package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ComplaintStatusOpen       = "open"
	ComplaintStatusInProgress = "in_progress"
	ComplaintStatusDone       = "done"
)

type Complaint struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	AssetID        uuid.UUID      `gorm:"type:uuid;index;not null" json:"asset_id"`
	ReportedBy     uuid.UUID      `gorm:"type:uuid;index;not null" json:"reported_by"`
	AssignedTo     *uuid.UUID     `gorm:"type:uuid;index" json:"assigned_to,omitempty"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Description    string         `gorm:"type:text;not null" json:"description"`
	ResolutionNote string         `gorm:"type:text" json:"resolution_note,omitempty"`
	Status         string         `gorm:"size:50;default:'open';index" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Asset    Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Reporter User  `gorm:"foreignKey:ReportedBy" json:"reporter,omitempty"`
	Assignee User  `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Tenant   Tenant
}

func IsValidComplaintStatus(status string) bool {
	switch status {
	case ComplaintStatusOpen, ComplaintStatusInProgress, ComplaintStatusDone:
		return true
	default:
		return false
	}
}
