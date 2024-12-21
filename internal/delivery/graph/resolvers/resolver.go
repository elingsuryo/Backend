package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/TrinityKnights/Backend/internal/service/event"
	"github.com/TrinityKnights/Backend/internal/service/payment"
	"github.com/TrinityKnights/Backend/internal/service/ticket"
	"github.com/TrinityKnights/Backend/internal/service/user"
	"github.com/TrinityKnights/Backend/internal/service/venue"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

type Resolver struct {
	UserService    user.UserService
	EventService   event.EventService
	TicketService  ticket.TicketService
	VenueService   venue.VenueService
	PaymentService payment.PaymentService
	helper         helper.ContextHelper
}

func NewResolver(userService user.UserService, eventService event.EventService, ticketService ticket.TicketService, venueService venue.VenueService, paymentService payment.PaymentService) *Resolver {
	return &Resolver{
		UserService:    userService,
		EventService:   eventService,
		TicketService:  ticketService,
		VenueService:   venueService,
		PaymentService: paymentService,
		helper:         *helper.NewContextHelper(),
	}
}
