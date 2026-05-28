package routes

import (
	"github.com/gin-gonic/gin"
)

func registerMaintenanceRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/maintenance-schedules", requirePerm(deps, "maintenance:create"), handlers.Maintenance.Create)
	v1.GET("/maintenance-schedules", requirePerm(deps, "maintenance:read"), handlers.Maintenance.GetAll)
	v1.GET("/maintenance-schedules/:id", requirePerm(deps, "maintenance:read"), handlers.Maintenance.GetByID)
	v1.PATCH("/maintenance-schedules/:id", requirePerm(deps, "maintenance:update"), handlers.Maintenance.Update)
	v1.PATCH("/maintenance-schedules/:id/complete", requirePerm(deps, "maintenance:complete"), handlers.Maintenance.Complete)
	v1.DELETE("/maintenance-schedules/:id", requirePerm(deps, "maintenance:delete"), handlers.Maintenance.Delete)
}

func registerMaintenanceDocRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/maintenance-schedules/:id/documents", requirePerm(deps, "maintenance_document:upload"), handlers.MaintenanceDoc.Upload)
	v1.GET("/maintenance-schedules/:id/documents", requirePerm(deps, "maintenance_document:read"), handlers.MaintenanceDoc.ListBySchedule)
	v1.GET("/maintenance-documents/:id/download", requirePerm(deps, "maintenance_document:read"), handlers.MaintenanceDoc.Download)
	v1.DELETE("/maintenance-documents/:id", requirePerm(deps, "maintenance_document:delete"), handlers.MaintenanceDoc.Delete)
}
