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

func TestReportUsecase_ExportAssets(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.AssetFilter{Status: "ready"}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockAssetRepository)
		wantErr string
	}{
		{
			name: "repo error",
			setup: func(repo *mocks.MockAssetRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return(nil, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockAssetRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return([]domain.Asset{{}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			complaintRepo := mocks.NewMockComplaintRepository(ctrl)
			maintenanceRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(assetRepo)
			}
			usecase := NewReportUsecase(assetRepo, complaintRepo, maintenanceRepo)
			_, err := usecase.ExportAssets(tenantID, filter)
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

func TestReportUsecase_ExportComplaints(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.ComplaintFilter{Status: "open"}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockComplaintRepository)
		wantErr string
	}{
		{
			name: "repo error",
			setup: func(repo *mocks.MockComplaintRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return(nil, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockComplaintRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return([]domain.Complaint{{}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			complaintRepo := mocks.NewMockComplaintRepository(ctrl)
			maintenanceRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(complaintRepo)
			}
			usecase := NewReportUsecase(assetRepo, complaintRepo, maintenanceRepo)
			_, err := usecase.ExportComplaints(tenantID, filter)
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

func TestReportUsecase_ExportMaintenance(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.MaintenanceScheduleFilter{Status: "scheduled"}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockMaintenanceScheduleRepository)
		wantErr string
	}{
		{
			name: "repo error",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return(nil, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockMaintenanceScheduleRepository) {
				repo.EXPECT().FindAllByTenantExport(tenantID, filter).Return([]domain.MaintenanceSchedule{{}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			complaintRepo := mocks.NewMockComplaintRepository(ctrl)
			maintenanceRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(maintenanceRepo)
			}
			usecase := NewReportUsecase(assetRepo, complaintRepo, maintenanceRepo)
			_, err := usecase.ExportMaintenance(tenantID, filter)
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
