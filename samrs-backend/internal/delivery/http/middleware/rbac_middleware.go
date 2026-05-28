package middleware

import (
	"fmt"

	rbacusecase "samrs-backend/internal/usecase/rbac"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(rbac rbacusecase.RBACUsecase, permissionSlug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDRaw, exists := c.Get("role_id")
		if !exists {
			util.ErrorResponseFromErr(c, "Role ID tidak ditemukan", util.ErrUnauthorized("role id missing"))
			c.Abort()
			return
		}

		roleID, ok := toInt(roleIDRaw)
		if !ok {
			util.ErrorResponseFromErr(c, "Format Role ID tidak valid", util.ErrUnauthorized("invalid role id"))
			c.Abort()
			return
		}

		allowed, err := rbac.Authorize(roleID, permissionSlug)
		if err != nil {
			util.ErrorResponseFromErr(c, "Role tidak ditemukan", util.ErrUnauthorized("role not found"))
			c.Abort()
			return
		}

		if !allowed {
			util.ErrorResponseFromErr(c, "Akses ditolak", util.ErrForbidden("access denied"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func toInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		var parsed int
		_, err := fmt.Sscanf(v, "%d", &parsed)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
