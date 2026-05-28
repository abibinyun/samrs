package config

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func newAuditStatement(t *testing.T, model interface{}) *gorm.Statement {
	t.Helper()
	stmt := &gorm.Statement{
		DB:      &gorm.DB{Config: &gorm.Config{}},
		Context: context.Background(),
	}
	parsed, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	stmt.Schema = parsed
	stmt.Table = parsed.Table
	stmt.Dest = model
	stmt.ReflectValue = reflect.ValueOf(model)
	return stmt
}

func TestBuildRecordID_EmptyWhenZero(t *testing.T) {
	room := &domain.Room{}
	stmt := newAuditStatement(t, room)
	if got := buildRecordID(stmt); got != "" {
		t.Fatalf("expected empty record id, got %s", got)
	}
}

func TestBuildRecordID_ReturnsValue(t *testing.T) {
	roomID := uuid.New()
	room := &domain.Room{ID: roomID}
	stmt := newAuditStatement(t, room)
	stmt.ReflectValue = reflect.ValueOf(room).Elem()
	if got := buildRecordID(stmt); got != roomID.String() {
		t.Fatalf("expected record id %s, got %s", roomID, got)
	}
}

func TestSanitizeData_UserRemovesPassword(t *testing.T) {
	user := &domain.User{
		ID:           uuid.New(),
		TenantID:     uuid.New(),
		Username:     "tester",
		PasswordHash: "secret",
	}
	stmt := newAuditStatement(t, user)
	sanitized := sanitizeData(user, stmt)
	data, ok := sanitized.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map data after sanitize")
	}
	if _, exists := data["password_hash"]; exists {
		t.Fatalf("password_hash should be removed from audit data")
	}
	if _, exists := data["PasswordHash"]; exists {
		t.Fatalf("PasswordHash should be removed from audit data")
	}
}

func TestExtractUUID(t *testing.T) {
	id := uuid.New()
	if got, ok := extractUUID(id.String()); !ok || got != id {
		t.Fatalf("expected uuid %s from string", id)
	}
}
