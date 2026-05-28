package usecase

import (
	"errors"
	"testing"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestStockOpnameUsecase_CreateSession(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		input   StockOpnameSessionInput
		setup   func(repo *mocks.MockStockOpnameRepository)
		wantErr string
	}{
		{
			name: "missing title",
			input: StockOpnameSessionInput{
				TenantID: tenantID,
				Title:    " ",
			},
			wantErr: "judul opname wajib diisi",
		},
		{
			name: "invalid opname date",
			input: StockOpnameSessionInput{
				TenantID: tenantID,
				Title:    "Opname",
				OpnameAt: "2026-13-01",
			},
			wantErr: "format opname_at harus YYYY-MM-DD",
		},
		{
			name: "create error",
			input: StockOpnameSessionInput{
				TenantID: tenantID,
				Title:    "Opname",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success default status",
			input: StockOpnameSessionInput{
				TenantID: tenantID,
				Title:    "Opname",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(session *domain.StockOpnameSession) error {
					if session.Status != domain.StockOpnameStatusDraft {
						t.Fatalf("expected default status draft, got %s", session.Status)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, err := usecase.CreateSession(tt.input)
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

func TestStockOpnameUsecase_UpdateSession(t *testing.T) {
	tenantID := uuid.New()
	sessionID := uint(1)

	tests := []struct {
		name    string
		input   StockOpnameSessionUpdateInput
		setup   func(repo *mocks.MockStockOpnameRepository)
		wantErr string
	}{
		{
			name: "session not found",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "missing title",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
				Title:    " ",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
			},
			wantErr: "judul opname wajib diisi",
		},
		{
			name: "invalid opname date",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
				Title:    "Opname",
				OpnameAt: "2026-13-01",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID, OpnameAt: time.Now()}, nil)
			},
			wantErr: "format opname_at harus YYYY-MM-DD",
		},
		{
			name: "invalid status",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
				Title:    "Opname",
				Status:   "invalid",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID, Status: domain.StockOpnameStatusDraft, OpnameAt: time.Now()}, nil)
			},
			wantErr: "status opname tidak valid",
		},
		{
			name: "update error",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
				Title:    "Opname",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				session := &domain.StockOpnameSession{ID: sessionID, Status: domain.StockOpnameStatusDraft, OpnameAt: time.Now()}
				repo.EXPECT().FindByID(tenantID, sessionID).Return(session, nil)
				repo.EXPECT().Update(session).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: StockOpnameSessionUpdateInput{
				TenantID: tenantID,
				ID:       sessionID,
				Title:    "Opname",
			},
			setup: func(repo *mocks.MockStockOpnameRepository) {
				session := &domain.StockOpnameSession{ID: sessionID, Status: domain.StockOpnameStatusDraft, OpnameAt: time.Now()}
				repo.EXPECT().FindByID(tenantID, sessionID).Return(session, nil)
				repo.EXPECT().Update(session).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, err := usecase.UpdateSession(tt.input)
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

func TestStockOpnameUsecase_CloseSession(t *testing.T) {
	tenantID := uuid.New()
	sessionID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockStockOpnameRepository)
		wantErr string
	}{
		{
			name: "session not found",
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "update error",
			setup: func(repo *mocks.MockStockOpnameRepository) {
				session := &domain.StockOpnameSession{ID: sessionID}
				repo.EXPECT().FindByID(tenantID, sessionID).Return(session, nil)
				repo.EXPECT().Update(session).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockStockOpnameRepository) {
				session := &domain.StockOpnameSession{ID: sessionID}
				repo.EXPECT().FindByID(tenantID, sessionID).Return(session, nil)
				repo.EXPECT().Update(session).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, err := usecase.CloseSession(tenantID, sessionID)
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

func TestStockOpnameUsecase_AddItem(t *testing.T) {
	tenantID := uuid.New()
	sessionID := uint(1)
	assetID := uuid.New()
	checkedBy := uuid.New()

	tests := []struct {
		name    string
		input   StockOpnameItemInput
		setup   func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository)
		wantErr string
	}{
		{
			name: "missing asset id",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				CheckedBy: checkedBy,
			},
			wantErr: "asset ID wajib diisi",
		},
		{
			name: "missing checked_by",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
			},
			wantErr: "checked_by wajib diisi",
		},
		{
			name: "session not found",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
				CheckedBy: checkedBy,
			},
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(nil, errors.New("not found"))
			},
			wantErr: "session tidak ditemukan atau akses ditolak",
		},
		{
			name: "asset not found",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
				CheckedBy: checkedBy,
			},
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "asset tidak ditemukan atau akses ditolak",
		},
		{
			name: "invalid condition",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
				CheckedBy: checkedBy,
				Condition: "invalid",
			},
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "condition tidak valid",
		},
		{
			name: "create error",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
				CheckedBy: checkedBy,
			},
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				itemRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success default condition",
			input: StockOpnameItemInput{
				TenantID:  tenantID,
				SessionID: sessionID,
				AssetID:   assetID,
				CheckedBy: checkedBy,
			},
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository, assetRepo *mocks.MockAssetRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				itemRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(item *domain.StockOpnameItem) error {
					if item.Condition != domain.StockOpnameConditionMatch {
						t.Fatalf("expected default condition match, got %s", item.Condition)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, itemRepo, assetRepo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, err := usecase.AddItem(tt.input)
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

func TestStockOpnameUsecase_ListItems(t *testing.T) {
	tenantID := uuid.New()
	sessionID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository)
		wantErr string
	}{
		{
			name: "session not found",
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockStockOpnameRepository, itemRepo *mocks.MockStockOpnameItemRepository) {
				repo.EXPECT().FindByID(tenantID, sessionID).Return(&domain.StockOpnameSession{ID: sessionID}, nil)
				itemRepo.EXPECT().ListBySession(tenantID, sessionID).Return([]domain.StockOpnameItem{{ID: 1}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, itemRepo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, err := usecase.ListItems(tenantID, sessionID)
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

func TestStockOpnameUsecase_GetAllSessions(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.StockOpnameFilter{Page: 1, PerPage: 10}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockStockOpnameRepository)
		wantErr string
	}{
		{
			name: "list error",
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindAllByTenant(tenantID, filter).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockStockOpnameRepository) {
				repo.EXPECT().FindAllByTenant(tenantID, filter).Return([]domain.StockOpnameSession{{ID: 1}}, int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockStockOpnameRepository(ctrl)
			itemRepo := mocks.NewMockStockOpnameItemRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewStockOpnameUsecase(repo, itemRepo, assetRepo)
			_, _, err := usecase.GetAllSessions(tenantID, filter)
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
