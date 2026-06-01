import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";
import { TableSkeleton } from "@/components/commons/data-table/TableSkeleton";
import MutationDetailPage from "./pages/MutationDetailPage";

export const mutationRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "assets/mutation",
    beforeLoad: () => requirePermission(SCOPES.ASSET_MUTATION.READ),
    component: lazyRouteComponent(() => import("./pages/MutationListPage")),
    pendingComponent: () => <TableSkeleton rows={10} isFilterBoxVisible />,
    pendingMs: 300,
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "assets/mutation/create",
    beforeLoad: () => requirePermission(SCOPES.ASSET_MUTATION.CREATE),
    component: lazyRouteComponent(() => import("./pages/MutationCreatePage")),
    pendingComponent: () => <TableSkeleton rows={10} />,
    pendingMs: 300,
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "assets/mutation/$id",
    beforeLoad: () => requirePermission(SCOPES.ASSET_MUTATION.READ),
    component: MutationDetailPage,
  }),
];
