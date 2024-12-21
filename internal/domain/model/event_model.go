package model

import (
	"time"

	"github.com/TrinityKnights/Backend/pkg/helper"
)

type EventResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date"`
	Time        helper.SQLTime `json:"time"`
	VenueID     uint           `json:"venue_id"`
}

type CreateEventRequest struct {
	Name        string `json:"name" validate:"required,lte=100"`
	Description string `json:"description" validate:"required,lte=255"`
	Date        string `json:"date" validate:"required" example:"2024-03-20"`
	Time        string `json:"time" validate:"required" example:"14:30:00"`
	VenueID     uint   `json:"venue_id" validate:"required"`
}

type UpdateEventRequest struct {
	ID          uint   `param:"id" validate:"required"`
	Name        string `json:"name" validate:"omitempty,lte=100"`
	Description string `json:"description" validate:"omitempty,lte=255"`
	Date        string `json:"date" validate:"omitempty" example:"2024-03-20"`
	Time        string `json:"time" validate:"omitempty" example:"14:30:00"`
	VenueID     uint   `json:"venue_id" validate:"omitempty"`
}

type GetEventRequest struct {
	ID uint `param:"id" validate:"required"`
}

type EventsRequest struct {
	Page  int    `query:"page" validate:"numeric,omitempty,gte=1"`
	Size  int    `query:"size" validate:"numeric,omitempty,gte=1,lte=100"`
	Sort  string `query:"sort" validate:"omitempty,oneof=id name description date time venue_id created_at updated_at"`
	Order string `query:"order" validate:"omitempty"`
}

type EventSearchRequest struct {
	Name        string `query:"name" validate:"omitempty,lte=100"`
	Description string `query:"description" validate:"omitempty,lte=255"`
	Date        string `query:"date" validate:"omitempty,datetime=2006-01-02"`
	Time        string `query:"time" validate:"omitempty,datetime=15:04:05"`
	VenueID     uint   `query:"venue_id" validate:"omitempty"`
	Page        int    `query:"page" validate:"numeric,omitempty,gte=1"`
	Size        int    `query:"size" validate:"numeric,omitempty,gte=1,lte=100"`
	Sort        string `query:"sort" validate:"omitempty,oneof=id name description date time venue_id created_at updated_at"`
	Order       string `query:"order" validate:"omitempty"`
}

type EventQueryOptions struct {
	ID          *uint
	Name        *string
	Description *string
	Date        *string
	Time        *string
	VenueID     *uint
	Page        int
	Size        int
	Sort        string
	Order       string
}
