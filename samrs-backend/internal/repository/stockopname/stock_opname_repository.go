package stockopnamerepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockOpnameFilter struct {
	Status   string
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PerPage  int
	SortBy   string
	SortDir  string
}

type StockOpnameRepository interface {
	Create(session *domain.StockOpnameSession) error
	FindAllByTenant(tenantID uuid.UUID, filter StockOpnameFilter) ([]domain.StockOpnameSession, int64, error)
	FindByID(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error)
	Update(session *domain.StockOpnameSession) error
	Delete(session *domain.StockOpnameSession) error
}

type StockOpnameItemRepository interface {
	Create(item *domain.StockOpnameItem) error
	ListBySession(tenantID uuid.UUID, sessionID uint) ([]domain.StockOpnameItem, error)
	Delete(tenantID uuid.UUID, sessionID uint, itemID uint) error
}

type stockOpnameRepository struct {
	db *gorm.DB
}

type stockOpnameItemRepository struct {
	db *gorm.DB
}

func NewStockOpnameRepository(db *gorm.DB) StockOpnameRepository {
	return &stockOpnameRepository{db}
}

func NewStockOpnameItemRepository(db *gorm.DB) StockOpnameItemRepository {
	return &stockOpnameItemRepository{db}
}

func (r *stockOpnameRepository) Create(session *domain.StockOpnameSession) error {
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameSession{})
	return db.Create(session).Error
}

func (r *stockOpnameRepository) FindAllByTenant(tenantID uuid.UUID, filter StockOpnameFilter) ([]domain.StockOpnameSession, int64, error) {
	var sessions []domain.StockOpnameSession
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameSession{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("opname_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("opname_at <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := mapStockOpnameSortBy(filter.SortBy)
	sortDir := repobase.MapSortDir(filter.SortDir)
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order(sortBy + " " + sortDir).
		Limit(filter.PerPage).Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

func (r *stockOpnameRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error) {
	var session domain.StockOpnameSession
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameSession{})
	if err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *stockOpnameRepository) Update(session *domain.StockOpnameSession) error {
	updates := map[string]interface{}{
		"title":      session.Title,
		"opname_at":  session.OpnameAt,
		"status":     session.Status,
		"notes":      session.Notes,
		"updated_at": time.Now(),
	}
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameSession{})
	return repobase.WithTenant(db, session.TenantID).
		Set("audit_record_id", session.ID).
		Where("id = ?", session.ID).
		Updates(updates).
		Error
}

func (r *stockOpnameRepository) Delete(session *domain.StockOpnameSession) error {
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameSession{})
	return repobase.WithTenant(db, session.TenantID).Delete(session).Error
}

func (r *stockOpnameItemRepository) Create(item *domain.StockOpnameItem) error {
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameItem{})
	return db.Create(item).Error
}

func (r *stockOpnameItemRepository) ListBySession(tenantID uuid.UUID, sessionID uint) ([]domain.StockOpnameItem, error) {
	var items []domain.StockOpnameItem
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameItem{})
	if err := repobase.WithTenant(db, tenantID).
		Where("session_id = ?", sessionID).
		Preload("Asset").
		Preload("Checker").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *stockOpnameItemRepository) Delete(tenantID uuid.UUID, sessionID uint, itemID uint) error {
	db := repobase.NewDB(r.db).Model(&domain.StockOpnameItem{})
	return repobase.WithTenant(db, tenantID).
		Where("id = ? AND session_id = ?", itemID, sessionID).
		Delete(&domain.StockOpnameItem{}).Error
}

func mapStockOpnameSortBy(sortBy string) string {
	switch sortBy {
	case "opname_at":
		return "opname_at"
	case "created_at":
		return "created_at"
	default:
		return "created_at"
	}
}
