package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null;size:255"        json:"name"`
	Description string         `gorm:"type:text"                json:"description"`
	Price       float64        `gorm:"not null"                 json:"price"`
	Stock       int            `gorm:"default:0"                json:"stock"`
	Category    string         `gorm:"size:100"                 json:"category"`
	CreatedAt   time.Time      `                                json:"created_at"`
	UpdatedAt   time.Time      `                                json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                    json:"-"`
}

type CreateProductInput struct {
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"       binding:"gte=0"`
	Category    string  `json:"category"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"  binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock"  binding:"omitempty,gte=0"`
	Category    *string  `json:"category"`
}
