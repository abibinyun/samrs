package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type VendorUsecase interface {
	CreateVendor(tenantID uuid.UUID, input VendorInput) (*domain.Vendor, error)
	GetAllVendors(tenantID uuid.UUID, filter repository.VendorFilter) ([]domain.Vendor, int64, error)
	GetVendorByID(tenantID uuid.UUID, id uint) (*domain.Vendor, error)
	UpdateVendor(tenantID uuid.UUID, id uint, input VendorInput) (*domain.Vendor, error)
	DeleteVendor(tenantID uuid.UUID, id uint) error
}

type VendorInput struct {
	Code        string
	Name        string
	ContactName string
	Phone       string
	Email       string
	Address     string
}

type vendorUsecase struct {
	repo repository.VendorRepository
}

func NewVendorUsecase(repo repository.VendorRepository) VendorUsecase {
	return &vendorUsecase{repo}
}

func (u *vendorUsecase) CreateVendor(tenantID uuid.UUID, input VendorInput) (*domain.Vendor, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama vendor wajib diisi")
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode vendor sudah digunakan")
	}

	vendor := &domain.Vendor{
		TenantID:    tenantID,
		Code:        code,
		Name:        name,
		ContactName: strings.TrimSpace(input.ContactName),
		Phone:       strings.TrimSpace(input.Phone),
		Email:       strings.TrimSpace(input.Email),
		Address:     strings.TrimSpace(input.Address),
	}

	if err := u.repo.Create(vendor); err != nil {
		return nil, err
	}
	return vendor, nil
}

func (u *vendorUsecase) GetAllVendors(tenantID uuid.UUID, filter repository.VendorFilter) ([]domain.Vendor, int64, error) {
	return u.repo.FindAllByTenant(tenantID, filter)
}

func (u *vendorUsecase) GetVendorByID(tenantID uuid.UUID, id uint) (*domain.Vendor, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *vendorUsecase) UpdateVendor(tenantID uuid.UUID, id uint, input VendorInput) (*domain.Vendor, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama vendor wajib diisi")
	}

	vendor, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = slug.Make(name)
	} else {
		code = slug.Make(code)
	}

	exists, err := u.repo.ExistsByCode(tenantID, code, &vendor.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("kode vendor sudah digunakan")
	}

	vendor.Code = code
	vendor.Name = name
	vendor.ContactName = strings.TrimSpace(input.ContactName)
	vendor.Phone = strings.TrimSpace(input.Phone)
	vendor.Email = strings.TrimSpace(input.Email)
	vendor.Address = strings.TrimSpace(input.Address)

	if err := u.repo.Update(vendor); err != nil {
		return nil, err
	}
	return vendor, nil
}

func (u *vendorUsecase) DeleteVendor(tenantID uuid.UUID, id uint) error {
	vendor, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(vendor)
}
