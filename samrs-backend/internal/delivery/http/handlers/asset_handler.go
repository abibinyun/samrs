package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type AssetHandler struct {
	usecase    assetusecase.AssetUsecase
	tenantRepo repository.TenantRepository
	db         *gorm.DB
}

func NewAssetHandler(u assetusecase.AssetUsecase, tr repository.TenantRepository, db *gorm.DB) *AssetHandler {
	return &AssetHandler{usecase: u, tenantRepo: tr, db: db}
}

func (h *AssetHandler) Create(c *gin.Context) {
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
		CategoryID   uint   `json:"category_id" binding:"required"`
		RoomID       string `json:"room_id"`
		BedID        uint   `json:"bed_id"`
		VendorID     uint   `json:"vendor_id"`
		BrandID      uint   `json:"brand_id"`
		ModelID      uint   `json:"model_id"`
		Code         string `json:"code" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Brand        string `json:"brand"`
		Model        string `json:"model"`
		Status       string `json:"status"`
		PurchaseDate string `json:"purchase_date"` // YYYY-MM-DD
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var roomID *uuid.UUID
	if strings.TrimSpace(input.RoomID) != "" {
		parsed, err := uuid.Parse(input.RoomID)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format Room ID tidak valid", err)
			return
		}
		roomID = &parsed
	}

	var bedID *uint
	if input.BedID != 0 {
		bedID = &input.BedID
	}
	var vendorID *uint
	if input.VendorID != 0 {
		vendorID = &input.VendorID
	}
	var brandID *uint
	if input.BrandID != 0 {
		brandID = &input.BrandID
	}
	var modelID *uint
	if input.ModelID != 0 {
		modelID = &input.ModelID
	}

	var asset *domain.Asset
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		assetRepo := repository.NewAssetRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		vendorRepo := repository.NewVendorRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelRepo := repository.NewAssetModelRepository(tx)
		statusRepo := repository.NewAssetStatusRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		assetUsecase := assetusecase.NewAssetUsecase(
			assetRepo,
			categoryRepo,
			roomRepo,
			bedRepo,
			vendorRepo,
			brandRepo,
			modelRepo,
			statusRepo,
		)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		created, err := assetUsecase.CreateAsset(assetusecase.CreateAssetInput{
			TenantID:     tenantID,
			CategoryID:   input.CategoryID,
			RoomID:       roomID,
			BedID:        bedID,
			VendorID:     vendorID,
			BrandID:      brandID,
			ModelID:      modelID,
			Code:         input.Code,
			Name:         input.Name,
			Brand:        input.Brand,
			Model:        input.Model,
			Status:       input.Status,
			PurchaseDate: input.PurchaseDate,
		})
		if err != nil {
			return err
		}
		asset = created
		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     asset.ID,
			UserID:      userID,
			EventType:   domain.AssetEventCreated,
			Description: "Asset dibuat",
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat aset", err)
		return
	}

	util.SuccessResponse(c, "Aset berhasil dibuat", asset)
}

func (h *AssetHandler) GetAll(c *gin.Context) {
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

	purchaseRange, err := httputil.ParseDateRange(c, "purchase_from", "purchase_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format purchase_date tidak valid", err)
		return
	}

	filter := repository.AssetFilter{
		Status:       strings.TrimSpace(c.Query("status")),
		Search:       strings.TrimSpace(c.Query("search")),
		DateFrom:     createdRange.From,
		DateTo:       createdRange.To,
		PurchaseFrom: purchaseRange.From,
		PurchaseTo:   purchaseRange.To,
		Page:         pagination.Page,
		PerPage:      pagination.PerPage,
		SortBy:       pagination.SortBy,
		SortDir:      pagination.SortDir,
	}

	if categoryStr := strings.TrimSpace(c.Query("category_id")); categoryStr != "" {
		id, err := strconv.ParseUint(categoryStr, 10, 32)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format category_id tidak valid", err)
			return
		}
		filter.CategoryID = uint(id)
	}

	if roomStr := strings.TrimSpace(c.Query("room_id")); roomStr != "" {
		roomID, err := uuid.Parse(roomStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format room_id tidak valid", err)
			return
		}
		filter.RoomID = roomID
	}

	if bedStr := strings.TrimSpace(c.Query("bed_id")); bedStr != "" {
		bedID, err := strconv.ParseUint(bedStr, 10, 32)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format bed_id tidak valid", err)
			return
		}
		filter.BedID = uint(bedID)
	}

	assets, total, err := h.usecase.GetAllAssets(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data aset ditemukan", assets, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *AssetHandler) GetByID(c *gin.Context) {
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

	asset, err := h.usecase.GetAssetByID(tenantID, assetID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil aset", err)
		return
	}

	util.SuccessResponse(c, "Aset ditemukan", asset)
}

func (h *AssetHandler) Update(c *gin.Context) {
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
	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format asset ID tidak valid", err)
		return
	}

	var input struct {
		CategoryID   uint   `json:"category_id" binding:"required"`
		RoomID       string `json:"room_id"`
		BedID        uint   `json:"bed_id"`
		VendorID     uint   `json:"vendor_id"`
		BrandID      uint   `json:"brand_id"`
		ModelID      uint   `json:"model_id"`
		Code         string `json:"code" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Brand        string `json:"brand"`
		Model        string `json:"model"`
		Status       string `json:"status"`
		PurchaseDate string `json:"purchase_date"` // YYYY-MM-DD
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var roomID *uuid.UUID
	if strings.TrimSpace(input.RoomID) != "" {
		parsed, err := uuid.Parse(input.RoomID)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format Room ID tidak valid", err)
			return
		}
		roomID = &parsed
	}

	var bedID *uint
	if input.BedID != 0 {
		bedID = &input.BedID
	}
	var vendorID *uint
	if input.VendorID != 0 {
		vendorID = &input.VendorID
	}
	var brandID *uint
	if input.BrandID != 0 {
		brandID = &input.BrandID
	}
	var modelID *uint
	if input.ModelID != 0 {
		modelID = &input.ModelID
	}

	var asset *domain.Asset
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		assetRepo := repository.NewAssetRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		vendorRepo := repository.NewVendorRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelRepo := repository.NewAssetModelRepository(tx)
		statusRepo := repository.NewAssetStatusRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		assetUsecase := assetusecase.NewAssetUsecase(
			assetRepo,
			categoryRepo,
			roomRepo,
			bedRepo,
			vendorRepo,
			brandRepo,
			modelRepo,
			statusRepo,
		)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		oldAsset, err := assetRepo.FindByIDAndTenant(assetID, tenantID)
		if err != nil {
			return err
		}
		updated, err := assetUsecase.UpdateAsset(assetusecase.UpdateAssetInput{
			TenantID:     tenantID,
			ID:           assetID,
			CategoryID:   input.CategoryID,
			RoomID:       roomID,
			BedID:        bedID,
			VendorID:     vendorID,
			BrandID:      brandID,
			ModelID:      modelID,
			Code:         input.Code,
			Name:         input.Name,
			Brand:        input.Brand,
			Model:        input.Model,
			Status:       input.Status,
			PurchaseDate: input.PurchaseDate,
		})
		if err != nil {
			return err
		}
		asset = updated

		desc := "Asset diperbarui"
		eventType := domain.AssetEventUpdated
		if strings.TrimSpace(oldAsset.Status) != strings.TrimSpace(asset.Status) {
			desc = "Status aset berubah dari " + oldAsset.Status + " ke " + asset.Status
			eventType = domain.AssetEventStatusChanged
		}
		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     asset.ID,
			UserID:      userID,
			EventType:   eventType,
			Description: desc,
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah aset", err)
		return
	}

	util.SuccessResponse(c, "Aset berhasil diubah", asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
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
	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format asset ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		assetRepo := repository.NewAssetRepository(tx)
		categoryRepo := repository.NewCategoryRepository(tx)
		roomRepo := repository.NewRoomRepository(tx)
		bedRepo := repository.NewBedRepository(tx)
		vendorRepo := repository.NewVendorRepository(tx)
		brandRepo := repository.NewAssetBrandRepository(tx)
		modelRepo := repository.NewAssetModelRepository(tx)
		statusRepo := repository.NewAssetStatusRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		assetUsecase := assetusecase.NewAssetUsecase(
			assetRepo,
			categoryRepo,
			roomRepo,
			bedRepo,
			vendorRepo,
			brandRepo,
			modelRepo,
			statusRepo,
		)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		_, err := eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     assetID,
			UserID:      userID,
			EventType:   domain.AssetEventDeleted,
			Description: "Asset dihapus",
		})
		if err != nil {
			return err
		}
		if err := assetUsecase.DeleteAsset(tenantID, assetID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus aset", err)
		return
	}

	util.SuccessResponse(c, "Aset berhasil dihapus", nil)
}

func (h *AssetHandler) GetQRCode(c *gin.Context) {
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

	asset, err := h.usecase.GetAssetByID(tenantID, assetID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Aset tidak ditemukan", err)
		return
	}

	tenant, err := h.tenantRepo.FindByID(tenantID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Tenant tidak ditemukan", err)
		return
	}

	baseURL := strings.TrimRight(os.Getenv("PUBLIC_ASSET_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = requestBaseURL(c)
	}

	publicURL := fmt.Sprintf("%s/public/tenants/%s/assets/%s", baseURL, tenant.Slug, asset.Code)
	size := httputil.ParseIntWithDefault(c.Query("size"), 256)
	if size < 128 {
		size = 128
	}
	if size > 512 {
		size = 512
	}

	png, err := qrcode.Encode(publicURL, qrcode.Medium, size)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat QR code", err)
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func requestBaseURL(c *gin.Context) string {
	proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))
	host := strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	if host == "" {
		host = c.Request.Host
	}
	return proto + "://" + host
}
