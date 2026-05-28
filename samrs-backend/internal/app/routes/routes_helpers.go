package routes

import (
	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func requirePerm(deps RouteDeps, perm string) gin.HandlerFunc {
	return middleware.RBACMiddleware(deps.RBAC, perm)
}
