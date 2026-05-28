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

type AssetBrandHandler struct {
	usecase masterdatausecase.AssetBrandUsecase
	db      *gorm.DB
}

func NewAssetBrandHandler(u masterdatausecase.AssetBrandUsecase, db *gorm.DB) *AssetBrandHandler {
	return &AssetBrandHandler{usecase: u, db: db}
}

func (h *AssetBrandHandler) Create(c *gin.Context) {
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

	var brand interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetBrandRepository(tx)
		brandUsecase := masterdatausecase.NewAssetBrandUsecase(repo)
		created, err := brandUsecase.CreateBrand(tenantID, masterdatausecase.AssetBrandInput{
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		brand = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat merek", err)
		return
	}

	util.SuccessResponse(c, "Merek berhasil dibuat", brand)
}

func (h *AssetBrandHandler) GetAll(c *gin.Context) {
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

	filter := repository.AssetBrandFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	brands, total, err := h.usecase.GetAllBrands(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data merek", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data merek ditemukan", brands, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AssetBrandHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format merek ID tidak valid", err)
		return
	}

	brand, err := h.usecase.GetBrandByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil merek", err)
		return
	}

	util.SuccessResponse(c, "Merek ditemukan", brand)
}

func (h *AssetBrandHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format merek ID tidak valid", err)
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

	var brand interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetBrandRepository(tx)
		brandUsecase := masterdatausecase.NewAssetBrandUsecase(repo)
		updated, err := brandUsecase.UpdateBrand(tenantID, uint(id), masterdatausecase.AssetBrandInput{
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		brand = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah merek", err)
		return
	}

	util.SuccessResponse(c, "Merek berhasil diubah", brand)
}

func (h *AssetBrandHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format merek ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewAssetBrandRepository(tx)
		brandUsecase := masterdatausecase.NewAssetBrandUsecase(repo)
		return brandUsecase.DeleteBrand(tenantID, uint(id))
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus merek", err)
		return
	}

	util.SuccessResponse(c, "Merek berhasil dihapus", nil)
}
