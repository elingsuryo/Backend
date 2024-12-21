package user

import (
	"errors"
	"net/http"

	"github.com/TrinityKnights/Backend/internal/delivery/http/handler"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/service/user"
	domainErrors "github.com/TrinityKnights/Backend/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandlerImpl struct {
	Log  *logrus.Logger
	User user.UserService
}

func NewUserHandler(log *logrus.Logger, userService user.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		Log:  log,
		User: userService,
	}
}

// Register function is a handler to register a new user
// @Summary Register a new user
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User data"
// @Success 201 {object} model.Response[model.UserResponse]
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users [post]
func (h *UserHandlerImpl) Register(ctx echo.Context) error {
	request := new(model.RegisterRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.Register(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to register: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrEmailAlreadyExists):
			return handler.HandleError(ctx, http.StatusConflict, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, model.NewResponse(response, nil))
}

// Login function is a handler to login user
// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "User data"
// @Success 200 {object} model.Response[model.TokenResponse]
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users/login [post]
func (h *UserHandlerImpl) Login(ctx echo.Context) error {
	request := new(model.LoginRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.Login(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to login: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrUnauthorized):
			return handler.HandleError(ctx, http.StatusUnauthorized, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// Profile function is a handler to get user profile
// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[model.UserResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /users [get]
func (h *UserHandlerImpl) Profile(ctx echo.Context) error {
	response, err := h.User.Profile(ctx.Request().Context())
	if err != nil {
		h.Log.Errorf("failed to get profile: %v", err)
		switch {
		case errors.Is(err, errors.New(http.StatusText(http.StatusBadRequest))):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// Update function is a handler to update user profile
// @Summary Update user profile
// @Description Update user profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.UpdateRequest true "User data"
// @Success 200 {object} model.Response[model.UserResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @security ApiKeyAuth
// @Router /users [put]
func (h *UserHandlerImpl) Update(ctx echo.Context) error {
	request := new(model.UpdateRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.Update(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to update: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// RefreshToken function is a handler to refresh token
// @Summary Refresh token
// @Description Refresh token
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.RefreshTokenRequest true "User data"
// @Success 200 {object} model.Response[model.TokenResponse]
// @Failure 400 {object} model.Error
// @Failure 401 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users/refresh [post]
func (h *UserHandlerImpl) RefreshToken(ctx echo.Context) error {
	request := new(model.RefreshTokenRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.RefreshToken(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to refresh token: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrUnauthorized):
			return handler.HandleError(ctx, http.StatusUnauthorized, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// RequestResetPassword function is a handler to request reset password via email
// @Summary Request reset password via email
// @Description Request reset password via email
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.ReqResetPasswordRequest true "User data"
// @Success 200 {object} model.Response[model.VerifyResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users/request-reset [post]
func (h *UserHandlerImpl) RequestReset(ctx echo.Context) error {
	request := new(model.ReqResetPasswordRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.RequestReset(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to request reset password: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// ResetPassword function is a handler to reset password
// @Summary Reset password
// @Description Reset password
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.ResetPasswordRequest true "User data"
// @Success 200 {object} model.Response[model.VerifyResponse]
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users/reset-password [post]
func (h *UserHandlerImpl) ResetPassword(ctx echo.Context) error {
	request := new(model.ResetPasswordRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	response, err := h.User.ResetPassword(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to reset password: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, model.NewResponse(response, nil))
}

// VerifyEmail function is a handler to verify email
// @Summary Verify email
// @Description Verify email
// @Tags user
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {string} string "Email verified"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /users/verify-email/{token} [get]
func (h *UserHandlerImpl) VerifyEmail(ctx echo.Context) error {
	request := new(model.VerifyRequest)
	if err := ctx.Bind(request); err != nil {
		h.Log.Errorf("failed to bind request: %v", err)
		return handler.HandleError(ctx, http.StatusBadRequest, domainErrors.ErrBadRequest)
	}

	_, err := h.User.VerifyEmail(ctx.Request().Context(), request)
	if err != nil {
		h.Log.Errorf("failed to verify email: %v", err)
		switch {
		case errors.Is(err, domainErrors.ErrBadRequest):
			return handler.HandleError(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domainErrors.ErrNotFound):
			return handler.HandleError(ctx, http.StatusNotFound, err)
		default:
			return handler.HandleError(ctx, http.StatusInternalServerError, err)
		}
	}

	return ctx.String(http.StatusOK, "Email verified")
}
