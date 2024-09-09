package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint    `gorm:"primary_key" json:"id"` // Use uint for ID
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UserID      uint    `json:"user_id"` // Use uint as foreign key to User
	User        User    `gorm:"foreignKey:UserID"`
}
