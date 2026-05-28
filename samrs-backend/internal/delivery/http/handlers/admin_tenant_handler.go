package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"encoding/json"
	"net/http"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	tenantusecase "samrs-backend/internal/usecase/tenant"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AdminTenantHandler struct {
	usecase tenantusecase.AdminTenantUsecase
	db      *gorm.DB
}

func NewAdminTenantHandler(u tenantusecase.AdminTenantUsecase, db *gorm.DB) *AdminTenantHandler {
	return &AdminTenantHandler{usecase: u, db: db}
}

func (h *AdminTenantHandler) GetAll(c *gin.Context) {
	pagination := httputil.ParsePagination(c)
	dateRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}

	filter := repository.TenantFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	items, total, err := h.usecase.ListTenantsWithRoles(filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data tenant", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data tenant ditemukan", items, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AdminTenantHandler) GetByID(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tenant ID tidak valid", err)
		return
	}

	userPage := httputil.ParseIntWithDefault(c.Query("user_page"), 1)
	userPerPage := httputil.ParseIntWithDefault(c.Query("user_per_page"), 20)
	if userPage < 1 {
		userPage = 1
	}
	if userPerPage < 1 {
		userPerPage = 1
	}
	if userPerPage > 100 {
		userPerPage = 100
	}

	userFilter := repository.UserFilter{
		Search:  strings.TrimSpace(c.Query("user_search")),
		Page:    userPage,
		PerPage: userPerPage,
	}

	rolePage := httputil.ParseIntWithDefault(c.Query("role_page"), 1)
	rolePerPage := httputil.ParseIntWithDefault(c.Query("role_per_page"), 20)
	if rolePage < 1 {
		rolePage = 1
	}
	if rolePerPage < 1 {
		rolePerPage = 1
	}
	if rolePerPage > 100 {
		rolePerPage = 100
	}

	roleFilter := repository.RoleFilter{
		Search:  strings.TrimSpace(c.Query("role_search")),
		Page:    rolePage,
		PerPage: rolePerPage,
	}

	detail, err := h.usecase.GetTenantDetail(tenantID, userFilter, roleFilter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil detail tenant", err)
		return
	}

	util.SuccessResponse(c, "Detail tenant ditemukan", detail)
}

func (h *AdminTenantHandler) Create(c *gin.Context) {
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	var input struct {
		Name       string `json:"name" binding:"required"`
		Slug       string `json:"slug"`
		Address    string `json:"address"`
		TenantType string `json:"type"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var tenant *domain.Tenant
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		tenantRepo := repository.NewTenantRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		auditRepo := repository.NewAuditTrailRepository(tx)

		adminUsecase := tenantusecase.NewAdminTenantUsecase(
			tenantRepo,
			roleRepo,
			userRepo,
			roomRepo,
			bedRepo,
			categoryRepo,
			assetRepo,
			auditRepo,
		)

		created, err := adminUsecase.CreateTenant(tenantusecase.CreateTenantInput{
			Name:       input.Name,
			Slug:       input.Slug,
			Address:    input.Address,
			TenantType: input.TenantType,
			Status:     input.Status,
		})
		if err != nil {
			return err
		}
		tenant = created

		newJSON, err := toJSON(tenantAuditData(tenant))
		if err != nil {
			return err
		}

		audit := &domain.AuditTrail{
			TenantID:  tenant.ID,
			UserID:    userID,
			Action:    "CREATE",
			TableName: "tenants",
			RecordID:  tenant.ID.String(),
			NewData:   newJSON,
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}
		return auditRepo.Create(audit)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat tenant", err)
		return
	}

	util.SuccessResponse(c, "Tenant berhasil dibuat", tenant)
}

func (h *AdminTenantHandler) Update(c *gin.Context) {
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tenant ID tidak valid", err)
		return
	}

	var input struct {
		Name       string `json:"name" binding:"required"`
		Slug       string `json:"slug"`
		Address    string `json:"address"`
		TenantType string `json:"type"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var tenant *domain.Tenant
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		tenantRepo := repository.NewTenantRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		auditRepo := repository.NewAuditTrailRepository(tx)

		adminUsecase := tenantusecase.NewAdminTenantUsecase(
			tenantRepo,
			roleRepo,
			userRepo,
			roomRepo,
			bedRepo,
			categoryRepo,
			assetRepo,
			auditRepo,
		)

		oldTenant, err := tenantRepo.FindByID(tenantID)
		if err != nil {
			return err
		}

		updated, err := adminUsecase.UpdateTenant(tenantusecase.UpdateTenantInput{
			ID:         tenantID,
			Name:       input.Name,
			Slug:       input.Slug,
			Address:    input.Address,
			TenantType: input.TenantType,
		})
		if err != nil {
			return err
		}
		tenant = updated

		oldJSON, err := toJSON(tenantAuditData(oldTenant))
		if err != nil {
			return err
		}
		newJSON, err := toJSON(tenantAuditData(tenant))
		if err != nil {
			return err
		}

		audit := &domain.AuditTrail{
			TenantID:  tenant.ID,
			UserID:    userID,
			Action:    "UPDATE",
			TableName: "tenants",
			RecordID:  tenant.ID.String(),
			OldData:   oldJSON,
			NewData:   newJSON,
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}
		return auditRepo.Create(audit)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah tenant", err)
		return
	}

	util.SuccessResponse(c, "Tenant berhasil diubah", tenant)
}

func (h *AdminTenantHandler) SetStatus(c *gin.Context) {
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tenant ID tidak valid", err)
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var tenant *domain.Tenant
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		tenantRepo := repository.NewTenantRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		auditRepo := repository.NewAuditTrailRepository(tx)

		adminUsecase := tenantusecase.NewAdminTenantUsecase(
			tenantRepo,
			roleRepo,
			userRepo,
			roomRepo,
			bedRepo,
			categoryRepo,
			assetRepo,
			auditRepo,
		)

		oldTenant, err := tenantRepo.FindByID(tenantID)
		if err != nil {
			return err
		}

		updated, err := adminUsecase.SetTenantStatus(tenantID, input.Status)
		if err != nil {
			return err
		}
		tenant = updated

		oldJSON, err := toJSON(tenantAuditData(oldTenant))
		if err != nil {
			return err
		}
		newJSON, err := toJSON(tenantAuditData(tenant))
		if err != nil {
			return err
		}

		audit := &domain.AuditTrail{
			TenantID:  tenant.ID,
			UserID:    userID,
			Action:    "UPDATE",
			TableName: "tenants",
			RecordID:  tenant.ID.String(),
			OldData:   oldJSON,
			NewData:   newJSON,
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}
		return auditRepo.Create(audit)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah status tenant", err)
		return
	}

	util.SuccessResponse(c, "Status tenant berhasil diubah", tenant)
}

func tenantAuditData(tenant *domain.Tenant) map[string]interface{} {
	if tenant == nil {
		return nil
	}
	return map[string]interface{}{
		"id":      tenant.ID.String(),
		"name":    tenant.Name,
		"slug":    tenant.Slug,
		"address": tenant.Address,
		"type":    tenant.TenantType,
		"status":  tenant.Status,
	}
}

func toJSON(payload interface{}) (datatypes.JSON, error) {
	if payload == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(bytes), nil
}
