package usecase

import (
	"errors"
	"testing"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestPublicAssetUsecase_GetPublicAsset(t *testing.T) {
	tenantID := uuid.New()
	roomID := uuid.New()
	purchaseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	tenant := &domain.Tenant{ID: tenantID, Name: "RS Utama", Slug: "rs-utama"}

	fullAsset := &domain.Asset{
		ID:           uuid.New(),
		Code:         "VENT-001",
		Name:         "Ventilator",
		Brand:        "Philips",
		Model:        "V60",
		Status:       "ready",
		PurchaseDate: &purchaseDate,
		Category:     domain.Category{ID: 1, Name: "Life Support", Slug: "life-support"},
		Room:         domain.Room{ID: roomID, Name: "ICU", Code: "ICU-1", Location: "L1"},
		Bed:          domain.Bed{ID: 2, Code: "BED-01"},
	}

	minAsset := &domain.Asset{
		ID:     uuid.New(),
		Code:   "ASSET-01",
		Name:   "Asset Minimal",
		Status: "ready",
	}

	tests := []struct {
		name       string
		tenantSlug string
		assetCode  string
		setup      func(tr *mocks.MockTenantRepository, ar *mocks.MockAssetRepository)
		wantErr    string
		assert     func(t *testing.T, resp *PublicAssetResponse)
	}{
		{
			name:       "missing tenant slug",
			tenantSlug: " ",
			assetCode:  "VENT-001",
			wantErr:    "tenant slug dan kode aset wajib diisi",
		},
		{
			name:       "missing asset code",
			tenantSlug: "rs-utama",
			assetCode:  "",
			wantErr:    "tenant slug dan kode aset wajib diisi",
		},
		{
			name:       "tenant not found",
			tenantSlug: "rs-utama",
			assetCode:  "VENT-001",
			setup: func(tr *mocks.MockTenantRepository, ar *mocks.MockAssetRepository) {
				tr.EXPECT().FindBySlug("rs-utama").Return(nil, errors.New("not found"))
			},
			wantErr: "tenant tidak ditemukan",
		},
		{
			name:       "asset not found",
			tenantSlug: "rs-utama",
			assetCode:  "VENT-001",
			setup: func(tr *mocks.MockTenantRepository, ar *mocks.MockAssetRepository) {
				tr.EXPECT().FindBySlug("rs-utama").Return(tenant, nil)
				ar.EXPECT().FindByCodeAndTenant("VENT-001", tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "aset tidak ditemukan",
		},
		{
			name:       "success with full relations",
			tenantSlug: " rs-utama ",
			assetCode:  " VENT-001 ",
			setup: func(tr *mocks.MockTenantRepository, ar *mocks.MockAssetRepository) {
				tr.EXPECT().FindBySlug("rs-utama").Return(tenant, nil)
				ar.EXPECT().FindByCodeAndTenant("VENT-001", tenantID).Return(fullAsset, nil)
			},
			assert: func(t *testing.T, resp *PublicAssetResponse) {
				if resp == nil {
					t.Fatalf("expected response, got nil")
				}
				if resp.ID != fullAsset.ID {
					t.Fatalf("expected id %s, got %s", fullAsset.ID, resp.ID)
				}
				if resp.Category == nil || resp.Category.ID != 1 {
					t.Fatalf("expected category to be set")
				}
				if resp.Room == nil || resp.Room.ID != roomID {
					t.Fatalf("expected room to be set")
				}
				if resp.Bed == nil || resp.Bed.ID != 2 {
					t.Fatalf("expected bed to be set")
				}
				if resp.Tenant.Slug != "rs-utama" {
					t.Fatalf("expected tenant slug rs-utama, got %s", resp.Tenant.Slug)
				}
			},
		},
		{
			name:       "success with minimal relations",
			tenantSlug: "rs-utama",
			assetCode:  "ASSET-01",
			setup: func(tr *mocks.MockTenantRepository, ar *mocks.MockAssetRepository) {
				tr.EXPECT().FindBySlug("rs-utama").Return(tenant, nil)
				ar.EXPECT().FindByCodeAndTenant("ASSET-01", tenantID).Return(minAsset, nil)
			},
			assert: func(t *testing.T, resp *PublicAssetResponse) {
				if resp == nil {
					t.Fatalf("expected response, got nil")
				}
				if resp.Category != nil || resp.Room != nil || resp.Bed != nil {
					t.Fatalf("expected no category/room/bed on minimal asset")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			tenantRepo := mocks.NewMockTenantRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(tenantRepo, assetRepo)
			}

			usecase := NewPublicAssetUsecase(tenantRepo, assetRepo)
			resp, err := usecase.GetPublicAsset(tt.tenantSlug, tt.assetCode)
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
			if tt.assert != nil {
				tt.assert(t, resp)
			}
		})
	}
}
