package repobase

import (
	"strings"

	"samrs-backend/internal/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func WithTenant(db *gorm.DB, tenantID uuid.UUID) *gorm.DB {
	return config.MarkTenantScope(db).Where("tenant_id = ?", tenantID)
}

func WithoutTenantScope(db *gorm.DB) *gorm.DB {
	return config.WithSkipTenantScope(db)
}

func MapSortDir(sortDir string) string {
	sortDir = strings.ToLower(strings.TrimSpace(sortDir))
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}
	return sortDir
}

// NewDB creates a fresh session while preserving audit metadata when present.
func NewDB(db *gorm.DB) *gorm.DB {
	if db == nil {
		return db
	}
	newDB := db.Session(&gorm.Session{NewDB: true})
	if meta, ok := db.Get("audit_meta"); ok {
		newDB = newDB.Set("audit_meta", meta)
	}
	if action, ok := db.Get("audit_action"); ok {
		newDB = newDB.Set("audit_action", action)
	}
	if recordID, ok := db.Get("audit_record_id"); ok {
		newDB = newDB.Set("audit_record_id", recordID)
	}
	if _, ok := db.Get("skip_audit"); ok {
		newDB = newDB.Set("skip_audit", true)
	}
	return newDB
}
