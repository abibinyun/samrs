package routes

import (
	"github.com/gin-gonic/gin"
)

func registerVendorRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/vendors", requirePerm(deps, "vendor:create"), handlers.Vendor.Create)
	v1.GET("/vendors", requirePerm(deps, "vendor:read"), handlers.Vendor.GetAll)
	v1.GET("/vendors/:id", requirePerm(deps, "vendor:read"), handlers.Vendor.GetByID)
	v1.PATCH("/vendors/:id", requirePerm(deps, "vendor:update"), handlers.Vendor.Update)
	v1.DELETE("/vendors/:id", requirePerm(deps, "vendor:delete"), handlers.Vendor.Delete)
}
