package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type AssetMutationUsecase interface {
	MoveAsset(input AssetMutationInput) (*domain.AssetMutation, *domain.Asset, error)
	GetAllMutations(tenantID uuid.UUID, filter repository.AssetMutationFilter) ([]domain.AssetMutation, int64, error)
	GetMutationByID(tenantID uuid.UUID, id uint) (*domain.AssetMutation, error)
}

type AssetMutationInput struct {
	TenantID uuid.UUID
	AssetID  uuid.UUID
	ToRoomID *uuid.UUID
	ToBedID  *uint
	Reason   string
	MovedBy  uuid.UUID
}

type assetMutationUsecase struct {
	repo      repository.AssetMutationRepository
	assetRepo repository.AssetRepository
	roomRepo  repository.RoomRepository
	bedRepo   repository.BedRepository
}

func NewAssetMutationUsecase(
	repo repository.AssetMutationRepository,
	ar repository.AssetRepository,
	rr repository.RoomRepository,
	br repository.BedRepository,
) AssetMutationUsecase {
	return &assetMutationUsecase{
		repo:      repo,
		assetRepo: ar,
		roomRepo:  rr,
		bedRepo:   br,
	}
}

func (u *assetMutationUsecase) MoveAsset(input AssetMutationInput) (*domain.AssetMutation, *domain.Asset, error) {
	if input.AssetID == uuid.Nil {
		return nil, nil, util.ErrValidation("asset ID wajib diisi")
	}
	if input.ToRoomID == nil && input.ToBedID == nil {
		return nil, nil, util.ErrValidation("tujuan ruangan atau bed wajib diisi")
	}

	asset, err := u.assetRepo.FindByIDAndTenant(input.AssetID, input.TenantID)
	if err != nil {
		return nil, nil, err
	}

	if input.ToRoomID != nil {
		if _, err := u.roomRepo.FindByIDAndTenant(*input.ToRoomID, input.TenantID); err != nil {
			return nil, nil, util.ErrNotFound("ruangan tujuan tidak ditemukan")
		}
	}

	if input.ToBedID != nil {
		bed, err := u.bedRepo.FindByIDAndTenant(*input.ToBedID, input.TenantID)
		if err != nil {
			return nil, nil, util.ErrNotFound("bed tujuan tidak ditemukan")
		}
		if input.ToRoomID != nil && bed.RoomID != *input.ToRoomID {
			return nil, nil, util.ErrValidation("bed tidak sesuai dengan ruangan tujuan")
		}
		if input.ToRoomID == nil {
			input.ToRoomID = &bed.RoomID
		}
	}

	mutation := &domain.AssetMutation{
		TenantID:   input.TenantID,
		AssetID:    input.AssetID,
		FromRoomID: asset.RoomID,
		FromBedID:  asset.BedID,
		ToRoomID:   input.ToRoomID,
		ToBedID:    input.ToBedID,
		Reason:     strings.TrimSpace(input.Reason),
		MovedBy:    input.MovedBy,
	}

	asset.RoomID = input.ToRoomID
	asset.BedID = input.ToBedID

	if err := u.repo.Create(mutation); err != nil {
		return nil, nil, err
	}
	if err := u.assetRepo.Update(asset); err != nil {
		return nil, nil, err
	}

	return mutation, asset, nil
}

func (u *assetMutationUsecase) GetAllMutations(tenantID uuid.UUID, filter repository.AssetMutationFilter) ([]domain.AssetMutation, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *assetMutationUsecase) GetMutationByID(tenantID uuid.UUID, id uint) (*domain.AssetMutation, error) {
	return u.repo.FindByID(tenantID, id)
}
