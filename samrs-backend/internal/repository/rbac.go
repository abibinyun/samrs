package repository

import (
	rbacrepo "samrs-backend/internal/repository/rbac"

	"gorm.io/gorm"
)

type RoleFilter = rbacrepo.RoleFilter

type RoleRepository = rbacrepo.RoleRepository

type PermissionRepository = rbacrepo.PermissionRepository

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return rbacrepo.NewRoleRepository(db)
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return rbacrepo.NewPermissionRepository(db)
}
