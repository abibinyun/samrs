package app

import (
	deliveryhandlers "samrs-backend/internal/delivery/http/handlers"
	"samrs-backend/internal/delivery/http/middleware"
	"samrs-backend/internal/app/routes"
	"samrs-backend/internal/repository"
	assetusecase "samrs-backend/internal/usecase/asset"
	auditusecase "samrs-backend/internal/usecase/audit"
	authusecase "samrs-backend/internal/usecase/auth"
	complaintusecase "samrs-backend/internal/usecase/complaint"
	documentusecase "samrs-backend/internal/usecase/document"
	masterdatausecase "samrs-backend/internal/usecase/masterdata"
	maintenanceusecase "samrs-backend/internal/usecase/maintenance"
	notificationusecase "samrs-backend/internal/usecase/notification"
	rbacusecase "samrs-backend/internal/usecase/rbac"
	reportusecase "samrs-backend/internal/usecase/report"
	stockopnameusecase "samrs-backend/internal/usecase/stockopname"
	tenantusecase "samrs-backend/internal/usecase/tenant"
	userusecase "samrs-backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
}

func New(db *gorm.DB, logger *zap.Logger) *App {
	// Repository
	permissionRepo := repository.NewPermissionRepository(db)
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bedRepo := repository.NewBedRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	assetEventRepo := repository.NewAssetEventRepository(db)
	assetMutationRepo := repository.NewAssetMutationRepository(db)
	maintenanceRepo := repository.NewMaintenanceScheduleRepository(db)
	maintenanceDocRepo := repository.NewMaintenanceDocumentRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	documentFileRepo := repository.NewDocumentFileRepository(db)
	stockOpnameRepo := repository.NewStockOpnameRepository(db)
	stockOpnameItemRepo := repository.NewStockOpnameItemRepository(db)
	vendorRepo := repository.NewVendorRepository(db)
	brandRepo := repository.NewAssetBrandRepository(db)
	modelRepo := repository.NewAssetModelRepository(db)
	statusRepo := repository.NewAssetStatusRepository(db)
	complaintRepo := repository.NewComplaintRepository(db)
	auditRepo := repository.NewAuditTrailRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	tenantRepo := repository.NewTenantRepository(db)

	// Usecase
	authUsecase := authusecase.NewAuthUsecase(userRepo, permissionRepo)
	roomUsecase := masterdatausecase.NewRoomUsecase(roomRepo)
	bedUsecase := masterdatausecase.NewBedUsecase(bedRepo, roomRepo)
	categoryUsecase := masterdatausecase.NewCategoryUsecase(categoryRepo)
	assetUsecase := assetusecase.NewAssetUsecase(
		assetRepo,
		categoryRepo,
		roomRepo,
		bedRepo,
		vendorRepo,
		brandRepo,
		modelRepo,
		statusRepo,
	)
	assetEventUsecase := assetusecase.NewAssetEventUsecase(assetEventRepo, assetRepo, userRepo)
	assetMutationUsecase := assetusecase.NewAssetMutationUsecase(assetMutationRepo, assetRepo, roomRepo, bedRepo)
	maintenanceUsecase := maintenanceusecase.NewMaintenanceScheduleUsecase(maintenanceRepo, assetRepo)
	maintenanceDocUsecase := maintenanceusecase.NewMaintenanceDocumentUsecase(maintenanceDocRepo, maintenanceRepo)
	documentUsecase := documentusecase.NewDocumentUsecase(documentRepo, documentFileRepo, db)
	stockOpnameUsecase := stockopnameusecase.NewStockOpnameUsecase(stockOpnameRepo, stockOpnameItemRepo, assetRepo)
	vendorUsecase := masterdatausecase.NewVendorUsecase(vendorRepo)
	brandUsecase := masterdatausecase.NewAssetBrandUsecase(brandRepo)
	modelUsecase := masterdatausecase.NewAssetModelUsecase(modelRepo, brandRepo)
	statusUsecase := masterdatausecase.NewAssetStatusUsecase(statusRepo)
	reportUsecase := reportusecase.NewReportUsecase(assetRepo, complaintRepo, maintenanceRepo)
	notificationUsecase := notificationusecase.NewNotificationUsecase([]notificationusecase.NotificationProvider{
		notificationusecase.NewNoopNotificationProvider("email"),
		notificationusecase.NewNoopNotificationProvider("whatsapp"),
		notificationusecase.NewNoopNotificationProvider("generic"),
	})
	complaintUsecase := complaintusecase.NewComplaintUsecase(complaintRepo, assetRepo, userRepo)
	rbacUsecase := rbacusecase.NewRBACUsecase(roleRepo, permissionRepo)
	permissionUsecase := rbacusecase.NewPermissionUsecase(permissionRepo, roleRepo)
	roleUsecase := rbacusecase.NewRoleUsecase(roleRepo)
	rolePermissionUsecase := rbacusecase.NewRolePermissionUsecase(roleRepo, permissionRepo)
	userUsecase := userusecase.NewUserUsecase(userRepo, roleRepo)
	auditUsecase := auditusecase.NewAuditTrailUsecase(auditRepo)
	publicAssetUsecase := assetusecase.NewPublicAssetUsecase(tenantRepo, assetRepo)
	adminTenantUsecase := tenantusecase.NewAdminTenantUsecase(
		tenantRepo,
		roleRepo,
		userRepo,
		roomRepo,
		bedRepo,
		categoryRepo,
		assetRepo,
		auditRepo,
	)

	// Handler
	authHandler := deliveryhandlers.NewAuthHandler(authUsecase)
	roomHandler := deliveryhandlers.NewRoomHandler(roomUsecase, db)
	bedHandler := deliveryhandlers.NewBedHandler(bedUsecase, db)
	categoryHandler := deliveryhandlers.NewCategoryHandler(categoryUsecase, db)
	assetHandler := deliveryhandlers.NewAssetHandler(assetUsecase, tenantRepo, db)
	assetEventHandler := deliveryhandlers.NewAssetEventHandler(assetEventUsecase)
	assetMutationHandler := deliveryhandlers.NewAssetMutationHandler(assetMutationUsecase, db)
	maintenanceHandler := deliveryhandlers.NewMaintenanceScheduleHandler(maintenanceUsecase, db)
	maintenanceDocHandler := deliveryhandlers.NewMaintenanceDocumentHandler(maintenanceDocUsecase, db)
	documentHandler := deliveryhandlers.NewDocumentHandler(documentUsecase, db)
	stockOpnameHandler := deliveryhandlers.NewStockOpnameHandler(stockOpnameUsecase, db)
	vendorHandler := deliveryhandlers.NewVendorHandler(vendorUsecase, db)
	brandHandler := deliveryhandlers.NewAssetBrandHandler(brandUsecase, db)
	modelHandler := deliveryhandlers.NewAssetModelHandler(modelUsecase, db)
	statusHandler := deliveryhandlers.NewAssetStatusHandler(statusUsecase, db)
	reportHandler := deliveryhandlers.NewReportHandler(reportUsecase)
	notificationHandler := deliveryhandlers.NewNotificationHandler(notificationUsecase)
	complaintHandler := deliveryhandlers.NewComplaintHandler(complaintUsecase, db)
	permissionHandler := deliveryhandlers.NewPermissionHandler(permissionUsecase)
	roleHandler := deliveryhandlers.NewRoleHandler(roleUsecase, rolePermissionUsecase, db)
	userHandler := deliveryhandlers.NewUserHandler(userUsecase, db)
	auditHandler := deliveryhandlers.NewAuditTrailHandler(auditUsecase)
	adminTenantHandler := deliveryhandlers.NewAdminTenantHandler(adminTenantUsecase, db)
	publicAssetHandler := deliveryhandlers.NewPublicAssetHandler(publicAssetUsecase)

	// Router
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.StructuredLogger(logger), gin.Recovery(), middleware.Metrics())
	handlers := routes.Handlers{
		Auth:           authHandler,
		Room:           roomHandler,
		Bed:            bedHandler,
		Category:       categoryHandler,
		Asset:          assetHandler,
		AssetEvent:     assetEventHandler,
		AssetMutation:  assetMutationHandler,
		Maintenance:    maintenanceHandler,
		MaintenanceDoc: maintenanceDocHandler,
		Document:       documentHandler,
		StockOpname:    stockOpnameHandler,
		Vendor:         vendorHandler,
		Brand:          brandHandler,
		Model:          modelHandler,
		Status:         statusHandler,
		Report:         reportHandler,
		Notification:   notificationHandler,
		Complaint:      complaintHandler,
		Permission:     permissionHandler,
		Role:           roleHandler,
		User:           userHandler,
		Audit:          auditHandler,
		AdminTenant:    adminTenantHandler,
		PublicAsset:    publicAssetHandler,
	}

	routes.Register(r, handlers, routes.RouteDeps{
		RBAC:     rbacUsecase,
		RoleRepo: roleRepo,
	})

	return &App{Router: r}
}
