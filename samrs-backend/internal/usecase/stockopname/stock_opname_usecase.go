package usecase

import (
	"strings"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type StockOpnameUsecase interface {
	CreateSession(input StockOpnameSessionInput) (*domain.StockOpnameSession, error)
	GetAllSessions(tenantID uuid.UUID, filter repository.StockOpnameFilter) ([]domain.StockOpnameSession, int64, error)
	GetSessionByID(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error)
	UpdateSession(input StockOpnameSessionUpdateInput) (*domain.StockOpnameSession, error)
	CloseSession(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error)
	AddItem(input StockOpnameItemInput) (*domain.StockOpnameItem, error)
	ListItems(tenantID uuid.UUID, sessionID uint) ([]domain.StockOpnameItem, error)
}

type StockOpnameSessionInput struct {
	TenantID uuid.UUID
	Title    string
	OpnameAt string
	Notes    string
}

type StockOpnameSessionUpdateInput struct {
	TenantID uuid.UUID
	ID       uint
	Title    string
	OpnameAt string
	Notes    string
	Status   string
}

type StockOpnameItemInput struct {
	TenantID  uuid.UUID
	SessionID uint
	AssetID   uuid.UUID
	Condition string
	Note      string
	CheckedBy uuid.UUID
}

type stockOpnameUsecase struct {
	repo      repository.StockOpnameRepository
	itemRepo  repository.StockOpnameItemRepository
	assetRepo repository.AssetRepository
}

func NewStockOpnameUsecase(
	repo repository.StockOpnameRepository,
	itemRepo repository.StockOpnameItemRepository,
	ar repository.AssetRepository,
) StockOpnameUsecase {
	return &stockOpnameUsecase{
		repo:      repo,
		itemRepo:  itemRepo,
		assetRepo: ar,
	}
}

func (u *stockOpnameUsecase) CreateSession(input StockOpnameSessionInput) (*domain.StockOpnameSession, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, util.ErrValidation("judul opname wajib diisi")
	}

	opnameAt := time.Now()
	if strings.TrimSpace(input.OpnameAt) != "" {
		parsed, err := time.Parse("2006-01-02", input.OpnameAt)
		if err != nil {
			return nil, util.ErrValidation("format opname_at harus YYYY-MM-DD")
		}
		opnameAt = parsed
	}

	session := &domain.StockOpnameSession{
		TenantID: input.TenantID,
		Title:    title,
		OpnameAt: opnameAt,
		Status:   domain.StockOpnameStatusDraft,
		Notes:    strings.TrimSpace(input.Notes),
	}

	if err := u.repo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *stockOpnameUsecase) GetAllSessions(tenantID uuid.UUID, filter repository.StockOpnameFilter) ([]domain.StockOpnameSession, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *stockOpnameUsecase) GetSessionByID(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *stockOpnameUsecase) UpdateSession(input StockOpnameSessionUpdateInput) (*domain.StockOpnameSession, error) {
	session, err := u.repo.FindByID(input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, util.ErrValidation("judul opname wajib diisi")
	}

	opnameAt := session.OpnameAt
	if strings.TrimSpace(input.OpnameAt) != "" {
		parsed, err := time.Parse("2006-01-02", input.OpnameAt)
		if err != nil {
			return nil, util.ErrValidation("format opname_at harus YYYY-MM-DD")
		}
		opnameAt = parsed
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = session.Status
	}
	if !domain.IsValidStockOpnameStatus(status) {
		return nil, util.ErrValidation("status opname tidak valid")
	}

	session.Title = title
	session.OpnameAt = opnameAt
	session.Notes = strings.TrimSpace(input.Notes)
	session.Status = status

	if err := u.repo.Update(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *stockOpnameUsecase) CloseSession(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error) {
	session, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}
	session.Status = domain.StockOpnameStatusClosed
	if err := u.repo.Update(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *stockOpnameUsecase) AddItem(input StockOpnameItemInput) (*domain.StockOpnameItem, error) {
	if input.AssetID == uuid.Nil {
		return nil, util.ErrValidation("asset ID wajib diisi")
	}
	if input.CheckedBy == uuid.Nil {
		return nil, util.ErrValidation("checked_by wajib diisi")
	}

	if _, err := u.repo.FindByID(input.TenantID, input.SessionID); err != nil {
		return nil, util.ErrNotFound("session tidak ditemukan atau akses ditolak")
	}
	if _, err := u.assetRepo.FindByIDAndTenant(input.AssetID, input.TenantID); err != nil {
		return nil, util.ErrNotFound("asset tidak ditemukan atau akses ditolak")
	}

	condition := strings.TrimSpace(input.Condition)
	if condition == "" {
		condition = domain.StockOpnameConditionMatch
	}
	if !domain.IsValidStockOpnameCondition(condition) {
		return nil, util.ErrValidation("condition tidak valid")
	}

	item := &domain.StockOpnameItem{
		TenantID:  input.TenantID,
		SessionID: input.SessionID,
		AssetID:   input.AssetID,
		Condition: condition,
		Note:      strings.TrimSpace(input.Note),
		CheckedBy: input.CheckedBy,
	}

	if err := u.itemRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *stockOpnameUsecase) ListItems(tenantID uuid.UUID, sessionID uint) ([]domain.StockOpnameItem, error) {
	if _, err := u.repo.FindByID(tenantID, sessionID); err != nil {
		return nil, err
	}
	return u.itemRepo.ListBySession(tenantID, sessionID)
}
