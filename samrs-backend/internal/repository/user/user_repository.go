package userrepo

import (
	repobase "samrs-backend/internal/repository/base"
	"strings"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.User, error)
	Create(user *domain.User) error
	FindAllByTenant(tenantID uuid.UUID) ([]domain.User, error)
	Update(user *domain.User) error
	Delete(user *domain.User) error
	CountByTenant(tenantID uuid.UUID) (int64, error)
	ListByTenant(tenantID uuid.UUID, filter UserFilter) ([]domain.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	// GORM akan otomatis Join dengan Tenant dan Role karena kita sudah setup relasinya di struct
	db := repobase.NewDB(r.db).Model(&domain.User{})
	err := repobase.WithoutTenantScope(db).Preload("Tenant").Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	db := repobase.NewDB(r.db).Model(&domain.User{})
	err := repobase.WithoutTenantScope(db).Preload("Tenant").Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.User, error) {
	var user domain.User
	db := repobase.NewDB(r.db).Model(&domain.User{})
	err := repobase.WithTenant(db, tenantID).Preload("Tenant").Preload("Role").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	db := repobase.NewDB(r.db).Model(&domain.User{})
	return db.Create(user).Error
}

func (r *userRepository) FindAllByTenant(tenantID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	db := repobase.NewDB(r.db).Model(&domain.User{})
	err := repobase.WithTenant(db, tenantID).Preload("Tenant").Preload("Role").Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *domain.User) error {
	updates := map[string]interface{}{
		"role_id":       user.RoleID,
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"is_active":     user.IsActive,
		"updated_at":    user.UpdatedAt,
	}
	db := repobase.NewDB(r.db).Model(&domain.User{})
	return repobase.WithTenant(db, user.TenantID).
		Set("audit_record_id", user.ID).
		Where("id = ?", user.ID).
		Updates(updates).
		Error
}

func (r *userRepository) Delete(user *domain.User) error {
	db := repobase.NewDB(r.db).Model(&domain.User{})
	return repobase.WithTenant(db, user.TenantID).Delete(user).Error
}

func (r *userRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	var count int64
	db := repobase.NewDB(r.db).Model(&domain.User{})
	if err := repobase.WithTenant(db, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type UserFilter struct {
	Search  string
	Page    int
	PerPage int
}

func (r *userRepository) ListByTenant(tenantID uuid.UUID, filter UserFilter) ([]domain.User, int64, error) {
	var users []domain.User
	db := repobase.NewDB(r.db).Model(&domain.User{})
	query := repobase.WithTenant(db, tenantID)

	if strings.TrimSpace(filter.Search) != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("username ILIKE ?", like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Preload("Tenant").Preload("Role").
		Order("created_at desc").Limit(filter.PerPage).Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
