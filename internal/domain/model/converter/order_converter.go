package converter

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/pkg/helper"
)

func OrderEntityToResponse(order *entity.Order) *model.OrderResponse {
	if order == nil {
		return nil
	}

	quantity := len(order.Tickets)
	eventID := uint(0)
	if len(order.Tickets) > 0 {
		eventID = order.Tickets[0].EventID
	}

	response := &model.OrderResponse{
		ID:         order.ID,
		EventID:    &eventID,
		UserID:     order.UserID,
		Quantity:   &quantity,
		TotalPrice: &order.TotalPrice,
		Date:       helper.FormatDate(order.Date),
	}

	// Only add tickets if they exist
	if len(order.Tickets) > 0 {
		tickets := make([]model.TicketResponse, len(order.Tickets))
		for i, ticket := range order.Tickets {
			tickets[i] = model.TicketResponse{
				ID:         ticket.ID,
				EventID:    ticket.EventID,
				OrderID:    helper.UintOrZero(ticket.OrderID),
				Price:      ticket.Price,
				Type:       ticket.Type,
				SeatNumber: ticket.SeatNumber,
			}
		}
		response.Tickets = &tickets
	}

	return response
}

func OrdersToResponses(orders []entity.Order) []*model.OrderResponse {
	orderResponses := make([]*model.OrderResponse, len(orders))
	for i := range orders {
		orderResponses[i] = OrderEntityToResponse(&orders[i])
	}
	return orderResponses
}

func OrdersToPaginatedResponse(orders []entity.Order, totalItems int64, page, size int) *model.Response[[]*model.OrderResponse] {
	ordersResponse := OrdersToResponses(orders)
	totalPages := (int(totalItems) + size - 1) / size

	return model.NewResponse(ordersResponse, &model.PageMetadata{
		Page:       page,
		Size:       size,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	})
}
