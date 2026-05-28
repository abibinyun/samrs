package routes

import (
	"github.com/gin-gonic/gin"
)

func registerAssetRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/assets", requirePerm(deps, "asset:create"), handlers.Asset.Create)
	v1.GET("/assets", requirePerm(deps, "asset:read"), handlers.Asset.GetAll)
	v1.GET("/assets/:id", requirePerm(deps, "asset:read"), handlers.Asset.GetByID)
	v1.GET("/assets/:id/timeline", requirePerm(deps, "asset:read"), handlers.AssetEvent.GetTimeline)
	v1.GET("/assets/:id/qr", requirePerm(deps, "asset:read"), handlers.Asset.GetQRCode)
	v1.PATCH("/assets/:id", requirePerm(deps, "asset:update"), handlers.Asset.Update)
	v1.DELETE("/assets/:id", requirePerm(deps, "asset:delete"), handlers.Asset.Delete)
}
