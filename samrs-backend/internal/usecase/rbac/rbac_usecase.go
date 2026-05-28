package usecase

import (
	"samrs-backend/internal/repository"
)

type RBACUsecase interface {
	Authorize(roleID int, permissionSlug string) (bool, error)
}

type rbacUsecase struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRBACUsecase(rr repository.RoleRepository, pr repository.PermissionRepository) RBACUsecase {
	return &rbacUsecase{rr, pr}
}

func (u *rbacUsecase) Authorize(roleID int, permissionSlug string) (bool, error) {
	role, err := u.roleRepo.FindByID(roleID)
	if err != nil {
		return false, err
	}

	if role.IsSystem {
		return true, nil
	}

	return u.permissionRepo.RoleHasPermission(roleID, permissionSlug)
}
