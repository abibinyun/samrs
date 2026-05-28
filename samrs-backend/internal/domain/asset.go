package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AssetStatusReady       = "ready"
	AssetStatusBroken      = "broken"
	AssetStatusMaintenance = "maintenance"
)

type Asset struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID     uuid.UUID      `gorm:"type:uuid;index;not null;index:idx_asset_tenant_code;index:idx_asset_tenant_created;index:idx_asset_tenant_status;index:idx_asset_tenant_category;index:idx_asset_tenant_room" json:"tenant_id"`
	CategoryID   uint           `gorm:"not null;index:idx_asset_tenant_category" json:"category_id"`
	RoomID       *uuid.UUID     `gorm:"type:uuid;index:idx_asset_tenant_room" json:"room_id,omitempty"`
	BedID        *uint          `gorm:"index:idx_asset_tenant_bed" json:"bed_id,omitempty"`
	VendorID     *uint          `gorm:"index:idx_asset_tenant_vendor" json:"vendor_id,omitempty"`
	BrandID      *uint          `gorm:"index:idx_asset_tenant_brand" json:"brand_id,omitempty"`
	ModelID      *uint          `gorm:"index:idx_asset_tenant_model" json:"model_id,omitempty"`
	Code         string         `gorm:"size:100;not null;index:idx_asset_tenant_code" json:"code"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Brand        string         `gorm:"size:100" json:"brand,omitempty"`
	Model        string         `gorm:"size:100" json:"model,omitempty"`
	Status       string         `gorm:"size:50;default:'ready';index:idx_asset_tenant_status" json:"status"`
	PurchaseDate *time.Time     `json:"purchase_date,omitempty"`
	CreatedAt    time.Time      `gorm:"index:idx_asset_tenant_created" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Category Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Room     Room       `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Bed      Bed        `gorm:"foreignKey:BedID" json:"bed,omitempty"`
	Vendor   Vendor     `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	BrandRef AssetBrand `gorm:"foreignKey:BrandID" json:"brand_master,omitempty"`
	ModelRef AssetModel `gorm:"foreignKey:ModelID" json:"model_master,omitempty"`
	Tenant   Tenant     `gorm:"foreignKey:TenantID" json:"-"`
}

func IsValidAssetStatus(status string) bool {
	switch status {
	case AssetStatusReady, AssetStatusBroken, AssetStatusMaintenance:
		return true
	default:
		return false
	}
}
