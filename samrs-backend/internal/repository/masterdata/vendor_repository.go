package masterdatarepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VendorRepository interface {
	Create(vendor *domain.Vendor) error
	FindAllByTenant(tenantID uuid.UUID, filter VendorFilter) ([]domain.Vendor, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.Vendor, error)
	Update(vendor *domain.Vendor) error
	Delete(vendor *domain.Vendor) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type vendorRepository struct {
	db *gorm.DB
}

type VendorFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return &vendorRepository{db}
}

func (r *vendorRepository) Create(vendor *domain.Vendor) error {
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	return db.Create(vendor).Error
}

func (r *vendorRepository) FindAllByTenant(tenantID uuid.UUID, filter VendorFilter) ([]domain.Vendor, int64, error) {
	var vendors []domain.Vendor
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
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

	sortBy := mapVendorSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&vendors).Error; err != nil {
		return nil, 0, err
	}
	return vendors, total, nil
}

func (r *vendorRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Vendor, error) {
	var vendor domain.Vendor
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (r *vendorRepository) Update(vendor *domain.Vendor) error {
	updates := map[string]interface{}{
		"code":        vendor.Code,
		"name":        vendor.Name,
		"contact_name": vendor.ContactName,
		"phone":       vendor.Phone,
		"email":       vendor.Email,
		"address":     vendor.Address,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	return repobase.WithTenant(db, vendor.TenantID).
		Set("audit_record_id", vendor.ID).
		Where("id = ?", vendor.ID).
		Updates(updates).
		Error
}

func (r *vendorRepository) Delete(vendor *domain.Vendor) error {
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	return repobase.WithTenant(db, vendor.TenantID).Delete(vendor).Error
}

func (r *vendorRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *vendorRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Vendor{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapVendorSortBy(sortBy string) string {
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
