package usecase

import (
	"strings"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type AssetUsecase interface {
	CreateAsset(input CreateAssetInput) (*domain.Asset, error)
	GetAllAssets(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, int64, error)
	GetAssetByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Asset, error)
	UpdateAsset(input UpdateAssetInput) (*domain.Asset, error)
	DeleteAsset(tenantID uuid.UUID, id uuid.UUID) error
}

type CreateAssetInput struct {
	TenantID     uuid.UUID
	CategoryID   uint
	RoomID       *uuid.UUID
	BedID        *uint
	VendorID     *uint
	BrandID      *uint
	ModelID      *uint
	Code         string
	Name         string
	Brand        string
	Model        string
	Status       string
	PurchaseDate string
}

type UpdateAssetInput struct {
	TenantID     uuid.UUID
	ID           uuid.UUID
	CategoryID   uint
	RoomID       *uuid.UUID
	BedID        *uint
	VendorID     *uint
	BrandID      *uint
	ModelID      *uint
	Code         string
	Name         string
	Brand        string
	Model        string
	Status       string
	PurchaseDate string
}

type assetUsecase struct {
	assetRepo    repository.AssetRepository
	categoryRepo repository.CategoryRepository
	roomRepo     repository.RoomRepository
	bedRepo      repository.BedRepository
	vendorRepo   repository.VendorRepository
	brandRepo    repository.AssetBrandRepository
	modelRepo    repository.AssetModelRepository
	statusRepo   repository.AssetStatusRepository
}

func NewAssetUsecase(
	ar repository.AssetRepository,
	cr repository.CategoryRepository,
	rr repository.RoomRepository,
	br repository.BedRepository,
	vr repository.VendorRepository,
	abr repository.AssetBrandRepository,
	amr repository.AssetModelRepository,
	asr repository.AssetStatusRepository,
) AssetUsecase {
	return &assetUsecase{
		assetRepo:    ar,
		categoryRepo: cr,
		roomRepo:     rr,
		bedRepo:      br,
		vendorRepo:   vr,
		brandRepo:    abr,
		modelRepo:    amr,
		statusRepo:   asr,
	}
}

func (u *assetUsecase) CreateAsset(input CreateAssetInput) (*domain.Asset, error) {
	if strings.TrimSpace(input.Code) == "" || strings.TrimSpace(input.Name) == "" {
		return nil, util.ErrValidation("kode dan nama aset wajib diisi")
	}

	if input.CategoryID == 0 {
		return nil, util.ErrValidation("kategori wajib diisi")
	}

	if _, err := u.categoryRepo.FindByID(input.TenantID, input.CategoryID); err != nil {
		return nil, util.ErrNotFound("kategori tidak ditemukan atau akses ditolak")
	}

	if input.RoomID != nil {
		if _, err := u.roomRepo.FindByIDAndTenant(*input.RoomID, input.TenantID); err != nil {
			return nil, util.ErrNotFound("ruangan tidak ditemukan atau akses ditolak")
		}
	}

	if input.BedID != nil {
		bed, err := u.bedRepo.FindByIDAndTenant(*input.BedID, input.TenantID)
		if err != nil {
			return nil, util.ErrNotFound("bed tidak ditemukan atau akses ditolak")
		}
		if input.RoomID != nil && bed.RoomID != *input.RoomID {
			return nil, util.ErrValidation("bed tidak sesuai dengan ruangan")
		}
		if input.RoomID == nil {
			input.RoomID = &bed.RoomID
		}
	}

	if input.VendorID != nil {
		if _, err := u.vendorRepo.FindByID(input.TenantID, *input.VendorID); err != nil {
			return nil, util.ErrNotFound("vendor tidak ditemukan atau akses ditolak")
		}
	}
	if input.BrandID != nil {
		if _, err := u.brandRepo.FindByID(input.TenantID, *input.BrandID); err != nil {
			return nil, util.ErrNotFound("merek tidak ditemukan atau akses ditolak")
		}
	}
	if input.ModelID != nil {
		model, err := u.modelRepo.FindByID(input.TenantID, *input.ModelID)
		if err != nil {
			return nil, util.ErrNotFound("tipe alat tidak ditemukan atau akses ditolak")
		}
		if input.BrandID != nil && model.BrandID != *input.BrandID {
			return nil, util.ErrValidation("tipe alat tidak sesuai dengan merek")
		}
		if input.BrandID == nil {
			input.BrandID = &model.BrandID
		}
	}

	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = domain.AssetStatusReady
	} else {
		status = strings.ToLower(status)
	}
	if err := u.validateStatus(input.TenantID, status); err != nil {
		return nil, err
	}

	var purchaseDate *time.Time
	if strings.TrimSpace(input.PurchaseDate) != "" {
		parsed, err := time.Parse("2006-01-02", input.PurchaseDate)
		if err != nil {
			return nil, util.ErrValidation("format purchase_date harus YYYY-MM-DD")
		}
		purchaseDate = &parsed
	}

	exists, err := u.assetRepo.ExistsByCode(input.TenantID, input.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode aset sudah digunakan")
	}

	asset := &domain.Asset{
		TenantID:     input.TenantID,
		CategoryID:   input.CategoryID,
		RoomID:       input.RoomID,
		BedID:        input.BedID,
		VendorID:     input.VendorID,
		BrandID:      input.BrandID,
		ModelID:      input.ModelID,
		Code:         input.Code,
		Name:         input.Name,
		Brand:        input.Brand,
		Model:        input.Model,
		Status:       status,
		PurchaseDate: purchaseDate,
	}

	if err := u.assetRepo.Create(asset); err != nil {
		return nil, err
	}

	// Fetch with relations populated
	return u.assetRepo.FindByIDAndTenant(asset.ID, input.TenantID)
}

