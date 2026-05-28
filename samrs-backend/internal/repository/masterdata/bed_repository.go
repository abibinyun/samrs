package masterdatarepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BedRepository interface {
	Create(bed *domain.Bed) error
	FindAllByTenant(tenantID uuid.UUID, filter BedFilter) ([]domain.Bed, int64, error)
	FindByIDAndTenant(id uint, tenantID uuid.UUID) (*domain.Bed, error)
	Update(bed *domain.Bed) error
	Delete(bed *domain.Bed) error
	ExistsByCode(tenantID uuid.UUID, roomID uuid.UUID, code string, excludeID *uint) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type bedRepository struct {
	db *gorm.DB
}

type BedFilter struct {
	Search   string
	RoomID   *uuid.UUID
	Status   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewBedRepository(db *gorm.DB) BedRepository {
	return &bedRepository{db}
}

func (r *bedRepository) Create(bed *domain.Bed) error {
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	return db.Create(bed).Error
}

func (r *bedRepository) FindAllByTenant(tenantID uuid.UUID, filter BedFilter) ([]domain.Bed, int64, error) {
	var beds []domain.Bed
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	query := repobase.WithTenant(db, tenantID)

	if filter.RoomID != nil {
		query = query.Where("room_id = ?", *filter.RoomID)
	}
	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("code ILIKE ?", like)
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

	sortBy := mapBedSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&beds).Error; err != nil {
		return nil, 0, err
	}
	return beds, total, nil
}

func (r *bedRepository) FindByIDAndTenant(id uint, tenantID uuid.UUID) (*domain.Bed, error) {
	var bed domain.Bed
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&bed).Error
	if err != nil {
		return nil, err
	}
	return &bed, nil
}

func (r *bedRepository) Update(bed *domain.Bed) error {
	updates := map[string]interface{}{
		"room_id":    bed.RoomID,
		"code":       bed.Code,
		"status":     bed.Status,
		"updated_at": time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	return repobase.WithTenant(db, bed.TenantID).
		Set("audit_record_id", bed.ID).
		Where("id = ?", bed.ID).
		Updates(updates).
		Error
}

func (r *bedRepository) Delete(bed *domain.Bed) error {
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	return repobase.WithTenant(db, bed.TenantID).Delete(bed).Error
}

func (r *bedRepository) ExistsByCode(tenantID uuid.UUID, roomID uuid.UUID, code string, excludeID *uint) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	query := repobase.WithTenant(db, tenantID).Where("room_id = ? AND code = ?", roomID, code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bedRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Bed{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapBedSortBy(sortBy string) string {
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
