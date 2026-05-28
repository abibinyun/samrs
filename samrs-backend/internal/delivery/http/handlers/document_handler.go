package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	documentusecase "samrs-backend/internal/usecase/document"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxGlobalDocumentSize = 10 << 20

type DocumentHandler struct {
	usecase documentusecase.DocumentUsecase
	db      *gorm.DB
}

func NewDocumentHandler(u documentusecase.DocumentUsecase, db *gorm.DB) *DocumentHandler {
	return &DocumentHandler{usecase: u, db: db}
}

func (h *DocumentHandler) Create(c *gin.Context) {
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

	title := strings.TrimSpace(c.PostForm("title"))
	docType := strings.TrimSpace(c.PostForm("doc_type"))
	description := strings.TrimSpace(c.PostForm("description"))

	files, storedPaths, err := saveUploadedFiles(c, "files", maxGlobalDocumentSize, globalDocumentDir())
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menyimpan dokumen", err)
		return
	}

	var (
		doc       *domain.Document
		docFiles  []domain.DocumentFile
		input     = documentusecase.DocumentInput{TenantID: tenantID, Title: title, DocType: docType, Description: description, UploadedBy: userID}
		fileInput = buildDocumentFileInputs(files, storedPaths)
	)

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewDocumentRepository(tx)
		fileRepo := repository.NewDocumentFileRepository(tx)
		docUsecase := documentusecase.NewDocumentUsecase(repo, fileRepo, tx)

		created, createdFiles, err := docUsecase.CreateDocument(input, fileInput)
		if err != nil {
			return err
		}
		doc = created
		docFiles = createdFiles
		return nil
	}); err != nil {
		for _, path := range storedPaths {
			_ = os.Remove(path)
		}
		util.ErrorResponseFromErr(c, "Gagal menyimpan dokumen", err)
		return
	}

	attachDocumentFileURLs(c, docFiles)
	doc.Files = docFiles
	util.SuccessResponse(c, "Dokumen berhasil dibuat", doc)
}

func (h *DocumentHandler) List(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	pagination := httputil.ParsePagination(c)
	dateRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}

	filter := repository.DocumentFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DocType:  strings.TrimSpace(c.Query("doc_type")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	items, total, err := h.usecase.ListDocuments(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil dokumen", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Dokumen ditemukan", items, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *DocumentHandler) GetByID(c *gin.Context) {
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

	files, err := h.usecase.ListFiles(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil file", err)
		return
	}
	attachDocumentFileURLs(c, files)
	doc.Files = files

	util.SuccessResponse(c, "Dokumen ditemukan", doc)
}

func (h *DocumentHandler) Update(c *gin.Context) {
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

	var input struct {
		Title       string `json:"title"`
		DocType     string `json:"doc_type"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	updated, err := h.usecase.UpdateDocument(tenantID, uint(id), documentusecase.DocumentUpdateInput{
		Title:       input.Title,
		DocType:     input.DocType,
		Description: input.Description,
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah dokumen", err)
		return
	}

	util.SuccessResponse(c, "Dokumen berhasil diubah", updated)
}

func (h *DocumentHandler) AddFiles(c *gin.Context) {
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

	docID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document ID tidak valid", err)
		return
	}

	files, storedPaths, err := saveUploadedFiles(c, "files", maxGlobalDocumentSize, globalDocumentDir())
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menambah file dokumen", err)
		return
	}

	docFiles, err := h.usecase.AddFiles(tenantID, uint(docID), userID, buildDocumentFileInputs(files, storedPaths))
	if err != nil {
		for _, path := range storedPaths {
			_ = os.Remove(path)
		}
		util.ErrorResponseFromErr(c, "Gagal menambah file dokumen", err)
		return
	}

	attachDocumentFileURLs(c, docFiles)
	util.SuccessResponse(c, "File dokumen berhasil ditambahkan", docFiles)
}

func (h *DocumentHandler) ListFiles(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	docID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document ID tidak valid", err)
		return
	}

	files, err := h.usecase.ListFiles(tenantID, uint(docID))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil file", err)
		return
	}

	attachDocumentFileURLs(c, files)
	util.SuccessResponse(c, "File dokumen ditemukan", files)
}

func (h *DocumentHandler) DownloadFile(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document file ID tidak valid", err)
		return
	}

	file, err := h.usecase.GetFileByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "File tidak ditemukan", err)
		return
	}

	c.FileAttachment(file.FilePath, file.Filename)
}

func (h *DocumentHandler) DeleteFile(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format document file ID tidak valid", err)
		return
	}

	file, err := h.usecase.DeleteFile(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus file", err)
		return
	}
	_ = os.Remove(file.FilePath)

	util.SuccessResponse(c, "File dokumen berhasil dihapus", nil)
}

func (h *DocumentHandler) Delete(c *gin.Context) {
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

	files, err := h.usecase.DeleteDocument(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus dokumen", err)
		return
	}
	for _, file := range files {
		_ = os.Remove(file.FilePath)
	}

	util.SuccessResponse(c, "Dokumen berhasil dihapus", nil)
}

func globalDocumentDir() string {
	uploadDir := strings.TrimSpace(os.Getenv("DOCUMENTS_DIR"))
	if uploadDir == "" {
		return "./uploads/documents"
	}
	return uploadDir
}

func saveUploadedFiles(c *gin.Context, field string, maxSize int64, uploadDir string) ([]*multipart.FileHeader, []string, error) {
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil, nil, util.ErrValidation("file dokumen wajib diunggah")
	}

	headers := form.File[field]
	if len(headers) == 0 {
		return nil, nil, util.ErrValidation("file dokumen wajib diunggah")
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("gagal menyiapkan folder upload")
	}

	storedPaths := make([]string, 0, len(headers))
	for _, fileHeader := range headers {
		if fileHeader.Size > maxSize {
			return nil, nil, util.ErrValidation("ukuran file terlalu besar (max 10MB)")
		}
		ext := filepath.Ext(fileHeader.Filename)
		storedName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		storedPath := filepath.Join(uploadDir, storedName)
		if err := c.SaveUploadedFile(fileHeader, storedPath); err != nil {
			return nil, nil, fmt.Errorf("gagal menyimpan file")
		}
		storedPaths = append(storedPaths, storedPath)
	}

	return headers, storedPaths, nil
}

func buildDocumentFileInputs(headers []*multipart.FileHeader, storedPaths []string) []documentusecase.DocumentFileInput {
	inputs := make([]documentusecase.DocumentFileInput, 0, len(headers))
	for i, fileHeader := range headers {
		inputs = append(inputs, documentusecase.DocumentFileInput{
			Filename: fileHeader.Filename,
			FilePath: storedPaths[i],
			MimeType: fileHeader.Header.Get("Content-Type"),
			Size:     fileHeader.Size,
		})
	}
	return inputs
}

func attachDocumentFileURLs(c *gin.Context, files []domain.DocumentFile) {
	for i := range files {
		files[i].FileURL = buildDocumentFileURL(c, files[i].ID)
	}
}

func buildDocumentFileURL(c *gin.Context, fileID uint) string {
	base := strings.TrimRight(os.Getenv("PUBLIC_FILE_BASE_URL"), "/")
	if base == "" {
		base = requestBaseURL(c)
	}
	return fmt.Sprintf("%s/api/v1/document-files/%d/download", base, fileID)
}
