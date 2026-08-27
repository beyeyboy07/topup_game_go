package routes

import (
	"net/http"

	"topup_games_go/controllers"

	"github.com/gin-gonic/gin"
)

// SetupGameRoutes digunakan untuk
// mendaftarkan seluruh endpoint Game.
func SetupGameRoutes(
	router *gin.Engine,
	controller *controllers.GameController,
) {

	// ========================================
	// Group /api/games
	// ========================================
	//
	// Semua endpoint di dalam group ini
	// otomatis memiliki prefix:
	//
	// /api/games
	gameRoutes := router.Group("/api/games")
	gameRoutes.Use(AuthMiddleware)

	// GET /api/games
	gameRoutes.GET(
		"",
		controller.GetGames,
	)

	// POST /api/games
	gameRoutes.POST(
		"",
		RequireAdmin,
		controller.CreateGame,
	)

	// ========================================
	// GET /api/games/:id
	// ========================================
	//
	// Mengambil satu game berdasarkan ID.
	gameRoutes.GET(
		"/:id",
		controller.GetGameByID,
	)

	// PUT /api/games/:id
	gameRoutes.PUT(
		"/:id",
		RequireAdmin,
		controller.UpdateGame,
	)

	// DELETE /api/games/:id
	gameRoutes.DELETE(
		"/:id",
		RequireAdmin,
		controller.DeleteGame,
	)

	// ========================================
	// Contoh route sederhana
	// ========================================
	//
	// Ini hanya untuk menunjukkan bahwa
	// Gin menangani HTTP method secara otomatis.
	_ = http.MethodGet
}
