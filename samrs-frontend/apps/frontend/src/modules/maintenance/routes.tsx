import { createRoute } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";
import { requirePermission } from "@/lib/routeGuard";
import { SCOPES } from "@/constants";
import { MaintenanceListPage } from "./pages/MaintenanceListPage";
import { MaintenanceCreatePage } from "./pages/MaintenanceCreatePage";

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

export const maintenanceRoutes = [maintenanceListRoute, maintenanceCreateRoute];
