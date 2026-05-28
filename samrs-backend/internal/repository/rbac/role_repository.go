package rbacrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByID(id int) (*domain.Role, error)
	FindByIDAndTenant(id int, tenantID uuid.UUID) (*domain.Role, error)
	FindAllByTenant(tenantID uuid.UUID) ([]domain.Role, error)
	FindAllByTenantIDs(tenantIDs []uuid.UUID) ([]domain.Role, error)
	CountByTenant(tenantID uuid.UUID) (int64, error)
	ListByTenant(tenantID uuid.UUID, filter RoleFilter) ([]domain.Role, int64, error)
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(role *domain.Role) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) FindByID(id int) (*domain.Role, error) {
	var role domain.Role
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if err := repobase.WithoutTenantScope(db).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByIDAndTenant(id int, tenantID uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if err := repobase.WithTenant(db, tenantID).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAllByTenant(tenantID uuid.UUID) ([]domain.Role, error) {
	var roles []domain.Role
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if err := repobase.WithTenant(db, tenantID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindAllByTenantIDs(tenantIDs []uuid.UUID) ([]domain.Role, error) {
	var roles []domain.Role
	if len(tenantIDs) == 0 {
		return roles, nil
	}
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if err := db.Where("tenant_id IN ?", tenantIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type RoleFilter struct {
	Search  string
	Page    int
	PerPage int
}

func (r *roleRepository) ListByTenant(tenantID uuid.UUID, filter RoleFilter) ([]domain.Role, int64, error) {
	var roles []domain.Role
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ?", like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Order("id desc").Limit(filter.PerPage).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *roleRepository) Create(role *domain.Role) error {
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	return db.Create(role).Error
}

func (r *roleRepository) Update(role *domain.Role) error {
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if role.TenantID == nil {
		updates := map[string]interface{}{
			"name":            role.Name,
			"is_system":       role.IsSystem,
			"is_tenant_admin": role.IsTenantAdmin,
		}
		return repobase.WithoutTenantScope(db).
			Set("audit_record_id", role.ID).
			Where("id = ?", role.ID).
			Updates(updates).
			Error
	}
	updates := map[string]interface{}{
		"name":            role.Name,
		"is_system":       role.IsSystem,
		"is_tenant_admin": role.IsTenantAdmin,
	}
	return repobase.WithTenant(db, *role.TenantID).
		Set("audit_record_id", role.ID).
		Where("id = ?", role.ID).
		Updates(updates).
		Error
}

func (r *roleRepository) Delete(role *domain.Role) error {
	db := repobase.NewDB(r.db).Model(&domain.Role{})
	if role.TenantID == nil {
		return repobase.WithoutTenantScope(db).Delete(role).Error
	}
	return repobase.WithTenant(db, *role.TenantID).Delete(role).Error
}
