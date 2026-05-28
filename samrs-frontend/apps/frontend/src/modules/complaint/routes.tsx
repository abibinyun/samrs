import { createRoute } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";
import { requirePermission } from "@/lib/routeGuard";
import { SCOPES } from "@/constants";
import { ComplaintListPage } from "./pages/ComplaintListPage";
import { ComplaintCreatePage } from "./pages/ComplaintCreatePage";

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

export const complaintRoutes = [complaintListRoute, complaintCreateRoute];
