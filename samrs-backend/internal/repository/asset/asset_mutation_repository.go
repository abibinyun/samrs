package assetrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetMutationFilter struct {
	AssetID  uuid.UUID
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

type AssetMutationRepository interface {
	Create(mutation *domain.AssetMutation) error
	FindAllByTenant(tenantID uuid.UUID, filter AssetMutationFilter) ([]domain.AssetMutation, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.AssetMutation, error)
}

type assetMutationRepository struct {
	db *gorm.DB
}

func NewAssetMutationRepository(db *gorm.DB) AssetMutationRepository {
	return &assetMutationRepository{db}
}

func (r *assetMutationRepository) Create(mutation *domain.AssetMutation) error {
	db := repobase.NewDB(r.db).Model(&domain.AssetMutation{})
	return db.Create(mutation).Error
}

func (r *assetMutationRepository) FindAllByTenant(tenantID uuid.UUID, filter AssetMutationFilter) ([]domain.AssetMutation, int64, error) {
	var mutations []domain.AssetMutation
	db := repobase.NewDB(r.db).Model(&domain.AssetMutation{})
	query := repobase.WithTenant(db, tenantID)

	if filter.AssetID != uuid.Nil {
		query = query.Where("asset_id = ?", filter.AssetID)
	}
	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("reason ILIKE ?", like)
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

	sortBy := mapAssetMutationSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Asset").
		Preload("Mover").
		Preload("FromRoom").
		Preload("ToRoom").
		Preload("FromBed").
		Preload("ToBed").
		Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&mutations).Error; err != nil {
		return nil, 0, err
	}
	return mutations, total, nil
}

func (r *assetMutationRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetMutation, error) {
	var mutation domain.AssetMutation
	db := repobase.NewDB(r.db).Model(&domain.AssetMutation{})
	if err := repobase.WithTenant(db, tenantID).
		Preload("Asset").
		Preload("Mover").
		Preload("FromRoom").
		Preload("ToRoom").
		Preload("FromBed").
		Preload("ToBed").
		Where("id = ?", id).
		First(&mutation).Error; err != nil {
		return nil, err
	}
	return &mutation, nil
}

func mapAssetMutationSortBy(sortBy string) string {
	switch sortBy {
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
