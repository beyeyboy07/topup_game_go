package controllers

import (
	// "encoding/json"
	"net/http"

	"topup_games_go/models"
	"topup_games_go/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// GameController menyimpan dependency
// yang dibutuhkan oleh Game Controller.
//
// Saat ini kita membutuhkan koneksi database.
type GameController struct {
	DB *gorm.DB
}

// CreateGame membuat game baru.
// POST /api/games
func (controller *GameController) CreateGame(c *gin.Context) {

	// Request body.
	var request struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	// Bind JSON dari request body ke struct.
	if err := c.ShouldBindJSON(&request); err != nil {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	// Validasi sederhana.
	if request.Name == "" || request.Code == "" {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Name and code are required",
		)

		return
	}

	// Buat object Game.
	game := models.Game{
		Name: request.Name,
		Code: request.Code,
	}

	// Simpan ke database.
	result := controller.DB.Create(&game)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Convert database model ke response model.
	gameResponse := game.ToResponse()

	// Response.
	utils.Success(
		c.Writer,
		http.StatusCreated,
		"Game created successfully",
		gameResponse,
	)
}

// GetGames mengambil semua game.
// GET /api/games
func (controller *GameController) GetGames(c *gin.Context) {

	// Menampung semua game dari database.
	var games []models.Game

	// Ambil semua data.
	result := controller.DB.Find(&games)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to get games",
		)

		return
	}

	// Convert Game -> GameResponse.
	gameResponses := make(
		[]models.GameResponse,
		0,
		len(games),
	)

	for _, game := range games {

		gameResponses = append(
			gameResponses,
			game.ToResponse(),
		)
	}

	// Response.
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Games retrieved successfully",
		gameResponses,
	)
}


// GET /api/games/1
func (controller *GameController) GetGameByID(c *gin.Context) {

	// ========================================
	// 1. Ambil ID dari parameter URL
	// ========================================
	id := c.Param("id")

	// ========================================
	// 2. Siapkan object Game
	// ========================================
	var game models.Game

	// ========================================
	// 3. Cari game berdasarkan ID
	// ========================================
	result := controller.DB.First(&game, id)

	// ========================================
	// 4. Jika game tidak ditemukan
	// ========================================
	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusNotFound,
			"Game not found",
		)

		return
	}

	// ========================================
	// 5. Convert Game menjadi GameResponse
	// ========================================
	gameResponse := game.ToResponse()

	// ========================================
	// 6. Response
	// ========================================
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Game retrieved successfully",
		gameResponse,
	)
}

// UpdateGame mengubah data game.
// PUT /api/games/:id
func (controller *GameController) UpdateGame(c *gin.Context) {

	// Ambil ID dari parameter URL.
	//
	// Contoh:
	// /api/games/1
	//
	// c.Param("id") = "1"
	id := c.Param("id")

	// Cari game berdasarkan ID.
	var game models.Game

	result := controller.DB.First(&game, id)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusNotFound,
			"Game not found",
		)

		return
	}

	// Request body.
	var request struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	// Decode JSON menggunakan Gin.
	if err := c.ShouldBindJSON(&request); err != nil {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	// Validasi.
	if request.Name == "" || request.Code == "" {

		utils.Error(
			c.Writer,
			http.StatusBadRequest,
			"Name and code are required",
		)

		return
	}

	// Update data.
	game.Name = request.Name
	game.Code = request.Code

	// Simpan perubahan.
	result = controller.DB.Save(&game)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Convert ke response model.
	gameResponse := game.ToResponse()

	// Response.
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Game updated successfully",
		gameResponse,
	)
}

// DeleteGame menghapus game berdasarkan ID.
// DELETE /api/games/:id
func (controller *GameController) DeleteGame(c *gin.Context) {

	// Ambil ID dari parameter URL.
	id := c.Param("id")

	// Cari game.
	var game models.Game

	result := controller.DB.First(&game, id)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusNotFound,
			"Game not found",
		)

		return
	}

	// Soft delete.
	result = controller.DB.Delete(&game)

	if result.Error != nil {

		utils.Error(
			c.Writer,
			http.StatusInternalServerError,
			result.Error.Error(),
		)

		return
	}

	// Response.
	utils.Success(
		c.Writer,
		http.StatusOK,
		"Game deleted successfully",
		nil,
	)
}