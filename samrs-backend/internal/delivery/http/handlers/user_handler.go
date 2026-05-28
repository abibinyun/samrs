package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/repository"
	userusecase "samrs-backend/internal/usecase/user"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	usecase userusecase.UserUsecase
	db      *gorm.DB
}

func NewUserHandler(u userusecase.UserUsecase, db *gorm.DB) *UserHandler {
	return &UserHandler{usecase: u, db: db}
}

func (h *UserHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleID   int    `json:"role_id" binding:"required"`
		IsActive *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	var createdUser interface{}
	err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		user, err := userUsecase.CreateUser(userusecase.CreateUserInput{
			TenantID: tenantID,
			RoleID:   input.RoleID,
			Username: input.Username,
			Password: input.Password,
			IsActive: isActive,
		})
		if err != nil {
			return err
		}

		createdUser = user
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat user", err)
		return
	}

	util.SuccessResponse(c, "User berhasil dibuat", createdUser)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	pagination := httputil.ParsePagination(c)
	filter := repository.UserFilter{
		Search:  strings.TrimSpace(c.Query("search")),
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	}

	users, total, err := h.usecase.ListUsers(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data user", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data user ditemukan", users, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format user ID tidak valid", err)
		return
	}

	var input struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var updatedUser interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		user, err := userUsecase.UpdateUserProfile(tenantID, userID, input.Username)
		if err != nil {
			return err
		}
		updatedUser = user
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah user", err)
		return
	}

	util.SuccessResponse(c, "User berhasil diubah", updatedUser)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format user ID tidak valid", err)
		return
	}

	var input struct {
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var updatedUser interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		user, err := userUsecase.UpdateUserRole(tenantID, userID, input.RoleID)
		if err != nil {
			return err
		}
		updatedUser = user
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah role user", err)
		return
	}

	util.SuccessResponse(c, "Role user berhasil diubah", updatedUser)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format user ID tidak valid", err)
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var updatedUser interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		user, err := userUsecase.ResetUserPassword(tenantID, userID, input.Password)
		if err != nil {
			return err
		}
		updatedUser = user
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal reset password", err)
		return
	}

	util.SuccessResponse(c, "Password user berhasil direset", updatedUser)
}

func (h *UserHandler) SetActive(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format user ID tidak valid", err)
		return
	}

	var input struct {
		IsActive *bool `json:"is_active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}
	if input.IsActive == nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", util.ErrValidation("is_active wajib diisi"))
		return
	}

	var updatedUser interface{}
	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		user, err := userUsecase.SetUserActive(tenantID, userID, *input.IsActive)
		if err != nil {
			return err
		}
		updatedUser = user
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah status user", err)
		return
	}

	util.SuccessResponse(c, "Status user berhasil diubah", updatedUser)
}

func (h *UserHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format user ID tidak valid", err)
		return
	}

	err = httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		userRepo := repository.NewUserRepository(tx)
		roleRepo := repository.NewRoleRepository(tx)
		userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)

		if err := userUsecase.DeleteUser(tenantID, userID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus user", err)
		return
	}

	util.SuccessResponse(c, "User berhasil dihapus", nil)
}
