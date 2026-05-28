package usecase

import (
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type RolePermissionUsecase interface {
	AssignPermissions(tenantID uuid.UUID, roleID int, permissionIDs []int) error
	RevokePermissions(tenantID uuid.UUID, roleID int, permissionIDs []int) error
}

type rolePermissionUsecase struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRolePermissionUsecase(rr repository.RoleRepository, pr repository.PermissionRepository) RolePermissionUsecase {
	return &rolePermissionUsecase{rr, pr}
}

func (u *rolePermissionUsecase) AssignPermissions(tenantID uuid.UUID, roleID int, permissionIDs []int) error {
	role, err := u.roleRepo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return util.ErrForbidden("role sistem tidak dapat diubah")
	}
	return u.permissionRepo.AssignPermissions(roleID, permissionIDs)
}

func (u *rolePermissionUsecase) RevokePermissions(tenantID uuid.UUID, roleID int, permissionIDs []int) error {
	role, err := u.roleRepo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return util.ErrForbidden("role sistem tidak dapat diubah")
	}
	return u.permissionRepo.RevokePermissions(roleID, permissionIDs)
}
