package repository

import (
	maintenancerepo "samrs-backend/internal/repository/maintenance"

	"gorm.io/gorm"
)

type MaintenanceScheduleFilter = maintenancerepo.MaintenanceScheduleFilter

type MaintenanceScheduleRepository = maintenancerepo.MaintenanceScheduleRepository

type MaintenanceDocumentRepository = maintenancerepo.MaintenanceDocumentRepository

func NewMaintenanceScheduleRepository(db *gorm.DB) MaintenanceScheduleRepository {
	return maintenancerepo.NewMaintenanceScheduleRepository(db)
}

func NewMaintenanceDocumentRepository(db *gorm.DB) MaintenanceDocumentRepository {
	return maintenancerepo.NewMaintenanceDocumentRepository(db)
}
