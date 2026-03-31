package models

import (
	"time"

	"gorm.io/gorm"
)

// Status transaksi
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusPaid      TransactionStatus = "paid"
	StatusCancelled TransactionStatus = "cancelled"
)

// Transaction adalah model untuk tabel transactions di database
type Transaction struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint              `gorm:"not null"                 json:"user_id"`
	TotalPrice  float64           `gorm:"not null"                 json:"total_price"`
	Status      TransactionStatus `gorm:"default:pending;size:20"  json:"status"`
	Notes       string            `gorm:"type:text"                json:"notes"`
	CreatedAt   time.Time         `                                json:"created_at"`
	UpdatedAt   time.Time         `                                json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index"                    json:"-"`

	// Relasi — transaksi milik satu user
	User  User              `gorm:"foreignKey:UserID"            json:"user,omitempty"`

	// Relasi — satu transaksi bisa punya banyak item
	Items []TransactionItem `gorm:"foreignKey:TransactionID"     json:"items,omitempty"`
}

// TransactionItem adalah detail produk di dalam transaksi
type TransactionItem struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	TransactionID uint    `gorm:"not null"                 json:"transaction_id"`
	ProductID     uint    `gorm:"not null"                 json:"product_id"`
	Quantity      int     `gorm:"not null"                 json:"quantity"`
	Price         float64 `gorm:"not null"                 json:"price"`        // harga saat transaksi terjadi
	Subtotal      float64 `gorm:"not null"                 json:"subtotal"`     // quantity * price

	// Relasi
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// CreateTransactionInput dipakai saat POST /transactions
type CreateTransactionInput struct {
	UserID uint                      `json:"user_id" binding:"required"`
	Notes  string                    `json:"notes"`
	Items  []CreateTransactionItemInput `json:"items"   binding:"required,min=1"`
}

// CreateTransactionItemInput detail item dalam transaksi
type CreateTransactionItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity"   binding:"required,gt=0"`
}

// UpdateTransactionStatusInput dipakai saat PATCH /transactions/:id/status
type UpdateTransactionStatusInput struct {
	Status TransactionStatus `json:"status" binding:"required,oneof=pending paid cancelled"`
}
