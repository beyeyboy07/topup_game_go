package models

import "gorm.io/gorm"

// User adalah akun yang dapat mengakses API.
type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(150);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Role         string `gorm:"type:varchar(20);not null;default:customer"`
}
