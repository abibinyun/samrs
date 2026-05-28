package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ScheduleTypeMaintenance = "maintenance"
	ScheduleTypeCalibration = "calibration"

	ScheduleStatusScheduled = "scheduled"
	ScheduleStatusDue       = "due"
	ScheduleStatusCompleted = "completed"
)

type MaintenanceSchedule struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	AssetID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"asset_id"`
	ScheduleType string         `gorm:"size:50;not null;index" json:"schedule_type"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	IntervalDays int            `gorm:"default:0" json:"interval_days"`
	NextDueDate  *time.Time     `json:"next_due_date,omitempty"`
	LastDoneAt   *time.Time     `json:"last_done_at,omitempty"`
	Status       string         `gorm:"size:50;default:'scheduled';index" json:"status"`
	Notes        string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

func IsValidScheduleType(value string) bool {
	switch value {
	case ScheduleTypeMaintenance, ScheduleTypeCalibration:
		return true
	default:
		return false
	}
}

func IsValidScheduleStatus(value string) bool {
	switch value {
	case ScheduleStatusScheduled, ScheduleStatusDue, ScheduleStatusCompleted:
		return true
	default:
		return false
	}
}
