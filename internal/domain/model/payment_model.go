package model

const (
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodQRIS         = "qris"
	// @TODO: Add other payment methods as needed
)

type PaymentRequest struct {
	OrderID       uint    `json:"order_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}

type PaymentResponse struct {
	ID            uint    `json:"id"`
	OrderID       uint    `json:"order_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	PaymentURL    string  `json:"payment_url"`
}
