package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestDocumentUsecase_CreateDocument(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		input   DocumentInput
		files   []DocumentFileInput
		wantErr string
	}{
		{
			name: "missing title",
			input: DocumentInput{
				TenantID: tenantID,
				Title:    " ",
			},
			files:   []DocumentFileInput{{Filename: "a.pdf", FilePath: "/a.pdf"}},
			wantErr: "judul dokumen wajib diisi",
		},
		{
			name: "missing files",
			input: DocumentInput{
				TenantID: tenantID,
				Title:    "Doc",
			},
			files:   nil,
			wantErr: "file dokumen wajib diunggah",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})

			_, _, err := usecase.CreateDocument(tt.input, tt.files)
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

func TestDocumentUsecase_AddFiles(t *testing.T) {
	tenantID := uuid.New()
	docID := uint(1)
	uploaderID := uuid.New()

	tests := []struct {
		name    string
		files   []DocumentFileInput
		setup   func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository)
		wantErr string
	}{
		{
			name:    "missing files",
			files:   nil,
			wantErr: "file dokumen wajib diunggah",
		},
		{
			name:  "document not found",
			files: []DocumentFileInput{{Filename: "a.pdf", FilePath: "/a.pdf"}},
			setup: func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:  "create error",
			files: []DocumentFileInput{{Filename: "a.pdf", FilePath: "/a.pdf"}},
			setup: func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(&domain.Document{ID: docID}, nil)
				fileRepo.EXPECT().CreateMany(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name:  "success",
			files: []DocumentFileInput{{Filename: "a.pdf", FilePath: "/a.pdf"}},
			setup: func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(&domain.Document{ID: docID}, nil)
				fileRepo.EXPECT().CreateMany(gomock.Any()).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, fileRepo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, err := usecase.AddFiles(tenantID, docID, uploaderID, tt.files)
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

func TestDocumentUsecase_ListDocuments(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.DocumentFilter{Page: 1, PerPage: 10}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockDocumentRepository)
		wantErr string
	}{
		{
			name: "list error",
			setup: func(repo *mocks.MockDocumentRepository) {
				repo.EXPECT().List(tenantID, filter).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockDocumentRepository) {
				repo.EXPECT().List(tenantID, filter).Return([]domain.Document{{ID: 1}}, int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, _, err := usecase.ListDocuments(tenantID, filter)
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

func TestDocumentUsecase_UpdateDocument(t *testing.T) {
	tenantID := uuid.New()
	docID := uint(1)

	tests := []struct {
		name    string
		input   DocumentUpdateInput
		setup   func(repo *mocks.MockDocumentRepository)
		wantErr string
	}{
		{
			name: "document not found",
			input: DocumentUpdateInput{
				Title: "Updated",
			},
			setup: func(repo *mocks.MockDocumentRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "invalid doc type",
			input: DocumentUpdateInput{
				DocType: "invalid",
			},
			setup: func(repo *mocks.MockDocumentRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(&domain.Document{ID: docID}, nil)
			},
			wantErr: "doc_type tidak valid",
		},
		{
			name: "update error",
			input: DocumentUpdateInput{
				Title:   "Updated",
				DocType: domain.DocumentTypeSOP,
			},
			setup: func(repo *mocks.MockDocumentRepository) {
				doc := &domain.Document{ID: docID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, docID).Return(doc, nil)
				repo.EXPECT().Update(doc).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: DocumentUpdateInput{
				Title:   "Updated",
				DocType: domain.DocumentTypeSOP,
			},
			setup: func(repo *mocks.MockDocumentRepository) {
				doc := &domain.Document{ID: docID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, docID).Return(doc, nil)
				repo.EXPECT().Update(doc).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, err := usecase.UpdateDocument(tenantID, docID, tt.input)
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

func TestDocumentUsecase_ListFiles(t *testing.T) {
	tenantID := uuid.New()
	docID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository)
		wantErr string
	}{
		{
			name: "document not found",
			setup: func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockDocumentRepository, fileRepo *mocks.MockDocumentFileRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(&domain.Document{ID: docID}, nil)
				fileRepo.EXPECT().ListByDocument(tenantID, docID).Return([]domain.DocumentFile{{ID: 1}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo, fileRepo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, err := usecase.ListFiles(tenantID, docID)
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

func TestDocumentUsecase_DeleteFile(t *testing.T) {
	tenantID := uuid.New()
	fileID := uint(1)

	tests := []struct {
		name    string
		setup   func(fileRepo *mocks.MockDocumentFileRepository)
		wantErr string
	}{
		{
			name: "file not found",
			setup: func(fileRepo *mocks.MockDocumentFileRepository) {
				fileRepo.EXPECT().FindByID(tenantID, fileID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(fileRepo *mocks.MockDocumentFileRepository) {
				file := &domain.DocumentFile{ID: fileID, TenantID: tenantID}
				fileRepo.EXPECT().FindByID(tenantID, fileID).Return(file, nil)
				fileRepo.EXPECT().Delete(file).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(fileRepo *mocks.MockDocumentFileRepository) {
				file := &domain.DocumentFile{ID: fileID, TenantID: tenantID}
				fileRepo.EXPECT().FindByID(tenantID, fileID).Return(file, nil)
				fileRepo.EXPECT().Delete(file).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(fileRepo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, err := usecase.DeleteFile(tenantID, fileID)
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

func TestDocumentUsecase_DeleteDocument(t *testing.T) {
	tenantID := uuid.New()
	docID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockDocumentRepository)
		wantErr string
	}{
		{
			name: "document not found",
			setup: func(repo *mocks.MockDocumentRepository) {
				repo.EXPECT().FindByID(tenantID, docID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockDocumentRepository(ctrl)
			fileRepo := mocks.NewMockDocumentFileRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewDocumentUsecase(repo, fileRepo, &gorm.DB{})
			_, err := usecase.DeleteDocument(tenantID, docID)
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
