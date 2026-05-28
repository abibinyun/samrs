package rbacrepo

import (
	"samrs-backend/internal/domain"
	repobase "samrs-backend/internal/repository/base"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	EnsurePermissions(perms []domain.Permission) error
	ListAll() ([]domain.Permission, error)
	ListByRoleID(roleID int) ([]domain.Permission, error)
	ListSlugsByRoleID(roleID int) ([]string, error)
	RoleHasPermission(roleID int, slug string) (bool, error)
	AssignPermissions(roleID int, permissionIDs []int) error
	RevokePermissions(roleID int, permissionIDs []int) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db}
}

func (r *permissionRepository) EnsurePermissions(perms []domain.Permission) error {
	for _, perm := range perms {
		db := repobase.NewDB(r.db).Model(&domain.Permission{})
		if err := db.Where("slug = ?", perm.Slug).FirstOrCreate(&perm).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *permissionRepository) ListAll() ([]domain.Permission, error) {
	var perms []domain.Permission
	db := repobase.NewDB(r.db).Model(&domain.Permission{})
	if err := db.Order("id ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) ListByRoleID(roleID int) ([]domain.Permission, error) {
	var perms []domain.Permission
	db := repobase.NewDB(r.db)
	err := db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.id ASC").
		Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) ListSlugsByRoleID(roleID int) ([]string, error) {
	var slugs []string
	db := repobase.NewDB(r.db)
	err := db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.id ASC").
		Pluck("permissions.slug", &slugs).Error
	return slugs, err
}

func (r *permissionRepository) RoleHasPermission(roleID int, slug string) (bool, error) {
	var count int64
	db := repobase.NewDB(r.db)
	err := db.Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.slug = ?", roleID, slug).
		Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) AssignPermissions(roleID int, permissionIDs []int) error {
	for _, pid := range permissionIDs {
		rp := domain.RolePermission{
			RoleID:       roleID,
			PermissionID: pid,
		}
		db := repobase.NewDB(r.db).Model(&domain.RolePermission{})
		if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *permissionRepository) RevokePermissions(roleID int, permissionIDs []int) error {
	db := repobase.NewDB(r.db).Model(&domain.RolePermission{})
	return db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&domain.RolePermission{}).Error
}
