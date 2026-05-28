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

type VendorHandler struct {
	usecase masterdatausecase.VendorUsecase
	db      *gorm.DB
}

func NewVendorHandler(u masterdatausecase.VendorUsecase, db *gorm.DB) *VendorHandler {
	return &VendorHandler{usecase: u, db: db}
}

func (h *VendorHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Address     string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var vendor interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewVendorRepository(tx)
		vendorUsecase := masterdatausecase.NewVendorUsecase(repo)
		created, err := vendorUsecase.CreateVendor(tenantID, masterdatausecase.VendorInput{
			Code:        input.Code,
			Name:        input.Name,
			ContactName: input.ContactName,
			Phone:       input.Phone,
			Email:       input.Email,
			Address:     input.Address,
		})
		if err != nil {
			return err
		}
		vendor = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat vendor", err)
		return
	}

	util.SuccessResponse(c, "Vendor berhasil dibuat", vendor)
}

func (h *VendorHandler) GetAll(c *gin.Context) {
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

	filter := repository.VendorFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	vendors, total, err := h.usecase.GetAllVendors(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data vendor", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data vendor ditemukan", vendors, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *VendorHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format vendor ID tidak valid", err)
		return
	}

	vendor, err := h.usecase.GetVendorByID(tenantID, id)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil vendor", err)
		return
	}

	util.SuccessResponse(c, "Vendor ditemukan", vendor)
}

func (h *VendorHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format vendor ID tidak valid", err)
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name" binding:"required"`
		ContactName string `json:"contact_name"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Address     string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var vendor interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewVendorRepository(tx)
		vendorUsecase := masterdatausecase.NewVendorUsecase(repo)
		updated, err := vendorUsecase.UpdateVendor(tenantID, id, masterdatausecase.VendorInput{
			Code:        input.Code,
			Name:        input.Name,
			ContactName: input.ContactName,
			Phone:       input.Phone,
			Email:       input.Email,
			Address:     input.Address,
		})
		if err != nil {
			return err
		}
		vendor = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah vendor", err)
		return
	}

	util.SuccessResponse(c, "Vendor berhasil diubah", vendor)
}

func (h *VendorHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := parseUintParam(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format vendor ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		repo := repository.NewVendorRepository(tx)
		vendorUsecase := masterdatausecase.NewVendorUsecase(repo)
		return vendorUsecase.DeleteVendor(tenantID, id)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus vendor", err)
		return
	}

	util.SuccessResponse(c, "Vendor berhasil dihapus", nil)
}

func parseUintParam(raw string) (uint, error) {
	parsed, err := httputil.ParseUint(raw)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}
