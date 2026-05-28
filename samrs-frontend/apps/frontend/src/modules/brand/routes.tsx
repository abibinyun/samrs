import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const brandRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "brands",
    beforeLoad: () => requirePermission(SCOPES.BRAND.READ),
    component: lazyRouteComponent(() => import("./pages/BrandListPage")),
    pendingMs: 300,
  }),
];
