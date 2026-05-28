package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetMutationHandler struct {
	usecase assetusecase.AssetMutationUsecase
	db      *gorm.DB
}

func NewAssetMutationHandler(u assetusecase.AssetMutationUsecase, db *gorm.DB) *AssetMutationHandler {
	return &AssetMutationHandler{usecase: u, db: db}
}

func (h *AssetMutationHandler) Create(c *gin.Context) {
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
		AssetID string `json:"asset_id" binding:"required"`
		RoomID  string `json:"room_id"`
		BedID   uint   `json:"bed_id"`
		Reason  string `json:"reason"`
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

	var roomID *uuid.UUID
	if strings.TrimSpace(input.RoomID) != "" {
		parsed, err := uuid.Parse(input.RoomID)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format room_id tidak valid", err)
			return
		}
		roomID = &parsed
	}

	var bedID *uint
	if input.BedID != 0 {
		bedID = &input.BedID
	}

	var mutation *domain.AssetMutation
	var asset *domain.Asset
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		mutationRepo := repository.NewAssetMutationRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		userRepo := repository.NewUserRepository(tx)

		mutationUsecase := assetusecase.NewAssetMutationUsecase(mutationRepo, assetRepo, roomRepo, bedRepo)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		created, updatedAsset, err := mutationUsecase.MoveAsset(assetusecase.AssetMutationInput{
			TenantID: tenantID,
			AssetID:  assetID,
			ToRoomID: roomID,
			ToBedID:  bedID,
			Reason:   input.Reason,
			MovedBy:  userID,
		})
		if err != nil {
			return err
		}
		mutation = created
		asset = updatedAsset

		desc := "Mutasi aset"
		if roomID != nil {
			desc = desc + " ke ruangan baru"
		}
		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     asset.ID,
			UserID:      userID,
			EventType:   domain.AssetEventMutation,
			Description: desc,
			Meta: map[string]interface{}{
				"from_room_id": mutation.FromRoomID,
				"from_bed_id":  mutation.FromBedID,
				"to_room_id":   mutation.ToRoomID,
				"to_bed_id":    mutation.ToBedID,
			},
		})
		return err
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mutasi asset", err)
		return
	}

	util.SuccessResponse(c, "Mutasi asset berhasil", gin.H{
		"mutation": mutation,
		"asset":    asset,
	})
}

func (h *AssetMutationHandler) GetAll(c *gin.Context) {
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

	filter := repository.AssetMutationFilter{
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
		Search:   strings.TrimSpace(c.Query("search")),
	}

	if assetStr := strings.TrimSpace(c.Query("asset_id")); assetStr != "" {
		assetID, err := uuid.Parse(assetStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
			return
		}
		filter.AssetID = assetID
	}

	items, total, err := h.usecase.GetAllMutations(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data mutasi ditemukan", items, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AssetMutationHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format mutation ID tidak valid", err)
		return
	}

	mutation, err := h.usecase.GetMutationByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil mutasi", err)
		return
	}

	util.SuccessResponse(c, "Mutasi ditemukan", mutation)
}
