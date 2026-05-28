package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditTrail struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Action    string         `gorm:"size:20;not null" json:"action"`
	TableName string         `gorm:"size:100;not null" json:"table_name"`
	RecordID  string         `gorm:"size:100;not null" json:"record_id"`
	OldData   datatypes.JSON `gorm:"type:jsonb" json:"old_data,omitempty"`
	NewData   datatypes.JSON `gorm:"type:jsonb" json:"new_data,omitempty"`
	IP        string         `gorm:"size:50" json:"ip_address"`
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
}
