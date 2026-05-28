package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetModelRepository interface {
	Create(model *domain.AssetModel) error
	FindAllByTenant(tenantID uuid.UUID, filter AssetModelFilter) ([]domain.AssetModel, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.AssetModel, error)
	Update(model *domain.AssetModel) error
	Delete(model *domain.AssetModel) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type assetModelRepository struct {
	db *gorm.DB
}

type AssetModelFilter struct {
	Search   string
	BrandID  uint
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewAssetModelRepository(db *gorm.DB) AssetModelRepository {
	return &assetModelRepository{db}
}

func (r *assetModelRepository) Create(model *domain.AssetModel) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	return db.Create(model).Error
}

func (r *assetModelRepository) FindAllByTenant(tenantID uuid.UUID, filter AssetModelFilter) ([]domain.AssetModel, int64, error) {
	var models []domain.AssetModel
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	query := repobase.WithTenant(db, tenantID)

	if filter.BrandID != 0 {
		query = query.Where("brand_id = ?", filter.BrandID)
	}
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

	sortBy := mapAssetModelSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Brand").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return models, total, nil
}

func (r *assetModelRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetModel, error) {
	var model domain.AssetModel
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	err := repobase.WithTenant(db, tenantID).Preload("Brand").Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *assetModelRepository) Update(model *domain.AssetModel) error {
	updates := map[string]interface{}{
		"brand_id":    model.BrandID,
		"code":        model.Code,
		"name":        model.Name,
		"description": model.Description,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	return repobase.WithTenant(db, model.TenantID).
		Set("audit_record_id", model.ID).
		Where("id = ?", model.ID).
		Updates(updates).
		Error
}

func (r *assetModelRepository) Delete(model *domain.AssetModel) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	return repobase.WithTenant(db, model.TenantID).Delete(model).Error
}

func (r *assetModelRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetModelRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetModel{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapAssetModelSortBy(sortBy string) string {
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
