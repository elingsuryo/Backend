package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository/payment"
	"github.com/TrinityKnights/Backend/pkg/cache"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/TrinityKnights/Backend/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	xendit "github.com/xendit/xendit-go/v6"
	invoice "github.com/xendit/xendit-go/v6/invoice"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	DB                *gorm.DB
	Cache             *cache.ImplCache
	Log               *logrus.Logger
	Validate          *validator.Validate
	PaymentRepository payment.PaymentRepository
	Xendit            *xendit.APIClient
	helper            *helper.ContextHelper
}

func NewPaymentServiceImpl(db *gorm.DB, cacheImpl *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, paymentRepository payment.PaymentRepository, xendit *xendit.APIClient) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		DB:                db,
		Cache:             cacheImpl,
		Log:               log,
		Validate:          validate,
		PaymentRepository: paymentRepository,
		Xendit:            xendit,
		helper:            helper.NewContextHelper(),
	}
}

func (s *PaymentServiceImpl) CreateInvoice(ctx context.Context, tx *gorm.DB, request *model.PaymentRequest) (*model.PaymentResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	// Get order details
	var order entity.Order
	if err := tx.Preload("User").First(&order, request.OrderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get order: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Verify amount matches order total
	if order.TotalPrice != request.Amount {
		return nil, domainErrors.ErrInvalidAmount
	}

	description := fmt.Sprintf("Payment for Order #%d", order.ID)
	currency := "IDR"
	shouldSendEmail := true

	// Create Xendit invoice
	createInvoiceRequest := invoice.CreateInvoiceRequest{
		ExternalId:      fmt.Sprintf("order_%d", order.ID),
		Amount:          float64(request.Amount),
		PayerEmail:      &order.User.Email,
		Description:     &description,
		Currency:        &currency,
		ShouldSendEmail: &shouldSendEmail,
	}

	invoice, _, err := s.Xendit.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()
	if err != nil {
		s.Log.Errorf("failed to create xendit invoice: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create payment record
	payment := &entity.Payment{
		OrderID:       order.ID,
		TransactionID: *invoice.Id,
		Method:        nil,
		Amount:        request.Amount,
		Status:        model.PaymentStatus(invoice.Status),
	}

	if err := tx.Create(payment).Error; err != nil {
		s.Log.Errorf("failed to create payment record: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.PaymentResponse{
		ID:         payment.ID,
		OrderID:    payment.OrderID,
		Amount:     request.Amount,
		Status:     string(invoice.Status),
		ExpiryDate: invoice.ExpiryDate.Format(time.RFC3339),
		PaymentURL: invoice.InvoiceUrl,
	}, nil
}

func (s *PaymentServiceImpl) Callback(ctx context.Context, request *model.PaymentCallbackRequest) (*model.PaymentCallbackResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Check if id is valid
	dataPayment, err := s.PaymentRepository.GetByTransactionID(ctx, request.ID)
	if err != nil {
		s.Log.Errorf("failed to get payment by transaction id: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	updatePayment := &model.PaymentUpdateRequest{
		ID:     dataPayment.ID,
		Method: *request.PaymentMethod,
		Status: model.PaymentStatus(request.Status),
	}

	// Update payment status
	if err := s.PaymentRepository.UpdatePaymentStatus(ctx, updatePayment); err != nil {
		s.Log.Errorf("failed to update payment status: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Errorf("failed to commit transaction: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.PaymentCallbackResponse{
		Status: string(updatePayment.Status),
	}, nil
}
