package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vendor struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index;index:idx_vendor_tenant_code;index:idx_vendor_tenant_created" json:"tenant_id"`
	Code        string         `gorm:"size:120;not null;index:idx_vendor_tenant_code" json:"code"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	ContactName string         `gorm:"size:255" json:"contact_name,omitempty"`
	Phone       string         `gorm:"size:50" json:"phone,omitempty"`
	Email       string         `gorm:"size:255" json:"email,omitempty"`
	Address     string         `gorm:"size:255" json:"address,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_vendor_tenant_created" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
