package usecase

import (
	"encoding/json"
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AssetEventUsecase interface {
	CreateEvent(input AssetEventInput) (*domain.AssetEvent, error)
	GetTimeline(tenantID uuid.UUID, assetID uuid.UUID, filter repository.AssetEventFilter) ([]domain.AssetEvent, int64, error)
}

type AssetEventInput struct {
	TenantID    uuid.UUID
	AssetID     uuid.UUID
	UserID      uuid.UUID
	EventType   string
	Description string
	Meta        map[string]interface{}
}

type assetEventUsecase struct {
	repo      repository.AssetEventRepository
	assetRepo repository.AssetRepository
	userRepo  repository.UserRepository
}

func NewAssetEventUsecase(er repository.AssetEventRepository, ar repository.AssetRepository, ur repository.UserRepository) AssetEventUsecase {
	return &assetEventUsecase{
		repo:      er,
		assetRepo: ar,
		userRepo:  ur,
	}
}

func (u *assetEventUsecase) CreateEvent(input AssetEventInput) (*domain.AssetEvent, error) {
	if input.AssetID == uuid.Nil {
		return nil, util.ErrValidation("asset ID wajib diisi")
	}
	if input.UserID == uuid.Nil {
		return nil, util.ErrValidation("user ID wajib diisi")
	}
	eventType := strings.TrimSpace(input.EventType)
	if eventType == "" {
		return nil, util.ErrValidation("event_type wajib diisi")
	}
	if _, err := u.assetRepo.FindByIDAndTenant(input.AssetID, input.TenantID); err != nil {
		return nil, util.ErrNotFound("asset tidak ditemukan atau akses ditolak")
	}
	if _, err := u.userRepo.FindByIDAndTenant(input.UserID, input.TenantID); err != nil {
		return nil, util.ErrNotFound("user tidak ditemukan")
	}

	metaJSON, err := toJSONMap(input.Meta)
	if err != nil {
		return nil, err
	}

	event := &domain.AssetEvent{
		TenantID:    input.TenantID,
		AssetID:     input.AssetID,
		UserID:      input.UserID,
		EventType:   eventType,
		Description: strings.TrimSpace(input.Description),
		Meta:        metaJSON,
	}

	if err := u.repo.Create(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (u *assetEventUsecase) GetTimeline(tenantID uuid.UUID, assetID uuid.UUID, filter repository.AssetEventFilter) ([]domain.AssetEvent, int64, error) {
	if _, err := u.assetRepo.FindByIDAndTenant(assetID, tenantID); err != nil {
		return nil, 0, err
	}
	return u.repo.ListByAsset(tenantID, assetID, filter)
}

func toJSONMap(payload map[string]interface{}) (datatypes.JSON, error) {
	if payload == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(bytes), nil
}
