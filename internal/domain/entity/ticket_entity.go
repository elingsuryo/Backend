package entity

import "gorm.io/gorm"

type Ticket struct {
	ID         string                 `json:"id" gorm:"primaryKey"`
	EventID    uint                   `json:"event_id" gorm:"not null"`
	OrderID    *uint                  `json:"order_id,omitempty" gorm:"null"`
	Price      float64                `json:"price" gorm:"not null"`
	Type       string                 `json:"type" gorm:"not null"`
	SeatNumber string                 `json:"seat_number"`
	Event      Event                  `json:"event" gorm:"foreignKey:EventID"`
	Order      Order                  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Metadata   map[string]interface{} `gorm:"-"`
	gorm.Model
}

func (t *Ticket) TableName() string {
	return "tickets"
}
