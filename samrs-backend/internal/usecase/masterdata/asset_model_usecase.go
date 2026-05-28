package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type AssetModelUsecase interface {
	CreateModel(tenantID uuid.UUID, input AssetModelInput) (*domain.AssetModel, error)
	GetAllModels(tenantID uuid.UUID, filter repository.AssetModelFilter) ([]domain.AssetModel, int64, error)
	GetModelByID(tenantID uuid.UUID, id uint) (*domain.AssetModel, error)
	UpdateModel(tenantID uuid.UUID, id uint, input AssetModelInput) (*domain.AssetModel, error)
	DeleteModel(tenantID uuid.UUID, id uint) error
}

type AssetModelInput struct {
	BrandID     uint
	Code        string
	Name        string
	Description string
}

type assetModelUsecase struct {
	repo      repository.AssetModelRepository
	brandRepo repository.AssetBrandRepository
}

func NewAssetModelUsecase(repo repository.AssetModelRepository, br repository.AssetBrandRepository) AssetModelUsecase {
	return &assetModelUsecase{repo: repo, brandRepo: br}
}

func (u *assetModelUsecase) CreateModel(tenantID uuid.UUID, input AssetModelInput) (*domain.AssetModel, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama tipe alat wajib diisi")
	}
	if input.BrandID == 0 {
		return nil, util.ErrValidation("merek wajib diisi")
	}
	if _, err := u.brandRepo.FindByID(tenantID, input.BrandID); err != nil {
		return nil, util.ErrNotFound("merek tidak ditemukan atau akses ditolak")
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
		return nil, util.ErrConflict("kode tipe alat sudah digunakan")
	}

	model := &domain.AssetModel{
		TenantID:    tenantID,
		BrandID:     input.BrandID,
		Code:        code,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	}

	if err := u.repo.Create(model); err != nil {
		return nil, err
	}
	return model, nil
}

func (u *assetModelUsecase) GetAllModels(tenantID uuid.UUID, filter repository.AssetModelFilter) ([]domain.AssetModel, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *assetModelUsecase) GetModelByID(tenantID uuid.UUID, id uint) (*domain.AssetModel, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *assetModelUsecase) UpdateModel(tenantID uuid.UUID, id uint, input AssetModelInput) (*domain.AssetModel, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama tipe alat wajib diisi")
	}
	if input.BrandID == 0 {
		return nil, util.ErrValidation("merek wajib diisi")
	}
	if _, err := u.brandRepo.FindByID(tenantID, input.BrandID); err != nil {
		return nil, util.ErrNotFound("merek tidak ditemukan atau akses ditolak")
	}

	model, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, &model.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode tipe alat sudah digunakan")
	}

	model.BrandID = input.BrandID
	model.Code = code
	model.Name = name
	model.Description = strings.TrimSpace(input.Description)

	if err := u.repo.Update(model); err != nil {
		return nil, err
	}
	return model, nil
}

func (u *assetModelUsecase) DeleteModel(tenantID uuid.UUID, id uint) error {
	model, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(model)
}
