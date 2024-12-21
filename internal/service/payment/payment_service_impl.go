package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/domain/model/converter"
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

func NewPaymentServiceImpl(db *gorm.DB, cacheImpl *cache.ImplCache, log *logrus.Logger, validate *validator.Validate, paymentRepository payment.PaymentRepository, x *xendit.APIClient) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		DB:                db,
		Cache:             cacheImpl,
		Log:               log,
		Validate:          validate,
		PaymentRepository: paymentRepository,
		Xendit:            x,
		helper:            helper.NewContextHelper(),
	}
}

func (s *PaymentServiceImpl) CreateInvoice(ctx context.Context, tx *gorm.DB, request *model.CreatePaymentRequest) (*model.CreatePaymentResponse, error) {
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

	i, _, err := s.Xendit.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()
	if err != nil {
		s.Log.Errorf("failed to create xendit invoice: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	// Create p record
	p := &entity.Payment{
		OrderID:       order.ID,
		TransactionID: *i.Id,
		Amount:        request.Amount,
		Status:        model.PaymentStatus(i.Status),
	}

	if err := tx.Create(p).Error; err != nil {
		s.Log.Errorf("failed to create payment record: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	return &model.CreatePaymentResponse{
		ID:         p.ID,
		OrderID:    p.OrderID,
		Amount:     request.Amount,
		Status:     string(i.Status),
		ExpiryDate: i.ExpiryDate.Format(time.RFC3339),
		PaymentURL: i.InvoiceUrl,
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

func (s *PaymentServiceImpl) GetPaymentByID(ctx context.Context, request *model.GetPaymentRequest) (*model.PaymentResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	key := fmt.Sprintf("payment:get:id:%d", request.ID)
	var cacheResponse *model.PaymentResponse
	if err := s.Cache.Get(key, &cacheResponse); err == nil {
		return cacheResponse, nil
	}

	tx := s.DB.WithContext(ctx)

	payments, err := s.PaymentRepository.Find(tx, &model.PaymentQueryOptions{
		ID: &request.ID,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get payment: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(payments) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.PaymentToEntityResponse(payments[0])

	if err := s.Cache.Set(key, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}

func (s *PaymentServiceImpl) GetPayments(ctx context.Context, request *model.PaymentsRequest) (*model.Response[[]*model.PaymentResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	opts := model.PaymentQueryOptions{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  request.Sort,
		Order: request.Order,
	}

	if opts.Size <= 0 {
		opts.Size = 10
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	cacheKey := fmt.Sprintf("payment:get:page:%d:size:%d:sort:%s:order:%s", opts.Page, opts.Size, opts.Sort, opts.Order)
	var cacheResponse model.Response[[]*model.PaymentResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	db := s.DB.WithContext(ctx)

	payments, err := s.PaymentRepository.Find(db, &opts)
	if err != nil {
		s.Log.Errorf("failed to get payments: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		return nil, domainErrors.ErrInternalServer
	}

	if len(payments) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.PaymentsToPaginatedResponse(payments, int64(len(payments)), opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}

func (s *PaymentServiceImpl) SearchPayments(ctx context.Context, request *model.PaymentSearchRequest) (*model.Response[[]*model.PaymentResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, domainErrors.ErrValidation
	}

	if request.Size <= 0 {
		request.Size = 10
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	cacheKey := fmt.Sprintf("payment:search:id:%d:order:%d:amount:%.2f:status:%s:page:%d:size:%d:sort:%s:order:%s",
		request.ID, request.OrderID, request.Amount, request.Status,
		request.Page, request.Size, request.Sort, request.Order)

	var cacheResponse model.Response[[]*model.PaymentResponse]
	if err := s.Cache.Get(cacheKey, &cacheResponse); err == nil {
		return &cacheResponse, nil
	}

	opts := model.PaymentQueryOptions{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  request.Sort,
		Order: request.Order,
	}

	if request.ID != 0 {
		opts.ID = &request.ID
	}
	if request.OrderID != 0 {
		opts.OrderID = &request.OrderID
	}
	if request.Amount != 0 {
		opts.Amount = &request.Amount
	}
	if request.Status != "" {
		opts.Status = &request.Status
	}

	db := s.DB.WithContext(ctx)

	payments, err := s.PaymentRepository.Find(db, &opts)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainErrors.ErrNotFound
		}
		s.Log.Errorf("failed to get payments: %v", err)
		return nil, domainErrors.ErrInternalServer
	}

	if len(payments) == 0 {
		return nil, domainErrors.ErrNotFound
	}

	response := converter.PaymentsToPaginatedResponse(payments, int64(len(payments)), opts.Page, opts.Size)

	if err := s.Cache.Set(cacheKey, response, 5*time.Minute); err != nil {
		s.Log.Errorf("failed to set cache: %v", err)
	}

	return response, nil
}
