package payment

import (
	"context"
	"gorm.io/gorm"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type PaymentRepository interface {
	GetByTransactionID(ctx context.Context, transactionID string) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, payment *model.PaymentUpdateRequest) error
	Find(db *gorm.DB, filter *model.PaymentQueryOptions) ([]*entity.Payment, error)
}
