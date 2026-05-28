import { redirect } from "@tanstack/react-router";
import { store } from "@/store";

/**
 * Check if user has permission or is Super Admin
 * Use this in route beforeLoad
 */
export const requirePermission = (permission: string | string[]) => {
  const { permissions, user } = store.getState().auth;
  
  // Debug log
  console.log('🔐 Route Guard Check:', { 
    user_is_system: user?.is_system, 
    user_role: user?.role,
    required_permission: permission 
  });
  
  // Super Admin bypass
  if (user?.is_system === true) {
    console.log('✅ Super Admin - Access Granted');
    return;
  }
  
  // Check permissions
  const hasPermission = Array.isArray(permission)
    ? permission.some(p => permissions?.includes(p))
    : permissions?.includes(permission);
  
  if (!hasPermission) {
    console.log('❌ Permission Denied - Redirecting to home');
    throw redirect({ to: "/" });
  }
  
  console.log('✅ Permission Granted');
};
