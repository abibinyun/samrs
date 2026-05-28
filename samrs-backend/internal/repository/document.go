package repository

import (
	documentrepo "samrs-backend/internal/repository/document"

	"gorm.io/gorm"
)

type DocumentFilter = documentrepo.DocumentFilter

type DocumentRepository = documentrepo.DocumentRepository

type DocumentFileRepository = documentrepo.DocumentFileRepository

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return documentrepo.NewDocumentRepository(db)
}

func NewDocumentFileRepository(db *gorm.DB) DocumentFileRepository {
	return documentrepo.NewDocumentFileRepository(db)
}
