package routes

import (
	"github.com/gin-gonic/gin"
)

func registerReportRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.GET("/reports/assets/export", requirePerm(deps, "report:export"), handlers.Report.ExportAssets)
	v1.GET("/reports/complaints/export", requirePerm(deps, "report:export"), handlers.Report.ExportComplaints)
	v1.GET("/reports/maintenance-schedules/export", requirePerm(deps, "report:export"), handlers.Report.ExportMaintenance)
}
