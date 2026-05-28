package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strconv"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	masterdatausecase "samrs-backend/internal/usecase/masterdata"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	usecase masterdatausecase.CategoryUsecase
	db      *gorm.DB
}

func NewCategoryHandler(u masterdatausecase.CategoryUsecase, db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{usecase: u, db: db}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var category *domain.Category
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		categoryRepo := repository.NewCategoryRepository(tx)
		categoryUsecase := masterdatausecase.NewCategoryUsecase(categoryRepo)

		created, err := categoryUsecase.CreateCategory(tenantID, input.Name, input.Description)
		if err != nil {
			return err
		}
		category = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat kategori", err)
		return
	}
	util.SuccessResponse(c, "Kategori berhasil dibuat", category)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
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

	filter := repository.CategoryFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	categories, total, err := h.usecase.GetAllCategories(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}
	util.SuccessResponseWithMeta(c, "Data kategori ditemukan", categories, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format category ID tidak valid", err)
		return
	}

	category, err := h.usecase.GetCategoryByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil kategori", err)
		return
	}

	util.SuccessResponse(c, "Kategori ditemukan", category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format category ID tidak valid", err)
		return
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var category *domain.Category
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		categoryRepo := repository.NewCategoryRepository(tx)
		categoryUsecase := masterdatausecase.NewCategoryUsecase(categoryRepo)

		updated, err := categoryUsecase.UpdateCategory(tenantID, uint(id), input.Name, input.Description)
		if err != nil {
			return err
		}
		category = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah kategori", err)
		return
	}

	util.SuccessResponse(c, "Kategori berhasil diubah", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format category ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		categoryRepo := repository.NewCategoryRepository(tx)
		categoryUsecase := masterdatausecase.NewCategoryUsecase(categoryRepo)

		if err := categoryUsecase.DeleteCategory(tenantID, uint(id)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus kategori", err)
		return
	}

	util.SuccessResponse(c, "Kategori berhasil dihapus", nil)
}
