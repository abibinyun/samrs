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

type assetMocks struct {
	assetRepo    *mocks.MockAssetRepository
	categoryRepo *mocks.MockCategoryRepository
	roomRepo     *mocks.MockRoomRepository
	bedRepo      *mocks.MockBedRepository
	vendorRepo   *mocks.MockVendorRepository
	brandRepo    *mocks.MockAssetBrandRepository
	modelRepo    *mocks.MockAssetModelRepository
	statusRepo   *mocks.MockAssetStatusRepository
}

func newAssetUsecase(ctrl *gomock.Controller) (AssetUsecase, assetMocks) {
	m := assetMocks{
		assetRepo:    mocks.NewMockAssetRepository(ctrl),
		categoryRepo: mocks.NewMockCategoryRepository(ctrl),
		roomRepo:     mocks.NewMockRoomRepository(ctrl),
		bedRepo:      mocks.NewMockBedRepository(ctrl),
		vendorRepo:   mocks.NewMockVendorRepository(ctrl),
		brandRepo:    mocks.NewMockAssetBrandRepository(ctrl),
		modelRepo:    mocks.NewMockAssetModelRepository(ctrl),
		statusRepo:   mocks.NewMockAssetStatusRepository(ctrl),
	}
	return NewAssetUsecase(
		m.assetRepo,
		m.categoryRepo,
		m.roomRepo,
		m.bedRepo,
		m.vendorRepo,
		m.brandRepo,
		m.modelRepo,
		m.statusRepo,
	), m
}

