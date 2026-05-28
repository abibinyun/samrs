package maintenancerepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceScheduleFilter struct {
	ScheduleType string
	Status       string
	AssetID      uuid.UUID
	DateFrom     *time.Time
	DateTo       *time.Time
	DueFrom      *time.Time
	DueTo        *time.Time
	Page         int
	PerPage      int
	SortBy       string
	SortDir      string
}

type MaintenanceScheduleRepository interface {
	Create(schedule *domain.MaintenanceSchedule) error
	FindAllByTenant(tenantID uuid.UUID, filter MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, int64, error)
	FindAllByTenantExport(tenantID uuid.UUID, filter MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceSchedule, error)
	Update(schedule *domain.MaintenanceSchedule) error
	Delete(schedule *domain.MaintenanceSchedule) error
}

type maintenanceScheduleRepository struct {
	db *gorm.DB
}

func NewMaintenanceScheduleRepository(db *gorm.DB) MaintenanceScheduleRepository {
	return &maintenanceScheduleRepository{db}
}

func (r *maintenanceScheduleRepository) Create(schedule *domain.MaintenanceSchedule) error {
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	return db.Create(schedule).Error
}

func (r *maintenanceScheduleRepository) FindAllByTenant(tenantID uuid.UUID, filter MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, int64, error) {
	var schedules []domain.MaintenanceSchedule
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.ScheduleType) != "" {
		query = query.Where("schedule_type = ?", filter.ScheduleType)
	}
	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssetID != uuid.Nil {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.DueFrom != nil {
		query = query.Where("next_due_date >= ?", *filter.DueFrom)
	}
	if filter.DueTo != nil {
		query = query.Where("next_due_date <= ?", *filter.DueTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := mapMaintenanceSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Asset").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&schedules).Error; err != nil {
		return nil, 0, err
	}
	return schedules, total, nil
}

func (r *maintenanceScheduleRepository) FindAllByTenantExport(tenantID uuid.UUID, filter MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, error) {
	var schedules []domain.MaintenanceSchedule
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.ScheduleType) != "" {
		query = query.Where("schedule_type = ?", filter.ScheduleType)
	}
	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssetID != uuid.Nil {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if filter.DueFrom != nil {
		query = query.Where("next_due_date >= ?", *filter.DueFrom)
	}
	if filter.DueTo != nil {
		query = query.Where("next_due_date <= ?", *filter.DueTo)
	}

	sortBy := mapMaintenanceSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	if err := query.Preload("Asset").
		Order(sortBy + " " + sortDir).
		Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *maintenanceScheduleRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceSchedule, error) {
	var schedule domain.MaintenanceSchedule
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	err := repobase.WithTenant(db, tenantID).Preload("Asset").Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *maintenanceScheduleRepository) Update(schedule *domain.MaintenanceSchedule) error {
	updates := map[string]interface{}{
		"asset_id":      schedule.AssetID,
		"schedule_type": schedule.ScheduleType,
		"title":         schedule.Title,
		"interval_days": schedule.IntervalDays,
		"next_due_date": schedule.NextDueDate,
		"last_done_at":  schedule.LastDoneAt,
		"status":        schedule.Status,
		"notes":         schedule.Notes,
		"updated_at":    time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	return repobase.WithTenant(db, schedule.TenantID).
		Set("audit_record_id", schedule.ID).
		Where("id = ?", schedule.ID).
		Updates(updates).
		Error
}

func (r *maintenanceScheduleRepository) Delete(schedule *domain.MaintenanceSchedule) error {
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceSchedule{})
	return repobase.WithTenant(db, schedule.TenantID).Delete(schedule).Error
}

func mapMaintenanceSortBy(sortBy string) string {
	switch sortBy {
	case "next_due_date":
		return "next_due_date"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
