package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentUsecase interface {
	CreateDocument(input DocumentInput, files []DocumentFileInput) (*domain.Document, []domain.DocumentFile, error)
	AddFiles(tenantID uuid.UUID, documentID uint, uploaderID uuid.UUID, files []DocumentFileInput) ([]domain.DocumentFile, error)
	ListDocuments(tenantID uuid.UUID, filter repository.DocumentFilter) ([]domain.Document, int64, error)
	GetDocumentByID(tenantID uuid.UUID, id uint) (*domain.Document, error)
	UpdateDocument(tenantID uuid.UUID, id uint, input DocumentUpdateInput) (*domain.Document, error)
	ListFiles(tenantID uuid.UUID, documentID uint) ([]domain.DocumentFile, error)
	GetFileByID(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error)
	DeleteFile(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error)
	DeleteDocument(tenantID uuid.UUID, id uint) ([]domain.DocumentFile, error)
}

type DocumentInput struct {
	TenantID    uuid.UUID
	Title       string
	DocType     string
	Description string
	UploadedBy  uuid.UUID
}

type DocumentUpdateInput struct {
	Title       string
	DocType     string
	Description string
}

type DocumentFileInput struct {
	Filename string
	FilePath string
	MimeType string
	Size     int64
}

type documentUsecase struct {
	repo     repository.DocumentRepository
	fileRepo repository.DocumentFileRepository
	db       *gorm.DB
}

func NewDocumentUsecase(
	repo repository.DocumentRepository,
	fileRepo repository.DocumentFileRepository,
	db *gorm.DB,
) DocumentUsecase {
	return &documentUsecase{repo: repo, fileRepo: fileRepo, db: db}
}

func (u *documentUsecase) CreateDocument(input DocumentInput, files []DocumentFileInput) (*domain.Document, []domain.DocumentFile, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, nil, util.ErrValidation("judul dokumen wajib diisi")
	}
	if !domain.IsValidDocumentType(input.DocType) {
		input.DocType = domain.DocumentTypeOther
	}
	if len(files) == 0 {
		return nil, nil, util.ErrValidation("file dokumen wajib diunggah")
	}

	doc := &domain.Document{
		TenantID:    input.TenantID,
		Title:       title,
		DocType:     input.DocType,
		Description: strings.TrimSpace(input.Description),
	}

	var createdFiles []domain.DocumentFile
	err := u.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewDocumentRepository(tx)
		txFileRepo := repository.NewDocumentFileRepository(tx)

		if err := txRepo.Create(doc); err != nil {
			return err
		}

		createdFiles = buildDocumentFiles(input.TenantID, input.UploadedBy, doc.ID, files)
		if err := txFileRepo.CreateMany(createdFiles); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return doc, createdFiles, nil
}

func (u *documentUsecase) AddFiles(tenantID uuid.UUID, documentID uint, uploaderID uuid.UUID, files []DocumentFileInput) ([]domain.DocumentFile, error) {
	if len(files) == 0 {
		return nil, util.ErrValidation("file dokumen wajib diunggah")
	}
	if _, err := u.repo.FindByID(tenantID, documentID); err != nil {
		return nil, err
	}

	createdFiles := buildDocumentFiles(tenantID, uploaderID, documentID, files)
	if err := u.fileRepo.CreateMany(createdFiles); err != nil {
		return nil, err
	}
	return createdFiles, nil
}

func (u *documentUsecase) ListDocuments(tenantID uuid.UUID, filter repository.DocumentFilter) ([]domain.Document, int64, error) {
	return u.repo.List(tenantID, filter)
}

func (u *documentUsecase) GetDocumentByID(tenantID uuid.UUID, id uint) (*domain.Document, error) {
	return u.repo.FindByID(tenantID, id)
}

func (u *documentUsecase) UpdateDocument(tenantID uuid.UUID, id uint, input DocumentUpdateInput) (*domain.Document, error) {
	doc, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title != "" {
		doc.Title = title
	}

	docType := strings.TrimSpace(input.DocType)
	if docType != "" {
		if !domain.IsValidDocumentType(docType) {
			return nil, util.ErrValidation("doc_type tidak valid")
		}
		doc.DocType = docType
	}

	if input.Description != "" {
		doc.Description = strings.TrimSpace(input.Description)
	}

	if err := u.repo.Update(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (u *documentUsecase) ListFiles(tenantID uuid.UUID, documentID uint) ([]domain.DocumentFile, error) {
	if _, err := u.repo.FindByID(tenantID, documentID); err != nil {
		return nil, err
	}
	return u.fileRepo.ListByDocument(tenantID, documentID)
}

func (u *documentUsecase) GetFileByID(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error) {
	return u.fileRepo.FindByID(tenantID, id)
}

func (u *documentUsecase) DeleteFile(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error) {
	file, err := u.fileRepo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := u.fileRepo.Delete(file); err != nil {
		return nil, err
	}
	return file, nil
}

func (u *documentUsecase) DeleteDocument(tenantID uuid.UUID, id uint) ([]domain.DocumentFile, error) {
	doc, err := u.repo.FindByID(tenantID, id)
	if err != nil {
		return nil, err
	}

	var files []domain.DocumentFile
	if err := u.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewDocumentRepository(tx)
		txFileRepo := repository.NewDocumentFileRepository(tx)

		list, err := txFileRepo.ListByDocument(tenantID, doc.ID)
		if err != nil {
			return err
		}
		files = list

		if err := txFileRepo.DeleteByDocument(tenantID, doc.ID); err != nil {
			return err
		}
		return txRepo.Delete(doc)
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func buildDocumentFiles(tenantID uuid.UUID, uploaderID uuid.UUID, documentID uint, inputs []DocumentFileInput) []domain.DocumentFile {
	files := make([]domain.DocumentFile, 0, len(inputs))
	for _, input := range inputs {
		files = append(files, domain.DocumentFile{
			TenantID:   tenantID,
			DocumentID: documentID,
			UploadedBy: uploaderID,
			Filename:   input.Filename,
			FilePath:   input.FilePath,
			MimeType:   input.MimeType,
			Size:       input.Size,
		})
	}
	return files
}
