package routes

import (
	"github.com/gin-gonic/gin"
)

func registerNotificationRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/notifications/send", requirePerm(deps, "notification:send"), handlers.Notification.Send)
}
