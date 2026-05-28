package routes

import (
	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuditRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.GET(
		"/audit-trails",
		requirePerm(deps, "audit:read"),
		middleware.TenantAdminOnly(deps.RoleRepo),
		handlers.Audit.GetAll,
	)
	v1.GET(
		"/audit-trails/:id",
		requirePerm(deps, "audit:read"),
		middleware.TenantAdminOnly(deps.RoleRepo),
		handlers.Audit.GetByID,
	)
}
