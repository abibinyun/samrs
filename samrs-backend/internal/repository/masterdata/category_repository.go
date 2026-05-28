package masterdatarepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *domain.Category) error
	FindAllByTenant(tenantID uuid.UUID, filter CategoryFilter) ([]domain.Category, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(category *domain.Category) error
	ExistsBySlug(tenantID uuid.UUID, slug string, excludeID *uint) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

type CategoryFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	return db.Create(category).Error
}

func (r *categoryRepository) FindAllByTenant(tenantID uuid.UUID, filter CategoryFilter) ([]domain.Category, int64, error) {
	var categories []domain.Category
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	query := repobase.WithTenant(db, tenantID)

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

	sortBy := mapCategorySortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *categoryRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Category, error) {
	var category domain.Category
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	updates := map[string]interface{}{
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	return repobase.WithTenant(db, category.TenantID).
		Set("audit_record_id", category.ID).
		Where("id = ?", category.ID).
		Updates(updates).
		Error
}

func (r *categoryRepository) Delete(category *domain.Category) error {
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	return repobase.WithTenant(db, category.TenantID).Delete(category).Error
}

func (r *categoryRepository) ExistsBySlug(tenantID uuid.UUID, slug string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	query := repobase.WithTenant(db, tenantID).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *categoryRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Category{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapCategorySortBy(sortBy string) string {
	switch sortBy {
	case "name":
		return "name"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
