package maintenancerepo

import (
	repobase "samrs-backend/internal/repository/base"
	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceDocumentRepository interface {
	Create(doc *domain.MaintenanceDocument) error
	ListBySchedule(tenantID uuid.UUID, scheduleID uint) ([]domain.MaintenanceDocument, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceDocument, error)
	Delete(doc *domain.MaintenanceDocument) error
}

type maintenanceDocumentRepository struct {
	db *gorm.DB
}

func NewMaintenanceDocumentRepository(db *gorm.DB) MaintenanceDocumentRepository {
	return &maintenanceDocumentRepository{db}
}

func (r *maintenanceDocumentRepository) Create(doc *domain.MaintenanceDocument) error {
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceDocument{})
	return db.Create(doc).Error
}

func (r *maintenanceDocumentRepository) ListBySchedule(tenantID uuid.UUID, scheduleID uint) ([]domain.MaintenanceDocument, error) {
	var docs []domain.MaintenanceDocument
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceDocument{})
	if err := repobase.WithTenant(db, tenantID).
		Where("schedule_id = ?", scheduleID).
		Preload("Uploader").
		Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *maintenanceDocumentRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceDocument, error) {
	var doc domain.MaintenanceDocument
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceDocument{})
	if err := repobase.WithTenant(db, tenantID).
		Preload("Uploader").
		Where("id = ?", id).
		First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *maintenanceDocumentRepository) Delete(doc *domain.MaintenanceDocument) error {
	db := repobase.NewDB(r.db).Model(&domain.MaintenanceDocument{})
	return repobase.WithTenant(db, doc.TenantID).Delete(doc).Error
}
