package domain

import "github.com/google/uuid"

type Role struct {
	ID            int        `gorm:"primaryKey" json:"id"`
	TenantID      *uuid.UUID `gorm:"type:uuid" json:"tenant_id,omitempty"`
	Name          string     `gorm:"not null" json:"name"`
	IsSystem      bool       `gorm:"default:false" json:"is_system"`
	IsTenantAdmin bool       `gorm:"default:false" json:"is_tenant_admin"`
}
