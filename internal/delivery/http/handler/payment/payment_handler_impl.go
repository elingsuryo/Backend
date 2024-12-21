package payment

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaymentHandlerImpl struct {
	Viper          *viper.Viper
	Log            *logrus.Logger
	PaymentService payment.PaymentService
}

func NewPaymentHandler(v *viper.Viper, log *logrus.Logger, paymentService payment.PaymentService) PaymentHandler {
	return &PaymentHandlerImpl{
		Viper:          v,
		Log:            log,
		PaymentService: paymentService,
	}
}

// @Summary Callback Payment
// @Description Callback Payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param payment body model.PaymentCallbackRequest true "Payment Callback Request"
// @Success 200 {object} model.PaymentCallbackResponse
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /payment/callback [post]
func (h *PaymentHandlerImpl) CallbackPayment(ctx echo.Context) error {
	request := new(model.PaymentCallbackRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	myCallbackToken := h.Viper.GetString("XENDIT_CALLBACK_TOKEN")
	callbackToken := ctx.Request().Header.Get("x-callback-token")

	if callbackToken != myCallbackToken {
		h.Log.Errorf("invalid webhook id or callback token")
		return handler.HandleError(ctx, http.StatusUnauthorized, errors.New("invalid webhook id or callback token"))
	}

	response, err := h.PaymentService.Callback(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to callback payment: %v", err)
		return handler.HandleError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// @Summary Get Payment By ID
// @Description Get Payment By ID
// @Tags Payment
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} model.Response[model.CreatePaymentResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Security ApiKeyAuth
// @Router /payment/{payment_id} [get]
func (h *PaymentHandlerImpl) GetPaymentByID(ctx echo.Context) error {
	request := new(model.GetPaymentRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.PaymentService.GetPaymentByID(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get payment by id: %v", err)
		switch err {
		case domainErrors.ErrBadRequest:
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case domainErrors.ErrNotFound:
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// @Summary Get Payments
// @Description Get Payments
// @Tags Payment
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(id,order_id,amount,status)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.PaymentResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Security ApiKeyAuth
// @Router /payment [get]
func (h *PaymentHandlerImpl) GetPayments(ctx echo.Context) error {
	request := new(model.PaymentsRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.PaymentService.GetPayments(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get payments: %v", err)
		switch err {
		case domainErrors.ErrBadRequest:
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case domainErrors.ErrNotFound:
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

// @Summary Search Payments
// @Description Search Payments
// @Tags Payment
// @Accept json
// @Produce json
// @Param id query int false "Payment ID"
// @Param order_id query int false "Order ID"
// @Param amount query float64 false "Amount"
// @Param status query string false "Status"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param sort query string false "Sort field" Enums(id,order_id,amount,status)
// @Param order query string false "Sort order"
// @Success 200 {object} model.Response[[]model.PaymentResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Security ApiKeyAuth
// @Router /payment/search [get]
func (h *PaymentHandlerImpl) SearchPayments(ctx echo.Context) error {
	request := new(model.PaymentSearchRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, err)
	}

	response, err := h.PaymentService.SearchPayments(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to search payments: %v", err)
		switch err {
		case domainErrors.ErrBadRequest:
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case domainErrors.ErrNotFound:
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, response)
}
