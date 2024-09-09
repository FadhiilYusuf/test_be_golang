package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID        uint    `gorm:"primary_key" json:"id"`
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"`
	UserID    uint    `json:"user_id"`
	User      User    `gorm:"foreignKey:UserID"`
	Status    string  `json:"status" gorm:"default:'diproses'"` // Default to 'diproses'
	Quantity  uint    `json:"quantity"`
	Total     float64 `json:"total"`
}

// GORM hook: BeforeCreate to set default status
func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if order.Status == "" {
		order.Status = "diproses" // Default status
	}
	return
}
