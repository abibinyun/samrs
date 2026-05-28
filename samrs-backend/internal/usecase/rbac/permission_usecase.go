package usecase

import (
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type PermissionUsecase interface {
	ListAllPermissions() ([]PermissionDTO, error)
	ListPermissionsByRole(tenantID uuid.UUID, roleID int) ([]PermissionDTO, error)
}

type PermissionDTO struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Module string `json:"module"`
}

type permissionUsecase struct {
	repo     repository.PermissionRepository
	roleRepo repository.RoleRepository
}

func NewPermissionUsecase(repo repository.PermissionRepository, rr repository.RoleRepository) PermissionUsecase {
	return &permissionUsecase{repo: repo, roleRepo: rr}
}

func (u *permissionUsecase) ListAllPermissions() ([]PermissionDTO, error) {
	perms, err := u.repo.ListAll()
	if err != nil {
		return nil, err
	}

	out := make([]PermissionDTO, 0, len(perms))
	for _, p := range perms {
		out = append(out, PermissionDTO{
			ID:     p.ID,
			Name:   p.Name,
			Slug:   p.Slug,
			Module: p.Module,
		})
	}
	return out, nil
}

func (u *permissionUsecase) ListPermissionsByRole(tenantID uuid.UUID, roleID int) ([]PermissionDTO, error) {
	role, err := u.roleRepo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return nil, util.ErrNotFound("role tidak ditemukan atau akses ditolak")
	}
	if role.IsSystem {
		return nil, util.ErrForbidden("role sistem tidak dapat diakses")
	}

	perms, err := u.repo.ListByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	out := make([]PermissionDTO, 0, len(perms))
	for _, p := range perms {
		out = append(out, PermissionDTO{
			ID:     p.ID,
			Name:   p.Name,
			Slug:   p.Slug,
			Module: p.Module,
		})
	}
	return out, nil
}
