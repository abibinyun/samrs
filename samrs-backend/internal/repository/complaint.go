package repository

import (
	complaintrepo "samrs-backend/internal/repository/complaint"

	"gorm.io/gorm"
)

type ComplaintFilter = complaintrepo.ComplaintFilter

type ComplaintRepository = complaintrepo.ComplaintRepository

func NewComplaintRepository(db *gorm.DB) ComplaintRepository {
	return complaintrepo.NewComplaintRepository(db)
}
