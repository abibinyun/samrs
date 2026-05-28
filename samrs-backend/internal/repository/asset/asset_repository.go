package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetFilter struct {
	Status       string
	CategoryID   uint
	RoomID       uuid.UUID
	BedID        uint
	Search       string
	DateFrom     *time.Time
	DateTo       *time.Time
	PurchaseFrom *time.Time
	PurchaseTo   *time.Time
	Page         int
	PerPage      int
	SortBy       string
	SortDir      string
}

type AssetRepository interface {
	Create(asset *domain.Asset) error
	FindAllByTenant(tenantID uuid.UUID, filter AssetFilter) ([]domain.Asset, int64, error)
	FindAllByTenantExport(tenantID uuid.UUID, filter AssetFilter) ([]domain.Asset, error)
	FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Asset, error)
	FindByCodeAndTenant(code string, tenantID uuid.UUID) (*domain.Asset, error)
	Update(asset *domain.Asset) error
	Delete(asset *domain.Asset) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db}
}

func (r *assetRepository) Create(asset *domain.Asset) error {
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	return db.Create(asset).Error
}

func (r *assetRepository) FindAllByTenant(tenantID uuid.UUID, filter AssetFilter) ([]domain.Asset, int64, error) {
	var assets []domain.Asset
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	query := repobase.WithTenant(db, tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.RoomID != uuid.Nil {
		query = query.Where("room_id = ?", filter.RoomID)
	}

	if filter.BedID != 0 {
		query = query.Where("bed_id = ?", filter.BedID)
	}

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR brand ILIKE ? OR model ILIKE ?", like, like, like, like)
	}

	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.PurchaseFrom != nil {
		query = query.Where("purchase_date >= ?", *filter.PurchaseFrom)
	}
	if filter.PurchaseTo != nil {
		query = query.Where("purchase_date <= ?", *filter.PurchaseTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := mapAssetSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	err := query.Preload("Category").Preload("Room").Preload("Bed").
		Preload("Vendor").Preload("BrandRef").Preload("ModelRef").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&assets).Error
	if err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *assetRepository) FindAllByTenantExport(tenantID uuid.UUID, filter AssetFilter) ([]domain.Asset, error) {
	var assets []domain.Asset
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	query := repobase.WithTenant(db, tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.RoomID != uuid.Nil {
		query = query.Where("room_id = ?", filter.RoomID)
	}
	if filter.BedID != 0 {
		query = query.Where("bed_id = ?", filter.BedID)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR brand ILIKE ? OR model ILIKE ?", like, like, like, like)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.PurchaseFrom != nil {
		query = query.Where("purchase_date >= ?", *filter.PurchaseFrom)
	}
	if filter.PurchaseTo != nil {
		query = query.Where("purchase_date <= ?", *filter.PurchaseTo)
	}

	sortBy := mapAssetSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	if err := query.Preload("Category").Preload("Room").Preload("Bed").
		Preload("Vendor").Preload("BrandRef").Preload("ModelRef").
		Order(sortBy + " " + sortDir).
		Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	err := repobase.WithTenant(db, tenantID).Preload("Category").Preload("Room").Preload("Bed").
		Preload("Vendor").Preload("BrandRef").Preload("ModelRef").
		Where("id = ?", id).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) FindByCodeAndTenant(code string, tenantID uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	err := repobase.WithTenant(db, tenantID).Preload("Category").Preload("Room").Preload("Bed").
		Preload("Vendor").Preload("BrandRef").Preload("ModelRef").
		Where("code = ?", code).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) Update(asset *domain.Asset) error {
	updates := map[string]interface{}{
		"category_id":   asset.CategoryID,
		"room_id":       asset.RoomID,
		"bed_id":        asset.BedID,
		"vendor_id":     asset.VendorID,
		"brand_id":      asset.BrandID,
		"model_id":      asset.ModelID,
		"code":          asset.Code,
		"name":          asset.Name,
		"brand":         asset.Brand,
		"model":         asset.Model,
		"status":        asset.Status,
		"purchase_date": asset.PurchaseDate,
		"updated_at":    time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	return repobase.WithTenant(db, asset.TenantID).
		Set("audit_record_id", asset.ID).
		Where("id = ?", asset.ID).
		Updates(updates).
		Error
}

func (r *assetRepository) Delete(asset *domain.Asset) error {
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	return repobase.WithTenant(db, asset.TenantID).Delete(asset).Error
}

func (r *assetRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Asset{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapAssetSortBy(sortBy string) string {
	switch sortBy {
	case "code":
		return "code"
	case "status":
		return "status"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
