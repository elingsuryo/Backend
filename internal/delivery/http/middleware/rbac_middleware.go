package middleware

import (
	"net/http"

	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func RBACMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errMessage := func(message string) error {
				return echo.NewHTTPError(http.StatusUnauthorized, model.NewErrorResponse[any](http.StatusUnauthorized, message))
			}
			claims := c.Get("claims").(*jwt.JWTClaims)
			userRole := claims.Role
			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return errMessage("Permission denied")
		}
	}
}
