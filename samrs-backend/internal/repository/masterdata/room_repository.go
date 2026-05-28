package masterdatarepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *domain.Room) error
	FindAllByTenant(tenantID uuid.UUID, filter RoomFilter) ([]domain.Room, int64, error)
	FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Room, error)
	Update(room *domain.Room) error
	Delete(room *domain.Room) error
	ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type RoomFilter struct {
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db}
}

func (r *roomRepository) Create(room *domain.Room) error {
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	return db.Create(room).Error
}

func (r *roomRepository) FindAllByTenant(tenantID uuid.UUID, filter RoomFilter) ([]domain.Room, int64, error) {
	var rooms []domain.Room
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ? OR location ILIKE ?", like, like, like)
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

	sortBy := mapRoomSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Room, error) {
	var room domain.Room
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(room *domain.Room) error {
	updates := map[string]interface{}{
		"name":       room.Name,
		"code":       room.Code,
		"location":   room.Location,
		"updated_at": time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	return repobase.WithTenant(db, room.TenantID).
		Set("audit_record_id", room.ID).
		Where("id = ?", room.ID).
		Updates(updates).
		Error
}

func (r *roomRepository) Delete(room *domain.Room) error {
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	return repobase.WithTenant(db, room.TenantID).Delete(room).Error
}

func (r *roomRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	query := repobase.WithTenant(db, tenantID).Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roomRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapRoomSortBy(sortBy string) string {
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
