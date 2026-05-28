package complaintrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComplaintFilter struct {
	Status     string
	AssetID    uuid.UUID
	AssignedTo uuid.UUID
	ReportedBy uuid.UUID
	Search     string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PerPage    int
	SortBy     string
	SortDir    string
}

type ComplaintRepository interface {
	Create(complaint *domain.Complaint) error
	FindAllByTenant(tenantID uuid.UUID, filter ComplaintFilter) ([]domain.Complaint, int64, error)
	FindAllByTenantExport(tenantID uuid.UUID, filter ComplaintFilter) ([]domain.Complaint, error)
	FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Complaint, error)
	Update(complaint *domain.Complaint) error
	Delete(complaint *domain.Complaint) error
}

type complaintRepository struct {
	db *gorm.DB
}

func NewComplaintRepository(db *gorm.DB) ComplaintRepository {
	return &complaintRepository{db}
}

func (r *complaintRepository) Create(complaint *domain.Complaint) error {
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	return db.Create(complaint).Error
}

func (r *complaintRepository) FindAllByTenant(tenantID uuid.UUID, filter ComplaintFilter) ([]domain.Complaint, int64, error) {
	var complaints []domain.Complaint
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssetID != uuid.Nil {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if filter.AssignedTo != uuid.Nil {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if filter.ReportedBy != uuid.Nil {
		query = query.Where("reported_by = ?", filter.ReportedBy)
	}
	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
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

	sortBy := mapComplaintSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Asset").Preload("Reporter").Preload("Assignee").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&complaints).Error; err != nil {
		return nil, 0, err
	}
	return complaints, total, nil
}

func (r *complaintRepository) FindAllByTenantExport(tenantID uuid.UUID, filter ComplaintFilter) ([]domain.Complaint, error) {
	var complaints []domain.Complaint
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssetID != uuid.Nil {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if filter.AssignedTo != uuid.Nil {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if filter.ReportedBy != uuid.Nil {
		query = query.Where("reported_by = ?", filter.ReportedBy)
	}
	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	sortBy := mapComplaintSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	if err := query.Preload("Asset").Preload("Reporter").Preload("Assignee").
		Order(sortBy + " " + sortDir).
		Find(&complaints).Error; err != nil {
		return nil, err
	}
	return complaints, nil
}

func (r *complaintRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Complaint, error) {
	var complaint domain.Complaint
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	err := repobase.WithTenant(db, tenantID).Preload("Asset").Preload("Reporter").Preload("Assignee").
		Where("id = ?", id).First(&complaint).Error
	if err != nil {
		return nil, err
	}
	return &complaint, nil
}

func (r *complaintRepository) Update(complaint *domain.Complaint) error {
	updates := map[string]interface{}{
		"asset_id":       complaint.AssetID,
		"reported_by":    complaint.ReportedBy,
		"assigned_to":    complaint.AssignedTo,
		"title":          complaint.Title,
		"description":    complaint.Description,
		"resolution_note": complaint.ResolutionNote,
		"status":         complaint.Status,
		"updated_at":     time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	return repobase.WithTenant(db, complaint.TenantID).
		Set("audit_record_id", complaint.ID).
		Where("id = ?", complaint.ID).
		Updates(updates).
		Error
}

func (r *complaintRepository) Delete(complaint *domain.Complaint) error {
	db := repobase.NewDB(r.db).Model(&domain.Complaint{})
	return repobase.WithTenant(db, complaint.TenantID).Delete(complaint).Error
}

func mapComplaintSortBy(sortBy string) string {
	switch sortBy {
	case "status":
		return "status"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
