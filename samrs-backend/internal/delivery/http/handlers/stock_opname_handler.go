package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	stockopnameusecase "samrs-backend/internal/usecase/stockopname"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockOpnameHandler struct {
	usecase stockopnameusecase.StockOpnameUsecase
	db      *gorm.DB
}

func NewStockOpnameHandler(u stockopnameusecase.StockOpnameUsecase, db *gorm.DB) *StockOpnameHandler {
	return &StockOpnameHandler{usecase: u, db: db}
}

func (h *StockOpnameHandler) CreateSession(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	var input struct {
		Title    string `json:"title" binding:"required"`
		OpnameAt string `json:"opname_at"`
		Notes    string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var session *domain.StockOpnameSession
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		sessionRepo := repository.NewStockOpnameRepository(tx)
		itemRepo := repository.NewStockOpnameItemRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		stockOpnameUsecase := stockopnameusecase.NewStockOpnameUsecase(sessionRepo, itemRepo, assetRepo)

		created, err := stockOpnameUsecase.CreateSession(stockopnameusecase.StockOpnameSessionInput{
			TenantID: tenantID,
			Title:    input.Title,
			OpnameAt: input.OpnameAt,
			Notes:    input.Notes,
		})
		if err != nil {
			return err
		}
		session = created
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat stock opname", err)
		return
	}

	util.SuccessResponse(c, "Stock opname berhasil dibuat", session)
}

func (h *StockOpnameHandler) GetAllSessions(c *gin.Context) {
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

	filter := repository.StockOpnameFilter{
		Status:   strings.TrimSpace(c.Query("status")),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	items, total, err := h.usecase.GetAllSessions(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data stock opname ditemukan", items, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *StockOpnameHandler) GetSessionByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format session ID tidak valid", err)
		return
	}

	session, err := h.usecase.GetSessionByID(tenantID, uint(id))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil session", err)
		return
	}

	util.SuccessResponse(c, "Session ditemukan", session)
}

func (h *StockOpnameHandler) UpdateSession(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format session ID tidak valid", err)
		return
	}

	var input struct {
		Title    string `json:"title" binding:"required"`
		OpnameAt string `json:"opname_at"`
		Notes    string `json:"notes"`
		Status   string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var session *domain.StockOpnameSession
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		sessionRepo := repository.NewStockOpnameRepository(tx)
		itemRepo := repository.NewStockOpnameItemRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		stockOpnameUsecase := stockopnameusecase.NewStockOpnameUsecase(sessionRepo, itemRepo, assetRepo)

		updated, err := stockOpnameUsecase.UpdateSession(stockopnameusecase.StockOpnameSessionUpdateInput{
			TenantID: tenantID,
			ID:       uint(id),
			Title:    input.Title,
			OpnameAt: input.OpnameAt,
			Notes:    input.Notes,
			Status:   input.Status,
		})
		if err != nil {
			return err
		}
		session = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah session", err)
		return
	}

	util.SuccessResponse(c, "Session berhasil diubah", session)
}

func (h *StockOpnameHandler) CloseSession(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	id, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format session ID tidak valid", err)
		return
	}

	var session *domain.StockOpnameSession
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		sessionRepo := repository.NewStockOpnameRepository(tx)
		itemRepo := repository.NewStockOpnameItemRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		stockOpnameUsecase := stockopnameusecase.NewStockOpnameUsecase(sessionRepo, itemRepo, assetRepo)

		updated, err := stockOpnameUsecase.CloseSession(tenantID, uint(id))
		if err != nil {
			return err
		}
		session = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menutup session", err)
		return
	}

	util.SuccessResponse(c, "Session berhasil ditutup", session)
}

func (h *StockOpnameHandler) AddItem(c *gin.Context) {
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

	sessionID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format session ID tidak valid", err)
		return
	}

	var input struct {
		AssetID   string `json:"asset_id" binding:"required"`
		Condition string `json:"condition"`
		Note      string `json:"note"`
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

	var item *domain.StockOpnameItem
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		sessionRepo := repository.NewStockOpnameRepository(tx)
		itemRepo := repository.NewStockOpnameItemRepository(tx)
		assetRepo := repository.NewAssetRepository(tx)
		eventRepo := repository.NewAssetEventRepository(tx)
		userRepo := repository.NewUserRepository(tx)
		stockOpnameUsecase := stockopnameusecase.NewStockOpnameUsecase(sessionRepo, itemRepo, assetRepo)
		eventUsecase := assetusecase.NewAssetEventUsecase(eventRepo, assetRepo, userRepo)

		created, err := stockOpnameUsecase.AddItem(stockopnameusecase.StockOpnameItemInput{
			TenantID:  tenantID,
			SessionID: uint(sessionID),
			AssetID:   assetID,
			Condition: input.Condition,
			Note:      input.Note,
			CheckedBy: userID,
		})
		if err != nil {
			return err
		}
		item = created

		_, err = eventUsecase.CreateEvent(assetusecase.AssetEventInput{
			TenantID:    tenantID,
			AssetID:     assetID,
			UserID:      userID,
			EventType:   domain.AssetEventStockOpname,
			Description: "Stock opname: " + created.Condition,
		})
		return err
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menambah item", err)
		return
	}

	util.SuccessResponse(c, "Item opname berhasil ditambah", item)
}

func (h *StockOpnameHandler) ListItems(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	sessionID, err := httputil.ParseUint(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format session ID tidak valid", err)
		return
	}

	items, err := h.usecase.ListItems(tenantID, uint(sessionID))
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil item", err)
		return
	}

	util.SuccessResponse(c, "Item ditemukan", items)
}
