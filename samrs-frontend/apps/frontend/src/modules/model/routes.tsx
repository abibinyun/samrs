import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const modelRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "models",
    beforeLoad: () => requirePermission(SCOPES.MODEL.READ),
    component: lazyRouteComponent(() => import("./pages/ModelListPage")),
    pendingMs: 300,
  }),
];
