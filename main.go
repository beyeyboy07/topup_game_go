package main

import (
	"fmt"

	"topup_games_go/config"
	"topup_games_go/controllers"
	_ "topup_games_go/docs"
	"topup_games_go/models"
	"topup_games_go/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Topup Games API
// @version 1.0
// @description REST API untuk mengelola game, product, provider, dan order top-up.
// @host localhost:9000
// @BasePath /
// @schemes http https

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
