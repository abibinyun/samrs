package tenantrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

type TenantRepository interface {
	ListAll(filter TenantFilter) ([]domain.Tenant, int64, error)
	FindByID(id uuid.UUID) (*domain.Tenant, error)
	FindBySlug(slug string) (*domain.Tenant, error)
	Create(tenant *domain.Tenant) error
	Update(tenant *domain.Tenant) error
	ExistsBySlug(slug string, excludeID *uuid.UUID) (bool, error)
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db}
}

func (r *tenantRepository) ListAll(filter TenantFilter) ([]domain.Tenant, int64, error) {
	var tenants []domain.Tenant
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	query := db

	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", like, like)
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

	sortBy := mapTenantSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

func (r *tenantRepository) FindByID(id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	if err := db.Where("id = ?", id).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) FindBySlug(slug string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	if err := db.Where("slug = ?", slug).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) Create(tenant *domain.Tenant) error {
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	return db.Create(tenant).Error
}

func (r *tenantRepository) Update(tenant *domain.Tenant) error {
	updates := map[string]interface{}{
		"name":        tenant.Name,
		"slug":        tenant.Slug,
		"address":     tenant.Address,
		"tenant_type": tenant.TenantType,
		"status":      tenant.Status,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	return db.
		Set("audit_record_id", tenant.ID).
		Where("id = ?", tenant.ID).
		Updates(updates).
		Error
}

func (r *tenantRepository) ExistsBySlug(slug string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Tenant{})
	query := db.Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func mapTenantSortBy(sortBy string) string {
	switch sortBy {
	case "name":
		return "name"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
