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
	orderRoutes.Use(AuthMiddleware)

	orderRoutes.POST(
		"",
		controller.CreateOrder,
	)

	orderRoutes.GET(
		"/:id",
		controller.GetOrderByID,
	)

	orderRoutes.GET(
		"",
		controller.GetOrders,
	)

	// PUT /api/orders/:id/status
	orderRoutes.PUT(
		"/:id/status",
		RequireAdmin,
		controller.UpdateOrderStatus,
	)
}
