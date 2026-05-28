package routes

import (
	deliveryhandlers "samrs-backend/internal/delivery/http/handlers"
	"samrs-backend/internal/repository"
	rbacusecase "samrs-backend/internal/usecase/rbac"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth           *deliveryhandlers.AuthHandler
	Room           *deliveryhandlers.RoomHandler
	Bed            *deliveryhandlers.BedHandler
	Category       *deliveryhandlers.CategoryHandler
	Asset          *deliveryhandlers.AssetHandler
	AssetEvent     *deliveryhandlers.AssetEventHandler
	AssetMutation  *deliveryhandlers.AssetMutationHandler
	Maintenance    *deliveryhandlers.MaintenanceScheduleHandler
	MaintenanceDoc *deliveryhandlers.MaintenanceDocumentHandler
	Document       *deliveryhandlers.DocumentHandler
	StockOpname    *deliveryhandlers.StockOpnameHandler
	Vendor         *deliveryhandlers.VendorHandler
	Brand          *deliveryhandlers.AssetBrandHandler
	Model          *deliveryhandlers.AssetModelHandler
	Status         *deliveryhandlers.AssetStatusHandler
	Report         *deliveryhandlers.ReportHandler
	Notification   *deliveryhandlers.NotificationHandler
	Complaint      *deliveryhandlers.ComplaintHandler
	Permission     *deliveryhandlers.PermissionHandler
	Role           *deliveryhandlers.RoleHandler
	User           *deliveryhandlers.UserHandler
	Audit          *deliveryhandlers.AuditTrailHandler
	AdminTenant    *deliveryhandlers.AdminTenantHandler
	PublicAsset    *deliveryhandlers.PublicAssetHandler
}

type RouteDeps struct {
	RBAC     rbacusecase.RBACUsecase
	RoleRepo repository.RoleRepository
}

func Register(r *gin.Engine, handlers Handlers, deps RouteDeps) {
	registerPublicRoutes(r, handlers)
	registerProtectedRoutes(r, handlers, deps)
}
