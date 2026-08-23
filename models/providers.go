package models

import "gorm.io/gorm"

// Provider adalah representasi tabel providers.
type Provider struct {
	gorm.Model

	// Nama provider.
	Name string `gorm:"type:varchar(100);not null"`

	// Kode unik provider.
	// Contoh: DIGIFLAZZ
	Code string `gorm:"type:varchar(50);uniqueIndex;not null"`

	// URL API provider.
	BaseURL string `gorm:"type:varchar(255);not null"`

	// API Key provider.
	ApiKey string `gorm:"type:varchar(255);not null"`

	// Status provider.
	// Contoh:
	// 1 = aktif
	// 0 = tidak aktif
	Status bool `gorm:"not null;default:true"`
}

// ProviderResponse digunakan untuk
// response API.
//
// ApiKey sengaja tidak dimasukkan ke response.
type ProviderResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	BaseURL string `json:"base_url"`
	Status  bool   `json:"status"`
}
