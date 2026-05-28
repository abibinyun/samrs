package routes

import (
	"github.com/gin-gonic/gin"
)

func registerComplaintRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/complaints", requirePerm(deps, "complaint:create"), handlers.Complaint.Create)
	v1.GET("/complaints", requirePerm(deps, "complaint:read"), handlers.Complaint.GetAll)
	v1.GET("/complaints/:id", requirePerm(deps, "complaint:read"), handlers.Complaint.GetByID)
	v1.PATCH("/complaints/:id", requirePerm(deps, "complaint:update"), handlers.Complaint.Update)
	v1.DELETE("/complaints/:id", requirePerm(deps, "complaint:delete"), handlers.Complaint.Delete)
}
