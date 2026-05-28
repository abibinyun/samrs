package middleware

import (
	"fmt"

	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenantScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawTenantID, exists := c.Get("tenant_id")
		if !exists {
			util.ErrorResponseFromErr(c, "Tenant ID tidak ditemukan", util.ErrUnauthorized("tenant id missing"))
			c.Abort()
			return
		}

		tenantID, err := uuid.Parse(fmt.Sprint(rawTenantID))
		if err != nil || tenantID == uuid.Nil {
			util.ErrorResponseFromErr(c, "Tenant ID tidak valid", util.ErrUnauthorized("invalid tenant id"))
			c.Abort()
			return
		}

		c.Set("tenant_id_uuid", tenantID)
		c.Next()
	}
}
