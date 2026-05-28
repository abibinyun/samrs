package repository

import (
	tenantrepo "samrs-backend/internal/repository/tenant"

	"gorm.io/gorm"
)

type TenantFilter = tenantrepo.TenantFilter

type TenantRepository = tenantrepo.TenantRepository

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return tenantrepo.NewTenantRepository(db)
}
