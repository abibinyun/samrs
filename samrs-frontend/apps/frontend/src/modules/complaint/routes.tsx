import { createRoute } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";
import { requirePermission } from "@/lib/routeGuard";
import { SCOPES } from "@/constants";
import { ComplaintListPage } from "./pages/ComplaintListPage";
import { ComplaintCreatePage } from "./pages/ComplaintCreatePage";
import { ComplaintDetailPage } from "./pages/ComplaintDetailPage";

export const complaintListRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "complaints",
  beforeLoad: () => requirePermission(SCOPES.COMPLAINT.READ),
  component: ComplaintListPage,
});

export const complaintCreateRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "complaints/create",
  beforeLoad: () => requirePermission(SCOPES.COMPLAINT.CREATE),
  component: ComplaintCreatePage,
});

export const complaintDetailRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "complaints/$complaintId",
  beforeLoad: () => requirePermission(SCOPES.COMPLAINT.READ),
  component: ComplaintDetailPage,
});

export const complaintRoutes = [complaintListRoute, complaintCreateRoute, complaintDetailRoute];
