package models

import "gorm.io/gorm"

// Game adalah representasi tabel "games" di database.
//
// Setiap object Game akan mewakili satu record/baris
// di dalam tabel games.
type Game struct {

	// ID adalah primary key dari tabel games.
	//
	// gorm.Model sebenarnya sudah menyediakan ID,
	// CreatedAt, UpdatedAt, dan DeletedAt.
	// Tapi di tahap awal kita gunakan gorm.Model
	// supaya memahami fitur bawaan GORM.
	gorm.Model

	// Name menyimpan nama game.
	//
	// Contoh:
	// "Mobile Legends"
	// "Free Fire"
	Name string `gorm:"not null"`

	// Code menyimpan kode unik game.
	//
	// Contoh:
	// "ML"
	// "FF"
	//
	// uniqueIndex berarti nilai Code tidak boleh
	// sama dengan game lainnya.
	Code string `gorm:"type:varchar(50);uniqueIndex;not null"`
}

// GameResponse digunakan khusus untuk
// response API.
//
// Kita tidak langsung mengirim Game
// karena Game memiliki field database
// seperti CreatedAt dan DeletedAt.

type GameResponse struct {

	// ID game
	ID uint `json:"id"`

	// Nama game
	Name string `json:"name"`

	// Kode game
	Code string `json:"code"`
}

// ToResponse mengubah object Game
// menjadi object GameResponse.
//
// Tujuannya agar field database yang
// tidak perlu dikirim ke client tidak
// ikut keluar.
func (game *Game) ToResponse() GameResponse {

	return GameResponse{
		ID:   game.ID,
		Name: game.Name,
		Code: game.Code,
	}
}
