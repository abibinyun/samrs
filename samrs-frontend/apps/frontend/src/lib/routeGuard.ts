import { redirect } from "@tanstack/react-router";
import { store } from "@/store";

/**
 * Check if user has permission or is Super Admin
 * Use this in route beforeLoad
 */
export const requirePermission = (permission: string | string[]) => {
  const { permissions, user } = store.getState().auth;
  
  // Super Admin bypass
  if (user?.is_system === true) {
    return;
  }
  
  // Check permissions
  const hasPermission = Array.isArray(permission)
    ? permission.some(p => permissions?.includes(p))
    : permissions?.includes(permission);
  
  if (!hasPermission) {
    const tenantSlug = user?.tenant_slug || "";
    throw redirect({ 
      to: "/$tenantId/forbidden", 
      params: { tenantId: tenantSlug } 
    });
  }
};
