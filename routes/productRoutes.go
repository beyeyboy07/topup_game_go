package routes

import (
	"topup_games_go/controllers"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes mendaftarkan
// seluruh endpoint Product.
func SetupProductRoutes(
	router *gin.Engine,
	controller *controllers.ProductController,
) {

	// Group /api/products
	productRoutes := router.Group("/api/products")
	productRoutes.Use(AuthMiddleware)

	// POST /api/products
	productRoutes.POST(
		"",
		RequireAdmin,
		controller.CreateProduct,
	)

	// POST /api/products
	productRoutes.GET(
		"",
		controller.GetProducts,
	)

	productRoutes.GET(
		"/:id",
		controller.GetProductByID,
	)

	productRoutes.PUT(
		"/:id",
		RequireAdmin,
		controller.UpdateProduct,
	)

	productRoutes.DELETE(
		"/:id",
		RequireAdmin,
		controller.DeleteProduct,
	)
}
