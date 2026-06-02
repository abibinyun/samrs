import { createRoute } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";
import { requirePermission } from "@/lib/routeGuard";
import { SCOPES } from "@/constants";
import { MaintenanceListPage } from "./pages/MaintenanceListPage";
import { MaintenanceCreatePage } from "./pages/MaintenanceCreatePage";
import { MaintenanceEditPage } from "./pages/MaintenanceEditPage";
import { MaintenanceDetailPage } from "./pages/MaintenanceDetailPage";
import { MaintenanceDocumentListPage } from "./pages/MaintenanceDocumentListPage";

export const maintenanceListRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "maintenance/schedules",
  beforeLoad: () => requirePermission(SCOPES.MAINTENANCE.READ),
  component: MaintenanceListPage,
});

export const maintenanceCreateRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "maintenance/schedules/create",
  beforeLoad: () => requirePermission(SCOPES.MAINTENANCE.CREATE),
  component: MaintenanceCreatePage,
});

export const maintenanceEditRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "maintenance/schedules/$scheduleId/edit",
  beforeLoad: () => requirePermission(SCOPES.MAINTENANCE.UPDATE),
  component: MaintenanceEditPage,
});

export const maintenanceDetailRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "maintenance/schedules/$scheduleId",
  beforeLoad: () => requirePermission(SCOPES.MAINTENANCE.READ),
  component: MaintenanceDetailPage,
});

export const maintenanceDocumentRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "maintenance/documents",
  beforeLoad: () => requirePermission(SCOPES.MAINTENANCE_DOCUMENT.READ),
  component: MaintenanceDocumentListPage,
});

export const maintenanceRoutes = [
  maintenanceListRoute,
  maintenanceCreateRoute,
  maintenanceEditRoute,
  maintenanceDetailRoute,
  maintenanceDocumentRoute,
];
