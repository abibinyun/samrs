package routes

import (
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/users", requirePerm(deps, "user:create"), handlers.User.Create)
	v1.GET("/users", requirePerm(deps, "user:read"), handlers.User.GetAll)
	v1.PATCH("/users/:id", requirePerm(deps, "user:update"), handlers.User.Update)
	v1.PATCH("/users/:id/role", requirePerm(deps, "user:update"), handlers.User.UpdateRole)
	v1.PATCH("/users/:id/status", requirePerm(deps, "user:update"), handlers.User.SetActive)
	v1.PATCH("/users/:id/password", requirePerm(deps, "user:reset_password"), handlers.User.ResetPassword)
	v1.DELETE("/users/:id", requirePerm(deps, "user:delete"), handlers.User.Delete)
}
