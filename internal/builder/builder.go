package builder

import (
	graphql "github.com/TrinityKnights/Backend/internal/delivery/graph/handler"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/event"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/order"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/payment"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/venue"
	"github.com/TrinityKnights/Backend/internal/delivery/http/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	App            *echo.Echo
	GraphQLHandler *graphql.GraphQLHandler
	UserHandler    *user.UserHandlerImpl
	VenueHandler   *venue.VenueHandlerImpl
	EventHandler   *event.EventHandlerImpl
	OrderHandler   *order.OrderHandlerImpl
	PaymentHandler *payment.PaymentHandlerImpl
	AuthMiddleware echo.MiddlewareFunc
	Routes         *route.Config
}

func (c *Config) BuildRoutes() {
	// Global middleware
	c.App.Use(middleware.Recover())

	// API group
	g := c.App.Group("/api/v1")
	g.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(60)))

	// Public routes
	for _, r := range c.Routes.PublicRoute() {
		g.Add(r.Method, r.Path, r.Handler)
	}

	// Private routes with auth middleware
	privateGroup := g.Group("", c.AuthMiddleware)
	for _, r := range c.Routes.PrivateRoute() {
		privateGroup.Add(r.Method, r.Path, r.Handler)
	}

	// GraphQL routes with auth middleware
	graphqlGroup := g.Group("/graphql")
	graphqlGroup.Use(c.AuthMiddleware)
	graphqlGroup.POST("", c.GraphQLHandler.GraphQLHandler)
	c.App.GET("/playground", c.GraphQLHandler.PlaygroundHandler)

	// Swagger routes
	c.Routes.SwaggerRoutes()
}
