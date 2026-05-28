package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestMaintenanceScheduleUsecase_CreateSchedule(t *testing.T) {
	tenantID := uuid.New()
	assetID := uuid.New()

	baseInput := func() MaintenanceScheduleInput {
		return MaintenanceScheduleInput{
			TenantID:     tenantID,
			AssetID:      assetID,
			ScheduleType: domain.ScheduleTypeMaintenance,
			Title:        "Monthly",
			IntervalDays: 30,
			NextDueDate:  "2026-02-01",
			Status:       domain.ScheduleStatusScheduled,
		}
	}

	tests := []struct {
		name    string
		input   MaintenanceScheduleInput
		setup   func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository)
		wantErr string
	}{
		{
			name: "missing asset id",
			input: MaintenanceScheduleInput{
				TenantID: tenantID,
				Title:    "Monthly",
			},
			wantErr: "asset ID wajib diisi",
		},
		{
			name:  "asset not found",
			input: baseInput(),
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "asset tidak ditemukan atau akses ditolak",
		},
		{
			name: "missing title",
			input: MaintenanceScheduleInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				NextDueDate: "2026-02-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "judul jadwal wajib diisi",
		},
		{
			name: "invalid schedule type",
			input: MaintenanceScheduleInput{
				TenantID:     tenantID,
				AssetID:      assetID,
				ScheduleType: "invalid",
				Title:        "Monthly",
				NextDueDate:  "2026-02-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "tipe jadwal tidak valid",
		},
		{
			name: "missing next due date",
			input: MaintenanceScheduleInput{
				TenantID: tenantID,
				AssetID:  assetID,
				Title:    "Monthly",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "next_due_date wajib diisi",
		},
		{
			name: "invalid next due date format",
			input: MaintenanceScheduleInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				Title:       "Monthly",
				NextDueDate: "2026-13-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "format next_due_date harus YYYY-MM-DD",
		},
		{
			name: "invalid status",
			input: MaintenanceScheduleInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				Title:       "Monthly",
				NextDueDate: "2026-02-01",
				Status:      "invalid",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
			},
			wantErr: "status jadwal tidak valid",
		},
		{
			name:  "create error",
			input: baseInput(),
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success default type/status",
			input: MaintenanceScheduleInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				Title:       "Monthly",
				NextDueDate: "2026-02-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository, assetRepo *mocks.MockAssetRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(schedule *domain.MaintenanceSchedule) error {
					if schedule.ScheduleType != domain.ScheduleTypeMaintenance {
						t.Fatalf("expected default schedule type, got %s", schedule.ScheduleType)
					}
					if schedule.Status != domain.ScheduleStatusScheduled {
						t.Fatalf("expected default status scheduled, got %s", schedule.Status)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, assetRepo)
			}
			usecase := NewMaintenanceScheduleUsecase(repo, assetRepo)
			_, err := usecase.CreateSchedule(tt.input)
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

func TestMaintenanceScheduleUsecase_UpdateSchedule(t *testing.T) {
	tenantID := uuid.New()
	scheduleID := uint(1)

	tests := []struct {
		name    string
		input   MaintenanceScheduleUpdateInput
		setup   func(repo *mocks.MockMaintenanceScheduleRepository)
		wantErr string
	}{
		{
			name: "schedule not found",
			input: MaintenanceScheduleUpdateInput{
				TenantID: tenantID,
				ID:       scheduleID,
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "missing title",
			input: MaintenanceScheduleUpdateInput{
				TenantID: tenantID,
				ID:       scheduleID,
				Title:    " ",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
			},
			wantErr: "judul jadwal wajib diisi",
		},
		{
			name: "invalid schedule type",
			input: MaintenanceScheduleUpdateInput{
				TenantID:     tenantID,
				ID:           scheduleID,
				Title:        "Monthly",
				ScheduleType: "invalid",
				NextDueDate:  "2026-02-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
			},
			wantErr: "tipe jadwal tidak valid",
		},
		{
			name: "invalid next due date format",
			input: MaintenanceScheduleUpdateInput{
				TenantID:    tenantID,
				ID:          scheduleID,
				Title:       "Monthly",
				NextDueDate: "2026-13-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
			},
			wantErr: "format next_due_date harus YYYY-MM-DD",
		},
		{
			name: "missing next due date",
			input: MaintenanceScheduleUpdateInput{
				TenantID: tenantID,
				ID:       scheduleID,
				Title:    "Monthly",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
			},
			wantErr: "next_due_date wajib diisi",
		},
		{
			name: "invalid status",
			input: MaintenanceScheduleUpdateInput{
				TenantID:    tenantID,
				ID:          scheduleID,
				Title:       "Monthly",
				NextDueDate: "2026-02-01",
				Status:      "invalid",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
			},
			wantErr: "status jadwal tidak valid",
		},
		{
			name: "update error",
			input: MaintenanceScheduleUpdateInput{
				TenantID:    tenantID,
				ID:          scheduleID,
				Title:       "Monthly",
				NextDueDate: "2026-02-01",
				Status:      domain.ScheduleStatusScheduled,
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Update(schedule).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success default type/status",
			input: MaintenanceScheduleUpdateInput{
				TenantID:    tenantID,
				ID:          scheduleID,
				Title:       "Monthly",
				NextDueDate: "2026-02-01",
			},
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Update(schedule).DoAndReturn(func(updated *domain.MaintenanceSchedule) error {
					if updated.ScheduleType != domain.ScheduleTypeMaintenance {
						t.Fatalf("expected default schedule type")
					}
					if updated.Status != domain.ScheduleStatusScheduled {
						t.Fatalf("expected default status scheduled")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewMaintenanceScheduleUsecase(repo, assetRepo)
			_, err := usecase.UpdateSchedule(tt.input)
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

func TestMaintenanceScheduleUsecase_DeleteSchedule(t *testing.T) {
	tenantID := uuid.New()
	scheduleID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockMaintenanceScheduleRepository)
		wantErr string
	}{
		{
			name: "not found",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Delete(schedule).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Delete(schedule).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewMaintenanceScheduleUsecase(repo, assetRepo)
			err := usecase.DeleteSchedule(tenantID, scheduleID)
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

func TestMaintenanceScheduleUsecase_CompleteSchedule(t *testing.T) {
	tenantID := uuid.New()
	scheduleID := uint(1)

	tests := []struct {
		name        string
		interval    int
		setup       func(repo *mocks.MockMaintenanceScheduleRepository)
		wantStatus  string
		wantNextDue bool
		wantErr     string
	}{
		{
			name: "not found",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:     "interval schedule",
			interval: 30,
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID, IntervalDays: 30}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Update(schedule).Return(nil)
			},
			wantStatus:  domain.ScheduleStatusScheduled,
			wantNextDue: true,
		},
		{
			name:     "non interval schedule",
			interval: 0,
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID, IntervalDays: 0}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Update(schedule).Return(nil)
			},
			wantStatus:  domain.ScheduleStatusCompleted,
			wantNextDue: false,
		},
		{
			name:     "update error",
			interval: 0,
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				schedule := &domain.MaintenanceSchedule{ID: scheduleID, IntervalDays: 0}
				repo.EXPECT().FindByID(tenantID, scheduleID).Return(schedule, nil)
				repo.EXPECT().Update(schedule).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewMaintenanceScheduleUsecase(repo, assetRepo)
			result, err := usecase.CompleteSchedule(tenantID, scheduleID, "done")
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
			if result.Status != tt.wantStatus {
				t.Fatalf("expected status %s, got %s", tt.wantStatus, result.Status)
			}
			if tt.wantNextDue && result.NextDueDate == nil {
				t.Fatalf("expected next due date set")
			}
			if !tt.wantNextDue && result.NextDueDate != nil && result.IntervalDays == 0 {
				t.Fatalf("expected no next due date")
			}
		})
	}
}
