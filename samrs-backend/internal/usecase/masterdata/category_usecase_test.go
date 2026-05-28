package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestCategoryUsecase_CreateCategory(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		input   struct{ name, desc string }
		setup   func(repo *mocks.MockCategoryRepository)
		wantErr string
	}{
		{
			name: "missing name",
			input: struct{ name, desc string }{
				name: " ",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().ExistsBySlug(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "nama kategori wajib diisi",
		},
		{
			name: "slug exists",
			input: struct{ name, desc string }{
				name: "Life Support",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", nil).Return(true, nil)
			},
			wantErr: "nama kategori sudah digunakan",
		},
		{
			name: "exists check error",
			input: struct{ name, desc string }{
				name: "Life Support",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", nil).Return(false, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: struct{ name, desc string }{
				name: "Life Support",
				desc: "desc",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", nil).Return(false, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(cat *domain.Category) error {
					if cat.Name != "Life Support" || cat.Slug != "life-support" {
						t.Fatalf("unexpected category: %+v", cat)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockCategoryRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewCategoryUsecase(repo)
			_, err := usecase.CreateCategory(tenantID, tt.input.name, tt.input.desc)
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

func TestCategoryUsecase_UpdateCategory(t *testing.T) {
	tenantID := uuid.New()
	categoryID := uint(1)

	tests := []struct {
		name    string
		input   struct{ name, desc string }
		setup   func(repo *mocks.MockCategoryRepository)
		wantErr string
	}{
		{
			name: "missing name",
			input: struct{ name, desc string }{
				name: " ",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "nama kategori wajib diisi",
		},
		{
			name: "not found",
			input: struct{ name, desc string }{
				name: "Life Support",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().FindByID(tenantID, categoryID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "slug exists",
			input: struct{ name, desc string }{
				name: "Life Support",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				category := &domain.Category{ID: categoryID, TenantID: tenantID, Name: "Old", Slug: "old"}
				repo.EXPECT().FindByID(tenantID, categoryID).Return(category, nil)
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", &categoryID).Return(true, nil)
			},
			wantErr: "nama kategori sudah digunakan",
		},
		{
			name: "update error",
			input: struct{ name, desc string }{
				name: "Life Support",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				category := &domain.Category{ID: categoryID, TenantID: tenantID, Name: "Old", Slug: "old"}
				repo.EXPECT().FindByID(tenantID, categoryID).Return(category, nil)
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", &categoryID).Return(false, nil)
				repo.EXPECT().Update(category).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: struct{ name, desc string }{
				name: "Life Support",
				desc: "desc",
			},
			setup: func(repo *mocks.MockCategoryRepository) {
				category := &domain.Category{ID: categoryID, TenantID: tenantID, Name: "Old", Slug: "old"}
				repo.EXPECT().FindByID(tenantID, categoryID).Return(category, nil)
				repo.EXPECT().ExistsBySlug(tenantID, "life-support", &categoryID).Return(false, nil)
				repo.EXPECT().Update(category).DoAndReturn(func(updated *domain.Category) error {
					if updated.Name != "Life Support" || updated.Slug != "life-support" {
						t.Fatalf("unexpected update: %+v", updated)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockCategoryRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewCategoryUsecase(repo)
			_, err := usecase.UpdateCategory(tenantID, categoryID, tt.input.name, tt.input.desc)
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

func TestCategoryUsecase_DeleteCategory(t *testing.T) {
	tenantID := uuid.New()
	categoryID := uint(1)

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockCategoryRepository)
		wantErr string
	}{
		{
			name: "not found",
			setup: func(repo *mocks.MockCategoryRepository) {
				repo.EXPECT().FindByID(tenantID, categoryID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(repo *mocks.MockCategoryRepository) {
				category := &domain.Category{ID: categoryID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, categoryID).Return(category, nil)
				repo.EXPECT().Delete(category).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockCategoryRepository) {
				category := &domain.Category{ID: categoryID, TenantID: tenantID}
				repo.EXPECT().FindByID(tenantID, categoryID).Return(category, nil)
				repo.EXPECT().Delete(category).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockCategoryRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewCategoryUsecase(repo)
			err := usecase.DeleteCategory(tenantID, categoryID)
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
