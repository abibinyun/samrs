package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func registerAuditTrailHooks(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("audit_trail:after_create", auditAfterCreate)
	db.Callback().Update().Before("gorm:update").Register("audit_trail:before_update", auditBeforeUpdate)
	db.Callback().Update().After("gorm:update").Register("audit_trail:after_update", auditAfterUpdate)
	db.Callback().Delete().Before("gorm:delete").Register("audit_trail:before_delete", auditBeforeDelete)
	db.Callback().Delete().After("gorm:delete").Register("audit_trail:after_delete", auditAfterDelete)
}

func auditAfterCreate(db *gorm.DB) {
	createAuditEntry(db, "CREATE", nil, db.Statement.Dest)
}

func auditBeforeUpdate(db *gorm.DB) {
	oldData, ok := fetchOldRecord(db)
	if ok {
		db.Statement.Settings.Store("audit_old_data", oldData)
	}
}

func auditAfterUpdate(db *gorm.DB) {
	oldData, _ := db.Statement.Settings.Load("audit_old_data")
	createAuditEntry(db, "UPDATE", oldData, db.Statement.Dest)
}

func auditBeforeDelete(db *gorm.DB) {
	oldData, ok := fetchOldRecord(db)
	if ok {
		db.Statement.Settings.Store("audit_old_data", oldData)
	}
}

func auditAfterDelete(db *gorm.DB) {
	oldData, _ := db.Statement.Settings.Load("audit_old_data")
	createAuditEntry(db, "DELETE", oldData, nil)
}

func createAuditEntry(db *gorm.DB, defaultAction string, oldData, newData interface{}) {
	if !shouldAudit(db) {
		return
	}

	meta, ok := getAuditMeta(db)
	if !ok {
		return
	}

	recordID := buildRecordID(db.Statement)
	if recordID == "" {
		if override, ok := db.Statement.Settings.Load("audit_record_id"); ok {
			recordID = fmt.Sprint(override)
		}
	}
	if recordID == "" {
		return
	}

	action := defaultAction
	if override, ok := db.Statement.Settings.Load("audit_action"); ok {
		if actionStr, ok := override.(string); ok && strings.TrimSpace(actionStr) != "" {
			action = actionStr
		}
	}

	oldJSON, err := toJSON(sanitizeData(oldData, db.Statement))
	if err != nil {
		db.AddError(err)
		return
	}

	newJSON, err := toJSON(sanitizeData(newData, db.Statement))
	if err != nil {
		db.AddError(err)
		return
	}

	audit := &domain.AuditTrail{
		TenantID:  meta.TenantID,
		UserID:    meta.UserID,
		Action:    action,
		TableName: db.Statement.Table,
		RecordID:  recordID,
		OldData:   oldJSON,
		NewData:   newJSON,
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
	}

	if err := db.Session(&gorm.Session{NewDB: true}).Set("skip_audit", true).Create(audit).Error; err != nil {
		db.AddError(err)
	}
}

type auditMeta struct {
	TenantID  uuid.UUID
	UserID    uuid.UUID
	IP        string
	UserAgent string
}

func getAuditMeta(db *gorm.DB) (auditMeta, bool) {
	raw, ok := db.Statement.Settings.Load("audit_meta")
	if !ok {
		return auditMeta{}, false
	}

	metaMap, ok := raw.(map[string]interface{})
	if !ok {
		return auditMeta{}, false
	}

	tenantID, ok := extractUUID(metaMap["tenant_id"])
	if !ok {
		return auditMeta{}, false
	}
	userID, ok := extractUUID(metaMap["user_id"])
	if !ok {
		return auditMeta{}, false
	}

	return auditMeta{
		TenantID:  tenantID,
		UserID:    userID,
		IP:        fmt.Sprint(metaMap["ip"]),
		UserAgent: fmt.Sprint(metaMap["user_agent"]),
	}, true
}

func extractUUID(value interface{}) (uuid.UUID, bool) {
	switch v := value.(type) {
	case uuid.UUID:
		if v == uuid.Nil {
			return uuid.UUID{}, false
		}
		return v, true
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil || parsed == uuid.Nil {
			return uuid.UUID{}, false
		}
		return parsed, true
	default:
		parsed, err := uuid.Parse(fmt.Sprint(v))
		if err != nil || parsed == uuid.Nil {
			return uuid.UUID{}, false
		}
		return parsed, true
	}
}

