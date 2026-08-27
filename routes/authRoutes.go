package routes

import (
	"topup_games_go/controllers"
	appmiddleware "topup_games_go/middleware"

	"github.com/gin-gonic/gin"
)

var AuthMiddleware gin.HandlerFunc = func(c *gin.Context) {
	c.Next()
}

var RequireAdmin gin.HandlerFunc = func(c *gin.Context) {
	c.Next()
}

func SetupAuthMiddleware(router *gin.Engine, auth gin.HandlerFunc) {
	AuthMiddleware = auth
	RequireAdmin = appmiddleware.RequireRole("admin")
}

func SetupAuthRoutes(router *gin.Engine, controller *controllers.AuthController) {
	authRoutes := router.Group("/api/auth")
	authRoutes.POST("/register", controller.Register)
	authRoutes.POST("/login", controller.Login)
}
