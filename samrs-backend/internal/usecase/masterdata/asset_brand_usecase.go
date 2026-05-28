package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type AssetBrandUsecase interface {
	CreateBrand(tenantID uuid.UUID, input AssetBrandInput) (*domain.AssetBrand, error)
	GetAllBrands(tenantID uuid.UUID, filter repository.AssetBrandFilter) ([]domain.AssetBrand, int64, error)
	GetBrandByID(tenantID uuid.UUID, id uint) (*domain.AssetBrand, error)
	UpdateBrand(tenantID uuid.UUID, id uint, input AssetBrandInput) (*domain.AssetBrand, error)
	DeleteBrand(tenantID uuid.UUID, id uint) error
}

type AssetBrandInput struct {
	Code        string
	Name        string
	Description string
}

type assetBrandUsecase struct {
	repo repository.AssetBrandRepository
}

func NewAssetBrandUsecase(repo repository.AssetBrandRepository) AssetBrandUsecase {
	return &assetBrandUsecase{repo}
}

func (u *assetBrandUsecase) CreateBrand(tenantID uuid.UUID, input AssetBrandInput) (*domain.AssetBrand, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama merek wajib diisi")
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
		return nil, util.ErrConflict("kode merek sudah digunakan")
	}

	brand := &domain.AssetBrand{
		TenantID:    tenantID,
		Code:        code,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	}

	if err := u.repo.Create(brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (u *assetBrandUsecase) GetAllBrands(tenantID uuid.UUID, filter repository.AssetBrandFilter) ([]domain.AssetBrand, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *assetBrandUsecase) GetBrandByID(tenantID uuid.UUID, id uint) (*domain.AssetBrand, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *assetBrandUsecase) UpdateBrand(tenantID uuid.UUID, id uint, input AssetBrandInput) (*domain.AssetBrand, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama merek wajib diisi")
	}

	brand, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, &brand.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode merek sudah digunakan")
	}

	brand.Code = code
	brand.Name = name
	brand.Description = strings.TrimSpace(input.Description)

	if err := u.repo.Update(brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (u *assetBrandUsecase) DeleteBrand(tenantID uuid.UUID, id uint) error {
	brand, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(brand)
}
