package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	CreateUser(input CreateUserInput) (*domain.User, error)
	ListUsers(tenantID uuid.UUID, filter repository.UserFilter) ([]domain.User, int64, error)
	UpdateUserProfile(tenantID uuid.UUID, userID uuid.UUID, username string) (*domain.User, error)
	UpdateUserRole(tenantID uuid.UUID, userID uuid.UUID, roleID int) (*domain.User, error)
	ResetUserPassword(tenantID uuid.UUID, userID uuid.UUID, newPassword string) (*domain.User, error)
	SetUserActive(tenantID uuid.UUID, userID uuid.UUID, isActive bool) (*domain.User, error)
	DeleteUser(tenantID uuid.UUID, userID uuid.UUID) error
}

type CreateUserInput struct {
	TenantID uuid.UUID
	RoleID   int
	Username string
	Password string
	IsActive bool
}

type userUsecase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserUsecase(ur repository.UserRepository, rr repository.RoleRepository) UserUsecase {
	return &userUsecase{ur, rr}
}

func (u *userUsecase) CreateUser(input CreateUserInput) (*domain.User, error) {
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, util.ErrValidation("username dan password wajib diisi")
	}

	role, err := u.roleRepo.FindByIDAndTenant(input.RoleID, input.TenantID)
	if err != nil {
		return nil, util.ErrNotFound("role tidak ditemukan atau akses ditolak")
	}
	if role.IsSystem {
		return nil, util.ErrForbidden("role sistem tidak dapat digunakan untuk user tenant")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		TenantID:     input.TenantID,
		RoleID:       role.ID,
		Username:     input.Username,
		PasswordHash: string(hash),
		IsActive:     input.IsActive,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) ListUsers(tenantID uuid.UUID, filter repository.UserFilter) ([]domain.User, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 20
	}
	if filter.PerPage > 100 {
		filter.PerPage = 100
	}
	return u.userRepo.ListByTenant(tenantID, filter)
}

func (u *userUsecase) UpdateUserProfile(tenantID uuid.UUID, userID uuid.UUID, username string) (*domain.User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, util.ErrValidation("username wajib diisi")
	}

	user, err := u.userRepo.FindByIDAndTenant(userID, tenantID)
	if err != nil {
		return nil, err
	}

	user.Username = username
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateUserRole(tenantID uuid.UUID, userID uuid.UUID, roleID int) (*domain.User, error) {
	user, err := u.userRepo.FindByIDAndTenant(userID, tenantID)
	if err != nil {
		return nil, err
	}

	role, err := u.roleRepo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return nil, util.ErrNotFound("role tidak ditemukan atau akses ditolak")
	}

	user.RoleID = role.ID
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) ResetUserPassword(tenantID uuid.UUID, userID uuid.UUID, newPassword string) (*domain.User, error) {
	if strings.TrimSpace(newPassword) == "" {
		return nil, util.ErrValidation("password wajib diisi")
	}

	user, err := u.userRepo.FindByIDAndTenant(userID, tenantID)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(hash)
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) SetUserActive(tenantID uuid.UUID, userID uuid.UUID, isActive bool) (*domain.User, error) {
	user, err := u.userRepo.FindByIDAndTenant(userID, tenantID)
	if err != nil {
		return nil, err
	}

	user.IsActive = isActive
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) DeleteUser(tenantID uuid.UUID, userID uuid.UUID) error {
	user, err := u.userRepo.FindByIDAndTenant(userID, tenantID)
	if err != nil {
		return err
	}
	return u.userRepo.Delete(user)
}
