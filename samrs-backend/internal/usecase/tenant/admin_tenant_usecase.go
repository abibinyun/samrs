package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type TenantWithRoles struct {
	Tenant domain.Tenant `json:"tenant"`
	Roles  []domain.Role `json:"roles"`
}

type AdminTenantUsecase interface {
	ListTenantsWithRoles(filter repository.TenantFilter) ([]TenantWithRoles, int64, error)
	GetTenantDetail(tenantID uuid.UUID, userFilter repository.UserFilter, roleFilter repository.RoleFilter) (*TenantDetail, error)
	CreateTenant(input CreateTenantInput) (*domain.Tenant, error)
	UpdateTenant(input UpdateTenantInput) (*domain.Tenant, error)
	SetTenantStatus(tenantID uuid.UUID, status string) (*domain.Tenant, error)
}

type adminTenantUsecase struct {
	tenantRepo   repository.TenantRepository
	roleRepo     repository.RoleRepository
	userRepo     repository.UserRepository
	roomRepo     repository.RoomRepository
	bedRepo      repository.BedRepository
	categoryRepo repository.CategoryRepository
	assetRepo    repository.AssetRepository
	auditRepo    repository.AuditTrailRepository
}

func NewAdminTenantUsecase(
	tr repository.TenantRepository,
	rr repository.RoleRepository,
	ur repository.UserRepository,
	roomRepo repository.RoomRepository,
	bedRepo repository.BedRepository,
	categoryRepo repository.CategoryRepository,
	assetRepo repository.AssetRepository,
	auditRepo repository.AuditTrailRepository,
) AdminTenantUsecase {
	return &adminTenantUsecase{
		tenantRepo:   tr,
		roleRepo:     rr,
		userRepo:     ur,
		roomRepo:     roomRepo,
		bedRepo:      bedRepo,
		categoryRepo: categoryRepo,
		assetRepo:    assetRepo,
		auditRepo:    auditRepo,
	}
}

func (u *adminTenantUsecase) ListTenantsWithRoles(filter repository.TenantFilter) ([]TenantWithRoles, int64, error) {
	tenants, total, err := u.tenantRepo.ListAll(filter)
	if err != nil {
		return nil, 0, err
	}

	tenantIDs := make([]uuid.UUID, 0, len(tenants))
	for _, t := range tenants {
		tenantIDs = append(tenantIDs, t.ID)
	}

	roles, err := u.roleRepo.FindAllByTenantIDs(tenantIDs)
	if err != nil {
		return nil, 0, err
	}

	rolesByTenant := make(map[uuid.UUID][]domain.Role)
	for _, r := range roles {
		if r.TenantID != nil {
			rolesByTenant[*r.TenantID] = append(rolesByTenant[*r.TenantID], r)
		}
	}

	result := make([]TenantWithRoles, 0, len(tenants))
	for _, t := range tenants {
		result = append(result, TenantWithRoles{
			Tenant: t,
			Roles:  rolesByTenant[t.ID],
		})
	}

	return result, total, nil
}

type TenantDetailCounts struct {
	Users      int64 `json:"users"`
	Roles      int64 `json:"roles"`
	Rooms      int64 `json:"rooms"`
	Beds       int64 `json:"beds"`
	Categories int64 `json:"categories"`
	Assets     int64 `json:"assets"`
	Audits     int64 `json:"audits"`
}

type TenantDetail struct {
	Tenant    domain.Tenant      `json:"tenant"`
	Roles     []domain.Role      `json:"roles"`
	Users     []domain.User      `json:"users"`
	UsersMeta UserListMeta       `json:"users_meta"`
	RolesMeta RoleListMeta       `json:"roles_meta"`
	Counts    TenantDetailCounts `json:"counts"`
}

type UserListMeta struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

type RoleListMeta struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

type CreateTenantInput struct {
	Name       string
	Slug       string
	Address    string
	TenantType string
	Status     string
}

type UpdateTenantInput struct {
	ID         uuid.UUID
	Name       string
	Slug       string
	Address    string
	TenantType string
}

