package documentrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentFilter struct {
	Search   string
	DocType  string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

type DocumentRepository interface {
	Create(doc *domain.Document) error
	Update(doc *domain.Document) error
	FindByID(tenantID uuid.UUID, id uint) (*domain.Document, error)
	List(tenantID uuid.UUID, filter DocumentFilter) ([]domain.Document, int64, error)
	Delete(doc *domain.Document) error
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(doc *domain.Document) error {
	db := repobase.NewDB(r.db).Model(&domain.Document{})
	return db.Create(doc).Error
}

func (r *documentRepository) Update(doc *domain.Document) error {
	updates := map[string]interface{}{
		"title":       doc.Title,
		"doc_type":    doc.DocType,
		"description": doc.Description,
		"updated_at":  time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Document{})
	return repobase.WithTenant(db, doc.TenantID).
		Set("audit_record_id", doc.ID).
		Where("id = ?", doc.ID).
		Updates(updates).
		Error
}

func (r *documentRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Document, error) {
	var doc domain.Document
	db := repobase.NewDB(r.db).Model(&domain.Document{})
	if err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) List(tenantID uuid.UUID, filter DocumentFilter) ([]domain.Document, int64, error) {
	var (
		docs  []domain.Document
		total int64
	)

	db := repobase.NewDB(r.db).Model(&domain.Document{})
	query := repobase.WithTenant(db, tenantID)

	if filter.Search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+filter.Search+"%")
	}
	if filter.DocType != "" {
		query = query.Where("doc_type = ?", filter.DocType)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	switch filter.SortBy {
	case "title", "created_at":
		sortBy = filter.SortBy
	}
	sortDir := "desc"
	if filter.SortDir == "asc" {
		sortDir = "asc"
	}
	offset := (filter.Page - 1) * filter.PerPage

	if err := query.Order(sortBy + " " + sortDir).Limit(filter.PerPage).Offset(offset).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func (r *documentRepository) Delete(doc *domain.Document) error {
	db := repobase.NewDB(r.db).Model(&domain.Document{})
	return repobase.WithTenant(db, doc.TenantID).Delete(doc).Error
}
