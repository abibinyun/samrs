package repository

import (
	auditrepo "samrs-backend/internal/repository/audit"

	"gorm.io/gorm"
)

type AuditTrailFilter = auditrepo.AuditTrailFilter

type AuditTrailRepository = auditrepo.AuditTrailRepository

func NewAuditTrailRepository(db *gorm.DB) AuditTrailRepository {
	return auditrepo.NewAuditTrailRepository(db)
}
