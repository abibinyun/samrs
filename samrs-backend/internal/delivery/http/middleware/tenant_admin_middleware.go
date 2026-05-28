package middleware

import (
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

func TenantAdminOnly(roleRepo repository.RoleRepository) gin.HandlerFunc {
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

		role, err := roleRepo.FindByID(roleID)
		if err != nil {
			util.ErrorResponseFromErr(c, "Role tidak ditemukan", util.ErrUnauthorized("role not found"))
			c.Abort()
			return
		}

		if role.IsSystem || role.IsTenantAdmin {
			c.Next()
			return
		}

		util.ErrorResponseFromErr(c, "Akses khusus tenant admin", util.ErrForbidden("tenant admin only"))
		c.Abort()
	}
}

func SuperAdminOnly(roleRepo repository.RoleRepository) gin.HandlerFunc {
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

		role, err := roleRepo.FindByID(roleID)
		if err != nil {
			util.ErrorResponseFromErr(c, "Role tidak ditemukan", util.ErrUnauthorized("role not found"))
			c.Abort()
			return
		}

		if role.IsSystem {
			c.Next()
			return
		}

		util.ErrorResponseFromErr(c, "Akses khusus super admin", util.ErrForbidden("super admin only"))
		c.Abort()
	}
}
