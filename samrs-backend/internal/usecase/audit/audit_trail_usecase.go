package usecase

import (
	"encoding/json"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditTrailUsecase interface {
	Log(action, tableName, recordID string, oldData interface{}, newData interface{}, meta AuditMeta) error
	ListByTenant(tenantID uuid.UUID, filter repository.AuditTrailFilter) ([]domain.AuditTrail, int64, error)
	GetByID(tenantID uuid.UUID, id uint) (*domain.AuditTrail, error)
}

type AuditMeta struct {
	TenantID  uuid.UUID
	UserID    uuid.UUID
	IP        string
	UserAgent string
}

type auditTrailUsecase struct {
	repo repository.AuditTrailRepository
}

func NewAuditTrailUsecase(repo repository.AuditTrailRepository) AuditTrailUsecase {
	return &auditTrailUsecase{repo}
}

func (u *auditTrailUsecase) Log(action, tableName, recordID string, oldData interface{}, newData interface{}, meta AuditMeta) error {
	if meta.TenantID == uuid.Nil || meta.UserID == uuid.Nil {
		return util.ErrValidation("metadata audit tidak lengkap")
	}

	oldJSON, err := toJSON(oldData)
	if err != nil {
		return err
	}

	newJSON, err := toJSON(newData)
	if err != nil {
		return err
	}

	audit := &domain.AuditTrail{
		TenantID:  meta.TenantID,
		UserID:    meta.UserID,
		Action:    action,
		TableName: tableName,
		RecordID:  recordID,
		OldData:   oldJSON,
		NewData:   newJSON,
		IP:        meta.IP,
		UserAgent: meta.UserAgent,
	}

	return u.repo.Create(audit)
}

func (u *auditTrailUsecase) ListByTenant(tenantID uuid.UUID, filter repository.AuditTrailFilter) ([]domain.AuditTrail, int64, error) {
	return u.repo.ListByTenant(tenantID, filter)
}

func (u *auditTrailUsecase) GetByID(tenantID uuid.UUID, id uint) (*domain.AuditTrail, error) {
	return u.repo.FindByID(tenantID, id)
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
