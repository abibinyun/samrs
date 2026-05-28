package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestBedUsecase_CreateBed(t *testing.T) {
	tenantID := uuid.New()
	roomID := uuid.New()

	tests := []struct {
		name    string
		code    string
		setup   func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository)
		wantErr string
	}{
		{
			name: "room not found",
			code: "BED-01",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "ruangan tidak ditemukan atau akses ditolak",
		},
		{
			name: "code exists",
			code: "BED-01",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", nil).Return(true, nil)
			},
			wantErr: "kode bed sudah digunakan",
		},
		{
			name: "create error",
			code: "BED-01",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", nil).Return(false, nil)
				bedRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			code: "BED-01",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", nil).Return(false, nil)
				bedRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(bed *domain.Bed) error {
					if bed.RoomID != roomID {
						t.Fatalf("expected room %v, got %v", roomID, bed.RoomID)
					}
					if bed.Code != "BED-01" {
						t.Fatalf("expected code BED-01, got %s", bed.Code)
					}
					if bed.Status != domain.BedStatusAvailable {
						t.Fatalf("expected status available, got %s", bed.Status)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			bedRepo := mocks.NewMockBedRepository(ctrl)
			roomRepo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(bedRepo, roomRepo)
			}
			usecase := NewBedUsecase(bedRepo, roomRepo)
			_, err := usecase.CreateBed(tenantID, roomID, tt.code)
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

func TestBedUsecase_UpdateBed(t *testing.T) {
	tenantID := uuid.New()
	roomID := uuid.New()
	bedID := uint(1)

	tests := []struct {
		name    string
		code    string
		status  string
		setup   func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository)
		wantErr string
	}{
		{
			name:   "missing code",
			code:   " ",
			status: "available",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "kode bed wajib diisi",
		},
		{
			name:   "room not found",
			code:   "BED-01",
			status: "available",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "ruangan tidak ditemukan atau akses ditolak",
		},
		{
			name:   "bed not found",
			code:   "BED-01",
			status: "available",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:   "code exists",
			code:   "BED-01",
			status: "available",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID, RoomID: roomID, Code: "OLD"}
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", &bedID).Return(true, nil)
			},
			wantErr: "kode bed sudah digunakan",
		},
		{
			name:   "invalid status",
			code:   "BED-01",
			status: "invalid",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID, RoomID: roomID, Code: "OLD"}
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", &bedID).Return(false, nil)
			},
			wantErr: "status bed tidak valid",
		},
		{
			name:   "update error",
			code:   "BED-01",
			status: "available",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID, RoomID: roomID, Code: "OLD"}
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", &bedID).Return(false, nil)
				bedRepo.EXPECT().Update(bed).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name:   "success",
			code:   "BED-01",
			status: "maintenance",
			setup: func(bedRepo *mocks.MockBedRepository, roomRepo *mocks.MockRoomRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID, RoomID: roomID, Code: "OLD"}
				roomRepo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(&domain.Room{ID: roomID, TenantID: tenantID}, nil)
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().ExistsByCode(tenantID, roomID, "BED-01", &bedID).Return(false, nil)
				bedRepo.EXPECT().Update(bed).DoAndReturn(func(updated *domain.Bed) error {
					if updated.Code != "BED-01" || updated.Status != domain.BedStatusMaintenance {
						t.Fatalf("unexpected bed update: %+v", updated)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			bedRepo := mocks.NewMockBedRepository(ctrl)
			roomRepo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(bedRepo, roomRepo)
			}
			usecase := NewBedUsecase(bedRepo, roomRepo)
			_, err := usecase.UpdateBed(tenantID, bedID, roomID, tt.code, tt.status)
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

func TestBedUsecase_DeleteBed(t *testing.T) {
	tenantID := uuid.New()
	bedID := uint(1)

	tests := []struct {
		name    string
		setup   func(bedRepo *mocks.MockBedRepository)
		wantErr string
	}{
		{
			name: "bed not found",
			setup: func(bedRepo *mocks.MockBedRepository) {
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(bedRepo *mocks.MockBedRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID}
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().Delete(bed).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(bedRepo *mocks.MockBedRepository) {
				bed := &domain.Bed{ID: bedID, TenantID: tenantID}
				bedRepo.EXPECT().FindByIDAndTenant(bedID, tenantID).Return(bed, nil)
				bedRepo.EXPECT().Delete(bed).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			bedRepo := mocks.NewMockBedRepository(ctrl)
			roomRepo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(bedRepo)
			}
			usecase := NewBedUsecase(bedRepo, roomRepo)
			err := usecase.DeleteBed(tenantID, bedID)
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
