package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestAssetEventUsecase_CreateEvent(t *testing.T) {
	tenantID := uuid.New()
	assetID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		input   AssetEventInput
		setup   func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name: "missing asset id",
			input: AssetEventInput{
				TenantID:  tenantID,
				UserID:    userID,
				EventType: "created",
			},
			wantErr: "asset ID wajib diisi",
		},
		{
			name: "missing user id",
			input: AssetEventInput{
				TenantID:  tenantID,
				AssetID:   assetID,
				EventType: "created",
			},
			wantErr: "user ID wajib diisi",
		},
		{
			name: "missing event type",
			input: AssetEventInput{
				TenantID: tenantID,
				AssetID:  assetID,
				UserID:   userID,
			},
			wantErr: "event_type wajib diisi",
		},
		{
			name: "asset not found",
			input: AssetEventInput{
				TenantID:  tenantID,
				AssetID:   assetID,
				UserID:    userID,
				EventType: "created",
			},
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "asset tidak ditemukan atau akses ditolak",
		},
		{
			name: "user not found",
			input: AssetEventInput{
				TenantID:  tenantID,
				AssetID:   assetID,
				UserID:    userID,
				EventType: "created",
			},
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "user tidak ditemukan",
		},
		{
			name: "invalid meta payload",
			input: AssetEventInput{
				TenantID:  tenantID,
				AssetID:   assetID,
				UserID:    userID,
				EventType: "created",
				Meta: map[string]interface{}{
					"bad": func() {},
				},
			},
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(&domain.User{ID: userID}, nil)
			},
			wantErr: "json: unsupported type: func()",
		},
		{
			name: "create error",
			input: AssetEventInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				UserID:      userID,
				EventType:   "created",
				Description: "desc",
			},
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(&domain.User{ID: userID}, nil)
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: AssetEventInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				UserID:      userID,
				EventType:   "created",
				Description: "desc",
				Meta: map[string]interface{}{
					"key": "value",
				},
			},
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(&domain.User{ID: userID}, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(event *domain.AssetEvent) error {
					if event.EventType != "created" {
						t.Fatalf("expected event_type created, got %s", event.EventType)
					}
					if event.Meta == nil {
						t.Fatalf("expected meta to be set")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eventRepo := mocks.NewMockAssetEventRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			if tt.setup != nil {
				tt.setup(eventRepo, assetRepo, userRepo)
			}
			usecase := NewAssetEventUsecase(eventRepo, assetRepo, userRepo)
			_, err := usecase.CreateEvent(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAssetEventUsecase_GetTimeline(t *testing.T) {
	tenantID := uuid.New()
	assetID := uuid.New()
	filter := repository.AssetEventFilter{Page: 1, PerPage: 10}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository)
		wantErr string
	}{
		{
			name: "asset not found",
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "list error",
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				repo.EXPECT().ListByAsset(tenantID, assetID, filter).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockAssetEventRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				repo.EXPECT().ListByAsset(tenantID, assetID, filter).Return([]domain.AssetEvent{{ID: 1}}, int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			eventRepo := mocks.NewMockAssetEventRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			if tt.setup != nil {
				tt.setup(eventRepo, assetRepo)
			}
			usecase := NewAssetEventUsecase(eventRepo, assetRepo, userRepo)
			_, _, err := usecase.GetTimeline(tenantID, assetID, filter)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