func shouldAudit(db *gorm.DB) bool {
	if db == nil || db.Statement == nil || db.Statement.Schema == nil {
		return false
	}
	if _, ok := db.Statement.Settings.Load("skip_audit"); ok {
		return false
	}
	table := strings.ToLower(db.Statement.Table)
	if table == "audit_trails" || db.Statement.Schema.Name == "AuditTrail" {
		return false
	}
	_, ok := db.Statement.Settings.Load("audit_meta")
	return ok
}

func buildRecordID(stmt *gorm.Statement) string {
	if stmt == nil || stmt.Schema == nil {
		return ""
	}
	if stmt.ReflectValue.IsValid() {
		kind := stmt.ReflectValue.Kind()
		if kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map {
			return ""
		}
	}
	pkFields := stmt.Schema.PrimaryFields
	if len(pkFields) == 0 {
		return ""
	}

	values := make([]string, 0, len(pkFields))
	for _, field := range pkFields {
		val, zero := field.ValueOf(stmt.Context, stmt.ReflectValue)
		if zero || isZeroValue(val) {
			if stmt.Dest == nil {
				return ""
			}
			destValue := reflect.ValueOf(stmt.Dest)
			if destValue.Kind() == reflect.Ptr {
				destValue = destValue.Elem()
			}
			if !destValue.IsValid() || destValue.Kind() != reflect.Struct {
				return ""
			}
			val, zero = field.ValueOf(stmt.Context, destValue)
			if zero || isZeroValue(val) {
				return ""
			}
		}
		if len(pkFields) == 1 {
			return fmt.Sprint(val)
		}
		values = append(values, fmt.Sprintf("%s=%v", field.DBName, val))
	}
	return strings.Join(values, ",")
}

func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}
	return reflect.ValueOf(value).IsZero()
}

func fetchOldRecord(db *gorm.DB) (interface{}, bool) {
	if db == nil || db.Statement == nil || db.Statement.Schema == nil {
		return nil, false
	}

	pkConditions, ok := primaryKeyConditions(db.Statement)
	if !ok {
		return nil, false
	}

	modelValue := reflect.New(db.Statement.Schema.ModelType).Interface()
	query := db.Session(&gorm.Session{NewDB: true}).Model(modelValue)

	tenantID, ok := tenantIDFromStatement(db.Statement, &db.Statement.Settings)
	if db.Statement.Schema.LookUpField("TenantID") != nil {
		if !ok {
			return nil, false
		}
		query = query.Where("tenant_id = ?", tenantID)
	}

	for column, value := range pkConditions {
		query = query.Where(fmt.Sprintf("%s = ?", column), value)
	}

	if err := query.First(modelValue).Error; err != nil {
		return nil, false
	}
	return modelValue, true
}

func primaryKeyConditions(stmt *gorm.Statement) (map[string]interface{}, bool) {
	if stmt == nil || stmt.Schema == nil {
		return nil, false
	}
	conditions := make(map[string]interface{}, len(stmt.Schema.PrimaryFields))
	for _, field := range stmt.Schema.PrimaryFields {
		val, zero := field.ValueOf(stmt.Context, stmt.ReflectValue)
		if zero {
			return nil, false
		}
		conditions[field.DBName] = val
	}
	return conditions, len(conditions) > 0
}

func tenantIDFromStatement(stmt *gorm.Statement, settings *sync.Map) (uuid.UUID, bool) {
	if stmt == nil || stmt.Schema == nil {
		return uuid.UUID{}, false
	}
	field := stmt.Schema.LookUpField("TenantID")
	if field != nil {
		val, ok := field.ValueOf(stmt.Context, stmt.ReflectValue)
		if ok {
			if tenantID, ok := extractUUID(val); ok {
				return tenantID, true
			}
		}
	}

	if settings != nil {
		if meta, ok := settings.Load("audit_meta"); ok {
			if metaMap, ok := meta.(map[string]interface{}); ok {
				if tenantID, ok := extractUUID(metaMap["tenant_id"]); ok {
					return tenantID, true
				}
			}
		}
	}
	return uuid.UUID{}, false
}

func sanitizeData(value interface{}, stmt *gorm.Statement) interface{} {
	if value == nil {
		return nil
	}
	data := toMap(value)
	if data == nil {
		return value
	}

	if stmt != nil && stmt.Schema != nil && stmt.Schema.Name == "User" {
		delete(data, "password_hash")
		delete(data, "PasswordHash")
		delete(data, "password")
	}
	return data
}

func toMap(value interface{}) map[string]interface{} {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	var out map[string]interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func toJSON(value interface{}) (datatypes.JSON, error) {
	if value == nil {
		return nil, nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}
