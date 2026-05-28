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

type AssetModelHandler struct {
	usecase masterdatausecase.AssetModelUsecase
	db      *gorm.DB
}

func NewAssetModelHandler(u masterdatausecase.AssetModelUsecase, db *gorm.DB) *AssetModelHandler {
	return &AssetModelHandler{usecase: u, db: db}
}

func (h *AssetModelHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	var input struct {
		BrandID     uint   `json:"brand_id" binding:"required"`
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var model interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		modelRepo := repository.NewAssetModelRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelUsecase := masterdatausecase.NewAssetModelUsecase(modelRepo, brandRepo)
		created, err := modelUsecase.CreateModel(tenantID, masterdatausecase.AssetModelInput{
			BrandID:     input.BrandID,
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		model = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat tipe alat", err)
		return
	}

	util.SuccessResponse(c, "Tipe alat berhasil dibuat", model)
}

func (h *AssetModelHandler) GetAll(c *gin.Context) {
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

	filter := repository.AssetModelFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	if brandStr := strings.TrimSpace(c.Query("brand_id")); brandStr != "" {
		brandID, err := httputil.ParseUint(brandStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format brand_id tidak valid", err)
			return
		}
		filter.BrandID = uint(brandID)
	}

	models, total, err := h.usecase.GetAllModels(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data tipe alat", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data tipe alat ditemukan", models, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AssetModelHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tipe alat ID tidak valid", err)
		return
	}

	model, err := h.usecase.GetModelByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil tipe alat", err)
		return
	}

	util.SuccessResponse(c, "Tipe alat ditemukan", model)
}

func (h *AssetModelHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tipe alat ID tidak valid", err)
		return
	}

	var input struct {
		BrandID     uint   `json:"brand_id" binding:"required"`
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var model interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		modelRepo := repository.NewAssetModelRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelUsecase := masterdatausecase.NewAssetModelUsecase(modelRepo, brandRepo)
		updated, err := modelUsecase.UpdateModel(tenantID, uint(id), masterdatausecase.AssetModelInput{
			BrandID:     input.BrandID,
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		model = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah tipe alat", err)
		return
	}

	util.SuccessResponse(c, "Tipe alat berhasil diubah", model)
}

func (h *AssetModelHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tipe alat ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		modelRepo := repository.NewAssetModelRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelUsecase := masterdatausecase.NewAssetModelUsecase(modelRepo, brandRepo)
		return modelUsecase.DeleteModel(tenantID, uint(id))
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus tipe alat", err)
		return
	}

	util.SuccessResponse(c, "Tipe alat berhasil dihapus", nil)
}
