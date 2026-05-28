package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"go.uber.org/mock/gomock"
)

func TestRBACUsecase_Authorize(t *testing.T) {
	tests := []struct {
		name          string
		roleID        int
		permission    string
		setup         func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository)
		wantAllowed   bool
		wantErr       bool
	}{
		{
			name:       "role not found",
			roleID:     1,
			permission: "room:read",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByID(1).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:       "system role bypass",
			roleID:     2,
			permission: "room:read",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByID(2).Return(&domain.Role{IsSystem: true}, nil)
			},
			wantAllowed: true,
		},
		{
			name:       "non system role with permission",
			roleID:     3,
			permission: "room:read",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByID(3).Return(&domain.Role{IsSystem: false}, nil)
				permRepo.EXPECT().RoleHasPermission(3, "room:read").Return(true, nil)
			},
			wantAllowed: true,
		},
		{
			name:       "non system role without permission",
			roleID:     4,
			permission: "room:read",
			setup: func(roleRepo *mocks.MockRoleRepository, permRepo *mocks.MockPermissionRepository) {
				roleRepo.EXPECT().FindByID(4).Return(&domain.Role{IsSystem: false}, nil)
				permRepo.EXPECT().RoleHasPermission(4, "room:read").Return(false, nil)
			},
			wantAllowed: false,
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

			usecase := NewRBACUsecase(roleRepo, permRepo)
			allowed, err := usecase.Authorize(tt.roleID, tt.permission)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allowed != tt.wantAllowed {
				t.Fatalf("expected allowed=%v, got %v", tt.wantAllowed, allowed)
			}
		})
	}
}
