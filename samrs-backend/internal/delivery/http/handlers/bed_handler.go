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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BedHandler struct {
	bedUsecase masterdatausecase.BedUsecase
	db         *gorm.DB
}

func NewBedHandler(bu masterdatausecase.BedUsecase, db *gorm.DB) *BedHandler {
	return &BedHandler{bedUsecase: bu, db: db}
}

func (h *BedHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	// 3. Bind JSON Input
	var input struct {
		RoomID string `json:"room_id" binding:"required"` // Gunakan string karena UUID datang dari JSON sebagai string
		Code   string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	// 4. Konversi RoomID dari string ke uuid.UUID
	roomUUID, err := uuid.Parse(input.RoomID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format Room ID tidak valid", err)
		return
	}

	var bed *domain.Bed
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		bedRepo := repository.NewBedRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedUsecase := masterdatausecase.NewBedUsecase(bedRepo, roomRepo)

		created, err := bedUsecase.CreateBed(tenantID, roomUUID, input.Code)
		if err != nil {
			return err
		}
		bed = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat bed", err)
		return
	}

	util.SuccessResponse(c, "Bed berhasil dibuat", bed)
}

func (h *BedHandler) GetAll(c *gin.Context) {
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

	var roomID *uuid.UUID
	if roomStr := strings.TrimSpace(c.Query("room_id")); roomStr != "" {
		parsed, err := uuid.Parse(roomStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format room_id tidak valid", err)
			return
		}
		roomID = &parsed
	}

	filter := repository.BedFilter{
		Search:   strings.TrimSpace(c.Query("search")),
		RoomID:   roomID,
		Status:   strings.TrimSpace(c.Query("status")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	beds, total, err := h.bedUsecase.GetBeds(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data bed", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data bed ditemukan", beds, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *BedHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	bedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format bed ID tidak valid", err)
		return
	}

	bed, err := h.bedUsecase.GetBedByID(tenantID, uint(bedID))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil bed", err)
		return
	}

	util.SuccessResponse(c, "Bed ditemukan", bed)
}

func (h *BedHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	bedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format bed ID tidak valid", err)
		return
	}

	var input struct {
		RoomID string `json:"room_id" binding:"required"`
		Code   string `json:"code" binding:"required"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	roomUUID, err := uuid.Parse(input.RoomID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format Room ID tidak valid", err)
		return
	}

	var bed *domain.Bed
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		bedRepo := repository.NewBedRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedUsecase := masterdatausecase.NewBedUsecase(bedRepo, roomRepo)

		updated, err := bedUsecase.UpdateBed(tenantID, uint(bedID), roomUUID, input.Code, input.Status)
		if err != nil {
			return err
		}
		bed = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah bed", err)
		return
	}

	util.SuccessResponse(c, "Bed berhasil diubah", bed)
}

func (h *BedHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	bedID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format bed ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		bedRepo := repository.NewBedRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedUsecase := masterdatausecase.NewBedUsecase(bedRepo, roomRepo)

		return bedUsecase.DeleteBed(tenantID, uint(bedID))
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus bed", err)
		return
	}

	util.SuccessResponse(c, "Bed berhasil dihapus", nil)
}
