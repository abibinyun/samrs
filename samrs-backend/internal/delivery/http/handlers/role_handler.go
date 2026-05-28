package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"fmt"
	"net/http"
	"strconv"

	"samrs-backend/internal/repository"
	rbacusecase "samrs-backend/internal/usecase/rbac"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	roleUsecase           rbacusecase.RoleUsecase
	rolePermissionUsecase rbacusecase.RolePermissionUsecase
	db                    *gorm.DB
}

func NewRoleHandler(ru rbacusecase.RoleUsecase, rpu rbacusecase.RolePermissionUsecase, db *gorm.DB) *RoleHandler {
	return &RoleHandler{
		roleUsecase:           ru,
		rolePermissionUsecase: rpu,
		db:                    db,
	}
}

func (h *RoleHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var createdRole interface{}
	err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx)
		roleUsecase := rbacusecase.NewRoleUsecase(roleRepo)

		role, err := roleUsecase.CreateRole(tenantID, input.Name)
		if err != nil {
			return err
		}
		createdRole = role
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat role", err)
		return
	}

	util.SuccessResponse(c, "Role berhasil dibuat", createdRole)
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	roles, err := h.roleUsecase.ListRoles(tenantID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data role", err)
		return
	}

	util.SuccessResponse(c, "Data role ditemukan", roles)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	role, err := h.roleUsecase.GetRoleByID(tenantID, roleID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil role", err)
		return
	}

	util.SuccessResponse(c, "Role ditemukan", role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var updatedRole interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx)
		roleUsecase := rbacusecase.NewRoleUsecase(roleRepo)

		role, err := roleUsecase.UpdateRole(tenantID, roleID, input.Name)
		if err != nil {
			return err
		}
		updatedRole = role
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah role", err)
		return
	}

	util.SuccessResponse(c, "Role berhasil diubah", updatedRole)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx)
		roleUsecase := rbacusecase.NewRoleUsecase(roleRepo)

		if err := roleUsecase.DeleteRole(tenantID, roleID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus role", err)
		return
	}

	util.SuccessResponse(c, "Role berhasil dihapus", nil)
}

func (h *RoleHandler) SetTenantAdmin(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	var input struct {
		IsTenantAdmin *bool `json:"is_tenant_admin" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}
	if input.IsTenantAdmin == nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", util.ErrValidation("is_tenant_admin wajib diisi"))
		return
	}

	var updatedRole interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roleRepo := repository.NewRoleRepository(tx)
		roleUsecase := rbacusecase.NewRoleUsecase(roleRepo)

		role, err := roleUsecase.SetTenantAdmin(roleID, *input.IsTenantAdmin)
		if err != nil {
			return err
		}
		updatedRole = role
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah tenant admin", err)
		return
	}

	util.SuccessResponse(c, "Role tenant admin berhasil diubah", updatedRole)
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	var input struct {
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		tx = tx.Set("audit_action", "ASSIGN")
		roleRepo := repository.NewRoleRepository(tx)
		permissionRepo := repository.NewPermissionRepository(tx)
		rolePermissionUsecase := rbacusecase.NewRolePermissionUsecase(roleRepo, permissionRepo)

		if err := rolePermissionUsecase.AssignPermissions(tenantID, roleID, input.PermissionIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal assign permission", err)
		return
	}

	util.SuccessResponse(c, "Permission berhasil di-assign", nil)
}

func (h *RoleHandler) RevokePermissions(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format role ID tidak valid", err)
		return
	}

	var input struct {
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		tx = tx.Set("audit_action", "REVOKE").
			Set("audit_record_id", fmt.Sprintf("role_id=%d,permission_ids=%v", roleID, input.PermissionIDs))
		roleRepo := repository.NewRoleRepository(tx)
		permissionRepo := repository.NewPermissionRepository(tx)
		rolePermissionUsecase := rbacusecase.NewRolePermissionUsecase(roleRepo, permissionRepo)

		if err := rolePermissionUsecase.RevokePermissions(tenantID, roleID, input.PermissionIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal revoke permission", err)
		return
	}

	util.SuccessResponse(c, "Permission berhasil di-revoke", nil)
}
