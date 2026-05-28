package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StockOpnameStatusDraft  = "draft"
	StockOpnameStatusClosed = "closed"

	StockOpnameConditionMatch   = "match"
	StockOpnameConditionMissing = "missing"
	StockOpnameConditionDamaged = "damaged"
)

type StockOpnameSession struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	OpnameAt  time.Time      `json:"opname_at"`
	Status    string         `gorm:"size:50;default:'draft';index" json:"status"`
	Notes     string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type StockOpnameItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	AssetID   uuid.UUID `gorm:"type:uuid;index;not null" json:"asset_id"`
	Condition string    `gorm:"size:50;not null;index" json:"condition"`
	Note      string    `gorm:"type:text" json:"note,omitempty"`
	CheckedBy uuid.UUID `gorm:"type:uuid;index;not null" json:"checked_by"`
	CreatedAt time.Time `json:"created_at"`

	Session StockOpnameSession `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	Asset   Asset              `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Checker User               `gorm:"foreignKey:CheckedBy" json:"checker,omitempty"`
}

func IsValidStockOpnameStatus(value string) bool {
	switch value {
	case StockOpnameStatusDraft, StockOpnameStatusClosed:
		return true
	default:
		return false
	}
}

func IsValidStockOpnameCondition(value string) bool {
	switch value {
	case StockOpnameConditionMatch, StockOpnameConditionMissing, StockOpnameConditionDamaged:
		return true
	default:
		return false
	}
}
