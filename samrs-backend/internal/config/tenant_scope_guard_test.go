package config

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func newStatement(t *testing.T, model interface{}) *gorm.Statement {
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
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	stmt.ReflectValue = value
	return stmt
}

func TestTenantScopeGuard(t *testing.T) {
	tenantID := uuid.New()
	room := &domain.Room{ID: uuid.New(), TenantID: tenantID}
	permission := &domain.Permission{ID: 1, Name: "View", Slug: "view"}

	tests := []struct {
		name     string
		model    interface{}
		clauses  map[string]clause.Clause
		settings func(stmt *gorm.Statement)
		wantErr  bool
	}{
		{
			name:    "missing tenant clause",
			model:   room,
			wantErr: true,
		},
		{
			name:  "tenant clause with eq",
			model: room,
			clauses: map[string]clause.Clause{
				"WHERE": {
					Expression: clause.Where{Exprs: []clause.Expression{
						clause.Eq{Column: "tenant_id", Value: tenantID},
					}},
				},
			},
		},
		{
			name:  "tenant clause with raw expr",
			model: room,
			clauses: map[string]clause.Clause{
				"WHERE": {
					Expression: clause.Where{Exprs: []clause.Expression{
						clause.Expr{SQL: "tenant_id = ?"},
					}},
				},
			},
		},
		{
			name:  "skip tenant scope flag",
			model: room,
			settings: func(stmt *gorm.Statement) {
				stmt.Settings.Store(skipTenantScopeKey, true)
			},
		},
		{
			name:  "tenant scope marked",
			model: room,
			settings: func(stmt *gorm.Statement) {
				stmt.Settings.Store(tenantScopeOKKey, true)
			},
		},
		{
			name:  "schema without tenant id",
			model: permission,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := newStatement(t, tt.model)
			stmt.Clauses = tt.clauses
			if tt.settings != nil {
				tt.settings(stmt)
			}
			db := &gorm.DB{Config: &gorm.Config{}, Statement: stmt}
			stmt.DB = db

			tenantScopeGuard(db)

			if tt.wantErr && db.Error == nil {
				t.Fatalf("expected tenant scope error, got nil")
			}
			if !tt.wantErr && db.Error != nil {
				t.Fatalf("expected no error, got %v", db.Error)
			}
		})
	}
}

func TestTenantScopeGuardNilDB(t *testing.T) {
	tenantScopeGuard(nil)
}
