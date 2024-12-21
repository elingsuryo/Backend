package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func PaymentToEntityResponse(payment *entity.Payment) *model.PaymentResponse {
	return &model.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Method:        payment.Method,
		TransactionID: payment.TransactionID,
		Amount:        payment.Amount,
		Status:        string(payment.Status),
		Order:         OrderEntityToResponse(&payment.Order),
	}
}

func CreatePaymentEntityToResponse(payment *entity.Payment) *model.CreatePaymentResponse {
	return &model.CreatePaymentResponse{
		ID:         payment.ID,
		OrderID:    payment.OrderID,
		Amount:     payment.Amount,
		Status:     string(payment.Status),
		PaymentURL: payment.TransactionID,
	}
}

func PaymentsToResponses(payments []*entity.Payment) []*model.PaymentResponse {
	responses := make([]*model.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = PaymentToEntityResponse(payment)
	}
	return responses
}

func PaymentsToPaginatedResponse(payments []*entity.Payment, totalItems int64, page, size int) *model.Response[[]*model.PaymentResponse] {
	paymentsResponse := PaymentsToResponses(payments)

	if len(payments) > 0 && payments[0].Metadata != nil {
		if total, ok := payments[0].Metadata["total_count"].(int64); ok {
			totalItems = total
		}
	}

	totalPages := 1
	if size > 0 {
		totalPages = (int(totalItems) + size - 1) / size
	}

	return model.NewResponse(paymentsResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
