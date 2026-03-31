package models

import (
	"time"

	"gorm.io/gorm"
)

// User adalah model untuk tabel users di database
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"not null;size:255"        json:"name"`
	Email     string         `gorm:"unique;not null;size:255" json:"email"`
	Password  string         `gorm:"not null"                 json:"-"` // json:"-" supaya password tidak tampil di response
	Phone     string         `gorm:"size:20"                  json:"phone"`
	Address   string         `gorm:"type:text"                json:"address"`
	CreatedAt time.Time      `                                json:"created_at"`
	UpdatedAt time.Time      `                                json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                    json:"-"`

	// Relasi — satu user bisa punya banyak transaksi
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

// CreateUserInput dipakai saat POST /users
type CreateUserInput struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// UpdateUserInput dipakai saat PATCH /users/:id
type UpdateUserInput struct {
	Name    *string `json:"name"`
	Email   *string `json:"email"  binding:"omitempty,email"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}