func (u *assetUsecase) GetAllAssets(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, int64, error) {
	return u.assetRepo.FindAllByTenant(tenantID, filter)
}

func (u *assetUsecase) GetAssetByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Asset, error) {
	return u.assetRepo.FindByIDAndTenant(id, tenantID)
}

func (u *assetUsecase) UpdateAsset(input UpdateAssetInput) (*domain.Asset, error) {
	if strings.TrimSpace(input.Code) == "" || strings.TrimSpace(input.Name) == "" {
		return nil, util.ErrValidation("kode dan nama aset wajib diisi")
	}

	asset, err := u.assetRepo.FindByIDAndTenant(input.ID, input.TenantID)
	if err != nil {
		return nil, err
	}

	if input.CategoryID == 0 {
		return nil, util.ErrValidation("kategori wajib diisi")
	}
	if _, err := u.categoryRepo.FindByID(input.TenantID, input.CategoryID); err != nil {
		return nil, util.ErrNotFound("kategori tidak ditemukan atau akses ditolak")
	}

	if input.RoomID != nil {
		if _, err := u.roomRepo.FindByIDAndTenant(*input.RoomID, input.TenantID); err != nil {
			return nil, util.ErrNotFound("ruangan tidak ditemukan atau akses ditolak")
		}
	}

	if input.BedID != nil {
		bed, err := u.bedRepo.FindByIDAndTenant(*input.BedID, input.TenantID)
		if err != nil {
			return nil, util.ErrNotFound("bed tidak ditemukan atau akses ditolak")
		}
		if input.RoomID != nil && bed.RoomID != *input.RoomID {
			return nil, util.ErrValidation("bed tidak sesuai dengan ruangan")
		}
		if input.RoomID == nil {
			input.RoomID = &bed.RoomID
		}
	}

	if input.VendorID != nil {
		if _, err := u.vendorRepo.FindByID(input.TenantID, *input.VendorID); err != nil {
			return nil, util.ErrNotFound("vendor tidak ditemukan atau akses ditolak")
		}
	}
	if input.BrandID != nil {
		if _, err := u.brandRepo.FindByID(input.TenantID, *input.BrandID); err != nil {
			return nil, util.ErrNotFound("merek tidak ditemukan atau akses ditolak")
		}
	}
	if input.ModelID != nil {
		model, err := u.modelRepo.FindByID(input.TenantID, *input.ModelID)
		if err != nil {
			return nil, util.ErrNotFound("tipe alat tidak ditemukan atau akses ditolak")
		}
		if input.BrandID != nil && model.BrandID != *input.BrandID {
			return nil, util.ErrValidation("tipe alat tidak sesuai dengan merek")
		}
		if input.BrandID == nil {
			input.BrandID = &model.BrandID
		}
	}

	exists, err := u.assetRepo.ExistsByCode(input.TenantID, input.Code, &asset.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode aset sudah digunakan")
	}

	var purchaseDate *time.Time
	if strings.TrimSpace(input.PurchaseDate) != "" {
		parsed, err := time.Parse("2006-01-02", input.PurchaseDate)
		if err != nil {
			return nil, util.ErrValidation("format purchase_date harus YYYY-MM-DD")
		}
		purchaseDate = &parsed
	}

	asset.CategoryID = input.CategoryID
	asset.RoomID = input.RoomID
	asset.BedID = input.BedID
	asset.VendorID = input.VendorID
	asset.BrandID = input.BrandID
	asset.ModelID = input.ModelID
	asset.Code = input.Code
	asset.Name = input.Name
	asset.Brand = input.Brand
	asset.Model = input.Model
	if strings.TrimSpace(input.Status) != "" {
		status := strings.ToLower(input.Status)
		if err := u.validateStatus(input.TenantID, status); err != nil {
			return nil, err
		}
		asset.Status = status
	}
	asset.PurchaseDate = purchaseDate

	if err := u.assetRepo.Update(asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (u *assetUsecase) DeleteAsset(tenantID uuid.UUID, id uuid.UUID) error {
	asset, err := u.assetRepo.FindByIDAndTenant(id, tenantID)
	if err != nil {
		return err
	}
	return u.assetRepo.Delete(asset)
}

func (u *assetUsecase) validateStatus(tenantID uuid.UUID, status string) error {
	hasAny, err := u.statusRepo.HasAny(tenantID)
	if err != nil {
		return err
	}
	if hasAny {
		exists, err := u.statusRepo.ExistsByCode(tenantID, status, nil)
		if err != nil {
			return err
		}
		if !exists {
			return util.ErrValidation("status aset tidak terdaftar")
		}
		return nil
	}
	if !domain.IsValidAssetStatus(status) {
		return util.ErrValidation("status aset tidak valid")
	}
	return nil
}
