package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_CreateUser(t *testing.T) {
	tenantID := uuid.New()
	roleID := 5

	tests := []struct {
		name    string
		input   CreateUserInput
		setup   func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository)
		wantErr string
	}{
		{
			name: "missing username",
			input: CreateUserInput{
				TenantID: tenantID,
				RoleID:   roleID,
				Username: " ",
				Password: "pass",
			},
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			wantErr: "username dan password wajib diisi",
		},
		{
			name: "missing password",
			input: CreateUserInput{
				TenantID: tenantID,
				RoleID:   roleID,
				Username: "user",
				Password: " ",
			},
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			wantErr: "username dan password wajib diisi",
		},
		{
			name: "role not found",
			input: CreateUserInput{
				TenantID: tenantID,
				RoleID:   roleID,
				Username: "user",
				Password: "pass",
				IsActive: true,
			},
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(nil, errors.New("not found"))
				userRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			wantErr: "role tidak ditemukan atau akses ditolak",
		},
		{
			name: "system role not allowed",
			input: CreateUserInput{
				TenantID: tenantID,
				RoleID:   roleID,
				Username: "user",
				Password: "pass",
				IsActive: true,
			},
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID, IsSystem: true}, nil)
				userRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			wantErr: "role sistem tidak dapat digunakan untuk user tenant",
		},
		{
			name: "success",
			input: CreateUserInput{
				TenantID: tenantID,
				RoleID:   roleID,
				Username: "user",
				Password: "pass",
				IsActive: true,
			},
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).
					Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				userRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *domain.User) error {
					if user.TenantID != tenantID {
						t.Fatalf("expected tenant %v, got %v", tenantID, user.TenantID)
					}
					if user.RoleID != roleID {
						t.Fatalf("expected role %d, got %d", roleID, user.RoleID)
					}
					if user.Username != "user" {
						t.Fatalf("expected username user, got %s", user.Username)
					}
					if user.PasswordHash == "" {
						t.Fatalf("expected password hash set")
					}
					if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("pass")); err != nil {
						t.Fatalf("expected password hash to match: %v", err)
					}
					if !user.IsActive {
						t.Fatalf("expected user active")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo, roleRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			_, err := usecase.CreateUser(tt.input)
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

func TestUserUsecase_UpdateUserProfile(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		input   string
		setup   func(userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name:  "empty username",
			input: " ",
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().Update(gomock.Any()).Times(0)
			},
			wantErr: "username wajib diisi",
		},
		{
			name:  "user not found",
			input: "newname",
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:  "update error",
			input: "newname",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, Username: "old"}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name:  "success",
			input: "newname",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, Username: "old"}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			_, err := usecase.UpdateUserProfile(tenantID, userID, tt.input)
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

func TestUserUsecase_UpdateUserRole(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	roleID := 7

	tests := []struct {
		name    string
		setup   func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository)
		wantErr string
	}{
		{
			name: "user not found",
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "role not found",
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, RoleID: 1}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "role tidak ditemukan atau akses ditolak",
		},
		{
			name: "update error",
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, RoleID: 1}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				userRepo.EXPECT().Update(user).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(userRepo *mocks.MockUserRepository, roleRepo *mocks.MockRoleRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, RoleID: 1}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				roleRepo.EXPECT().FindByIDAndTenant(roleID, tenantID).Return(&domain.Role{ID: roleID, TenantID: &tenantID}, nil)
				userRepo.EXPECT().Update(user).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo, roleRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			_, err := usecase.UpdateUserRole(tenantID, userID, roleID)
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

func TestUserUsecase_ResetUserPassword(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name     string
		password string
		setup    func(userRepo *mocks.MockUserRepository)
		wantErr  string
	}{
		{
			name:     "empty password",
			password: " ",
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
				userRepo.EXPECT().Update(gomock.Any()).Times(0)
			},
			wantErr: "password wajib diisi",
		},
		{
			name:     "user not found",
			password: "newpass",
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:     "update error",
			password: "newpass",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, PasswordHash: "old"}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name:     "success",
			password: "newpass",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, PasswordHash: "old"}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).DoAndReturn(func(updated *domain.User) error {
					if updated.PasswordHash == "" {
						t.Fatalf("expected password hash set")
					}
					if err := bcrypt.CompareHashAndPassword([]byte(updated.PasswordHash), []byte("newpass")); err != nil {
						t.Fatalf("expected password hash to match: %v", err)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			_, err := usecase.ResetUserPassword(tenantID, userID, tt.password)
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

func TestUserUsecase_SetUserActive(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		active  bool
		setup   func(userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name:   "user not found",
			active: true,
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name:   "update error",
			active: false,
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, IsActive: true}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name:   "success",
			active: false,
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID, IsActive: true}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Update(user).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			_, err := usecase.SetUserActive(tenantID, userID, tt.active)
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

func TestUserUsecase_DeleteUser(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func(userRepo *mocks.MockUserRepository)
		wantErr string
	}{
		{
			name: "user not found",
			setup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Delete(user).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: userID, TenantID: tenantID}
				userRepo.EXPECT().FindByIDAndTenant(userID, tenantID).Return(user, nil)
				userRepo.EXPECT().Delete(user).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			roleRepo := mocks.NewMockRoleRepository(ctrl)
			if tt.setup != nil {
				tt.setup(userRepo)
			}

			usecase := NewUserUsecase(userRepo, roleRepo)
			err := usecase.DeleteUser(tenantID, userID)
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
