package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	AssetEventCreated           = "created"
	AssetEventUpdated           = "updated"
	AssetEventDeleted           = "deleted"
	AssetEventComplaintReported = "complaint_reported"
	AssetEventComplaintProgress = "complaint_in_progress"
	AssetEventComplaintResolved = "complaint_done"
	AssetEventStatusChanged     = "status_changed"
	AssetEventMaintenanceDone   = "maintenance_done"
	AssetEventCalibrationDone   = "calibration_done"
	AssetEventMutation          = "mutation"
	AssetEventStockOpname       = "stock_opname"
)

type AssetEvent struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	AssetID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"asset_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	EventType   string         `gorm:"size:100;not null;index" json:"event_type"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Meta        datatypes.JSON `gorm:"type:jsonb" json:"meta,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`

	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
