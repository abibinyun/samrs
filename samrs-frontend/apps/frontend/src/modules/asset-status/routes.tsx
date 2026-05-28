import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const assetStatusRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "asset-statuses",
    beforeLoad: () => requirePermission(SCOPES.ASSET_STATUS.READ),
    component: lazyRouteComponent(() => import("./pages/AssetStatusListPage")),
    pendingMs: 300,
  }),
];
