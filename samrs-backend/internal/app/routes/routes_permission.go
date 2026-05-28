package routes

import (
	"github.com/gin-gonic/gin"
)

func registerPermissionRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.GET("/permissions", requirePerm(deps, "permission:read"), handlers.Permission.GetAll)
	v1.GET("/roles/:id/permissions", requirePerm(deps, "role:read_permissions"), handlers.Permission.GetByRole)
}
