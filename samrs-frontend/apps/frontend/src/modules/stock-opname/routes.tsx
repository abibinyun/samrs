import { tenantRoute } from "@/routes";
import { createRoute } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { requirePermission } from "@/lib/routeGuard";
import { StockOpnameListPage } from "./pages/StockOpnameListPage";
import { StockOpnameCreatePage } from "./pages/StockOpnameCreatePage";
import { StockOpnameDetailPage } from "./pages/StockOpnameDetailPage";
import { StockOpnameEditPage } from "./pages/StockOpnameEditPage";

export const stockOpnameRoutes = [
  createRoute({
    getParentRoute: () => tenantRoute,
    path: "stock-opname",
    beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.READ),
    component: StockOpnameListPage,
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "stock-opname/create",
    beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.CREATE),
    component: StockOpnameCreatePage,
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "stock-opname/$sessionId",
    beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.READ),
    component: StockOpnameDetailPage,
    parseParams: (params) => ({ sessionId: Number(params.sessionId) }),
  }),

  createRoute({
    getParentRoute: () => tenantRoute,
    path: "stock-opname/$sessionId/edit",
    beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.UPDATE),
    component: StockOpnameEditPage,
    parseParams: (params) => ({ sessionId: Number(params.sessionId) }),
  }),
];
