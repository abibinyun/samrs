package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type MaintenanceDocumentUsecase interface {
	CreateDocument(input MaintenanceDocumentInput) (*domain.MaintenanceDocument, error)
	ListDocuments(tenantID uuid.UUID, scheduleID uint) ([]domain.MaintenanceDocument, error)
	GetDocumentByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceDocument, error)
	DeleteDocument(tenantID uuid.UUID, id uint) error
}

type MaintenanceDocumentInput struct {
	TenantID   uuid.UUID
	ScheduleID uint
	UploadedBy uuid.UUID
	DocType    string
	Filename   string
	FilePath   string
	MimeType   string
	Size       int64
}

type maintenanceDocumentUsecase struct {
	repo         repository.MaintenanceDocumentRepository
	scheduleRepo repository.MaintenanceScheduleRepository
}

func NewMaintenanceDocumentUsecase(
	repo repository.MaintenanceDocumentRepository,
	sr repository.MaintenanceScheduleRepository,
) MaintenanceDocumentUsecase {
	return &maintenanceDocumentUsecase{repo: repo, scheduleRepo: sr}
}

func (u *maintenanceDocumentUsecase) CreateDocument(input MaintenanceDocumentInput) (*domain.MaintenanceDocument, error) {
	if input.ScheduleID == 0 {
		return nil, util.ErrValidation("schedule ID wajib diisi")
	}
	if input.UploadedBy == uuid.Nil {
		return nil, util.ErrValidation("uploaded_by wajib diisi")
	}
	if strings.TrimSpace(input.Filename) == "" || strings.TrimSpace(input.FilePath) == "" {
		return nil, util.ErrValidation("file wajib diisi")
	}

	docType := strings.TrimSpace(input.DocType)
	if docType == "" {
		docType = domain.MaintenanceDocOther
	}
	switch docType {
	case domain.MaintenanceDocCertificate, domain.MaintenanceDocReport, domain.MaintenanceDocOther:
	default:
		return nil, util.ErrValidation("doc_type tidak valid")
	}

	schedule, err := u.scheduleRepo.FindByID(input.TenantID, input.ScheduleID)
	if err != nil {
		return nil, util.ErrNotFound("jadwal tidak ditemukan atau akses ditolak")
	}

	doc := &domain.MaintenanceDocument{
		TenantID:   input.TenantID,
		ScheduleID: input.ScheduleID,
		AssetID:    schedule.AssetID,
		UploadedBy: input.UploadedBy,
		DocType:    docType,
		Filename:   strings.TrimSpace(input.Filename),
		FilePath:   strings.TrimSpace(input.FilePath),
		MimeType:   strings.TrimSpace(input.MimeType),
		Size:       input.Size,
	}

	if err := u.repo.Create(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (u *maintenanceDocumentUsecase) ListDocuments(tenantID uuid.UUID, scheduleID uint) ([]domain.MaintenanceDocument, error) {
	_, err := u.scheduleRepo.FindByID(tenantID, scheduleID)
	if err != nil {
		return nil, err
	}
	return u.repo.ListBySchedule(tenantID, scheduleID)
}

func (u *maintenanceDocumentUsecase) GetDocumentByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceDocument, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *maintenanceDocumentUsecase) DeleteDocument(tenantID uuid.UUID, id uint) error {
	doc, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(doc)
}
