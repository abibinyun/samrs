package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Slug       string    `gorm:"unique;not null" json:"slug"`
	Address    string    `gorm:"size:255" json:"address,omitempty"`
	TenantType string    `gorm:"size:50" json:"type,omitempty"`
	Status     string    `gorm:"default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
