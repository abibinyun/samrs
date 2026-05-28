package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetStatus struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index;index:idx_status_tenant_code;index:idx_status_tenant_created" json:"tenant_id"`
	Code        string         `gorm:"size:120;not null;index:idx_status_tenant_code" json:"code"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_status_tenant_created" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
