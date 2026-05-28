package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetEventHandler struct {
	usecase assetusecase.AssetEventUsecase
}

func NewAssetEventHandler(u assetusecase.AssetEventUsecase) *AssetEventHandler {
	return &AssetEventHandler{usecase: u}
}

func (h *AssetEventHandler) GetTimeline(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format asset ID tidak valid", err)
		return
	}

	pagination := httputil.ParsePagination(c)
	dateRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}

	filter := repository.AssetEventFilter{
		EventType: strings.TrimSpace(c.Query("event_type")),
		DateFrom:  dateRange.From,
		DateTo:    dateRange.To,
		Page:      pagination.Page,
		PerPage:   pagination.PerPage,
		SortBy:    pagination.SortBy,
		SortDir:   pagination.SortDir,
	}

	events, total, err := h.usecase.GetTimeline(tenantID, assetID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil timeline", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Timeline aset ditemukan", events, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}
