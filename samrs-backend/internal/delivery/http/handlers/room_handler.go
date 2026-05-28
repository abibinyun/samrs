package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"net/http"
	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	masterdatausecase "samrs-backend/internal/usecase/masterdata"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomHandler struct {
	roomUsecase masterdatausecase.RoomUsecase
	db          *gorm.DB
}

func NewRoomHandler(u masterdatausecase.RoomUsecase, db *gorm.DB) *RoomHandler {
	return &RoomHandler{roomUsecase: u, db: db}
}

func (h *RoomHandler) Create(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	// 2. Bind JSON Input
	var input struct {
		Name     string `json:"name" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Location string `json:"location"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	// 3. Mapping ke Domain
	room := &domain.Room{
		TenantID: tenantID,
		Name:     input.Name,
		Code:     input.Code,
		Location: input.Location,
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roomRepo := repository.NewRoomRepository(tx)
		roomUsecase := masterdatausecase.NewRoomUsecase(roomRepo)
		return roomUsecase.CreateRoom(room)
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal membuat ruangan", err)
		return
	}

	util.SuccessResponse(c, "Ruangan berhasil dibuat", room)
}

func (h *RoomHandler) GetAll(c *gin.Context) {
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

	filter := repository.RoomFilter{
		Search:   c.Query("search"),
		DateFrom: dateRange.From,
		DateTo:   dateRange.To,
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		SortBy:   pagination.SortBy,
		SortDir:  pagination.SortDir,
	}

	rooms, total, err := h.roomUsecase.GetAllRooms(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}

	util.SuccessResponseWithMeta(c, "Data ruangan ditemukan", rooms, httputil.PaginationMeta{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

func (h *RoomHandler) GetByID(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format room ID tidak valid", err)
		return
	}

	room, err := h.roomUsecase.GetRoomByID(tenantID, roomID)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil ruangan", err)
		return
	}

	util.SuccessResponse(c, "Ruangan ditemukan", room)
}

func (h *RoomHandler) Update(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format room ID tidak valid", err)
		return
	}

	var input struct {
		Name     string `json:"name" binding:"required"`
		Code     string `json:"code" binding:"required"`
		Location string `json:"location"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ErrorResponseFromErr(c, "Input tidak valid", err)
		return
	}

	var room *domain.Room
	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roomRepo := repository.NewRoomRepository(tx)
		roomUsecase := masterdatausecase.NewRoomUsecase(roomRepo)

		updated, err := roomUsecase.UpdateRoom(tenantID, roomID, input.Name, input.Code, input.Location)
		if err != nil {
			return err
		}

		room = updated
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengubah ruangan", err)
		return
	}

	util.SuccessResponse(c, "Ruangan berhasil diubah", room)
}

func (h *RoomHandler) Delete(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		util.ErrorResponseFromErr(c, "Format room ID tidak valid", err)
		return
	}

	if err := httputil.WithTx(c, h.db, func(tx *gorm.DB) error {
		roomRepo := repository.NewRoomRepository(tx)
		roomUsecase := masterdatausecase.NewRoomUsecase(roomRepo)

		if err := roomUsecase.DeleteRoom(tenantID, roomID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ErrorResponseFromErr(c, "Gagal menghapus ruangan", err)
		return
	}

	util.SuccessResponse(c, "Ruangan berhasil dihapus", nil)
}
