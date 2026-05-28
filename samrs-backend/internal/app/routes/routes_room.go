package routes

import (
	"github.com/gin-gonic/gin"
)

func registerRoomRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/rooms", requirePerm(deps, "room:create"), handlers.Room.Create)
	v1.GET("/rooms", requirePerm(deps, "room:read"), handlers.Room.GetAll)
	v1.GET("/rooms/:id", requirePerm(deps, "room:read"), handlers.Room.GetByID)
	v1.PATCH("/rooms/:id", requirePerm(deps, "room:update"), handlers.Room.Update)
	v1.DELETE("/rooms/:id", requirePerm(deps, "room:delete"), handlers.Room.Delete)
}
