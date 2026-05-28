package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(v1 *gin.RouterGroup, handlers Handlers) {
	v1.GET("/auth/me", handlers.Auth.Me)

	v1.GET("/secure-ping", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")

		c.JSON(http.StatusOK, gin.H{
			"message":   "Authorized!",
			"user_id":   userID,
			"tenant_id": tenantID,
		})
	})
}
