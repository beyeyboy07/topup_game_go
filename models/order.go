package models

import (
	"time"

	"gorm.io/gorm"
)

// Order adalah representasi tabel orders.
type Order struct {
	gorm.Model

	// User yang melakukan order.
	UserID uint `gorm:"not null;index"`

	// Product yang dibeli.
	ProductID uint `gorm:"not null;index"`

	// ID player game.
	PlayerID string `gorm:"type:varchar(100);not null"`

	// Server ID game.
	ServerID string `gorm:"type:varchar(100)"`

	// Nomor order dari sistem kita.
	OrderNumber string `gorm:"type:varchar(50);uniqueIndex;not null"`

	// Nama product saat order dibuat.
	ProductName string `gorm:"type:varchar(100);not null"`

	// Harga beli dari provider.
	BuyPrice int64 `gorm:"not null"`

	// Harga jual ke user.
	SellPrice int64 `gorm:"not null"`

	// Status order.
	//
	// PENDING
	// PROCESSING
	// SUCCESS
	// FAILED
	Status string `gorm:"type:varchar(30);not null;index"`

	// ID transaksi dari provider.
	ProviderTransactionID string `gorm:"type:varchar(100)"`

	// Waktu order expired.
	ExpiredAt *time.Time

	// Waktu pembayaran.
	PaidAt *time.Time
}
