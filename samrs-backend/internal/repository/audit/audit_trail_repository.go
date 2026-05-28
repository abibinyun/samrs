package auditrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditTrailRepository interface {
	Create(audit *domain.AuditTrail) error
	ListByTenant(tenantID uuid.UUID, filter AuditTrailFilter) ([]domain.AuditTrail, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.AuditTrail, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
}

type auditTrailRepository struct {
	db *gorm.DB
}

type AuditTrailFilter struct {
	Action   string
	Table    string
	UserID   uuid.UUID
	RecordID string
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

func NewAuditTrailRepository(db *gorm.DB) AuditTrailRepository {
	return &auditTrailRepository{db}
}

func (r *auditTrailRepository) Create(audit *domain.AuditTrail) error {
	db := repobase.NewDB(r.db).Model(&domain.AuditTrail{})
	return db.Create(audit).Error
}

func (r *auditTrailRepository) ListByTenant(tenantID uuid.UUID, filter AuditTrailFilter) ([]domain.AuditTrail, int64, error) {
	var items []domain.AuditTrail
	db := repobase.NewDB(r.db).Model(&domain.AuditTrail{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Action) != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if strings.TrimSpace(filter.Table) != "" {
		query = query.Where("table_name = ?", filter.Table)
	}
	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if strings.TrimSpace(filter.RecordID) != "" {
		query = query.Where("record_id = ?", filter.RecordID)
	}
	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("table_name ILIKE ? OR record_id ILIKE ? OR action ILIKE ?", like, like, like)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := mapAuditSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *auditTrailRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AuditTrail, error) {
	var audit domain.AuditTrail
	db := repobase.NewDB(r.db).Model(&domain.AuditTrail{})
	err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&audit).Error
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

func (r *auditTrailRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.AuditTrail{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func mapAuditSortBy(sortBy string) string {
	switch sortBy {
	case "action":
		return "action"
	case "table_name":
		return "table_name"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
