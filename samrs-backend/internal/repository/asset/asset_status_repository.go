package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetStatusRepository interface {
	Create(status *domain.AssetStatus) error
	FindAllByTenant(tenantID uuid.UUID, filter AssetStatusFilter) ([]domain.AssetStatus, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.AssetStatus, error)
	Update(status *domain.AssetStatus) error
	Delete(status *domain.AssetStatus) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error)
	HasAny(tenantID uuid.UUID) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type assetStatusRepository struct {
	db *gorm.DB
}

type AssetStatusFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewAssetStatusRepository(db *gorm.DB) AssetStatusRepository {
	return &assetStatusRepository{db}
}

func (r *assetStatusRepository) Create(status *domain.AssetStatus) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	return db.Create(status).Error
}

func (r *assetStatusRepository) FindAllByTenant(tenantID uuid.UUID, filter AssetStatusFilter) ([]domain.AssetStatus, int64, error) {
	var statuses []domain.AssetStatus
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
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

	sortBy := mapAssetStatusSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&statuses).Error; err != nil {
		return nil, 0, err
	}
	return statuses, total, nil
}

func (r *assetStatusRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetStatus, error) {
	var status domain.AssetStatus
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *assetStatusRepository) Update(status *domain.AssetStatus) error {
	updates := map[string]interface{}{
		"code":       status.Code,
		"name":       status.Name,
		"updated_at": time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	return repobase.WithTenant(db, status.TenantID).
		Set("audit_record_id", status.ID).
		Where("id = ?", status.ID).
		Updates(updates).
		Error
}

func (r *assetStatusRepository) Delete(status *domain.AssetStatus) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	return repobase.WithTenant(db, status.TenantID).Delete(status).Error
}

func (r *assetStatusRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetStatusRepository) HasAny(tenantID uuid.UUID) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	if err := repobase.WithTenant(db, tenantID).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *assetStatusRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AssetStatus{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapAssetStatusSortBy(sortBy string) string {
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
