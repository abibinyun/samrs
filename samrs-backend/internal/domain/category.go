package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;index;index:idx_category_tenant_slug;index:idx_category_tenant_created" json:"tenant_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"size:120;not null;index:idx_category_tenant_slug" json:"slug"`
	Description string         `gorm:"text" json:"description"`
	CreatedAt   time.Time      `gorm:"index:idx_category_tenant_created" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
