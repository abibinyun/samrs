package usecase

import (
	"strings"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/pkg/util"

	"github.com/google/uuid"
)

type ComplaintUsecase interface {
	CreateComplaint(input CreateComplaintInput) (*domain.Complaint, error)
	GetAllComplaints(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, int64, error)
	GetComplaintByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Complaint, error)
	UpdateComplaint(input UpdateComplaintInput) (*domain.Complaint, error)
	DeleteComplaint(tenantID uuid.UUID, id uuid.UUID) error
}

type CreateComplaintInput struct {
	TenantID    uuid.UUID
	AssetID     uuid.UUID
	ReportedBy  uuid.UUID
	AssignedTo  *uuid.UUID
	Title       string
	Description string
}

type UpdateComplaintInput struct {
	TenantID       uuid.UUID
	ID             uuid.UUID
	Title          string
	Description    string
	Status         string
	AssignedTo     *uuid.UUID
	AssignedToSet  bool
	ResolutionNote string
}

type complaintUsecase struct {
	complaintRepo repository.ComplaintRepository
	assetRepo     repository.AssetRepository
	userRepo      repository.UserRepository
}

func NewComplaintUsecase(cr repository.ComplaintRepository, ar repository.AssetRepository, ur repository.UserRepository) ComplaintUsecase {
	return &complaintUsecase{cr, ar, ur}
}

func (u *complaintUsecase) CreateComplaint(input CreateComplaintInput) (*domain.Complaint, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, util.ErrValidation("judul complaint wajib diisi")
	}
	if strings.TrimSpace(input.Description) == "" {
		return nil, util.ErrValidation("deskripsi complaint wajib diisi")
	}

	if _, err := u.assetRepo.FindByIDAndTenant(input.AssetID, input.TenantID); err != nil {
		return nil, util.ErrNotFound("aset tidak ditemukan atau akses ditolak")
	}
	if _, err := u.userRepo.FindByIDAndTenant(input.ReportedBy, input.TenantID); err != nil {
		return nil, util.ErrNotFound("user pelapor tidak ditemukan")
	}

	if input.AssignedTo != nil {
		if _, err := u.userRepo.FindByIDAndTenant(*input.AssignedTo, input.TenantID); err != nil {
			return nil, util.ErrNotFound("user teknisi tidak ditemukan")
		}
	}

	complaint := &domain.Complaint{
		TenantID:    input.TenantID,
		AssetID:     input.AssetID,
		ReportedBy:  input.ReportedBy,
		AssignedTo:  input.AssignedTo,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Status:      domain.ComplaintStatusOpen,
	}

	if err := u.complaintRepo.Create(complaint); err != nil {
		return nil, err
	}
	return complaint, nil
}

func (u *complaintUsecase) GetAllComplaints(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, int64, error) {
	return u.complaintRepo.FindAllByTenant(tenantID, filter)
}

func (u *complaintUsecase) GetComplaintByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Complaint, error) {
	return u.complaintRepo.FindByIDAndTenant(id, tenantID)
}

func (u *complaintUsecase) UpdateComplaint(input UpdateComplaintInput) (*domain.Complaint, error) {
	complaint, err := u.complaintRepo.FindByIDAndTenant(input.ID, input.TenantID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Title) != "" {
		complaint.Title = strings.TrimSpace(input.Title)
	}
	if strings.TrimSpace(input.Description) != "" {
		complaint.Description = strings.TrimSpace(input.Description)
	}
	if strings.TrimSpace(input.Status) != "" {
		status := strings.ToLower(strings.TrimSpace(input.Status))
		if !domain.IsValidComplaintStatus(status) {
			return nil, util.ErrValidation("status complaint tidak valid")
		}
		complaint.Status = status
	}
	if input.AssignedToSet {
		if input.AssignedTo != nil {
			if _, err := u.userRepo.FindByIDAndTenant(*input.AssignedTo, input.TenantID); err != nil {
				return nil, util.ErrNotFound("user teknisi tidak ditemukan")
			}
		}
		complaint.AssignedTo = input.AssignedTo
	}
	if strings.TrimSpace(input.ResolutionNote) != "" {
		complaint.ResolutionNote = strings.TrimSpace(input.ResolutionNote)
	}

	if err := u.complaintRepo.Update(complaint); err != nil {
		return nil, err
	}
	return complaint, nil
}

func (u *complaintUsecase) DeleteComplaint(tenantID uuid.UUID, id uuid.UUID) error {
	complaint, err := u.complaintRepo.FindByIDAndTenant(id, tenantID)
	if err != nil {
		return err
	}
	return u.complaintRepo.Delete(complaint)
}
