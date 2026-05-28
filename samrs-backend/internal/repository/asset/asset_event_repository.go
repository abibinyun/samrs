package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetEventFilter struct {
	EventType string
	DateFrom  *time.Time
	DateTo    *time.Time
	Page      int
	PerPage   int
	SortBy    string
	SortDir   string
}

type AssetEventRepository interface {
	Create(event *domain.AssetEvent) error
	ListByAsset(tenantID uuid.UUID, assetID uuid.UUID, filter AssetEventFilter) ([]domain.AssetEvent, int64, error)
}

type assetEventRepository struct {
	db *gorm.DB
}

func NewAssetEventRepository(db *gorm.DB) AssetEventRepository {
	return &assetEventRepository{db}
}

func (r *assetEventRepository) Create(event *domain.AssetEvent) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetEvent{})
	return db.Create(event).Error
}

func (r *assetEventRepository) ListByAsset(tenantID uuid.UUID, assetID uuid.UUID, filter AssetEventFilter) ([]domain.AssetEvent, int64, error) {
	var events []domain.AssetEvent
	db := repobase.NewDB(r.db).Model(&domain.AssetEvent{})
	query := repobase.WithTenant(db, tenantID).Where("asset_id = ?", assetID)

	if strings.TrimSpace(filter.EventType) != "" {
		query = query.Where("event_type = ?", filter.EventType)
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

	sortBy := mapAssetEventSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("User").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func mapAssetEventSortBy(sortBy string) string {
	switch sortBy {
	case "created_at":
		return "created_at"
	case "event_type":
		return "event_type"
	default:
		return "created_at"
	}
}
