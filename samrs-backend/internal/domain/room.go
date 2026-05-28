package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;index;not null;index:idx_room_tenant_code;index:idx_room_tenant_created" json:"tenant_id"`
	Name      string         `gorm:"not null" json:"name"`
	Code      string         `gorm:"not null;index:idx_room_tenant_code" json:"code"` // Misal: R001, VIP-01
	Location  string         `json:"location"`                                        // Misal: Gedung A Lantai 2
	CreatedAt time.Time      `gorm:"index:idx_room_tenant_created" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
}
