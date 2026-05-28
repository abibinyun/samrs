package routes

import (
	"github.com/gin-gonic/gin"
)

func registerBedRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/beds", requirePerm(deps, "bed:create"), handlers.Bed.Create)
	v1.GET("/beds", requirePerm(deps, "bed:read"), handlers.Bed.GetAll)
	v1.GET("/beds/:id", requirePerm(deps, "bed:read"), handlers.Bed.GetByID)
	v1.PATCH("/beds/:id", requirePerm(deps, "bed:update"), handlers.Bed.Update)
	v1.DELETE("/beds/:id", requirePerm(deps, "bed:delete"), handlers.Bed.Delete)
}
