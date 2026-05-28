package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/repository"
	masterdatausecase "samrs-backend/internal/usecase/masterdata"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetStatusHandler struct {
	usecase masterdatausecase.AssetStatusUsecase
	db      *gorm.DB
}

func NewAssetStatusHandler(u masterdatausecase.AssetStatusUsecase, db *gorm.DB) *AssetStatusHandler {
	return &AssetStatusHandler{usecase: u, db: db}
}

func (h *AssetStatusHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var status interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetStatusRepository(tx)
		statusUsecase := masterdatausecase.NewAssetStatusUsecase(repo)
		created, err := statusUsecase.CreateStatus(tenantID, masterdatausecase.AssetStatusInput{
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		status = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat status aset", err)
		return
	}

	util.SuccessResponse(c, "Status aset berhasil dibuat", status)
}

func (h *AssetStatusHandler) GetAll(c *gin.Context) {
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

	filter := repository.AssetStatusFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	statuses, total, err := h.usecase.GetAllStatuses(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data status aset", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data status aset ditemukan", statuses, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AssetStatusHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format status ID tidak valid", err)
		return
	}

	status, err := h.usecase.GetStatusByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil status aset", err)
		return
	}

	util.SuccessResponse(c, "Status aset ditemukan", status)
}

func (h *AssetStatusHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format status ID tidak valid", err)
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var status interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetStatusRepository(tx)
		statusUsecase := masterdatausecase.NewAssetStatusUsecase(repo)
		updated, err := statusUsecase.UpdateStatus(tenantID, uint(id), masterdatausecase.AssetStatusInput{
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		status = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah status aset", err)
		return
	}

	util.SuccessResponse(c, "Status aset berhasil diubah", status)
}

func (h *AssetStatusHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format status ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetStatusRepository(tx)
		statusUsecase := masterdatausecase.NewAssetStatusUsecase(repo)
		return statusUsecase.DeleteStatus(tenantID, uint(id))
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus status aset", err)
		return
	}

	util.SuccessResponse(c, "Status aset berhasil dihapus", nil)
}
