package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MaintenanceDocCertificate = "certificate"
	MaintenanceDocReport      = "report"
	MaintenanceDocOther       = "other"
)

type MaintenanceDocument struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	ScheduleID uint      `gorm:"index;not null" json:"schedule_id"`
	AssetID    uuid.UUID `gorm:"type:uuid;index;not null" json:"asset_id"`
	UploadedBy uuid.UUID `gorm:"type:uuid;index;not null" json:"uploaded_by"`
	DocType    string    `gorm:"size:50;not null;index" json:"doc_type"`
	Filename   string    `gorm:"size:255;not null" json:"filename"`
	FilePath   string    `gorm:"size:500;not null" json:"file_path"`
	MimeType   string    `gorm:"size:100" json:"mime_type,omitempty"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`

	Schedule MaintenanceSchedule `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Asset    Asset               `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Uploader User                `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`

	FileURL string `gorm:"-" json:"file_url,omitempty"`
}
