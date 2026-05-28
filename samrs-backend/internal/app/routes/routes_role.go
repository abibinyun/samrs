package routes

import (
	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func registerRoleRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/roles", requirePerm(deps, "role:create"), handlers.Role.Create)
	v1.GET("/roles", requirePerm(deps, "role:read"), handlers.Role.GetAll)
	v1.GET("/roles/:id", requirePerm(deps, "role:read"), handlers.Role.GetByID)
	v1.PATCH("/roles/:id", requirePerm(deps, "role:update"), handlers.Role.Update)
	v1.PATCH(
		"/roles/:id/tenant-admin",
		requirePerm(deps, "role:update"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.Role.SetTenantAdmin,
	)
	v1.DELETE("/roles/:id", requirePerm(deps, "role:delete"), handlers.Role.Delete)
	v1.POST("/roles/:id/permissions", requirePerm(deps, "role:assign_permissions"), handlers.Role.AssignPermissions)
	v1.DELETE("/roles/:id/permissions", requirePerm(deps, "role:revoke_permissions"), handlers.Role.RevokePermissions)
}
