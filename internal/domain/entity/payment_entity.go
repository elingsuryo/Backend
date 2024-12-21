package entity

import (
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID       uint                   `json:"order_id" gorm:"not null"`
	Method        string                 `json:"method" gorm:"null"`
	TransactionID string                 `json:"transaction_id" gorm:"not null"`
	Amount        float64                `json:"amount" gorm:"null"`
	Status        model.PaymentStatus    `json:"status" gorm:"null"`
	Order         Order                  `json:"order" gorm:"foreignKey:OrderID"`
	Metadata      map[string]interface{} `gorm:"-"`
	gorm.Model
}

func (p *Payment) TableName() string {
	return "payments"
}
