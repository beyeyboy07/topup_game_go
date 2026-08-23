package main

import (
	"fmt"

	"topup_games_go/config"
	"topup_games_go/controllers"
	"topup_games_go/models"
	"topup_games_go/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// ========================================
	// 1. Connect ke database
	// ========================================
	db, err := config.ConnectDatabase()

	if err != nil {

		fmt.Println(
			"Database connection failed:",
			err,
		)

		return
	}

	// ========================================
	// 2. Auto Migration
	// ========================================
	err = db.AutoMigrate(
		&models.Game{},
		&models.Product{},
		&models.Provider{},
		&models.Order{},
	)

	if err != nil {

		fmt.Println(
			"Migration failed:",
			err,
		)

		return
	}

	fmt.Println("Database migration success")

	// ========================================
	// 3. Buat Gin Router
	// ========================================
	router := gin.Default()

	// ========================================
	// 4. Buat Game Controller
	// ========================================
	gameController := &controllers.GameController{
		DB: db,
	}

	productContoller := &controllers.ProductController{
		DB: db,
	}

	ProviderController := &controllers.ProviderController{
		DB: db,
	}

	OrderController := &controllers.OrderController{
		DB: db,
	}

	// ========================================
	// 5. Register routes
	// ========================================
	routes.SetupGameRoutes(
		router,
		gameController,
	)

	routes.SetupProductRoutes(
		router,
		productContoller,
	)

	routes.SetupProviderRoutes(
		router,
		ProviderController,
	)

	routes.SetupOrderRoutes(
		router,
		OrderController,
	)

	// ========================================
	// 6. Jalankan server
	// ========================================
	fmt.Println(
		"Server running on http://localhost:9000",
	)

	err = router.Run(":9000")

	if err != nil {

		fmt.Println(
			"Server error:",
			err,
		)
	}
}
