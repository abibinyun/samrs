package routes

import (
	"github.com/gin-gonic/gin"
)

func registerBrandRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/brands", requirePerm(deps, "brand:create"), handlers.Brand.Create)
	v1.GET("/brands", requirePerm(deps, "brand:read"), handlers.Brand.GetAll)
	v1.GET("/brands/:id", requirePerm(deps, "brand:read"), handlers.Brand.GetByID)
	v1.PATCH("/brands/:id", requirePerm(deps, "brand:update"), handlers.Brand.Update)
	v1.DELETE("/brands/:id", requirePerm(deps, "brand:delete"), handlers.Brand.Delete)
}
