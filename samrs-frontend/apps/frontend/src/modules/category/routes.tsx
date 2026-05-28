import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const categoryRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "categories",
    beforeLoad: () => requirePermission(SCOPES.CATEGORY.READ),
    component: lazyRouteComponent(() => import("./pages/CategoryListPage")),
    pendingMs: 300,
  }),
];
