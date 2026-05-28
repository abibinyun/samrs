package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type AssetStatusUsecase interface {
	CreateStatus(tenantID uuid.UUID, input AssetStatusInput) (*domain.AssetStatus, error)
	GetAllStatuses(tenantID uuid.UUID, filter repository.AssetStatusFilter) ([]domain.AssetStatus, int64, error)
	GetStatusByID(tenantID uuid.UUID, id uint) (*domain.AssetStatus, error)
	UpdateStatus(tenantID uuid.UUID, id uint, input AssetStatusInput) (*domain.AssetStatus, error)
	DeleteStatus(tenantID uuid.UUID, id uint) error
}

type AssetStatusInput struct {
	Code        string
	Name        string
	Description string
}

type assetStatusUsecase struct {
	repo repository.AssetStatusRepository
}

func NewAssetStatusUsecase(repo repository.AssetStatusRepository) AssetStatusUsecase {
	return &assetStatusUsecase{repo}
}

func (u *assetStatusUsecase) CreateStatus(tenantID uuid.UUID, input AssetStatusInput) (*domain.AssetStatus, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama status wajib diisi")
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode status sudah digunakan")
	}

	status := &domain.AssetStatus{
		TenantID:    tenantID,
		Code:        code,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	}

	if err := u.repo.Create(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *assetStatusUsecase) GetAllStatuses(tenantID uuid.UUID, filter repository.AssetStatusFilter) ([]domain.AssetStatus, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *assetStatusUsecase) GetStatusByID(tenantID uuid.UUID, id uint) (*domain.AssetStatus, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *assetStatusUsecase) UpdateStatus(tenantID uuid.UUID, id uint, input AssetStatusInput) (*domain.AssetStatus, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama status wajib diisi")
	}

	status, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, &status.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode status sudah digunakan")
	}

	status.Code = code
	status.Name = name
	status.Description = strings.TrimSpace(input.Description)

	if err := u.repo.Update(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (u *assetStatusUsecase) DeleteStatus(tenantID uuid.UUID, id uint) error {
	status, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(status)
}
