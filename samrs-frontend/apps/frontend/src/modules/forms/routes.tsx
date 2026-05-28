import { tenantRoute } from "@/routes";
import { createRoute, lazyRouteComponent, redirect } from "@tanstack/react-router";
import { SCOPES } from "@/constants";
import { store } from "@/store";
import { TableSkeleton } from "@/components/commons/data-table/TableSkeleton";

export const formRoutes = [
  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "forms", 
    beforeLoad: () => {
      const { permissions } = store.getState().auth
      if (!permissions.includes(SCOPES.ASSET.READ)) {
        throw redirect({ to: '/' })
      }
    },
    component: lazyRouteComponent(() =>
      import("./pages/FormTemplatesPage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} isFilterBoxVisible />,
    pendingMs: 300,
  }),

  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "forms/create", 
    beforeLoad: () => {
      const { permissions } = store.getState().auth
      if (!permissions.includes(SCOPES.ASSET.CREATE)) {
        throw redirect({ to: '/' })
      }
    },
    component: lazyRouteComponent(() =>
      import("./pages/CreateTemplatesPage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} />,
    pendingMs: 300
  }),
  
  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "forms/edit/$id", 
    beforeLoad: () => {
      const { permissions } = store.getState().auth;
      if (!permissions.includes(SCOPES.ASSET.CREATE)) {
        throw redirect({ to: "/" });
      }
    },
    component: lazyRouteComponent(() =>
      import("./pages/EditTemplatesPage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} />,
    pendingMs: 300
  }),

  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "forms/fill/$id", 
    beforeLoad: () => {
      const { permissions } = store.getState().auth;
      if (!permissions.includes(SCOPES.ASSET.READ)) {
        throw redirect({ to: "/" });
      }
    },
    component: lazyRouteComponent(() =>
      import("./pages/FormFillPage")
    ),
    pendingComponent: () => <TableSkeleton rows={10} />,
    pendingMs: 300
  }),
];

