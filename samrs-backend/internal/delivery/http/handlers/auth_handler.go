package handlers

import (
	"fmt"
	"net/http"
	authusecase "samrs-backend/internal/usecase/auth"
	"samrs-backend/pkg/util" // Import helper baru

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUsecase authusecase.AuthUsecase
}

func NewAuthHandler(au authusecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{au}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Username dan password wajib diisi", nil)
		return
	}

	token, err := h.authUsecase.Login(input.Username, input.Password)
	if err != nil {
		util.ErrorResponseFromErr(c, "Login gagal", err)
		return
	}

	util.SuccessResponse(c, "Login berhasil", gin.H{"token": token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format User ID tidak valid", util.ErrUnauthorized("invalid user id"))
		return
	}

	user, permissions, err := h.authUsecase.GetMe(userID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data user", err)
		return
	}

	util.SuccessResponse(c, "Data user ditemukan", gin.H{
		"user":        user,
		"permissions": permissions,
	})
}
