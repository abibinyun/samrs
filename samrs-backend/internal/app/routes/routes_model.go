package routes

import (
	"github.com/gin-gonic/gin"
)

func registerModelRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/models", requirePerm(deps, "model:create"), handlers.Model.Create)
	v1.GET("/models", requirePerm(deps, "model:read"), handlers.Model.GetAll)
	v1.GET("/models/:id", requirePerm(deps, "model:read"), handlers.Model.GetByID)
	v1.PATCH("/models/:id", requirePerm(deps, "model:update"), handlers.Model.Update)
	v1.DELETE("/models/:id", requirePerm(deps, "model:delete"), handlers.Model.Delete)
}
