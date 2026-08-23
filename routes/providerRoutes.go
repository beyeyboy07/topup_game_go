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

	providerRoutes.POST(
		"",
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
		controller.UpdateProvider,
	)

	providerRoutes.DELETE(
		"/:id",
		controller.DeleteProvider,
	)
}
