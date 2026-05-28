package routes

import (
	"github.com/gin-gonic/gin"
)

func registerAssetMutationRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/asset-mutations", requirePerm(deps, "asset_mutation:create"), handlers.AssetMutation.Create)
	v1.GET("/asset-mutations", requirePerm(deps, "asset_mutation:read"), handlers.AssetMutation.GetAll)
	v1.GET("/asset-mutations/:id", requirePerm(deps, "asset_mutation:read"), handlers.AssetMutation.GetByID)
}
