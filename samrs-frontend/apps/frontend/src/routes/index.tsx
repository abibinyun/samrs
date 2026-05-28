// import { lazy } from "react";
import { createRoute, createRouter, redirect, lazyRouteComponent } from "@tanstack/react-router";
import { rootRoute } from "./__root";
import { store } from "@/store";
import { setCredentials } from "@/modules/auth/actions/slice";

// const LoginPage = lazy(() => import('../modules/auth/pages/LoginPage'))
// const DashboardShell = lazy(() => import('../modules/dashboard/layouts/DashboardShell'))
// const DashboardPage = lazy(() => import('../modules/dashboard/pages/DashboardPage'))

// Import semua route modul di sini
import { assetRoutes } from "@/modules/asset/routes";
import { roomRoutes } from "@/modules/room/routes";
import { notificationRoutes } from "@/modules/notification/routes";
import { mutationRoutes } from "@/modules/mutation/routes";
import { formRoutes } from "@/modules/forms/routes";
import { complaintRoutes } from "@/modules/complaint/routes";
import { maintenanceRoutes } from "@/modules/maintenance/routes";
import { stockOpnameRoutes } from "@/modules/stock-opname/routes";
import { categoryRoutes } from "@/modules/category/routes";
import { vendorRoutes } from "@/modules/vendor/routes";
import { brandRoutes } from "@/modules/brand/routes";
import { modelRoutes } from "@/modules/model/routes";
import { assetStatusRoutes } from "@/modules/asset-status/routes";
import { bedRoutes } from "@/modules/bed/routes";

// Public Routes
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  // component: LoginPage,
  component: lazyRouteComponent(() =>
    import("../modules/auth/pages/LoginPage")
  ),
});

// Index Route dengan Proteksi
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: () => {
    const { user } = store.getState().auth;
    const token = localStorage.getItem("samrs_token");
    
    if (!token) throw redirect({ to: "/login" });
    if (user?.tenant_slug) {
      throw redirect({ 
        to: "/$tenantId", 
        params: { tenantId: user.tenant_slug } 
      });
    }
  },
});

// Protected Tenant Route
export const tenantRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/$tenantId",
  // Di dalam tenantRoute
  beforeLoad: async ({ params }) => {
    const state = store.getState();
    const token = localStorage.getItem("samrs_token");

    if (!token) throw redirect({ to: "/login" });

    let currentUser = state?.auth?.user

    // Jika user refresh halaman, kita fetch ulang datanya SECARA SINKRON
    if (!currentUser) {
      try {
        // Gunakan trigger manual dari RTK Query atau fetch biasa
        const response = await fetch(`/api/v1/auth/me`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const resData = await response.json();
        const userData = resData.data.user;
        
        currentUser = { 
          id: userData.id,
          tenant_slug: userData.tenant.slug, 
          role: userData.role.name,
          name: userData.username,
          is_system: userData.role.is_system, // IMPORTANT: Super Admin flag
        };
        
        // Update store agar tidak fetch terus menerus
        store.dispatch(
          setCredentials({ 
            token, 
            user: currentUser, 
            permissions: resData.data.permissions 
          }));
      
      } catch {
        throw redirect({ to: "/login" });
      }
    }

    // Bandingkan URL vs Data Asli
    if (params.tenantId !== currentUser.tenant_slug) {
      throw redirect({ to: "/$tenantId", params: { tenantId: currentUser.tenant_slug } });
    }
  },
  // component: DashboardShell,
  component: lazyRouteComponent(() =>
    import("../modules/dashboard/layouts/DashboardShell")
  ),
});

// Dashboard Route
const dashboardRoute = createRoute({
  getParentRoute: () => tenantRoute,
  path: "/",
  // component: DashboardPage,
  component: lazyRouteComponent(() =>
    import("../modules/dashboard/pages/DashboardPage")
  ),
});

// Join all routes
const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  tenantRoute.addChildren([
    dashboardRoute,
    ...roomRoutes,
    ...notificationRoutes,
    ...assetRoutes,
    ...mutationRoutes,
    ...formRoutes,
    ...complaintRoutes,
    ...maintenanceRoutes,
    ...stockOpnameRoutes,
    ...categoryRoutes,
    ...vendorRoutes,
    ...brandRoutes,
    ...modelRoutes,
    ...assetStatusRoutes,
    ...bedRoutes,
    // Daftarkan 100 rute modul di sini
  ]),
]);

export const router = createRouter({ 
  routeTree,
  // Opsional: berikan context store ke router agar bisa diakses di beforeLoad jika perlu
  context: {
    store,
  }
});