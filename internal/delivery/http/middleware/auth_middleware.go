package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/jwt"
	"github.com/labstack/echo/v4"
)

const contextKey = "claims"

func AuthMiddleware(jwtService jwt.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/api/v1/graphql" {
				return next(c)
			}

			errMessage := func(message string) error {
				return echo.NewHTTPError(http.StatusUnauthorized, model.NewErrorResponse[any](http.StatusUnauthorized, message))
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errMessage("Missing authorization header")
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				return errMessage("Invalid authorization header")
			}

			claims, err := jwtService.ValidateToken(bearerToken[1])
			if err != nil {
				return errMessage("Invalid token")
			}

			c.Set(contextKey, claims)

			ctx := context.WithValue(c.Request().Context(), contextKey, claims)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