func (u *adminTenantUsecase) GetTenantDetail(tenantID uuid.UUID, userFilter repository.UserFilter, roleFilter repository.RoleFilter) (*TenantDetail, error) {
	tenant, err := u.tenantRepo.FindByID(tenantID)
	if err != nil {
		return nil, err
	}

	roles, roleTotal, err := u.roleRepo.ListByTenant(tenantID, roleFilter)
	if err != nil {
		return nil, err
	}

	users, userTotal, err := u.userRepo.ListByTenant(tenantID, userFilter)
	if err != nil {
		return nil, err
	}

	userCount, err := u.userRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	roleCount, err := u.roleRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	roomCount, err := u.roomRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	bedCount, err := u.bedRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	categoryCount, err := u.categoryRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	assetCount, err := u.assetRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	auditCount, err := u.auditRepo.CountByTenant(tenantID)
	if err != nil {
		return nil, err
	}

	return &TenantDetail{
		Tenant: *tenant,
		Roles:  roles,
		Users:  users,
		UsersMeta: UserListMeta{
			Total:   userTotal,
			Page:    userFilter.Page,
			PerPage: userFilter.PerPage,
		},
		RolesMeta: RoleListMeta{
			Total:   roleTotal,
			Page:    roleFilter.Page,
			PerPage: roleFilter.PerPage,
		},
		Counts: TenantDetailCounts{
			Users:      userCount,
			Roles:      roleCount,
			Rooms:      roomCount,
			Beds:       bedCount,
			Categories: categoryCount,
			Assets:     assetCount,
			Audits:     auditCount,
		},
	}, nil
}

func (u *adminTenantUsecase) CreateTenant(input CreateTenantInput) (*domain.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama tenant wajib diisi")
	}

	tenantSlug := strings.TrimSpace(input.Slug)
	if tenantSlug == "" {
		tenantSlug = slug.Make(name)
	} else {
		tenantSlug = slug.Make(tenantSlug)
	}

	exists, err := u.tenantRepo.ExistsBySlug(tenantSlug, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("slug tenant sudah digunakan")
	}

	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "active"
	}
	if !isValidTenantStatus(status) {
		return nil, util.ErrValidation("status tenant tidak valid")
	}

	tenant := &domain.Tenant{
		Name:       name,
		Slug:       tenantSlug,
		Address:    strings.TrimSpace(input.Address),
		TenantType: strings.TrimSpace(input.TenantType),
		Status:     status,
	}

	if err := u.tenantRepo.Create(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (u *adminTenantUsecase) UpdateTenant(input UpdateTenantInput) (*domain.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, util.ErrValidation("nama tenant wajib diisi")
	}

	tenant, err := u.tenantRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	tenantSlug := strings.TrimSpace(input.Slug)
	if tenantSlug == "" {
		tenantSlug = slug.Make(name)
	} else {
		tenantSlug = slug.Make(tenantSlug)
	}

	exists, err := u.tenantRepo.ExistsBySlug(tenantSlug, &tenant.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, util.ErrConflict("slug tenant sudah digunakan")
	}

	tenant.Name = name
	tenant.Slug = tenantSlug
	tenant.Address = strings.TrimSpace(input.Address)
	tenant.TenantType = strings.TrimSpace(input.TenantType)

	if err := u.tenantRepo.Update(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (u *adminTenantUsecase) SetTenantStatus(tenantID uuid.UUID, status string) (*domain.Tenant, error) {
	tenant, err := u.tenantRepo.FindByID(tenantID)
	if err != nil {
		return nil, err
	}

	normalized := strings.TrimSpace(status)
	if !isValidTenantStatus(normalized) {
		return nil, util.ErrValidation("status tenant tidak valid")
	}

	tenant.Status = normalized
	if err := u.tenantRepo.Update(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func isValidTenantStatus(status string) bool {
	switch status {
	case "active", "suspended":
		return true
	default:
		return false
	}
}
