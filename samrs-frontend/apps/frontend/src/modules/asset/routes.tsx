import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";
import { TableSkeleton } from "@/components/commons/data-table/TableSkeleton";
import { AssetDetailPage } from "./pages/AssetDetailPage";
import { AssetUpdatePage } from "./pages/AssetUpdatePage";

export const assetRoutes = [
  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "assets", 
    beforeLoad: () => requirePermission(SCOPES.ASSET.READ),
    component: lazyRouteComponent(() =>
      import("./pages/AssetListPage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} isFilterBoxVisible />,
    pendingMs: 300,
  }),

  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "assets/create", 
    beforeLoad: () => requirePermission(SCOPES.ASSET.CREATE),
    component: lazyRouteComponent(() =>
      import("./pages/AssetCreatePage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} />,
    pendingMs: 300
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "assets/$id",
    beforeLoad: () => requirePermission(SCOPES.ASSET.READ),
    component: AssetDetailPage,
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "assets/$id/edit",
    beforeLoad: () => requirePermission(SCOPES.ASSET.UPDATE),
    component: AssetUpdatePage,
  }),
];
