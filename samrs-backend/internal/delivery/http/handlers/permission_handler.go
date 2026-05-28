package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strconv"

	rbacusecase "samrs-backend/internal/usecase/rbac"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	usecase rbacusecase.PermissionUsecase
}

func NewPermissionHandler(u rbacusecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{u}
}

func (h *PermissionHandler) GetAll(c *gin.Context) {
	perms, err := h.usecase.ListAllPermissions()
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil permissions", err)
		return
	}
	util.SuccessResponse(c, "Data permissions ditemukan", perms)
}

func (h *PermissionHandler) GetByRole(c *gin.Context) {
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

	perms, err := h.usecase.ListPermissionsByRole(tenantID, roleID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil permissions", err)
		return
	}
	util.SuccessResponse(c, "Data permissions ditemukan", perms)
}
