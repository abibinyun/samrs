package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	complaintusecase "samrs-backend/internal/usecase/complaint"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComplaintHandler struct {
	usecase complaintusecase.ComplaintUsecase
	db      *gorm.DB
}

func NewComplaintHandler(u complaintusecase.ComplaintUsecase, db *gorm.DB) *ComplaintHandler {
	return &ComplaintHandler{usecase: u, db: db}
}

func (h *ComplaintHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	reporterID, ok := httputil.UserIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	var input struct {
		AssetID     string  `json:"asset_id" binding:"required"`
		AssignedTo  *string `json:"assigned_to"`
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
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

	var assignedTo *uuid.UUID
	if input.AssignedTo != nil && strings.TrimSpace(*input.AssignedTo) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*input.AssignedTo))
		if err != nil {
			util.ErrorResponseFromErr(c, "Format assigned_to tidak valid", err)
			return
		}
		assignedTo = &parsed
	}

	var complaint interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		complaintRepo := repository.NewComplaintRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		complaintUsecase := complaintusecase.NewComplaintUsecase(complaintRepo, assetRepo, userRepo)
		eventRepo := repository.NewAssetEventRepository(tx)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		created, err := complaintUsecase.CreateComplaint(complaintusecase.CreateComplaintInput{
			TenantID:    tenantID,
			AssetID:     assetID,
			ReportedBy:  reporterID,
			AssignedTo:  assignedTo,
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			return err
		}
		complaint = created
		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     assetID,
			UserID:      reporterID,
			EventType:   domain.AssetEventComplaintReported,
			Description: "Complaint dibuat: " + input.Title,
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat complaint", err)
		return
	}

	util.SuccessResponse(c, "Complaint berhasil dibuat", complaint)
}

func (h *ComplaintHandler) GetAll(c *gin.Context) {
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

	filter := repository.ComplaintFilter{
		Status:   strings.TrimSpace(c.Query("status")),
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: createdRange.From,
		DateTo:   createdRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	if assetStr := strings.TrimSpace(c.Query("asset_id")); assetStr != "" {
		assetID, err := uuid.Parse(assetStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
			return
		}
		filter.AssetID = assetID
	}
	if assignedStr := strings.TrimSpace(c.Query("assigned_to")); assignedStr != "" {
		assignedID, err := uuid.Parse(assignedStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format assigned_to tidak valid", err)
			return
		}
		filter.AssignedTo = assignedID
	}
	if reportedStr := strings.TrimSpace(c.Query("reported_by")); reportedStr != "" {
		reportedID, err := uuid.Parse(reportedStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format reported_by tidak valid", err)
			return
		}
		filter.ReportedBy = reportedID
	}

	complaints, total, err := h.usecase.GetAllComplaints(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data complaint ditemukan", complaints, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *ComplaintHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	complaintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format complaint ID tidak valid", err)
		return
	}

	complaint, err := h.usecase.GetComplaintByID(tenantID, complaintID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil complaint", err)
		return
	}

	util.SuccessResponse(c, "Complaint ditemukan", complaint)
}

func (h *ComplaintHandler) Update(c *gin.Context) {
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

	complaintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format complaint ID tidak valid", err)
		return
	}

	var input struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		Status         string  `json:"status"`
		AssignedTo     *string `json:"assigned_to"`
		ResolutionNote string  `json:"resolution_note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var (
		assignedTo    *uuid.UUID
		assignedToSet bool
	)
	if input.AssignedTo != nil {
		assignedToSet = true
		trimmed := strings.TrimSpace(*input.AssignedTo)
		if trimmed != "" {
			parsed, err := uuid.Parse(trimmed)
			if err != nil {
				util.ErrorResponseFromErr(c, "Format assigned_to tidak valid", err)
				return
			}
			assignedTo = &parsed
		}
	}

	var complaint interface{}
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		complaintRepo := repository.NewComplaintRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		complaintUsecase := complaintusecase.NewComplaintUsecase(complaintRepo, assetRepo, userRepo)
		eventRepo := repository.NewAssetEventRepository(tx)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		oldComplaint, err := complaintRepo.FindByIDAndTenant(complaintID, tenantID)
		if err != nil {
			return err
		}
		updated, err := complaintUsecase.UpdateComplaint(complaintusecase.UpdateComplaintInput{
			TenantID:       tenantID,
			ID:             complaintID,
			Title:          input.Title,
			Description:    input.Description,
			Status:         input.Status,
			AssignedTo:     assignedTo,
			AssignedToSet:  assignedToSet,
			ResolutionNote: input.ResolutionNote,
		})
		if err != nil {
			return err
		}
		complaint = updated

		if strings.TrimSpace(input.Status) != "" && oldComplaint.Status != updated.Status {
			eventType := domain.AssetEventComplaintProgress
			desc := "Complaint status: " + updated.Status
			if updated.Status == domain.ComplaintStatusDone {
				eventType = domain.AssetEventComplaintResolved
				desc = "Complaint diselesaikan"
			}
			_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
				TenantID:    tenantID,
				AssetID:     updated.AssetID,
				UserID:      userID,
				EventType:   eventType,
				Description: desc,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah complaint", err)
		return
	}

	util.SuccessResponse(c, "Complaint berhasil diubah", complaint)
}

func (h *ComplaintHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	complaintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format complaint ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		complaintRepo := repository.NewComplaintRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		complaintUsecase := complaintusecase.NewComplaintUsecase(complaintRepo, assetRepo, userRepo)

		return complaintUsecase.DeleteComplaint(tenantID, complaintID)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus complaint", err)
		return
	}

	util.SuccessResponse(c, "Complaint berhasil dihapus", nil)
}
