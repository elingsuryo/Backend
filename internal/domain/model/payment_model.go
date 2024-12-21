package model

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusExpired PaymentStatus = "EXPIRED"
)

type PaymentResponse struct {
	ID            uint           `json:"id"`
	OrderID       uint           `json:"order_id"`
	TransactionID string         `json:"transaction_id"`
	Amount        float64        `json:"amount"`
	Method        string         `json:"method"`
	Status        string         `json:"status"`
	Order         *OrderResponse `json:"order,omitempty"`
}

type CreatePaymentRequest struct {
	OrderID uint    `json:"order_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
}

type CreatePaymentResponse struct {
	ID         uint    `json:"id"`
	OrderID    uint    `json:"order_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	ExpiryDate string  `json:"expiry_date"`
	PaymentURL string  `json:"payment_url"`
}

type PaymentUpdateRequest struct {
	ID     uint          `json:"id" validate:"required"`
	Method string        `json:"method" validate:"required"`
	Status PaymentStatus `json:"status" validate:"required"`
}

type PaymentCallbackRequest struct {
	ID                 string  `json:"id"`
	ExternalID         string  `json:"external_id"`
	UserID             string  `json:"user_id"`
	IsHigh             bool    `json:"is_high"`
	PaymentMethod      *string `json:"payment_method,omitempty"`
	Status             string  `json:"status"`
	MerchantName       string  `json:"merchant_name"`
	Amount             int     `json:"amount"`
	BankCode           *string `json:"bank_code,omitempty"`
	PaidAmount         int     `json:"paid_amount"`
	PaidAt             *string `json:"paid_at,omitempty"`
	PayerEmail         *string `json:"payer_email,omitempty"`
	Description        string  `json:"description"`
	Created            string  `json:"created"`
	Updated            string  `json:"updated"`
	Currency           *string `json:"currency,omitempty"`
	PaymentChannel     *string `json:"payment_channel,omitempty"`
	PaymentDestination *string `json:"payment_destination,omitempty"`
	SuccessRedirectURL *string `json:"success_redirect_url,omitempty"`
	FailedRedirectURL  *string `json:"failed_redirect_url,omitempty"`
}

type PaymentCallbackResponse struct {
	Status string `json:"status"`
}

type PaymentQueryOptions struct {
	ID            *uint    `query:"id,omitempty"`
	OrderID       *uint    `query:"order_id,omitempty"`
	Method        *string  `query:"method,omitempty"`
	TransactionID *string  `query:"transaction_id,omitempty"`
	Amount        *float64 `query:"amount,omitempty"`
	Status        *string  `query:"status,omitempty"`
	Page          int      `query:"page,omitempty" validate:"omitempty,gte=1"`
	Size          int      `query:"size,omitempty" validate:"omitempty,gte=1,lte=100"`
	Sort          string   `query:"sort,omitempty"`
	Order         string   `query:"order,omitempty"`
}

type PaymentsRequest struct {
	Page  int    `query:"page" validate:"numeric,omitempty,gte=1"`
	Size  int    `query:"size" validate:"numeric,omitempty,gte=1,lte=100"`
	Sort  string `query:"sort" validate:"omitempty"`
	Order string `query:"order" validate:"omitempty"`
}

type PaymentSearchRequest struct {
	ID      uint    `query:"id" validate:"omitempty"`
	OrderID uint    `query:"order_id" validate:"omitempty"`
	Amount  float64 `query:"amount" validate:"omitempty"`
	Status  string  `query:"status" validate:"omitempty"`
	Page    int     `query:"page" validate:"numeric,omitempty,gte=1"`
	Size    int     `query:"size" validate:"numeric,omitempty,gte=1,lte=100"`
	Sort    string  `query:"sort" validate:"omitempty"`
	Order   string  `query:"order" validate:"omitempty"`
}

type GetPaymentRequest struct {
	ID uint `param:"id" validate:"required"`
}
