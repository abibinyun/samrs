import { useAppSelector } from "@/store/hooks";

export const useAuth = () => {
  const { permissions } = useAppSelector((state) => state.auth);
  const hasPermission = (perm: string) => permissions.includes(perm);
  return { permissions, hasPermission };
};
