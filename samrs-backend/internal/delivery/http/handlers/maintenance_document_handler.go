package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	maintenanceusecase "samrs-backend/internal/usecase/maintenance"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxDocumentSize = 10 << 20

type MaintenanceDocumentHandler struct {
	usecase maintenanceusecase.MaintenanceDocumentUsecase
	db      *gorm.DB
}

func NewMaintenanceDocumentHandler(u maintenanceusecase.MaintenanceDocumentUsecase, db *gorm.DB) *MaintenanceDocumentHandler {
	return &MaintenanceDocumentHandler{usecase: u, db: db}
}

func (h *MaintenanceDocumentHandler) Upload(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	scheduleID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		util.ErrorResponseFromErr(c, "File wajib diisi", util.ErrValidation("file wajib diisi"))
		return
	}
	if fileHeader.Size > maxDocumentSize {
		util.ErrorResponseFromErr(c, "Ukuran file terlalu besar", util.ErrValidation("max 10MB"))
		return
	}

	docType := strings.TrimSpace(c.PostForm("doc_type"))
	if docType == "" {
		docType = domain.MaintenanceDocOther
	}

	uploadDir := os.Getenv("MAINTENANCE_DOCS_DIR")
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = "./uploads/maintenance"
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menyiapkan folder upload", err)
		return
	}

	ext := filepath.Ext(fileHeader.Filename)
	storedName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	storedPath := filepath.Join(uploadDir, storedName)

	if err := c.SaveUploadedFile(fileHeader, storedPath); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menyimpan file", err)
		return
	}

	var doc *domain.MaintenanceDocument
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		docRepo := repository.NewMaintenanceDocumentRepository(tx)
		scheduleRepo := repository.NewMaintenanceScheduleRepository(tx)
		docUsecase := maintenanceusecase.NewMaintenanceDocumentUsecase(docRepo, scheduleRepo)

		created, err := docUsecase.CreateDocument(maintenanceusecase.MaintenanceDocumentInput{
			TenantID:   tenantID,
			ScheduleID: uint(scheduleID),
			UploadedBy: userID,
			DocType:    docType,
			Filename:   fileHeader.Filename,
			FilePath:   storedPath,
			MimeType:   fileHeader.Header.Get("Content-Type"),
			Size:       fileHeader.Size,
		})
		if err != nil {
			return err
		}
		doc = created
		return nil
	}); err != nil {
		_ = os.Remove(storedPath)
		util.ErrorResponseFromErr(c, "Gagal menyimpan dokumen", err)
		return
	}

	doc.FileURL = buildFileURL(c, doc.ID)
	util.SuccessResponse(c, "Dokumen berhasil diunggah", doc)
}

func (h *MaintenanceDocumentHandler) ListBySchedule(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	scheduleID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	docs, err := h.usecase.ListDocuments(tenantID, uint(scheduleID))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil dokumen", err)
		return
	}

	for i := range docs {
		docs[i].FileURL = buildFileURL(c, docs[i].ID)
	}

	util.SuccessResponse(c, "Dokumen ditemukan", docs)
}

func (h *MaintenanceDocumentHandler) Download(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document ID tidak valid", err)
		return
	}

	doc, err := h.usecase.GetDocumentByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Dokumen tidak ditemukan", err)
		return
	}

	c.FileAttachment(doc.FilePath, doc.Filename)
}

func (h *MaintenanceDocumentHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document ID tidak valid", err)
		return
	}

	doc, err := h.usecase.GetDocumentByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Dokumen tidak ditemukan", err)
		return
	}

	if err := h.usecase.DeleteDocument(tenantID, uint(id)); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus dokumen", err)
		return
	}
	_ = os.Remove(doc.FilePath)

	util.SuccessResponse(c, "Dokumen berhasil dihapus", nil)
}

func buildFileURL(c *gin.Context, docID uint) string {
	base := strings.TrimRight(os.Getenv("PUBLIC_FILE_BASE_URL"), "/")
	if base == "" {
		base = requestBaseURL(c)
	}
	return fmt.Sprintf("%s/api/v1/maintenance-documents/%d/download", base, docID)
}
