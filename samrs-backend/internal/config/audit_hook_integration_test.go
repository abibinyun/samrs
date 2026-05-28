package config

import (
	"os"
	"testing"
	"time"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuditHookWritesAuditTrail(t *testing.T) {
	dsn := os.Getenv("SAMRS_TEST_DSN")
	if dsn == "" {
		t.Skip("SAMRS_TEST_DSN tidak diset, skip integration test")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}

	if err := db.AutoMigrate(&domain.Room{}, &domain.AuditTrail{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	registerTenantScopeGuard(db)
	registerAuditTrailHooks(db)

	tenantID := uuid.New()
	userID := uuid.New()
	roomID := uuid.New()
	meta := map[string]interface{}{
		"tenant_id":  tenantID,
		"user_id":    userID,
		"ip":         "127.0.0.1",
		"user_agent": "test-agent",
	}

	defer func() {
		db.Where("tenant_id = ?", tenantID).Delete(&domain.AuditTrail{})
		db.Where("tenant_id = ?", tenantID).Delete(&domain.Room{})
	}()

	if err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Set("audit_meta", meta)
		room := domain.Room{
			ID:       roomID,
			TenantID: tenantID,
			Name:     "Ruang A",
			Code:     "R-001",
		}
		return tx.Create(&room).Error
	}); err != nil {
		t.Fatalf("failed to create room: %v", err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Set("audit_meta", meta).Set("audit_record_id", roomID)
		room := domain.Room{
			ID:       roomID,
			TenantID: tenantID,
		}
		return tx.Model(&room).
			Where("tenant_id = ? AND id = ?", tenantID, roomID).
			Updates(map[string]interface{}{"name": "Ruang A-Updated"}).Error
	}); err != nil {
		t.Fatalf("failed to update room: %v", err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Set("audit_meta", meta)
		room := domain.Room{
			ID:       roomID,
			TenantID: tenantID,
		}
		return tx.Where("tenant_id = ? AND id = ?", tenantID, roomID).Delete(&room).Error
	}); err != nil {
		t.Fatalf("failed to delete room: %v", err)
	}

	var count int64
	if err := db.Model(&domain.AuditTrail{}).
		Where("tenant_id = ? AND table_name = ? AND record_id = ?", tenantID, "rooms", roomID.String()).
		Count(&count).Error; err != nil {
		t.Fatalf("failed to count audit: %v", err)
	}

	if count < 3 {
		t.Fatalf("expected >= 3 audit entries, got %d", count)
	}

	var latest domain.AuditTrail
	if err := db.Where("tenant_id = ? AND table_name = ? AND record_id = ?", tenantID, "rooms", roomID.String()).
		Order("created_at desc").First(&latest).Error; err != nil {
		t.Fatalf("failed to fetch latest audit: %v", err)
	}

	if latest.UserID != userID {
		t.Fatalf("expected user_id %s, got %s", userID, latest.UserID)
	}
	if latest.IP != "127.0.0.1" {
		t.Fatalf("expected ip 127.0.0.1, got %s", latest.IP)
	}
	if latest.UserAgent != "test-agent" {
		t.Fatalf("expected user_agent test-agent, got %s", latest.UserAgent)
	}
	if time.Since(latest.CreatedAt) > time.Minute {
		t.Fatalf("audit timestamp too old")
	}
}
