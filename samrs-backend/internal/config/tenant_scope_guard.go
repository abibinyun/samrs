package config

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	skipTenantScopeKey  = "skip_tenant_scope"
	tenantScopeOKKey    = "tenant_scope_ok"
)

func WithSkipTenantScope(db *gorm.DB) *gorm.DB {
	return db.Set(skipTenantScopeKey, true)
}

func MarkTenantScope(db *gorm.DB) *gorm.DB {
	return db.Set(tenantScopeOKKey, true)
}

func registerTenantScopeGuard(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("tenant_scope_guard:query", tenantScopeGuard)
	db.Callback().Update().Before("gorm:update").Register("tenant_scope_guard:update", tenantScopeGuard)
	db.Callback().Delete().Before("gorm:delete").Register("tenant_scope_guard:delete", tenantScopeGuard)
}

func tenantScopeGuard(db *gorm.DB) {
	if db == nil || db.Statement == nil || db.Statement.Schema == nil {
		return
	}
	if _, ok := db.Statement.Settings.Load(skipTenantScopeKey); ok {
		return
	}
	if _, ok := db.Statement.Settings.Load(tenantScopeOKKey); ok {
		return
	}
	if db.Statement.Schema.LookUpField("TenantID") == nil {
		return
	}
	if hasTenantCondition(db.Statement.Clauses) {
		return
	}
	db.AddError(errors.New("tenant scope is required for this query"))
}

func hasTenantCondition(clauses map[string]clause.Clause) bool {
	whereClause, ok := clauses["WHERE"]
	if !ok {
		return false
	}
	where, ok := whereClause.Expression.(clause.Where)
	if !ok {
		return false
	}
	return exprsContainTenant(where.Exprs)
}

func exprsContainTenant(exprs []clause.Expression) bool {
	for _, expr := range exprs {
		if exprContainsTenant(expr) {
			return true
		}
	}
	return false
}

func exprContainsTenant(expr clause.Expression) bool {
	switch e := expr.(type) {
	case clause.Eq:
		return isTenantColumn(e.Column)
	case clause.Neq:
		return isTenantColumn(e.Column)
	case clause.IN:
		return isTenantColumn(e.Column)
	case clause.AndConditions:
		return exprsContainTenant(e.Exprs)
	case clause.OrConditions:
		return exprsContainTenant(e.Exprs)
	case clause.NotConditions:
		return exprsContainTenant(e.Exprs)
	case clause.Expr:
		return strings.Contains(strings.ToLower(e.SQL), "tenant_id")
	default:
		return false
	}
}

func isTenantColumn(column interface{}) bool {
	switch col := column.(type) {
	case string:
		return strings.Contains(strings.ToLower(col), "tenant_id")
	case clause.Column:
		return strings.Contains(strings.ToLower(col.Name), "tenant_id")
	default:
		return false
	}
}
