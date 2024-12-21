package model

type VenueResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
}

type CreateVenueRequest struct {
	Name     string `json:"name" validate:"required,lte=100"`
	Address  string `json:"address" validate:"required,lte=255"`
	Capacity int    `json:"capacity" validate:"omitempty"`
	City     string `json:"city" validate:"required,lte=100"`
	State    string `json:"state" validate:"required,lte=100"`
	Zip      string `json:"zip" validate:"required,lte=10"`
}

type UpdateVenueRequest struct {
	ID       uint   `param:"id" validate:"required"`
	Name     string `json:"name" validate:"omitempty,lte=100"`
	Address  string `json:"address" validate:"omitempty,lte=255"`
	Capacity int    `json:"capacity" validate:"omitempty"`
	City     string `json:"city" validate:"omitempty,lte=100"`
	State    string `json:"state" validate:"omitempty,lte=100"`
	Zip      string `json:"zip" validate:"omitempty,lte=10"`
}

type GetVenueRequest struct {
	ID uint `param:"id" validate:"required"`
}

type VenuesRequest struct {
	Page  int    `query:"page" validate:"numeric,min=1"`
	Size  int    `query:"size" validate:"numeric,min=1,max=100"`
	Sort  string `query:"sort" validate:"omitempty,oneof=id name address capacity city state zip created_at updated_at"`
	Order string `query:"order" validate:"omitempty"`
}

type VenueSearchRequest struct {
	Name     string `query:"name" validate:"omitempty,lte=100"`
	Address  string `query:"address" validate:"omitempty,lte=255"`
	Capacity int    `query:"capacity" validate:"omitempty"`
	City     string `query:"city" validate:"omitempty,lte=100"`
	State    string `query:"state" validate:"omitempty,lte=100"`
	Zip      string `query:"zip" validate:"omitempty,lte=10"`
	Page     int    `query:"page" validate:"numeric,min=1"`
	Size     int    `query:"size" validate:"numeric,min=1,max=100"`
	Sort     string `query:"sort" validate:"omitempty,oneof=id name address capacity city state zip created_at updated_at"`
	Order    string `query:"order" validate:"omitempty"`
}

type VenueQueryOptions struct {
	ID       *uint
	Name     *string
	Address  *string
	Capacity *int
	City     *string
	State    *string
	Zip      *string
	Page     int
	Size     int
	Sort     string
	Order    string
}
