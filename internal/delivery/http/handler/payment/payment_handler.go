package payment

import (
	"github.com/labstack/echo/v4"
)

type PaymentHandler interface {
	CallbackPayment(ctx echo.Context) error
	GetPaymentByID(ctx echo.Context) error
	GetPayments(ctx echo.Context) error
	SearchPayments(ctx echo.Context) error
}
