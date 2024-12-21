package route

import (
	"net/http"

	graphql "github.com/TrinityKnights/Backend/internal/delivery/graph/handler"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/event"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/order"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/payment"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/ticket"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/user"
	"github.com/TrinityKnights/Backend/internal/delivery/http/handler/venue"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/route"
	"github.com/labstack/echo/v4"
	swagger "github.com/swaggo/echo-swagger"
)

type RouteConfig interface {
	PublicRoute() []route.Route
	PrivateRoute() []route.Route
	SwaggerRoutes()
}

type Config struct {
	App            *echo.Echo
	GraphQLHandler *graphql.GraphQLHandler
	UserHandler    *user.UserHandlerImpl
	VenueHandler   *venue.VenueHandlerImpl
	EventHandler   *event.EventHandlerImpl
	TicketHandler  *ticket.TicketHandlerImpl
	OrderHandler   *order.OrderHandlerImpl
	PaymentHandler *payment.PaymentHandlerImpl
}

func (c Config) PublicRoute() []route.Route {
	return []route.Route{
		{
			Method:  echo.POST,
			Path:    "/users",
			Handler: c.UserHandler.Register,
		},
		{
			Method:  echo.POST,
			Path:    "/users/login",
			Handler: c.UserHandler.Login,
		},
		{
			Method:  echo.POST,
			Path:    "/users/refresh",
			Handler: c.UserHandler.RefreshToken,
		},
		{
			Method:  echo.POST,
			Path:    "/users/request-reset",
			Handler: c.UserHandler.RequestReset,
		},
		{
			Method:  echo.POST,
			Path:    "/users/reset-password/:token",
			Handler: c.UserHandler.ResetPassword,
		},
		{
			Method:  echo.GET,
			Path:    "/users/verify-email/:token",
			Handler: c.UserHandler.VerifyEmail,
		},
		{
			Method:  echo.GET,
			Path:    "/events/:id",
			Handler: c.EventHandler.GetEventByID,
		},
		{
			Method:  echo.GET,
			Path:    "/events",
			Handler: c.EventHandler.GetAllEvents,
		},
		{
			Method:  echo.GET,
			Path:    "/events/search",
			Handler: c.EventHandler.SearchEvents,
		},
		{
			Method:  echo.GET,
			Path:    "/tickets/:id",
			Handler: c.TicketHandler.GetTicketByID,
		},
		{
			Method:  echo.GET,
			Path:    "/tickets",
			Handler: c.TicketHandler.GetAllTickets,
		},
		{
			Method:  echo.GET,
			Path:    "/tickets/search",
			Handler: c.TicketHandler.SearchTickets,
		},
		{
			Method:  echo.POST,
			Path:    "/payment/callback",
			Handler: c.PaymentHandler.CallbackPayment,
		},
	}
}

func (c Config) PrivateRoute() []route.Route {
	return []route.Route{
		{
			Method:  echo.GET,
			Path:    "/users",
			Handler: c.UserHandler.Profile,
			Roles:   []string{"buyer", "admin"},
		},
		{
			Method:  echo.PUT,
			Path:    "/users",
			Handler: c.UserHandler.Update,
			Roles:   []string{"buyer", "admin"},
		},
		{
			Method:  echo.POST,
			Path:    "/venues",
			Handler: c.VenueHandler.CreateVenue,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/venues",
			Handler: c.VenueHandler.GetAllVenues,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/venues/:id",
			Handler: c.VenueHandler.GetVenueByID,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.PUT,
			Path:    "/venues/:id",
			Handler: c.VenueHandler.UpdateVenue,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/venues/search",
			Handler: c.VenueHandler.SearchVenues,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.POST,
			Path:    "/events",
			Handler: c.EventHandler.CreateEvent,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.PUT,
			Path:    "/events/:id",
			Handler: c.EventHandler.UpdateEvent,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.POST,
			Path:    "/orders",
			Handler: c.OrderHandler.CreateOrder,
			Roles:   []string{"buyer", "admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/orders/:id",
			Handler: c.OrderHandler.GetOrderByID,
			Roles:   []string{"buyer", "admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/orders",
			Handler: c.OrderHandler.GetAllOrders,
			Roles:   []string{"buyer", "admin"},
		},
		{
			Method:  echo.POST,
			Path:    "/tickets",
			Handler: c.TicketHandler.CreateTicket,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.PUT,
			Path:    "/tickets/:id",
			Handler: c.TicketHandler.UpdateTicket,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/payment/:id",
			Handler: c.PaymentHandler.GetPaymentByID,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/payment",
			Handler: c.PaymentHandler.GetPayments,
			Roles:   []string{"admin"},
		},
		{
			Method:  echo.GET,
			Path:    "/payment/search",
			Handler: c.PaymentHandler.SearchPayments,
			Roles:   []string{"admin"},
		},
	}
}

func (c Config) SwaggerRoutes() {
	c.App.GET("/swagger/*", swagger.WrapHandler)

	c.App.GET("/", func(ctx echo.Context) error {
		return ctx.Redirect(301, "/swagger/index.html")
	})
}

func (c Config) NotFoundRoute() {
	c.App.Any("*", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusNotFound, model.NewErrorResponse[any](http.StatusNotFound, "Route not found"))
	})
}
