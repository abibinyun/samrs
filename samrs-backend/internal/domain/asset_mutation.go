package domain

import (
	"time"

	"github.com/google/uuid"
)

type AssetMutation struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TenantID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"tenant_id"`
	AssetID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"asset_id"`
	FromRoomID *uuid.UUID `gorm:"type:uuid;index" json:"from_room_id,omitempty"`
	FromBedID  *uint      `gorm:"index" json:"from_bed_id,omitempty"`
	ToRoomID   *uuid.UUID `gorm:"type:uuid;index" json:"to_room_id,omitempty"`
	ToBedID    *uint      `gorm:"index" json:"to_bed_id,omitempty"`
	Reason     string     `gorm:"type:text" json:"reason,omitempty"`
	MovedBy    uuid.UUID  `gorm:"type:uuid;index;not null" json:"moved_by"`
	CreatedAt  time.Time  `json:"created_at"`

	Asset    Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Mover    User  `gorm:"foreignKey:MovedBy" json:"mover,omitempty"`
	FromRoom Room  `gorm:"foreignKey:FromRoomID" json:"from_room,omitempty"`
	ToRoom   Room  `gorm:"foreignKey:ToRoomID" json:"to_room,omitempty"`
	FromBed  Bed   `gorm:"foreignKey:FromBedID" json:"from_bed,omitempty"`
	ToBed    Bed   `gorm:"foreignKey:ToBedID" json:"to_bed,omitempty"`
}
