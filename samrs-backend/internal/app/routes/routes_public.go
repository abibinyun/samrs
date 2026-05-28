package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"samrs-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

func registerPublicRoutes(r *gin.Engine, handlers Handlers) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "online",
			"message": "SAMRS API Foundation is ready",
		})
	})

	r.GET("/swagger", func(c *gin.Context) {
		c.File("./docs/swagger.html")
	})
	r.GET("/swagger-ui", func(c *gin.Context) {
		c.File("./docs/swagger-ui.html")
	})
	r.GET("/swagger.yaml", func(c *gin.Context) {
		c.File("./docs/swagger.yaml")
	})
	r.GET("/metrics", metricsHandler())
	r.GET("/public/tenants/:slug/assets/:code", handlers.PublicAsset.GetByTenantAndCode)

	loginLimiter := middleware.NewRateLimiter(
		getEnvLimit("LOGIN_RATE_LIMIT", 2),
		getEnvIntLogin("LOGIN_RATE_BURST", 5),
		time.Duration(getEnvIntLogin("LOGIN_RATE_WINDOW_MINUTES", 5))*time.Minute,
	)
	r.POST("/api/v1/auth/login", loginLimiter.Middleware(), handlers.Auth.Login)
}

func getEnvIntLogin(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if val, err := strconv.Atoi(raw); err == nil {
		return val
	}
	return fallback
}

func getEnvLimit(key string, fallback int) rate.Limit {
	return rate.Limit(getEnvIntLogin(key, fallback))
}

func metricsHandler() gin.HandlerFunc {
	token := strings.TrimSpace(os.Getenv("METRICS_TOKEN"))
	handler := promhttp.Handler()
	if token == "" {
		return gin.WrapH(handler)
	}
	return func(c *gin.Context) {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if auth == "" || auth != "Bearer "+token {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
