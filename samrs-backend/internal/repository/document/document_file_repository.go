package documentrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentFileRepository interface {
	CreateMany(files []domain.DocumentFile) error
	ListByDocument(tenantID uuid.UUID, documentID uint) ([]domain.DocumentFile, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error)
	Delete(file *domain.DocumentFile) error
	DeleteByDocument(tenantID uuid.UUID, documentID uint) error
}

type documentFileRepository struct {
	db *gorm.DB
}

func NewDocumentFileRepository(db *gorm.DB) DocumentFileRepository {
	return &documentFileRepository{db: db}
}

func (r *documentFileRepository) CreateMany(files []domain.DocumentFile) error {
	if len(files) == 0 {
		return nil
	}
	db := repobase.NewDB(r.db).Model(&domain.DocumentFile{})
	return db.Create(&files).Error
}

func (r *documentFileRepository) ListByDocument(tenantID uuid.UUID, documentID uint) ([]domain.DocumentFile, error) {
	var files []domain.DocumentFile
	db := repobase.NewDB(r.db).Model(&domain.DocumentFile{})
	if err := repobase.WithTenant(db, tenantID).Where("document_id = ?", documentID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *documentFileRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error) {
	var file domain.DocumentFile
	db := repobase.NewDB(r.db).Model(&domain.DocumentFile{})
	if err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *documentFileRepository) Delete(file *domain.DocumentFile) error {
	db := repobase.NewDB(r.db).Model(&domain.DocumentFile{})
	return repobase.WithTenant(db, file.TenantID).Delete(file).Error
}

func (r *documentFileRepository) DeleteByDocument(tenantID uuid.UUID, documentID uint) error {
	db := repobase.NewDB(r.db).Model(&domain.DocumentFile{})
	return repobase.WithTenant(db, tenantID).Where("document_id = ?", documentID).Delete(&domain.DocumentFile{}).Error
}
