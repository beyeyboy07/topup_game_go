package routes

import (
	"topup_games_go/controllers"

	"github.com/gin-gonic/gin"
)

// SetupGameRoutes digunakan untuk
// mendaftarkan seluruh endpoint Game.

func SetupOrderRoutes(router *gin.Engine, controller *controllers.OrderController) {

	// Group /api/orders
	orderRoutes := router.Group("/api/orders")

	orderRoutes.POST(
		"",
		controller.CreateOrder,
	)
}