func TestAssetUsecase_CreateAsset(t *testing.T) {
	tenantID := uuid.New()
	categoryID := uint(1)
	roomID := uuid.New()
	bedID := uint(2)
	vendorID := uint(3)
	brandID := uint(4)
	modelID := uint(5)

	baseInput := func() CreateAssetInput {
		return CreateAssetInput{
			TenantID:     tenantID,
			CategoryID:   categoryID,
			RoomID:       &roomID,
			BedID:        &bedID,
			VendorID:     &vendorID,
			BrandID:      &brandID,
			ModelID:      &modelID,
			Code:         "AST-1",
			Name:         "Asset One",
			Status:       "",
			PurchaseDate: "2024-01-15",
		}
	}

	tests := []struct {
		name    string
		input   CreateAssetInput
		setup   func(m assetMocks)
		wantErr string
	}{
		{
			name: "missing code",
			input: CreateAssetInput{
				TenantID: tenantID,
				Code:     " ",
				Name:     "Asset",
			},
			wantErr: "kode dan nama aset wajib diisi",
		},
		{
			name: "missing category",
			input: CreateAssetInput{
				TenantID: tenantID,
				Code:     "AST-1",
				Name:     "Asset",
			},
			wantErr: "kategori wajib diisi",
		},
		{
			name:  "category not found",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(nil, errors.New("not found"))
			},
			wantErr: "kategori tidak ditemukan atau akses ditolak",
		},
		{
			name:  "room not found",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "ruangan tidak ditemukan atau akses ditolak",
		},
		{
			name: "bed not found",
			input: func() CreateAssetInput {
				in := baseInput()
				in.RoomID = nil
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "bed tidak ditemukan atau akses ditolak",
		},
		{
			name:  "bed room mismatch",
			input: baseInput(),
			setup: func(m assetMocks) {
				bedRoom := uuid.New()
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(&domain.Bed{ID: bedID, RoomID: bedRoom}, nil)
			},
			wantErr: "bed tidak sesuai dengan ruangan",
		},
		{
			name:  "vendor not found",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(&domain.Bed{ID: bedID, RoomID: roomID}, nil)
				m.vendorRepo.EXPECT().FindByID(tenantID, vendorID).Return(nil, errors.New("not found"))
			},
			wantErr: "vendor tidak ditemukan atau akses ditolak",
		},
		{
			name:  "model brand mismatch",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(&domain.Bed{ID: bedID, RoomID: roomID}, nil)
				m.vendorRepo.EXPECT().FindByID(tenantID, vendorID).Return(&domain.Vendor{ID: vendorID}, nil)
				m.brandRepo.EXPECT().FindByID(tenantID, brandID).Return(&domain.AssetBrand{ID: brandID}, nil)
				m.modelRepo.EXPECT().FindByID(tenantID, modelID).Return(&domain.AssetModel{ID: modelID, BrandID: 999}, nil)
			},
			wantErr: "tipe alat tidak sesuai dengan merek",
		},
		{
			name:  "status invalid without master",
			input: func() CreateAssetInput {
				in := baseInput()
				in.Status = "invalid"
				in.VendorID = nil
				in.BrandID = nil
				in.ModelID = nil
				in.BedID = nil
				in.RoomID = nil
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
			},
			wantErr: "status aset tidak valid",
		},
		{
			name:  "status not registered",
			input: func() CreateAssetInput {
				in := baseInput()
				in.Status = "ready"
				in.VendorID = nil
				in.BrandID = nil
				in.ModelID = nil
				in.BedID = nil
				in.RoomID = nil
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(true, nil)
				m.statusRepo.EXPECT().ExistsByCode(tenantID, "ready", nil).Return(false, nil)
			},
			wantErr: "status aset tidak terdaftar",
		},
		{
			name:  "invalid purchase date",
			input: func() CreateAssetInput {
				in := baseInput()
				in.VendorID = nil
				in.BrandID = nil
				in.ModelID = nil
				in.BedID = nil
				in.RoomID = nil
				in.PurchaseDate = "2024-13-01"
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
			},
			wantErr: "format purchase_date harus YYYY-MM-DD",
		},
		{
			name:  "code exists",
			input: func() CreateAssetInput {
				in := baseInput()
				in.VendorID = nil
				in.BrandID = nil
				in.ModelID = nil
				in.BedID = nil
				in.RoomID = nil
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", nil).Return(true, nil)
			},
			wantErr: "kode aset sudah digunakan",
		},
		{
			name: "success with bed sets room",
			input: func() CreateAssetInput {
				in := baseInput()
				in.RoomID = nil
				in.VendorID = nil
				in.BrandID = nil
				in.ModelID = nil
				in.Status = ""
				in.PurchaseDate = ""
				return in
			}(),
			setup: func(m assetMocks) {
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(&domain.Bed{ID: bedID, RoomID: roomID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", nil).Return(false, nil)
				m.assetRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(asset *domain.Asset) error {
					if asset.RoomID == nil || *asset.RoomID != roomID {
						t.Fatalf("expected room to be set from bed")
					}
					if asset.Status != domain.AssetStatusReady {
						t.Fatalf("expected default status ready, got %s", asset.Status)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase, m := newAssetUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(m)
			}
			_, err := usecase.CreateAsset(tt.input)
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

func TestAssetUsecase_UpdateAsset(t *testing.T) {
	tenantID := uuid.New()
	categoryID := uint(1)
	roomID := uuid.New()
	bedID := uint(2)
	assetID := uuid.New()

	baseInput := func() UpdateAssetInput {
		return UpdateAssetInput{
			TenantID:   tenantID,
			ID:         assetID,
			CategoryID: categoryID,
			Code:       "AST-1",
			Name:       "Asset One",
		}
	}

	tests := []struct {
		name    string
		input   UpdateAssetInput
		setup   func(m assetMocks)
		wantErr string
	}{
		{
			name: "missing code",
			input: UpdateAssetInput{
				TenantID: tenantID,
				ID:       assetID,
				Code:     " ",
				Name:     "Asset",
			},
			wantErr: "kode dan nama aset wajib diisi",
		},
		{
			name:  "asset not found",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "category missing",
			input: func() UpdateAssetInput {
				in := baseInput()
				in.CategoryID = 0
				return in
			}(),
			setup: func(m assetMocks) {
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "kategori wajib diisi",
		},
		{
			name:  "category not found",
			input: baseInput(),
			setup: func(m assetMocks) {
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(nil, errors.New("not found"))
			},
			wantErr: "kategori tidak ditemukan atau akses ditolak",
		},
		{
			name: "bed mismatch",
			input: func() UpdateAssetInput {
				in := baseInput()
				in.RoomID = &roomID
				in.BedID = &bedID
				return in
			}(),
			setup: func(m assetMocks) {
				bedRoom := uuid.New()
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID}, nil)
				m.bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(&domain.Bed{ID: bedID, RoomID: bedRoom}, nil)
			},
			wantErr: "bed tidak sesuai dengan ruangan",
		},
		{
			name: "status invalid",
			input: func() UpdateAssetInput {
				in := baseInput()
				in.Status = "invalid"
				return in
			}(),
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", &assetID).Return(false, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
			},
			wantErr: "status aset tidak valid",
		},
		{
			name:  "code exists",
			input: baseInput(),
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", &assetID).Return(true, nil)
			},
			wantErr: "kode aset sudah digunakan",
		},
		{
			name: "invalid purchase date",
			input: func() UpdateAssetInput {
				in := baseInput()
				in.PurchaseDate = "2024-13-01"
				return in
			}(),
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", &assetID).Return(false, nil)
			},
			wantErr: "format purchase_date harus YYYY-MM-DD",
		},
		{
			name:  "update error",
			input: baseInput(),
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", &assetID).Return(false, nil)
				m.assetRepo.EXPECT().Update(asset).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: func() UpdateAssetInput {
				in := baseInput()
				in.Status = "ready"
				in.PurchaseDate = "2024-01-15"
				return in
			}(),
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.categoryRepo.EXPECT().FindByID(tenantID, categoryID).Return(&domain.Category{ID: categoryID}, nil)
				m.statusRepo.EXPECT().HasAny(tenantID).Return(false, nil)
				m.assetRepo.EXPECT().ExistsByCode(tenantID, "AST-1", &assetID).Return(false, nil)
				m.assetRepo.EXPECT().Update(asset).DoAndReturn(func(updated *domain.Asset) error {
					if updated.Status != "ready" {
						t.Fatalf("expected status ready, got %s", updated.Status)
					}
					if updated.PurchaseDate == nil || !updated.PurchaseDate.Equal(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)) {
						t.Fatalf("unexpected purchase date: %v", updated.PurchaseDate)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase, m := newAssetUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(m)
			}
			_, err := usecase.UpdateAsset(tt.input)
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

func TestAssetUsecase_DeleteAsset(t *testing.T) {
	tenantID := uuid.New()
	assetID := uuid.New()

	tests := []struct {
		name    string
		setup   func(m assetMocks)
		wantErr string
	}{
		{
			name: "asset not found",
			setup: func(m assetMocks) {
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.assetRepo.EXPECT().Delete(asset).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(m assetMocks) {
				asset := &domain.Asset{ID: assetID, TenantID: tenantID}
				m.assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(asset, nil)
				m.assetRepo.EXPECT().Delete(asset).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase, m := newAssetUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(m)
			}
			err := usecase.DeleteAsset(tenantID, assetID)
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
