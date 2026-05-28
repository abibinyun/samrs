import { useAppSelector } from "@/store/hooks";

export const usePermission = () => {
  const { permissions, user } = useAppSelector((state) => state.auth);

  const hasPermission = (permission: string | string[]) => {
    // Super Admin (is_system) bypass all permission checks
    if (user?.is_system) return true;
    
    if (!permissions) return false;
    
    // Jika input berupa array, cek apakah salah satu izin terpenuhi (OR logic)
    if (Array.isArray(permission)) {
      return permission.some((p) => permissions.includes(p));
    }
    
    return permissions.includes(permission);
  };

  return { hasPermission };
};