package httputil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TenantIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	if cached, ok := c.Get("tenant_id_uuid"); ok {
		if tenantID, ok := cached.(uuid.UUID); ok {
			return tenantID, true
		}
	}

	tenantIDRaw, exists := c.Get("tenant_id")
	if !exists {
		return uuid.UUID{}, false
	}

	tenantID, err := uuid.Parse(fmt.Sprint(tenantIDRaw))
	if err != nil {
		return uuid.UUID{}, false
	}
	return tenantID, true
}

func UserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	if cached, ok := c.Get("user_id_uuid"); ok {
		if userID, ok := cached.(uuid.UUID); ok {
			return userID, true
		}
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return uuid.UUID{}, false
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return uuid.UUID{}, false
	}
	return userID, true
}
