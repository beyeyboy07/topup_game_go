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

	// POST /api/products
	productRoutes.POST(
		"",
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
		controller.UpdateProduct,
	)

	productRoutes.DELETE(
		"/:id",
		controller.DeleteProduct,
	)
}
