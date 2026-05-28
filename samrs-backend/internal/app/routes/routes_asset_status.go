package routes

import (
	"github.com/gin-gonic/gin"
)

func registerStatusRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/asset-statuses", requirePerm(deps, "asset_status:create"), handlers.Status.Create)
	v1.GET("/asset-statuses", requirePerm(deps, "asset_status:read"), handlers.Status.GetAll)
	v1.GET("/asset-statuses/:id", requirePerm(deps, "asset_status:read"), handlers.Status.GetByID)
	v1.PATCH("/asset-statuses/:id", requirePerm(deps, "asset_status:update"), handlers.Status.Update)
	v1.DELETE("/asset-statuses/:id", requirePerm(deps, "asset_status:delete"), handlers.Status.Delete)
}
