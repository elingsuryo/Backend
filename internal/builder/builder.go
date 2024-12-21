package builder

import (
	graphql "github.com/TrinityKnights/Backend/internal/delivery/graph/handler"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/event"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/order"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/payment"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/ticket"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/venue"
	rbac "github.com/TrinityKnights/Backend/internal/delivery/http/middleware"
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
	TicketHandler  *ticket.TicketHandlerImpl
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

	// Private routes with auth and rbac middleware
	privateGroup := g.Group("", c.AuthMiddleware)
	for _, r := range c.Routes.PrivateRoute() {
		rbacMiddleware := rbac.RBACMiddleware(r.Roles)
		privateGroup.Add(r.Method, r.Path, rbacMiddleware(r.Handler))
	}

	// GraphQL routes
	privateGroup.POST("/graphql", c.GraphQLHandler.GraphQLHandler)

	// GraphQL Playground (you might want to restrict this in production)
	c.App.GET("/playground", c.GraphQLHandler.PlaygroundHandler)

	// Swagger routes
	c.Routes.SwaggerRoutes()

	// Not found route
	c.Routes.NotFoundRoute()
}
