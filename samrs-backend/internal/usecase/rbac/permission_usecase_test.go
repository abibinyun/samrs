package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestPermissionUsecase_ListPermissionsByRole(t *testing.T) {
	tenantID := uuid.New()
	roleID := 1

	tests := []struct {
		name        string
		setup       func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository)
		wantErr     string
		wantEntries int
	}{
		{
			name: "role not found",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "role tidak ditemukan atau akses ditolak",
		},
		{
			name: "system role",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID, IsSystem: true}, nil)
			},
			wantErr: "role sistem tidak dapat diakses",
		},
		{
			name: "success",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				permRepo.EXPECT().ListByRoleID(roleID).Return([]domain.Permission{
					{ID: 1, Name: "View", Slug: "room:read", Module: "room"},
					{ID: 2, Name: "Create", Slug: "room:create", Module: "room"},
				}, nil)
			},
			wantEntries: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			permRepo := mocks.NewMockPermissionRepository(ctrl)
			if tt.setup != nil {
				tt.setup(roleRepo, permRepo)
			}

			usecase := NewPermissionUsecase(permRepo, roleRepo)
			perms, err := usecase.ListPermissionsByRole(tenantID, roleID)
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
			if len(perms) != tt.wantEntries {
				t.Fatalf("expected %d permissions, got %d", tt.wantEntries, len(perms))
			}
		})
	}
}
