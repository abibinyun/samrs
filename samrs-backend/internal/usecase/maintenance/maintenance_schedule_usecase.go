package usecase

import (
	"strings"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type MaintenanceScheduleUsecase interface {
	CreateSchedule(input MaintenanceScheduleInput) (*domain.MaintenanceSchedule, error)
	GetAllSchedules(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, int64, error)
	GetScheduleByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceSchedule, error)
	UpdateSchedule(input MaintenanceScheduleUpdateInput) (*domain.MaintenanceSchedule, error)
	DeleteSchedule(tenantID uuid.UUID, id uint) error
	CompleteSchedule(tenantID uuid.UUID, id uint, note string) (*domain.MaintenanceSchedule, error)
}

type MaintenanceScheduleInput struct {
	TenantID     uuid.UUID
	AssetID      uuid.UUID
	ScheduleType string
	Title        string
	IntervalDays int
	NextDueDate  string
	Status       string
	Notes        string
}

type MaintenanceScheduleUpdateInput struct {
	TenantID     uuid.UUID
	ID           uint
	ScheduleType string
	Title        string
	IntervalDays int
	NextDueDate  string
	Status       string
	Notes        string
}

type maintenanceScheduleUsecase struct {
	repo      repository.MaintenanceScheduleRepository
	assetRepo repository.AssetRepository
}

func NewMaintenanceScheduleUsecase(repo repository.MaintenanceScheduleRepository, ar repository.AssetRepository) MaintenanceScheduleUsecase {
	return &maintenanceScheduleUsecase{repo: repo, assetRepo: ar}
}

func (u *maintenanceScheduleUsecase) CreateSchedule(input MaintenanceScheduleInput) (*domain.MaintenanceSchedule, error) {
	if input.AssetID == uuid.Nil {
		return nil, util.ErrValidation("asset ID wajib diisi")
	}
	if _, err := u.assetRepo.FindByIDAndTenant(input.AssetID, input.TenantID); err != nil {
		return nil, util.ErrNotFound("asset tidak ditemukan atau akses ditolak")
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, util.ErrValidation("judul jadwal wajib diisi")
	}

	scheduleType := strings.TrimSpace(input.ScheduleType)
	if scheduleType == "" {
		scheduleType = domain.ScheduleTypeMaintenance
	}
	if !domain.IsValidScheduleType(scheduleType) {
		return nil, util.ErrValidation("tipe jadwal tidak valid")
	}

	var nextDue *time.Time
	if strings.TrimSpace(input.NextDueDate) != "" {
		parsed, err := time.Parse("2006-01-02", input.NextDueDate)
		if err != nil {
			return nil, util.ErrValidation("format next_due_date harus YYYY-MM-DD")
		}
		nextDue = &parsed
	}
	if nextDue == nil {
		return nil, util.ErrValidation("next_due_date wajib diisi")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = domain.ScheduleStatusScheduled
	}
	if !domain.IsValidScheduleStatus(status) {
		return nil, util.ErrValidation("status jadwal tidak valid")
	}

	schedule := &domain.MaintenanceSchedule{
		TenantID:     input.TenantID,
		AssetID:      input.AssetID,
		ScheduleType: scheduleType,
		Title:        title,
		IntervalDays: input.IntervalDays,
		NextDueDate:  nextDue,
		Status:       status,
		Notes:        strings.TrimSpace(input.Notes),
	}

	if err := u.repo.Create(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (u *maintenanceScheduleUsecase) GetAllSchedules(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *maintenanceScheduleUsecase) GetScheduleByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceSchedule, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *maintenanceScheduleUsecase) UpdateSchedule(input MaintenanceScheduleUpdateInput) (*domain.MaintenanceSchedule, error) {
	schedule, err := u.repo.FindByID(input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, util.ErrValidation("judul jadwal wajib diisi")
	}

	scheduleType := strings.TrimSpace(input.ScheduleType)
	if scheduleType == "" {
		scheduleType = domain.ScheduleTypeMaintenance
	}
	if !domain.IsValidScheduleType(scheduleType) {
		return nil, util.ErrValidation("tipe jadwal tidak valid")
	}

	var nextDue *time.Time
	if strings.TrimSpace(input.NextDueDate) != "" {
		parsed, err := time.Parse("2006-01-02", input.NextDueDate)
		if err != nil {
			return nil, util.ErrValidation("format next_due_date harus YYYY-MM-DD")
		}
		nextDue = &parsed
	}
	if nextDue == nil {
		return nil, util.ErrValidation("next_due_date wajib diisi")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = domain.ScheduleStatusScheduled
	}
	if !domain.IsValidScheduleStatus(status) {
		return nil, util.ErrValidation("status jadwal tidak valid")
	}

	schedule.ScheduleType = scheduleType
	schedule.Title = title
	schedule.IntervalDays = input.IntervalDays
	schedule.NextDueDate = nextDue
	schedule.Status = status
	schedule.Notes = strings.TrimSpace(input.Notes)

	if err := u.repo.Update(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (u *maintenanceScheduleUsecase) DeleteSchedule(tenantID uuid.UUID, id uint) error {
	schedule, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(schedule)
}

func (u *maintenanceScheduleUsecase) CompleteSchedule(tenantID uuid.UUID, id uint, note string) (*domain.MaintenanceSchedule, error) {
	schedule, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	schedule.LastDoneAt = &now
	schedule.Notes = strings.TrimSpace(note)

	if schedule.IntervalDays > 0 {
		next := now.AddDate(0, 0, schedule.IntervalDays)
		schedule.NextDueDate = &next
		schedule.Status = domain.ScheduleStatusScheduled
	} else {
		schedule.Status = domain.ScheduleStatusCompleted
	}

	if err := u.repo.Update(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}
