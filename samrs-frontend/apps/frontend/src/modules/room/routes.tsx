import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";

export const roomRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "rooms",
    beforeLoad: () => requirePermission(SCOPES.ROOM.READ),
    component: lazyRouteComponent(() => import("./pages/RoomListPage")),
    pendingMs: 300,
  }),
];
