package routes

import (
	"topup_games_go/controllers"

	"github.com/gin-gonic/gin"
)

// SetupProviderRoutes mendaftarkan
// endpoint Provider.

func SetupProviderRoutes(router *gin.Engine, controller *controllers.ProviderController) {

	// Group /api/providers

	providerRoutes := router.Group("/api/providers")
	providerRoutes.Use(AuthMiddleware)

	providerRoutes.POST(
		"",
		RequireAdmin,
		controller.CreateProvider,
	)

	providerRoutes.GET(
		"",
		controller.GetProviders,
	)

	providerRoutes.GET(
		"/:id",
		controller.GetProviderByID,
	)

	providerRoutes.PUT(
		"/:id",
		RequireAdmin,
		controller.UpdateProvider,
	)

	providerRoutes.DELETE(
		"/:id",
		RequireAdmin,
		controller.DeleteProvider,
	)
}
