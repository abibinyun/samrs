package config

import (
	"fmt"
	"log"
	"os"
	"samrs-backend/internal/domain"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() *gorm.DB {
	// Load .env (optional in Docker)
	_ = godotenv.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             1 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	fmt.Println("Berhasil terkoneksi ke Database SAMRS!")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatal("Gagal memastikan uuid-ossp extension:", err)
	}

	// Jalankan AutoMigrate via GORM
	err = db.AutoMigrate(
		&domain.Role{},
		&domain.Tenant{},
		&domain.User{},
		&domain.Permission{},
		&domain.RolePermission{},
		&domain.AuditTrail{},
		&domain.Room{},
		&domain.Bed{},
		&domain.Category{},
		&domain.Asset{},
		&domain.AssetEvent{},
		&domain.AssetMutation{},
		&domain.MaintenanceSchedule{},
		&domain.StockOpnameSession{},
		&domain.StockOpnameItem{},
		&domain.MaintenanceDocument{},
		&domain.Document{},
		&domain.DocumentFile{},
		&domain.Vendor{},
		&domain.AssetBrand{},
		&domain.AssetModel{},
		&domain.AssetStatus{},
		&domain.Complaint{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	if err := ensurePartialUniqueIndexes(db); err != nil {
		log.Fatal("Gagal membuat partial unique index:", err)
	}

	registerTenantScopeGuard(db)
	registerAuditTrailHooks(db)

	return db
}

func ensurePartialUniqueIndexes(db *gorm.DB) error {
	statements := []string{
		`DROP INDEX IF EXISTS idx_room_code_tenant`,
		`DROP INDEX IF EXISTS idx_bed_code_room_tenant`,
		`DROP INDEX IF EXISTS idx_category_slug_tenant`,
		`DROP INDEX IF EXISTS idx_asset_code_tenant`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_rooms_tenant_code_active ON rooms (tenant_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_beds_tenant_room_code_active ON beds (tenant_id, room_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_tenant_slug_active ON categories (tenant_id, slug) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_assets_tenant_code_active ON assets (tenant_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_vendors_tenant_code_active ON vendors (tenant_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_brands_tenant_code_active ON asset_brands (tenant_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_models_tenant_code_active ON asset_models (tenant_id, code) WHERE deleted_at IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_statuses_tenant_code_active ON asset_statuses (tenant_id, code) WHERE deleted_at IS NULL`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}
