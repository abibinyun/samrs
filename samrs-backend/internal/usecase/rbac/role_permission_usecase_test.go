package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestRolePermissionUsecase_AssignPermissions(t *testing.T) {
	tenantID := uuid.New()
	roleID := 10
	permIDs := []int{1, 2}

	tests := []struct {
		name    string
		setup   func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository)
		wantErr string
	}{
		{
			name: "role not found",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "system role",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID, IsSystem: true}, nil)
				permRepo.EXPECT().AssignPermissions(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat diubah",
		},
		{
			name: "success",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				permRepo.EXPECT().AssignPermissions(roleID, permIDs).Return(nil)
			},
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
			usecase := NewRolePermissionUsecase(roleRepo, permRepo)
			err := usecase.AssignPermissions(tenantID, roleID, permIDs)
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

func TestRolePermissionUsecase_RevokePermissions(t *testing.T) {
	tenantID := uuid.New()
	roleID := 10
	permIDs := []int{1, 2}

	tests := []struct {
		name    string
		setup   func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository)
		wantErr string
	}{
		{
			name: "role not found",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "system role",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID, IsSystem: true}, nil)
				permRepo.EXPECT().RevokePermissions(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat diubah",
		},
		{
			name: "success",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				permRepo.EXPECT().RevokePermissions(roleID, permIDs).Return(nil)
			},
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
			usecase := NewRolePermissionUsecase(roleRepo, permRepo)
			err := usecase.RevokePermissions(tenantID, roleID, permIDs)
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
