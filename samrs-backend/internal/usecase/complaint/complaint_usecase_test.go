package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestComplaintUsecase_CreateComplaint(t *testing.T) {
	tenantID := uuid.New()
	assetID := uuid.New()
	reporterID := uuid.New()
	assigneeID := uuid.New()

	tests := []struct {
		name    string
		input   CreateComplaintInput
		setup   func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name: "missing title",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				Title:       " ",
				Description: "desc",
			},
			wantErr: "judul complaint wajib diisi",
		},
		{
			name: "missing description",
			input: CreateComplaintInput{
				TenantID:   tenantID,
				AssetID:    assetID,
				ReportedBy: reporterID,
				Title:      "title",
			},
			wantErr: "deskripsi complaint wajib diisi",
		},
		{
			name: "asset not found",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				Title:       "title",
				Description: "desc",
			},
			setup: func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "aset tidak ditemukan atau akses ditolak",
		},
		{
			name: "reporter not found",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				Title:       "title",
				Description: "desc",
			},
			setup: func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(reporterID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "user pelapor tidak ditemukan",
		},
		{
			name: "assignee not found",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				AssignedTo:  &assigneeID,
				Title:       "title",
				Description: "desc",
			},
			setup: func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(reporterID, tenantID).Return(&domain.User{ID: reporterID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(assigneeID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "user teknisi tidak ditemukan",
		},
		{
			name: "create error",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				Title:       "title",
				Description: "desc",
			},
			setup: func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(reporterID, tenantID).Return(&domain.User{ID: reporterID}, nil)
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: CreateComplaintInput{
				TenantID:    tenantID,
				AssetID:     assetID,
				ReportedBy:  reporterID,
				Title:       "title",
				Description: "desc",
				AssignedTo:  &assigneeID,
			},
			setup: func(repo *mocks.MockComplaintRepository, assetRepo *mocks.MockAssetRepository, userRepo *mocks.MockUserRepository) {
				assetRepo.EXPECT().FindByIDAndTenant(assetID, tenantID).Return(&domain.Asset{ID: assetID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(reporterID, tenantID).Return(&domain.User{ID: reporterID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(assigneeID, tenantID).Return(&domain.User{ID: assigneeID}, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(c *domain.Complaint) error {
					if c.Status != domain.ComplaintStatusOpen {
						t.Fatalf("expected status open, got %s", c.Status)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockComplaintRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, assetRepo, userRepo)
			}
			usecase := NewComplaintUsecase(repo, assetRepo, userRepo)
			_, err := usecase.CreateComplaint(tt.input)
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

func TestComplaintUsecase_UpdateComplaint(t *testing.T) {
	tenantID := uuid.New()
	complaintID := uuid.New()
	assigneeID := uuid.New()

	tests := []struct {
		name    string
		input   UpdateComplaintInput
		setup   func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name: "complaint not found",
			input: UpdateComplaintInput{
				TenantID: tenantID,
				ID:       complaintID,
			},
			setup: func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository) {
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "invalid status",
			input: UpdateComplaintInput{
				TenantID: tenantID,
				ID:       complaintID,
				Status:   "invalid",
			},
			setup: func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository) {
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).
					Return(&domain.Complaint{ID: complaintID}, nil)
			},
			wantErr: "status complaint tidak valid",
		},
		{
			name: "assignee not found",
			input: UpdateComplaintInput{
				TenantID:      tenantID,
				ID:            complaintID,
				AssignedTo:    &assigneeID,
				AssignedToSet: true,
			},
			setup: func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository) {
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).
					Return(&domain.Complaint{ID: complaintID}, nil)
				userRepo.EXPECT().FindByIDAndTenant(assigneeID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "user teknisi tidak ditemukan",
		},
		{
			name: "update error",
			input: UpdateComplaintInput{
				TenantID: tenantID,
				ID:       complaintID,
				Title:    "updated",
			},
			setup: func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository) {
				complaint := &domain.Complaint{ID: complaintID}
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(complaint, nil)
				repo.EXPECT().Update(complaint).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: UpdateComplaintInput{
				TenantID:       tenantID,
				ID:             complaintID,
				Title:          "updated",
				Description:    "desc",
				Status:         domain.ComplaintStatusInProgress,
				ResolutionNote: "note",
			},
			setup: func(repo *mocks.MockComplaintRepository, userRepo *mocks.MockUserRepository) {
				complaint := &domain.Complaint{ID: complaintID}
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(complaint, nil)
				repo.EXPECT().Update(complaint).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockComplaintRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, userRepo)
			}
			usecase := NewComplaintUsecase(repo, assetRepo, userRepo)
			_, err := usecase.UpdateComplaint(tt.input)
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

func TestComplaintUsecase_DeleteComplaint(t *testing.T) {
	tenantID := uuid.New()
	complaintID := uuid.New()

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockComplaintRepository)
		wantErr string
	}{
		{
			name: "not found",
			setup: func(repo *mocks.MockComplaintRepository) {
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(repo *mocks.MockComplaintRepository) {
				complaint := &domain.Complaint{ID: complaintID, TenantID: tenantID}
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(complaint, nil)
				repo.EXPECT().Delete(complaint).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockComplaintRepository) {
				complaint := &domain.Complaint{ID: complaintID, TenantID: tenantID}
				repo.EXPECT().FindByIDAndTenant(complaintID, tenantID).Return(complaint, nil)
				repo.EXPECT().Delete(complaint).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockComplaintRepository(ctrl)
			assetRepo := mocks.NewMockAssetRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewComplaintUsecase(repo, assetRepo, userRepo)
			err := usecase.DeleteComplaint(tenantID, complaintID)
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
