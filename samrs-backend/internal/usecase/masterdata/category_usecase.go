package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CategoryUsecase interface {
	CreateCategory(tenantID uuid.UUID, name, desc string) (*domain.Category, error)
	GetAllCategories(tenantID uuid.UUID, filter repository.CategoryFilter) ([]domain.Category, int64, error)
	GetCategoryByID(tenantID uuid.UUID, id uint) (*domain.Category, error)
	UpdateCategory(tenantID uuid.UUID, id uint, name, desc string) (*domain.Category, error)
	DeleteCategory(tenantID uuid.UUID, id uint) error
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo}
}

func (u *categoryUsecase) CreateCategory(tenantID uuid.UUID, name, desc string) (*domain.Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, util.ErrValidation("nama kategori wajib diisi")
	}

	category := &domain.Category{
		TenantID:    tenantID,
		Name:        name,
		Slug:        slug.Make(name), // Otomatis jadi "life-support"
		Description: desc,
	}

	exists, err := u.repo.ExistsBySlug(tenantID, category.Slug, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("nama kategori sudah digunakan")
	}

	if err := u.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (u *categoryUsecase) GetAllCategories(tenantID uuid.UUID, filter repository.CategoryFilter) ([]domain.Category, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *categoryUsecase) GetCategoryByID(tenantID uuid.UUID, id uint) (*domain.Category, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *categoryUsecase) UpdateCategory(tenantID uuid.UUID, id uint, name, desc string) (*domain.Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, util.ErrValidation("nama kategori wajib diisi")
	}

	category, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	newSlug := slug.Make(name)
	exists, err := u.repo.ExistsBySlug(tenantID, newSlug, &category.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("nama kategori sudah digunakan")
	}

	category.Name = name
	category.Slug = newSlug
	category.Description = desc
	if err := u.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (u *categoryUsecase) DeleteCategory(tenantID uuid.UUID, id uint) error {
	category, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(category)
}
