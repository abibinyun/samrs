package routes

import (
	"github.com/gin-gonic/gin"
)

func registerDocumentRoutes(v1 *gin.RouterGroup, handlers Handlers, deps RouteDeps) {
	v1.POST("/documents", requirePerm(deps, "document:create"), handlers.Document.Create)
	v1.GET("/documents", requirePerm(deps, "document:read"), handlers.Document.List)
	v1.GET("/documents/:id", requirePerm(deps, "document:read"), handlers.Document.GetByID)
	v1.PATCH("/documents/:id", requirePerm(deps, "document:update"), handlers.Document.Update)
	v1.DELETE("/documents/:id", requirePerm(deps, "document:delete"), handlers.Document.Delete)
	v1.POST("/documents/:id/files", requirePerm(deps, "document:update"), handlers.Document.AddFiles)
	v1.GET("/documents/:id/files", requirePerm(deps, "document:read"), handlers.Document.ListFiles)
	v1.GET("/document-files/:id/download", requirePerm(deps, "document:read"), handlers.Document.DownloadFile)
	v1.DELETE("/document-files/:id", requirePerm(deps, "document:delete"), handlers.Document.DeleteFile)
}
