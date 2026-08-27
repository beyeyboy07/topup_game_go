package middleware

import (
	"net/http"
	"strings"

	"topup_games_go/utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			utils.Error(c.Writer, http.StatusUnauthorized, "Authorization bearer token is required")
			c.Abort()
			return
		}
		userID, role, err := utils.ParseToken(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			utils.Error(c.Writer, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.Error(c.Writer, http.StatusUnauthorized, "Authentication is required")
			c.Abort()
			return
		}
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}
		utils.Error(c.Writer, http.StatusForbidden, "You do not have permission for this action")
		c.Abort()
	}
}
