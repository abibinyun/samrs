package seed

import (
	"errors"
	"log"
	"os"
	"strings"

	"samrs-backend/internal/config"
	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	permissionRepo := repository.NewPermissionRepository(db.Session(&gorm.Session{NewDB: true}))

	log.Println("Menjalankan seeder inti jika diperlukan...")

	// a. Create Default Role (System Role)
	adminRole := domain.Role{Name: "Super Admin", IsSystem: true}
	roleDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.Role{})
	if err := roleDB.Where("name = ? AND tenant_id IS NULL AND is_system = ?", adminRole.Name, true).
		FirstOrCreate(&adminRole).Error; err != nil {
		return err
	}

	// b. Create Default Tenant
	tenant := domain.Tenant{
		Name:   getenvDefault("SEED_TENANT_NAME", "RS Pusat Utama"),
		Slug:   getenvDefault("SEED_TENANT_SLUG", "rs-pusat"),
		Status: getenvDefault("SEED_TENANT_STATUS", "active"),
	}
	tenantDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.Tenant{})
	if err := tenantDB.Where(domain.Tenant{Slug: tenant.Slug}).FirstOrCreate(&tenant).Error; err != nil {
		return err
	}

	// c. Create Tenant Admin Role (Scoped)
	tenantAdminRole := domain.Role{
		Name:          "Tenant Admin",
		TenantID:      &tenant.ID,
		IsSystem:      false,
		IsTenantAdmin: true,
	}
	tenantRoleDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.Role{})
	if err := tenantRoleDB.Where("name = ? AND tenant_id = ?", tenantAdminRole.Name, tenant.ID).
		FirstOrCreate(&tenantAdminRole).Error; err != nil {
		return err
	}
	if !tenantAdminRole.IsTenantAdmin {
		tenantAdminRole.IsTenantAdmin = true
		updates := map[string]interface{}{
			"is_tenant_admin": true,
		}
		if err := db.Session(&gorm.Session{NewDB: true}).
			Model(&domain.Role{}).
			Where("id = ?", tenantAdminRole.ID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	// c2. Seed Asset Status (Scoped)
	defaultAssetStatuses := []domain.AssetStatus{
		{TenantID: tenant.ID, Code: "ready", Name: "Ready"},
		{TenantID: tenant.ID, Code: "broken", Name: "Broken"},
		{TenantID: tenant.ID, Code: "maintenance", Name: "Maintenance"},
	}
	for _, status := range defaultAssetStatuses {
		statusDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.AssetStatus{})
		if err := statusDB.Where("tenant_id = ? AND code = ?", status.TenantID, status.Code).
			FirstOrCreate(&status).Error; err != nil {
			return err
		}
	}

	// d. Create Admin User (Password: password123)
	var existingAdmin domain.User
	adminUsername := getenvDefault("SEED_ADMIN_USERNAME", "admin")
	adminPassword := getenvDefault("SEED_ADMIN_PASSWORD", "password123")
	userDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.User{})
	err := config.WithSkipTenantScope(userDB).Where("username = ?", adminUsername).First(&existingAdmin).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		hpwd, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 10)
		if err != nil {
			return err
		}
		adminUser := domain.User{
			TenantID:     tenant.ID,
			RoleID:       adminRole.ID,
			Username:     adminUsername,
			PasswordHash: string(hpwd),
			IsActive:     true,
		}

		userCreateDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.User{})
		if err := userCreateDB.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Printf("Seeder berhasil: User '%s' siap digunakan.", adminUsername)
	}

	// e. Seed Permissions (RBAC)
	defaultPermissions := []domain.Permission{
		{Name: "Create Room", Slug: "room:create", Module: "room"},
		{Name: "Read Room", Slug: "room:read", Module: "room"},
		{Name: "Update Room", Slug: "room:update", Module: "room"},
		{Name: "Delete Room", Slug: "room:delete", Module: "room"},
		{Name: "Create Bed", Slug: "bed:create", Module: "bed"},
		{Name: "Read Bed", Slug: "bed:read", Module: "bed"},
		{Name: "Update Bed", Slug: "bed:update", Module: "bed"},
		{Name: "Delete Bed", Slug: "bed:delete", Module: "bed"},
		{Name: "Create Category", Slug: "category:create", Module: "category"},
		{Name: "Read Category", Slug: "category:read", Module: "category"},
		{Name: "Update Category", Slug: "category:update", Module: "category"},
		{Name: "Delete Category", Slug: "category:delete", Module: "category"},
		{Name: "Create Asset", Slug: "asset:create", Module: "asset"},
		{Name: "Read Asset", Slug: "asset:read", Module: "asset"},
		{Name: "Update Asset", Slug: "asset:update", Module: "asset"},
		{Name: "Delete Asset", Slug: "asset:delete", Module: "asset"},
		{Name: "Create Vendor", Slug: "vendor:create", Module: "vendor"},
		{Name: "Read Vendor", Slug: "vendor:read", Module: "vendor"},
		{Name: "Update Vendor", Slug: "vendor:update", Module: "vendor"},
		{Name: "Delete Vendor", Slug: "vendor:delete", Module: "vendor"},
		{Name: "Create Brand", Slug: "brand:create", Module: "brand"},
		{Name: "Read Brand", Slug: "brand:read", Module: "brand"},
		{Name: "Update Brand", Slug: "brand:update", Module: "brand"},
		{Name: "Delete Brand", Slug: "brand:delete", Module: "brand"},
		{Name: "Create Model", Slug: "model:create", Module: "model"},
		{Name: "Read Model", Slug: "model:read", Module: "model"},
		{Name: "Update Model", Slug: "model:update", Module: "model"},
		{Name: "Delete Model", Slug: "model:delete", Module: "model"},
		{Name: "Create Asset Status", Slug: "asset_status:create", Module: "asset_status"},
		{Name: "Read Asset Status", Slug: "asset_status:read", Module: "asset_status"},
		{Name: "Update Asset Status", Slug: "asset_status:update", Module: "asset_status"},
		{Name: "Delete Asset Status", Slug: "asset_status:delete", Module: "asset_status"},
		{Name: "Create Maintenance Schedule", Slug: "maintenance:create", Module: "maintenance"},
		{Name: "Read Maintenance Schedule", Slug: "maintenance:read", Module: "maintenance"},
		{Name: "Update Maintenance Schedule", Slug: "maintenance:update", Module: "maintenance"},
		{Name: "Delete Maintenance Schedule", Slug: "maintenance:delete", Module: "maintenance"},
		{Name: "Complete Maintenance Schedule", Slug: "maintenance:complete", Module: "maintenance"},
		{Name: "Export Report", Slug: "report:export", Module: "report"},
		{Name: "Upload Maintenance Document", Slug: "maintenance_document:upload", Module: "maintenance_document"},
		{Name: "Read Maintenance Document", Slug: "maintenance_document:read", Module: "maintenance_document"},
		{Name: "Delete Maintenance Document", Slug: "maintenance_document:delete", Module: "maintenance_document"},
		{Name: "Create Document", Slug: "document:create", Module: "document"},
		{Name: "Read Document", Slug: "document:read", Module: "document"},
		{Name: "Update Document", Slug: "document:update", Module: "document"},
		{Name: "Delete Document", Slug: "document:delete", Module: "document"},
		{Name: "Send Notification", Slug: "notification:send", Module: "notification"},
		{Name: "Create Asset Mutation", Slug: "asset_mutation:create", Module: "asset_mutation"},
		{Name: "Read Asset Mutation", Slug: "asset_mutation:read", Module: "asset_mutation"},
		{Name: "Create Stock Opname", Slug: "stock_opname:create", Module: "stock_opname"},
		{Name: "Read Stock Opname", Slug: "stock_opname:read", Module: "stock_opname"},
		{Name: "Update Stock Opname", Slug: "stock_opname:update", Module: "stock_opname"},
		{Name: "Close Stock Opname", Slug: "stock_opname:close", Module: "stock_opname"},
		{Name: "Create Complaint", Slug: "complaint:create", Module: "complaint"},
		{Name: "Read Complaint", Slug: "complaint:read", Module: "complaint"},
		{Name: "Update Complaint", Slug: "complaint:update", Module: "complaint"},
		{Name: "Delete Complaint", Slug: "complaint:delete", Module: "complaint"},
		{Name: "Read Permission", Slug: "permission:read", Module: "permission"},
		{Name: "Create Role", Slug: "role:create", Module: "role"},
		{Name: "Read Role", Slug: "role:read", Module: "role"},
		{Name: "Read Role Permissions", Slug: "role:read_permissions", Module: "role"},
		{Name: "Update Role", Slug: "role:update", Module: "role"},
		{Name: "Delete Role", Slug: "role:delete", Module: "role"},
		{Name: "Assign Role Permission", Slug: "role:assign_permissions", Module: "role"},
		{Name: "Revoke Role Permission", Slug: "role:revoke_permissions", Module: "role"},
		{Name: "Create User", Slug: "user:create", Module: "user"},
		{Name: "Read User", Slug: "user:read", Module: "user"},
		{Name: "Update User", Slug: "user:update", Module: "user"},
		{Name: "Delete User", Slug: "user:delete", Module: "user"},
		{Name: "Reset User Password", Slug: "user:reset_password", Module: "user"},
		{Name: "Read Audit Trail", Slug: "audit:read", Module: "audit"},
		{Name: "Read Tenant", Slug: "tenant:read", Module: "tenant"},
		{Name: "Create Tenant", Slug: "tenant:create", Module: "tenant"},
		{Name: "Update Tenant", Slug: "tenant:update", Module: "tenant"},
		{Name: "Update Tenant Status", Slug: "tenant:update_status", Module: "tenant"},
	}
	if err := permissionRepo.EnsurePermissions(defaultPermissions); err != nil {
		return err
	}

	// f. Assign all permissions to Tenant Admin
	if tenantAdminRole.ID != 0 {
		var permissionIDs []int
		permDB := db.Session(&gorm.Session{NewDB: true}).Model(&domain.Permission{})
		if err := permDB.Pluck("id", &permissionIDs).Error; err == nil {
			if err := permissionRepo.AssignPermissions(tenantAdminRole.ID, permissionIDs); err != nil {
				return err
			}
		}
	}

	return nil
}

func getenvDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
