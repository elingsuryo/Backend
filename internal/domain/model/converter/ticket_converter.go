package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

func ticketOrderToResponse(order *entity.Order, eventID *uint) *model.OrderResponse {
	if order == nil || order.ID == 0 {
		return nil
	}

	quantity := 1
	return &model.OrderResponse{
		ID:         order.ID,
		EventID:    eventID,
		UserID:     order.UserID,
		Quantity:   &quantity,
		TotalPrice: &order.TotalPrice,
	}
}

func TicketEntityToResponse(ticket *entity.Ticket) *model.TicketResponse {
	return &model.TicketResponse{
		ID:         ticket.ID,
		EventID:    ticket.EventID,
		OrderID:    helper.UintOrZero(ticket.OrderID),
		Price:      ticket.Price,
		Type:       ticket.Type,
		SeatNumber: ticket.SeatNumber,
		Event:      EventEntityToResponse(&ticket.Event),
		Order:      ticketOrderToResponse(&ticket.Order, &ticket.EventID),
	}
}

func TicketsToResponses(tickets []*entity.Ticket) []*model.TicketResponse {
	ticketsResponses := make([]*model.TicketResponse, len(tickets))
	for i := range tickets {
		ticketsResponses[i] = TicketEntityToResponse(tickets[i])
	}
	return ticketsResponses
}

func TicketsToPaginatedResponse(tickets []*entity.Ticket, totalItems int64, page, size int) *model.Response[[]*model.TicketResponse] {
	ticketsResponse := TicketsToResponses(tickets)

	if len(tickets) > 0 && tickets[0].Metadata != nil {
		if total, ok := tickets[0].Metadata["total_count"].(int64); ok {
			totalItems = total
		}
	}

	totalPages := 1
	if size > 0 {
		totalPages = (int(totalItems) + size - 1) / size
	}

	return model.NewResponse(ticketsResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
