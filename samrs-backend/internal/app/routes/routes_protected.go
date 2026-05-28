package routes

import (
	"os"
	"strconv"
	"time"

	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func registerProtectedRoutes(r *gin.Engine, handlers Handlers, deps RouteDeps) {
	v1 := r.Group("/api/v1")
	apiLimiter := middleware.NewRateLimiter(
		getEnvLimitAPI("API_RATE_LIMIT", 10),
		getEnvIntAPI("API_RATE_BURST", 20),
		time.Duration(getEnvIntAPI("API_RATE_WINDOW_MINUTES", 5))*time.Minute,
	)
	v1.Use(middleware.AuthMiddleware(), middleware.TenantScope(), apiLimiter.Middleware())

	registerAuthRoutes(v1, handlers)
	registerRoomRoutes(v1, handlers, deps)
	registerBedRoutes(v1, handlers, deps)
	registerCategoryRoutes(v1, handlers, deps)
	registerVendorRoutes(v1, handlers, deps)
	registerBrandRoutes(v1, handlers, deps)
	registerModelRoutes(v1, handlers, deps)
	registerStatusRoutes(v1, handlers, deps)
	registerMaintenanceRoutes(v1, handlers, deps)
	registerMaintenanceDocRoutes(v1, handlers, deps)
	registerDocumentRoutes(v1, handlers, deps)
	registerAssetMutationRoutes(v1, handlers, deps)
	registerStockOpnameRoutes(v1, handlers, deps)
	registerReportRoutes(v1, handlers, deps)
	registerNotificationRoutes(v1, handlers, deps)
	registerAssetRoutes(v1, handlers, deps)
	registerComplaintRoutes(v1, handlers, deps)
	registerPermissionRoutes(v1, handlers, deps)
	registerRoleRoutes(v1, handlers, deps)
	registerUserRoutes(v1, handlers, deps)
	registerAuditRoutes(v1, handlers, deps)
	registerAdminTenantRoutes(v1, handlers, deps)
}

func getEnvIntAPI(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if val, err := strconv.Atoi(raw); err == nil {
		return val
	}
	return fallback
}

func getEnvLimitAPI(key string, fallback int) rate.Limit {
	return rate.Limit(getEnvIntAPI(key, fallback))
}
