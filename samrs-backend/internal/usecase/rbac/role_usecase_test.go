package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestRoleUsecase_CreateRole(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name      string
		inputName string
		setup     func(repo *mocks.MockRoleRepository)
		wantErr   string
	}{
		{
			name:      "empty name",
			inputName: "  ",
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().Create(gomock.Any()).Times(0)
			},
			wantErr: "nama role wajib diisi",
		},
		{
			name:      "success",
			inputName: "Operator",
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(role *domain.Role) error {
					if role.Name != "Operator" {
						t.Fatalf("expected name Operator, got %s", role.Name)
					}
					if role.TenantID == nil || *role.TenantID != tenantID {
						t.Fatalf("expected tenant %v, got %v", tenantID, role.TenantID)
					}
					if role.IsSystem {
						t.Fatalf("expected IsSystem false")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(roleRepo)
			}
			usecase := NewRoleUsecase(roleRepo)
			role, err := usecase.CreateRole(tenantID, tt.inputName)
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
			if role == nil || role.Name != tt.inputName {
				t.Fatalf("expected role name %q, got %+v", tt.inputName, role)
			}
		})
	}
}

func TestRoleUsecase_UpdateRole(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name      string
		role      *domain.Role
		inputName string
		setup     func(repo *mocks.MockRoleRepository)
		wantErr   string
	}{
		{
			name:      "empty name",
			role:      &domain.Role{ID: 1, TenantID: &tenantID},
			inputName: "",
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByIDAndTenant(1, tenantID).Times(0)
			},
			wantErr: "nama role wajib diisi",
		},
		{
			name:      "role not found",
			role:      &domain.Role{ID: 2, TenantID: &tenantID},
			inputName: "Updater",
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByIDAndTenant(2, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:      "system role",
			role:      &domain.Role{ID: 3, TenantID: &tenantID, IsSystem: true},
			inputName: "Updater",
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByIDAndTenant(3, tenantID).Return(&domain.Role{ID: 3, TenantID: &tenantID, IsSystem: true}, nil)
				repo.EXPECT().Update(gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat diubah",
		},
		{
			name:      "success",
			role:      &domain.Role{ID: 4, TenantID: &tenantID, Name: "Old"},
			inputName: "New",
			setup: func(repo *mocks.MockRoleRepository) {
				role := &domain.Role{ID: 4, TenantID: &tenantID, Name: "Old"}
				repo.EXPECT().FindByIDAndTenant(4, tenantID).Return(role, nil)
				repo.EXPECT().Update(role).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(roleRepo)
			}
			usecase := NewRoleUsecase(roleRepo)
			updated, err := usecase.UpdateRole(tenantID, tt.role.ID, tt.inputName)
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
			if updated == nil || updated.Name != tt.inputName {
				t.Fatalf("expected updated name %q, got %+v", tt.inputName, updated)
			}
		})
	}
}

func TestRoleUsecase_DeleteRole(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockRoleRepository)
		roleID  int
		wantErr string
	}{
		{
			name:   "role not found",
			roleID: 1,
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByIDAndTenant(1, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:   "system role",
			roleID: 2,
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByIDAndTenant(2, tenantID).Return(&domain.Role{ID: 2, TenantID: &tenantID, IsSystem: true}, nil)
				repo.EXPECT().Delete(gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat dihapus",
		},
		{
			name:   "success",
			roleID: 3,
			setup: func(repo *mocks.MockRoleRepository) {
				role := &domain.Role{ID: 3, TenantID: &tenantID, Name: "Operator"}
				repo.EXPECT().FindByIDAndTenant(3, tenantID).Return(role, nil)
				repo.EXPECT().Delete(role).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(roleRepo)
			}
			usecase := NewRoleUsecase(roleRepo)
			err := usecase.DeleteRole(tenantID, tt.roleID)
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

func TestRoleUsecase_SetTenantAdmin(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name       string
		roleID     int
		isAdmin    bool
		setup      func(repo *mocks.MockRoleRepository)
		wantErr    string
		wantResult bool
	}{
		{
			name:   "role not found",
			roleID: 1,
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByID(1).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:   "system role",
			roleID: 2,
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByID(2).Return(&domain.Role{ID: 2, IsSystem: true}, nil)
				repo.EXPECT().Update(gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat diubah",
		},
		{
			name:   "role without tenant",
			roleID: 3,
			setup: func(repo *mocks.MockRoleRepository) {
				repo.EXPECT().FindByID(3).Return(&domain.Role{ID: 3, TenantID: nil}, nil)
				repo.EXPECT().Update(gomock.Any()).Times(0)
			},
			wantErr: "role tidak terkait tenant",
		},
		{
			name:    "success",
			roleID:  4,
			isAdmin: true,
			setup: func(repo *mocks.MockRoleRepository) {
				role := &domain.Role{ID: 4, TenantID: &tenantID}
				repo.EXPECT().FindByID(4).Return(role, nil)
				repo.EXPECT().Update(role).Return(nil)
			},
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(roleRepo)
			}
			usecase := NewRoleUsecase(roleRepo)
			role, err := usecase.SetTenantAdmin(tt.roleID, tt.isAdmin)
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
			if role == nil || role.IsTenantAdmin != tt.wantResult {
				t.Fatalf("expected is_tenant_admin=%v, got %+v", tt.wantResult, role)
			}
		})
	}
}
