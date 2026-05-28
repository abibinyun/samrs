package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type BedUsecase interface {
	CreateBed(tenantID uuid.UUID, roomID uuid.UUID, code string) (*domain.Bed, error)
	GetBeds(tenantID uuid.UUID, filter repository.BedFilter) ([]domain.Bed, int64, error)
	GetBedByID(tenantID uuid.UUID, bedID uint) (*domain.Bed, error)
	UpdateBed(tenantID uuid.UUID, bedID uint, roomID uuid.UUID, code, status string) (*domain.Bed, error)
	DeleteBed(tenantID uuid.UUID, bedID uint) error
}

type bedUsecase struct {
	bedRepo  repository.BedRepository
	roomRepo repository.RoomRepository // Kita butuh roomRepo untuk validasi kepemilikan
}

func NewBedUsecase(br repository.BedRepository, rr repository.RoomRepository) BedUsecase {
	return &bedUsecase{br, rr}
}

func (u *bedUsecase) CreateBed(tenantID uuid.UUID, roomID uuid.UUID, code string) (*domain.Bed, error) {
	if _, err := u.roomRepo.FindByIDAndTenant(roomID, tenantID); err != nil {
		return nil, util.ErrNotFound("ruangan tidak ditemukan atau akses ditolak")
	}

	exists, err := u.bedRepo.ExistsByCode(tenantID, roomID, code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode bed sudah digunakan")
	}

	bed := &domain.Bed{
		TenantID: tenantID,
		RoomID:   roomID,
		Code:     code,
		Status:   domain.BedStatusAvailable,
	}

	if err := u.bedRepo.Create(bed); err != nil {
		return nil, err
	}
	return bed, nil
}

func (u *bedUsecase) GetBeds(tenantID uuid.UUID, filter repository.BedFilter) ([]domain.Bed, int64, error) {
	return u.bedRepo.FindAllByTenant(tenantID, filter)
}

func (u *bedUsecase) GetBedByID(tenantID uuid.UUID, bedID uint) (*domain.Bed, error) {
	return u.bedRepo.FindByIDAndTenant(bedID, tenantID)
}

func (u *bedUsecase) UpdateBed(tenantID uuid.UUID, bedID uint, roomID uuid.UUID, code, status string) (*domain.Bed, error) {
	if strings.TrimSpace(code) == "" {
		return nil, util.ErrValidation("kode bed wajib diisi")
	}

	if _, err := u.roomRepo.FindByIDAndTenant(roomID, tenantID); err != nil {
		return nil, util.ErrNotFound("ruangan tidak ditemukan atau akses ditolak")
	}

	bed, err := u.bedRepo.FindByIDAndTenant(bedID, tenantID)
	if err != nil {
		return nil, err
	}

	exists, err := u.bedRepo.ExistsByCode(tenantID, roomID, code, &bed.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode bed sudah digunakan")
	}

	bed.RoomID = roomID
	bed.Code = code
	if strings.TrimSpace(status) != "" {
		status = strings.ToLower(status)
		if !domain.IsValidBedStatus(status) {
			return nil, util.ErrValidation("status bed tidak valid")
		}
		bed.Status = status
	}

	if err := u.bedRepo.Update(bed); err != nil {
		return nil, err
	}
	return bed, nil
}

func (u *bedUsecase) DeleteBed(tenantID uuid.UUID, bedID uint) error {
	bed, err := u.bedRepo.FindByIDAndTenant(bedID, tenantID)
	if err != nil {
		return err
	}
	return u.bedRepo.Delete(bed)
}
