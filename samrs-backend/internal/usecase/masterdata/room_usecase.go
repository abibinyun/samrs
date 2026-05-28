package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type RoomUsecase interface {
	CreateRoom(room *domain.Room) error
	GetAllRooms(tenantID uuid.UUID, filter repository.RoomFilter) ([]domain.Room, int64, error)
	GetRoomByID(tenantID uuid.UUID, roomID uuid.UUID) (*domain.Room, error)
	UpdateRoom(tenantID uuid.UUID, roomID uuid.UUID, name, code, location string) (*domain.Room, error)
	DeleteRoom(tenantID uuid.UUID, roomID uuid.UUID) error
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(repo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{repo}
}

func (u *roomUsecase) CreateRoom(room *domain.Room) error {
	if strings.TrimSpace(room.Name) == "" || strings.TrimSpace(room.Code) == "" {
		return util.ErrValidation("nama dan kode ruangan wajib diisi")
	}

	exists, err := u.roomRepo.ExistsByCode(room.TenantID, room.Code, nil)
	if err != nil {
		return err
	}
	if exists {
		return util.ErrConflict("kode ruangan sudah digunakan")
	}

	return u.roomRepo.Create(room)
}

func (u *roomUsecase) GetAllRooms(tenantID uuid.UUID, filter repository.RoomFilter) ([]domain.Room, int64, error) {
	return u.roomRepo.FindAllByTenant(tenantID, filter)
}

func (u *roomUsecase) GetRoomByID(tenantID uuid.UUID, roomID uuid.UUID) (*domain.Room, error) {
	return u.roomRepo.FindByIDAndTenant(roomID, tenantID)
}

func (u *roomUsecase) UpdateRoom(tenantID uuid.UUID, roomID uuid.UUID, name, code, location string) (*domain.Room, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(code) == "" {
		return nil, util.ErrValidation("nama dan kode ruangan wajib diisi")
	}

	room, err := u.roomRepo.FindByIDAndTenant(roomID, tenantID)
	if err != nil {
		return nil, err
	}

	exists, err := u.roomRepo.ExistsByCode(tenantID, code, &room.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode ruangan sudah digunakan")
	}

	room.Name = name
	room.Code = code
	room.Location = location
	if err := u.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (u *roomUsecase) DeleteRoom(tenantID uuid.UUID, roomID uuid.UUID) error {
	room, err := u.roomRepo.FindByIDAndTenant(roomID, tenantID)
	if err != nil {
		return err
	}
	return u.roomRepo.Delete(room)
}
