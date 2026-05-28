package usecase

import (
	"strings"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type PublicAssetUsecase interface {
	GetPublicAsset(tenantSlug, assetCode string) (*PublicAssetResponse, error)
}

type PublicAssetResponse struct {
	ID           uuid.UUID              `json:"id"`
	Code         string                 `json:"code"`
	Name         string                 `json:"name"`
	Brand        string                 `json:"brand,omitempty"`
	Model        string                 `json:"model,omitempty"`
	Status       string                 `json:"status"`
	PurchaseDate *time.Time             `json:"purchase_date,omitempty"`
	Category     *PublicAssetCategory   `json:"category,omitempty"`
	Room         *PublicAssetRoom       `json:"room,omitempty"`
	Bed          *PublicAssetBed        `json:"bed,omitempty"`
	Tenant       PublicAssetTenantBrief `json:"tenant"`
}

type PublicAssetCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PublicAssetRoom struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	Location string    `json:"location,omitempty"`
}

type PublicAssetBed struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

type PublicAssetTenantBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type publicAssetUsecase struct {
	tenantRepo repository.TenantRepository
	assetRepo  repository.AssetRepository
}

func NewPublicAssetUsecase(tr repository.TenantRepository, ar repository.AssetRepository) PublicAssetUsecase {
	return &publicAssetUsecase{
		tenantRepo: tr,
		assetRepo:  ar,
	}
}

func (u *publicAssetUsecase) GetPublicAsset(tenantSlug, assetCode string) (*PublicAssetResponse, error) {
	slug := strings.TrimSpace(tenantSlug)
	code := strings.TrimSpace(assetCode)
	if slug == "" || code == "" {
		return nil, util.ErrValidation("tenant slug dan kode aset wajib diisi")
	}

	tenant, err := u.tenantRepo.FindBySlug(slug)
	if err != nil {
		return nil, util.ErrNotFound("tenant tidak ditemukan")
	}

	asset, err := u.assetRepo.FindByCodeAndTenant(code, tenant.ID)
	if err != nil {
		return nil, util.ErrNotFound("aset tidak ditemukan")
	}

	return mapPublicAsset(asset, tenant), nil
}

func mapPublicAsset(asset *domain.Asset, tenant *domain.Tenant) *PublicAssetResponse {
	response := &PublicAssetResponse{
		ID:           asset.ID,
		Code:         asset.Code,
		Name:         asset.Name,
		Brand:        asset.Brand,
		Model:        asset.Model,
		Status:       asset.Status,
		PurchaseDate: asset.PurchaseDate,
		Tenant: PublicAssetTenantBrief{
			ID:   tenant.ID,
			Name: tenant.Name,
			Slug: tenant.Slug,
		},
	}

	if asset.Category.ID != 0 {
		response.Category = &PublicAssetCategory{
			ID:   asset.Category.ID,
			Name: asset.Category.Name,
			Slug: asset.Category.Slug,
		}
	}

	if asset.Room.ID != uuid.Nil {
		response.Room = &PublicAssetRoom{
			ID:       asset.Room.ID,
			Name:     asset.Room.Name,
			Code:     asset.Room.Code,
			Location: asset.Room.Location,
		}
	}

	if asset.Bed.ID != 0 {
		response.Bed = &PublicAssetBed{
			ID:   asset.Bed.ID,
			Code: asset.Bed.Code,
		}
	}

	return response
}
