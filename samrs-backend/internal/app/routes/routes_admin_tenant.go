package routes

import (
	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func registerAdminTenantRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.GET(
		"/admin/tenants",
		requirePerm(deps, "tenant:read"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.AdminTenant.GetAll,
	)
	v1.POST(
		"/admin/tenants",
		requirePerm(deps, "tenant:create"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.AdminTenant.Create,
	)
	v1.GET(
		"/admin/tenants/:id",
		requirePerm(deps, "tenant:read"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.AdminTenant.GetByID,
	)
	v1.PATCH(
		"/admin/tenants/:id",
		requirePerm(deps, "tenant:update"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.AdminTenant.Update,
	)
	v1.PATCH(
		"/admin/tenants/:id/status",
		requirePerm(deps, "tenant:update_status"),
		middleware.SuperAdminOnly(deps.RoleRepo),
		handlers.AdminTenant.SetStatus,
	)
}
