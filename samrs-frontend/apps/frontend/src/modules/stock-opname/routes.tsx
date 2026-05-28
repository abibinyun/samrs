import { createRoute } from "@tanstack/react-router";
import { tenantRoute } from "@/routes";
import { requirePermission } from "@/lib/routeGuard";
import { SCOPES } from "@/constants";
import { StockOpnameListPage } from "./pages/StockOpnameListPage";
import { StockOpnameCreatePage } from "./pages/StockOpnameCreatePage";

export const stockOpnameListRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "assets/stock-opname",
  beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.READ),
  component: StockOpnameListPage,
});

export const stockOpnameCreateRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "assets/stock-opname/create",
  beforeLoad: () => requirePermission(SCOPES.STOCK_OPNAME.CREATE),
  component: StockOpnameCreatePage,
});

export const stockOpnameRoutes = [stockOpnameListRoute, stockOpnameCreateRoute];
