package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestMaintenanceDocumentUsecase_CreateDocument(t *testing.T) {
	tenantID := uuid.New()
	scheduleID := uint(1)
	assetID := uuid.New()
	uploaderID := uuid.New()

	baseInput := func() MaintenanceDocumentInput {
		return MaintenanceDocumentInput{
			TenantID:   tenantID,
			ScheduleID: scheduleID,
			UploadedBy: uploaderID,
			DocType:    domain.MaintenanceDocReport,
			Filename:   "report.pdf",
			FilePath:   "/tmp/report.pdf",
			MimeType:   "application/pdf",
			Size:       10,
		}
	}

	tests := []struct {
		name    string
		input   MaintenanceDocumentInput
		setup   func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository)
		wantErr string
	}{
		{
			name: "missing schedule id",
			input: MaintenanceDocumentInput{
				TenantID: tenantID,
			},
			wantErr: "schedule ID wajib diisi",
		},
		{
			name:  "missing uploaded by",
			input: MaintenanceDocumentInput{TenantID: tenantID, ScheduleID: scheduleID},
			wantErr: "uploaded_by wajib diisi",
		},
		{
			name: "missing file",
			input: MaintenanceDocumentInput{
				TenantID:   tenantID,
				ScheduleID: scheduleID,
				UploadedBy: uploaderID,
			},
			wantErr: "file wajib diisi",
		},
		{
			name: "invalid doc type",
			input: MaintenanceDocumentInput{
				TenantID:   tenantID,
				ScheduleID: scheduleID,
				UploadedBy: uploaderID,
				DocType:    "invalid",
				Filename:   "report.pdf",
				FilePath:   "/tmp/report.pdf",
			},
			wantErr: "doc_type tidak valid",
		},
		{
			name:  "schedule not found",
			input: baseInput(),
			setup: func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository) {
				scheduleRepo.EXPECT().FindByID(tenantID, scheduleID).Return(nil, errors.New("not found"))
			},
			wantErr: "jadwal tidak ditemukan atau akses ditolak",
		},
		{
			name:  "create error",
			input: baseInput(),
			setup: func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository) {
				scheduleRepo.EXPECT().FindByID(tenantID, scheduleID).
					Return(&domain.MaintenanceSchedule{ID: scheduleID, AssetID: assetID}, nil)
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success default doc type",
			input: func() MaintenanceDocumentInput {
				in := baseInput()
				in.DocType = ""
				return in
			}(),
			setup: func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository) {
				scheduleRepo.EXPECT().FindByID(tenantID, scheduleID).
					Return(&domain.MaintenanceSchedule{ID: scheduleID, AssetID: assetID}, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(doc *domain.MaintenanceDocument) error {
					if doc.DocType != domain.MaintenanceDocOther {
						t.Fatalf("expected default doc type other")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceDocumentRepository(ctrl)
			scheduleRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, scheduleRepo)
			}
			usecase := NewMaintenanceDocumentUsecase(repo, scheduleRepo)
			_, err := usecase.CreateDocument(tt.input)
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

func TestMaintenanceDocumentUsecase_ListDocuments(t *testing.T) {
	tenantID := uuid.New()
	scheduleID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository)
		wantErr string
	}{
		{
			name: "schedule not found",
			setup: func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository) {
				scheduleRepo.EXPECT().FindByID(tenantID, scheduleID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockMaintenanceDocumentRepository, scheduleRepo *mocks.MockMaintenanceScheduleRepository) {
				scheduleRepo.EXPECT().FindByID(tenantID, scheduleID).
					Return(&domain.MaintenanceSchedule{ID: scheduleID}, nil)
				repo.EXPECT().ListBySchedule(tenantID, scheduleID).Return([]domain.MaintenanceDocument{{ID: 1}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceDocumentRepository(ctrl)
			scheduleRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, scheduleRepo)
			}
			usecase := NewMaintenanceDocumentUsecase(repo, scheduleRepo)
			_, err := usecase.ListDocuments(tenantID, scheduleID)
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

func TestMaintenanceDocumentUsecase_DeleteDocument(t *testing.T) {
	tenantID := uuid.New()
	docID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockMaintenanceDocumentRepository)
		wantErr string
	}{
		{
			name: "not found",
			setup: func(repo *mocks.MockMaintenanceDocumentRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(repo *mocks.MockMaintenanceDocumentRepository) {
				doc := &domain.MaintenanceDocument{ID: docID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, docID).Return(doc, nil)
				repo.EXPECT().Delete(doc).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockMaintenanceDocumentRepository) {
				doc := &domain.MaintenanceDocument{ID: docID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, docID).Return(doc, nil)
				repo.EXPECT().Delete(doc).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockMaintenanceDocumentRepository(ctrl)
			scheduleRepo := mocks.NewMockMaintenanceScheduleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewMaintenanceDocumentUsecase(repo, scheduleRepo)
			err := usecase.DeleteDocument(tenantID, docID)
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
