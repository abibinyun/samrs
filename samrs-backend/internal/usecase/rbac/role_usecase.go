package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type RoleUsecase interface {
	CreateRole(tenantID uuid.UUID, name string) (*domain.Role, error)
	ListRoles(tenantID uuid.UUID) ([]domain.Role, error)
	GetRoleByID(tenantID uuid.UUID, roleID int) (*domain.Role, error)
	UpdateRole(tenantID uuid.UUID, roleID int, name string) (*domain.Role, error)
	DeleteRole(tenantID uuid.UUID, roleID int) error
	SetTenantAdmin(roleID int, isTenantAdmin bool) (*domain.Role, error)
}

type roleUsecase struct {
	repo repository.RoleRepository
}

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{repo}
}

func (u *roleUsecase) CreateRole(tenantID uuid.UUID, name string) (*domain.Role, error) {
	if strings.TrimSpace(name) == "" {
		return nil, util.ErrValidation("nama role wajib diisi")
	}

	role := &domain.Role{
		Name:     name,
		TenantID: &tenantID,
		IsSystem: false,
	}
	if err := u.repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) ListRoles(tenantID uuid.UUID) ([]domain.Role, error) {
	return u.repo.FindAllByTenant(tenantID)
}

func (u *roleUsecase) GetRoleByID(tenantID uuid.UUID, roleID int) (*domain.Role, error) {
	role, err := u.repo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return nil, err
	}
	if role.IsSystem {
		return nil, util.ErrForbidden("role sistem tidak dapat diakses")
	}
	return role, nil
}

func (u *roleUsecase) UpdateRole(tenantID uuid.UUID, roleID int, name string) (*domain.Role, error) {
	if strings.TrimSpace(name) == "" {
		return nil, util.ErrValidation("nama role wajib diisi")
	}

	role, err := u.repo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return nil, err
	}
	if role.IsSystem {
		return nil, util.ErrForbidden("role sistem tidak dapat diubah")
	}

	role.Name = name
	if err := u.repo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (u *roleUsecase) DeleteRole(tenantID uuid.UUID, roleID int) error {
	role, err := u.repo.FindByIDAndTenant(roleID, tenantID)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return util.ErrForbidden("role sistem tidak dapat dihapus")
	}
	return u.repo.Delete(role)
}

func (u *roleUsecase) SetTenantAdmin(roleID int, isTenantAdmin bool) (*domain.Role, error) {
	role, err := u.repo.FindByID(roleID)
	if err != nil {
		return nil, err
	}
	if role.IsSystem {
		return nil, util.ErrForbidden("role sistem tidak dapat diubah")
	}
	if role.TenantID == nil {
		return nil, util.ErrForbidden("role tidak terkait tenant")
	}

	role.IsTenantAdmin = isTenantAdmin
	if err := u.repo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}
