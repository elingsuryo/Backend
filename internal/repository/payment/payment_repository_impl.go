package payment

import (
	"context"
	"fmt"
	"strings"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewPaymentRepository(db *gorm.DB, log *logrus.Logger) PaymentRepository {
	return &PaymentRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *PaymentRepositoryImpl) GetByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePaymentStatus(ctx context.Context, payment *model.PaymentUpdateRequest) error {
	return r.DB.Model(&entity.Payment{}).Where("id = ?", payment.ID).Updates(payment).Error
}

func (r *PaymentRepositoryImpl) Find(db *gorm.DB, opts *model.PaymentQueryOptions) ([]*entity.Payment, error) {
	// First, get total count
	var totalCount int64
	countQuery := db.Model(&entity.Payment{})

	// Apply filters to count query
	if opts.ID != nil {
		countQuery = countQuery.Where("id = ?", *opts.ID)
	}
	if opts.OrderID != nil {
		countQuery = countQuery.Where("order_id = ?", *opts.OrderID)
	}
	if opts.Amount != nil {
		countQuery = countQuery.Where("amount = ?", *opts.Amount)
	}
	if opts.Status != nil {
		countQuery = countQuery.Where("status = ?", *opts.Status)
	}

	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Main query with preloads
	query := db.Model(&entity.Payment{}).
		Preload("Order") // Add Order preload

	// Apply same filters to main query
	if opts.ID != nil {
		query = query.Where("id = ?", *opts.ID)
	}
	if opts.OrderID != nil {
		query = query.Where("order_id = ?", *opts.OrderID)
	}
	if opts.Amount != nil {
		query = query.Where("amount = ?", *opts.Amount)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}

	// Apply pagination
	if opts.Page > 0 && opts.Size > 0 {
		offset := (opts.Page - 1) * opts.Size
		query = query.Offset(offset).Limit(opts.Size)
	}

	// Validate and apply sorting
	if opts.Sort != "" && opts.Order != "" {
		// Whitelist of allowed columns for sorting
		allowedColumns := map[string]bool{
			"id":         true,
			"order_id":   true,
			"amount":     true,
			"status":     true,
			"created_at": true,
			// Add other allowed columns as needed
		}

		// Whitelist of allowed order directions
		allowedOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		// Validate sort column and order
		if !allowedColumns[opts.Sort] {
			return nil, fmt.Errorf("invalid sort column: %s", opts.Sort)
		}
		if !allowedOrders[strings.ToLower(opts.Order)] {
			return nil, fmt.Errorf("invalid sort order: %s", opts.Order)
		}

		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, opts.Order))
	}

	var payments []*entity.Payment
	if err := query.Find(&payments).Error; err != nil {
		return nil, err
	}

	// Store total count in the first payment's metadata
	if len(payments) > 0 {
		payments[0].Metadata = map[string]interface{}{
			"total_count": totalCount,
		}
	}

	return payments, nil
}
