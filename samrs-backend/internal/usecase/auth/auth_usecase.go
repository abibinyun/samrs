package usecase

import (
	"os"
	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(username, password string) (string, error)
	GetMe(userID uuid.UUID) (*domain.User, []string, error)
}

type authUsecase struct {
	userRepo       repository.UserRepository
	permissionRepo repository.PermissionRepository
}

func NewAuthUsecase(ur repository.UserRepository, pr repository.PermissionRepository) AuthUsecase {
	return &authUsecase{ur, pr}
}

func (u *authUsecase) Login(username, password string) (string, error) {
	// 1. Cari user di DB
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", util.ErrUnauthorized("username atau password salah")
	}

	// 2. Cek apakah user aktif
	if !user.IsActive {
		return "", util.ErrForbidden("akun anda tidak aktif")
	}

	// 3. Verifikasi Password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", util.ErrUnauthorized("username atau password salah")
	}

	// 4. Generate JWT Token
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role_id":   user.RoleID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func (u *authUsecase) GetMe(userID uuid.UUID) (*domain.User, []string, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	perms, err := u.permissionRepo.ListSlugsByRoleID(user.RoleID)
	if err != nil {
		return nil, nil, err
	}

	return user, perms, nil
}
