package handlers

import (
	"net/http"
	"strings"

	notificationusecase "samrs-backend/internal/usecase/notification"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	usecase notificationusecase.NotificationUsecase
}

func NewNotificationHandler(u notificationusecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{usecase: u}
}

func (h *NotificationHandler) Send(c *gin.Context) {
	var input notificationusecase.NotificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	channel := strings.TrimSpace(input.Channel)
	to := strings.TrimSpace(input.To)
	message := strings.TrimSpace(input.Message)
	template := strings.TrimSpace(input.Template)
	if channel == "" || to == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "channel dan to wajib diisi", nil)
		return
	}
	if message == "" && template == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "message atau template wajib diisi", nil)
		return
	}

	if err := h.usecase.Send(c.Request.Context(), input); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengirim notifikasi", err)
		return
	}

	util.SuccessResponse(c, "Notifikasi diterima untuk diproses", gin.H{
		"channel": channel,
		"to":      to,
		"status":  "queued",
	})
}
