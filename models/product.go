package models

import "gorm.io/gorm"

type Product struct {

	// ID, CreatedAt, UpdatedAt, DeletedAt
	// dari GORM.
	gorm.Model

	// ID game yang memiliki product ini.
	GameID uint `gorm:"not null;index"`

	// Nama product.
	Name string `gorm:"type:varchar(100);not null"`

	// Kode product.
	Code string `gorm:"type:varchar(50);uniqueIndex;not null"`

	// Harga beli dari provider.
	BuyPrice int64 `gorm:"not null"`

	// Harga jual ke user.
	SellPrice int64 `gorm:"not null"`

	// Relasi ke Game.
	Game Game `gorm:"foreignKey:GameID"`
}

type ProductResponse struct {
	ID        uint   `json:"id"`
	GameID    uint   `json:"game_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	BuyPrice  int64  `json:"buy_price"`
	SellPrice int64  `json:"sell_price"`
}
