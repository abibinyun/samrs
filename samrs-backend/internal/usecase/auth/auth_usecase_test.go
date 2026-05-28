package usecase

import (
	"errors"
	"os"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Login(t *testing.T) {
	if err := os.Setenv("JWT_SECRET", "test-secret"); err != nil {
		t.Fatalf("failed to set JWT_SECRET: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv("JWT_SECRET")
	})

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), 10)
	activeUser := domain.User{
		ID:           uuid.New(),
		TenantID:     uuid.New(),
		RoleID:       1,
		Username:     "alice",
		PasswordHash: string(hash),
		IsActive:     true,
	}
	inactiveUser := activeUser
	inactiveUser.Username = "bob"
	inactiveUser.IsActive = false

	tests := []struct {
		name     string
		username string
		password string
		setup    func(repo *mocks.MockUserRepository)
		wantErr  string
	}{
		{
			name:     "success",
			username: "alice",
			password: "secret",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByUsername("alice").Return(&activeUser, nil)
			},
		},
		{
			name:     "user not found",
			username: "missing",
			password: "secret",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByUsername("missing").Return(nil, errors.New("not found"))
			},
			wantErr: "username atau password salah",
		},
		{
			name:     "inactive user",
			username: "bob",
			password: "secret",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByUsername("bob").Return(&inactiveUser, nil)
			},
			wantErr: "akun anda tidak aktif",
		},
		{
			name:     "wrong password",
			username: "alice",
			password: "wrong",
			setup: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().GetByUsername("alice").Return(&activeUser, nil)
			},
			wantErr: "username atau password salah",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			permRepo := mocks.NewMockPermissionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(userRepo)
			}

			usecase := NewAuthUsecase(userRepo, permRepo)
			token, err := usecase.Login(tt.username, tt.password)
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
				t.Fatalf("expected success, got error: %v", err)
			}
			if token == "" {
				t.Fatalf("expected token, got empty string")
			}
		})
	}
}

func TestAuthUsecase_GetMe(t *testing.T) {
	user := domain.User{
		ID:       uuid.New(),
		TenantID: uuid.New(),
		RoleID:   2,
		Username: "carol",
		IsActive: true,
	}

	tests := []struct {
		name       string
		setup      func(userRepo *mocks.MockUserRepository, permRepo *mocks.MockPermissionRepository)
		wantUserID uuid.UUID
		wantPerms  int
		wantErr    bool
	}{
		{
			name: "success",
			setup: func(userRepo *mocks.MockUserRepository, permRepo *mocks.MockPermissionRepository) {
				userRepo.EXPECT().FindByID(user.ID).Return(&user, nil)
				permRepo.EXPECT().ListSlugsByRoleID(user.RoleID).Return([]string{"asset:read", "asset:update"}, nil)
			},
			wantUserID: user.ID,
			wantPerms:  2,
		},
		{
			name: "user not found",
			setup: func(userRepo *mocks.MockUserRepository, permRepo *mocks.MockPermissionRepository) {
				userRepo.EXPECT().FindByID(user.ID).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "permission repo error",
			setup: func(userRepo *mocks.MockUserRepository, permRepo *mocks.MockPermissionRepository) {
				userRepo.EXPECT().FindByID(user.ID).Return(&user, nil)
				permRepo.EXPECT().ListSlugsByRoleID(user.RoleID).Return(nil, errors.New("db down"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := mocks.NewMockUserRepository(ctrl)
			permRepo := mocks.NewMockPermissionRepository(ctrl)

			if tt.setup != nil {
				tt.setup(userRepo, permRepo)
			}

			usecase := NewAuthUsecase(userRepo, permRepo)
			gotUser, perms, err := usecase.GetMe(user.ID)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
			if gotUser == nil || gotUser.ID != tt.wantUserID {
				t.Fatalf("expected user %v, got %+v", tt.wantUserID, gotUser)
			}
			if len(perms) != tt.wantPerms {
				t.Fatalf("expected %d permissions, got %d", tt.wantPerms, len(perms))
			}
		})
	}
}
