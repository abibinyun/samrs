package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetBrandRepository interface {
	Create(brand *domain.AssetBrand) error
	FindAllByTenant(tenantID uuid.UUID, filter AssetBrandFilter) ([]domain.AssetBrand, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.AssetBrand, error)
	Update(brand *domain.AssetBrand) error
	Delete(brand *domain.AssetBrand) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type assetBrandRepository struct {
	db *gorm.DB
}

type AssetBrandFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewAssetBrandRepository(db *gorm.DB) AssetBrandRepository {
	return &assetBrandRepository{db}
}

func (r *assetBrandRepository) Create(brand *domain.AssetBrand) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	return db.Create(brand).Error
}

func (r *assetBrandRepository) FindAllByTenant(tenantID uuid.UUID, filter AssetBrandFilter) ([]domain.AssetBrand, int64, error) {
	var brands []domain.AssetBrand
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := mapAssetBrandSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&brands).Error; err != nil {
		return nil, 0, err
	}
	return brands, total, nil
}

func (r *assetBrandRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetBrand, error) {
	var brand domain.AssetBrand
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *assetBrandRepository) Update(brand *domain.AssetBrand) error {
	updates := map[string]interface{}{
		"code":        brand.Code,
		"name":        brand.Name,
		"description": brand.Description,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	return repobase.WithTenant(db, brand.TenantID).
		Set("audit_record_id", brand.ID).
		Where("id = ?", brand.ID).
		Updates(updates).
		Error
}

func (r *assetBrandRepository) Delete(brand *domain.AssetBrand) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	return repobase.WithTenant(db, brand.TenantID).Delete(brand).Error
}

func (r *assetBrandRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetBrandRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetBrand{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapAssetBrandSortBy(sortBy string) string {
	switch sortBy {
	case "name":
		return "name"
	case "code":
		return "code"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
