package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	BedStatusAvailable   = "available"
	BedStatusOccupied    = "occupied"
	BedStatusMaintenance = "maintenance"
)

type Bed struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index;index:idx_bed_tenant_room_code;index:idx_bed_tenant_room;index:idx_bed_tenant_created;index:idx_bed_tenant_status" json:"tenant_id"`
	RoomID    uuid.UUID      `gorm:"type:uuid;not null;index;index:idx_bed_tenant_room_code;index:idx_bed_tenant_room" json:"room_id"` // Sesuaikan ke UUID
	Code      string         `gorm:"size:50;not null;index:idx_bed_tenant_room_code" json:"code"`
	Status    string         `gorm:"size:20;default:'available';index:idx_bed_tenant_status" json:"status"`
	CreatedAt time.Time      `gorm:"index:idx_bed_tenant_created" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi ke Room
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

func IsValidBedStatus(status string) bool {
	switch status {
	case BedStatusAvailable, BedStatusOccupied, BedStatusMaintenance:
		return true
	default:
		return false
	}
}
