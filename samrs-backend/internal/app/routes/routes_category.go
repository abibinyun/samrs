package routes

import (
	"github.com/gin-gonic/gin"
)

func registerCategoryRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/categories", requirePerm(deps, "category:create"), handlers.Category.Create)
	v1.GET("/categories", requirePerm(deps, "category:read"), handlers.Category.GetAll)
	v1.GET("/categories/:id", requirePerm(deps, "category:read"), handlers.Category.GetByID)
	v1.PATCH("/categories/:id", requirePerm(deps, "category:update"), handlers.Category.Update)
	v1.DELETE("/categories/:id", requirePerm(deps, "category:delete"), handlers.Category.Delete)
}
