package httputil

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithTx(c *gin.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if c != nil {
			tenantID, tenantOK := TenantIDFromContext(c)
			userID, userOK := UserIDFromContext(c)
			if tenantOK && userOK {
				meta := map[string]interface{}{
					"tenant_id":  tenantID,
					"user_id":    userID,
					"ip":         c.ClientIP(),
					"user_agent": c.GetHeader("User-Agent"),
				}
				tx = tx.Set("audit_meta", meta)
			}
		}
		return fn(tx)
	})
}
