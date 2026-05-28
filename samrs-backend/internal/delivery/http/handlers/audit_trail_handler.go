package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strconv"
	"strings"

	"samrs-backend/internal/repository"
	auditusecase "samrs-backend/internal/usecase/audit"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditTrailHandler struct {
	usecase auditusecase.AuditTrailUsecase
}

func NewAuditTrailHandler(u auditusecase.AuditTrailUsecase) *AuditTrailHandler {
	return &AuditTrailHandler{u}
}

func (h *AuditTrailHandler) GetAll(c *gin.Context) {
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

	filter := repository.AuditTrailFilter{
		Action:   strings.TrimSpace(c.Query("action")),
		Table:    strings.TrimSpace(c.Query("table")),
		RecordID: strings.TrimSpace(c.Query("record_id")),
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	if userStr := strings.TrimSpace(c.Query("user_id")); userStr != "" {
		userID, err := uuid.Parse(userStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format user_id tidak valid", err)
			return
		}
		filter.UserID = userID
	}

	items, total, err := h.usecase.ListByTenant(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil audit trail", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data audit trail ditemukan", items, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AuditTrailHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format audit ID tidak valid", err)
		return
	}

	audit, err := h.usecase.GetByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil audit trail", err)
		return
	}

	util.SuccessResponse(c, "Audit trail ditemukan", audit)
}
