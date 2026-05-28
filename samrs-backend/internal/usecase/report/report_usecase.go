package usecase

import (
	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"

	"github.com/google/uuid"
)

type ReportUsecase interface {
	ExportAssets(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, error)
	ExportComplaints(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, error)
	ExportMaintenance(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, error)
}

type reportUsecase struct {
	assetRepo       repository.AssetRepository
	complaintRepo   repository.ComplaintRepository
	maintenanceRepo repository.MaintenanceScheduleRepository
}

func NewReportUsecase(
	ar repository.AssetRepository,
	cr repository.ComplaintRepository,
	mr repository.MaintenanceScheduleRepository,
) ReportUsecase {
	return &reportUsecase{
		assetRepo:       ar,
		complaintRepo:   cr,
		maintenanceRepo: mr,
	}
}

func (u *reportUsecase) ExportAssets(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, error) {
	return u.assetRepo.FindAllByTenantExport(tenantID, filter)
}

func (u *reportUsecase) ExportComplaints(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, error) {
	return u.complaintRepo.FindAllByTenantExport(tenantID, filter)
}

func (u *reportUsecase) ExportMaintenance(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, error) {
	return u.maintenanceRepo.FindAllByTenantExport(tenantID, filter)
}
