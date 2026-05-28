package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DocumentTypeSOP         = "sop"
	DocumentTypeManual      = "manual"
	DocumentTypeCertificate = "certificate"
	DocumentTypeOther       = "other"
)

type Document struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;index;not null;index:idx_document_tenant_created" json:"tenant_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	DocType     string         `gorm:"size:50;not null;index" json:"doc_type"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"index:idx_document_tenant_created" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Files []DocumentFile `gorm:"foreignKey:DocumentID" json:"files,omitempty"`
}

type DocumentFile struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TenantID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
	DocumentID uint           `gorm:"index;not null" json:"document_id"`
	UploadedBy uuid.UUID      `gorm:"type:uuid;index;not null" json:"uploaded_by"`
	Filename   string         `gorm:"size:255;not null" json:"filename"`
	FilePath   string         `gorm:"size:500;not null" json:"file_path"`
	MimeType   string         `gorm:"size:100" json:"mime_type,omitempty"`
	Size       int64          `json:"size"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Document Document `gorm:"foreignKey:DocumentID" json:"-"`
	Uploader User     `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`

	FileURL string `gorm:"-" json:"file_url,omitempty"`
}

func IsValidDocumentType(docType string) bool {
	switch docType {
	case DocumentTypeSOP, DocumentTypeManual, DocumentTypeCertificate, DocumentTypeOther:
		return true
	default:
		return false
	}
}
