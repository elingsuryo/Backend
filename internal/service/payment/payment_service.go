package payment

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/model"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreateInvoice(ctx context.Context, tx *gorm.DB, request *model.CreatePaymentRequest) (*model.CreatePaymentResponse, error)
	Callback(ctx context.Context, request *model.PaymentCallbackRequest) (*model.PaymentCallbackResponse, error)
	GetPaymentByID(ctx context.Context, request *model.GetPaymentRequest) (*model.PaymentResponse, error)
	GetPayments(ctx context.Context, request *model.PaymentsRequest) (*model.Response[[]*model.PaymentResponse], error)
	SearchPayments(ctx context.Context, request *model.PaymentSearchRequest) (*model.Response[[]*model.PaymentResponse], error)
}
