package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	maintenanceusecase "samrs-backend/internal/usecase/maintenance"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceScheduleHandler struct {
	usecase maintenanceusecase.MaintenanceScheduleUsecase
	db      *gorm.DB
}

func NewMaintenanceScheduleHandler(u maintenanceusecase.MaintenanceScheduleUsecase, db *gorm.DB) *MaintenanceScheduleHandler {
	return &MaintenanceScheduleHandler{usecase: u, db: db}
}

func (h *MaintenanceScheduleHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	var input struct {
		AssetID      string `json:"asset_id" binding:"required"`
		ScheduleType string `json:"schedule_type"`
		Title        string `json:"title" binding:"required"`
		IntervalDays int    `json:"interval_days"`
		NextDueDate  string `json:"next_due_date" binding:"required"`
		Status       string `json:"status"`
		Notes        string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	assetID, err := uuid.Parse(input.AssetID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
		return
	}

	var schedule interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		scheduleRepo := repository.NewMaintenanceScheduleRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		scheduleUsecase := maintenanceusecase.NewMaintenanceScheduleUsecase(scheduleRepo, assetRepo)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		created, err := scheduleUsecase.CreateSchedule(maintenanceusecase.MaintenanceScheduleInput{
			TenantID:     tenantID,
			AssetID:      assetID,
			ScheduleType: input.ScheduleType,
			Title:        input.Title,
			IntervalDays: input.IntervalDays,
			NextDueDate:  input.NextDueDate,
			Status:       input.Status,
			Notes:        input.Notes,
		})
		if err != nil {
			return err
		}
		schedule = created

		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     assetID,
			UserID:      userID,
			EventType:   domain.AssetEventUpdated,
			Description: "Jadwal " + created.ScheduleType + " dibuat",
		})
		return err
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat jadwal", err)
		return
	}

	util.SuccessResponse(c, "Jadwal berhasil dibuat", schedule)
}

func (h *MaintenanceScheduleHandler) GetAll(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	pagination := httputil.ParsePagination(c)
	createdRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}
	dueRange, err := httputil.ParseDateRange(c, "due_from", "due_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format due date tidak valid", err)
		return
	}

	filter := repository.MaintenanceScheduleFilter{
		ScheduleType: strings.TrimSpace(c.Query("schedule_type")),
		Status:       strings.TrimSpace(c.Query("status")),
		DateFrom:     createdRange.From,
		DateTo:       createdRange.To,
		DueFrom:      dueRange.From,
		DueTo:        dueRange.To,
		Page:         pagination.Page,
		PerPage:      pagination.PerPage,
		SortBy:       pagination.SortBy,
		SortDir:      pagination.SortDir,
	}

	if assetStr := strings.TrimSpace(c.Query("asset_id")); assetStr != "" {
		assetID, err := uuid.Parse(assetStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
			return
		}
		filter.AssetID = assetID
	}

	schedules, total, err := h.usecase.GetAllSchedules(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data jadwal", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data jadwal ditemukan", schedules, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *MaintenanceScheduleHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	schedule, err := h.usecase.GetScheduleByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil jadwal", err)
		return
	}

	util.SuccessResponse(c, "Jadwal ditemukan", schedule)
}

func (h *MaintenanceScheduleHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	var input struct {
		ScheduleType string `json:"schedule_type"`
		Title        string `json:"title" binding:"required"`
		IntervalDays int    `json:"interval_days"`
		NextDueDate  string `json:"next_due_date" binding:"required"`
		Status       string `json:"status"`
		Notes        string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var schedule interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		scheduleRepo := repository.NewMaintenanceScheduleRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		scheduleUsecase := maintenanceusecase.NewMaintenanceScheduleUsecase(scheduleRepo, assetRepo)
		updated, err := scheduleUsecase.UpdateSchedule(maintenanceusecase.MaintenanceScheduleUpdateInput{
			TenantID:     tenantID,
			ID:           uint(id),
			ScheduleType: input.ScheduleType,
			Title:        input.Title,
			IntervalDays: input.IntervalDays,
			NextDueDate:  input.NextDueDate,
			Status:       input.Status,
			Notes:        input.Notes,
		})
		if err != nil {
			return err
		}
		schedule = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah jadwal", err)
		return
	}

	util.SuccessResponse(c, "Jadwal berhasil diubah", schedule)
}

func (h *MaintenanceScheduleHandler) Complete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	userID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	var input struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var schedule *domain.MaintenanceSchedule
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		scheduleRepo := repository.NewMaintenanceScheduleRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		scheduleUsecase := maintenanceusecase.NewMaintenanceScheduleUsecase(scheduleRepo, assetRepo)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		updated, err := scheduleUsecase.CompleteSchedule(tenantID, uint(id), input.Note)
		if err != nil {
			return err
		}
		schedule = updated

		eventType := domain.AssetEventMaintenanceDone
		if updated.ScheduleType == domain.ScheduleTypeCalibration {
			eventType = domain.AssetEventCalibrationDone
		}
		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     updated.AssetID,
			UserID:      userID,
			EventType:   eventType,
			Description: "Jadwal " + updated.ScheduleType + " diselesaikan",
		})
		return err
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menyelesaikan jadwal", err)
		return
	}

	util.SuccessResponse(c, "Jadwal berhasil diselesaikan", schedule)
}

func (h *MaintenanceScheduleHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format schedule ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		scheduleRepo := repository.NewMaintenanceScheduleRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		scheduleUsecase := maintenanceusecase.NewMaintenanceScheduleUsecase(scheduleRepo, assetRepo)
		return scheduleUsecase.DeleteSchedule(tenantID, uint(id))
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus jadwal", err)
		return
	}

	util.SuccessResponse(c, "Jadwal berhasil dihapus", nil)
}
