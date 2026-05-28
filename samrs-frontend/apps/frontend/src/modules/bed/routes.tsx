import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const bedRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "beds",
    beforeLoad: () => requirePermission(SCOPES.BED.READ),
    component: lazyRouteComponent(() => import("./pages/BedListPage")),
    pendingMs: 300,
  }),
];
