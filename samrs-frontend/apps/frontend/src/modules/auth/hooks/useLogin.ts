import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { useDispatch } from "react-redux";
import type { ILoginForm, ILoginResponse, IMeResponse } from "../types";
import { authApi } from "../actions/api";
import { setCredentials } from "../actions/slice";

export const useLogin = () => {
  const [loginMutation, { isLoading: isLoggingIn }] = authApi.useLoginMutation();
  const [meQueryTrigger] = authApi.useLazyMeQuery();
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogin = async (formData: ILoginForm) => {
    try {
      // 1️⃣ Login untuk dapat token
      const loginRes = await loginMutation(formData).unwrap() as ILoginResponse;
      const token = loginRes.data.token;

      // Simpan token sementara ke Redux supaya request me punya header Authorization
      dispatch(setCredentials({ token }));

      // 2️⃣ Hit API 'me' untuk dapat user & permissions
      const meRes = await meQueryTrigger().unwrap() as IMeResponse;
      const userData = meRes.data.user;
      const permissions = meRes.data.permissions;

      // 3️⃣ Update Redux dengan token + user + permissions
      dispatch(setCredentials({
        token,
        user: {
          id: userData.id,
          tenant_slug: userData.tenant.slug,
          role: userData.role.name,
          name: userData.username,
          is_system: userData.role.is_system, // Super Admin flag from role
        },
        permissions,
      }));

      // 4️⃣ Toast & redirect
      toast.success(`Selamat datang, ${userData.username}`);
      
      // Debug log
      console.log('🔐 Login Success:', {
        tenant_slug: userData.tenant.slug,
        is_system: userData.role.is_system,
        permissions: permissions.length
      });
      
      // Navigate to tenant dashboard
      const tenantSlug = userData.tenant?.slug;
      if (!tenantSlug) {
        throw new Error('Tenant slug tidak ditemukan');
      }
      
      navigate({
        to: "/$tenantId",
        params: { tenantId: tenantSlug },
      });

    } catch (err: any) {
      toast.error(err.data?.message || "Login atau verifikasi profil gagal");
    }
  };

  return { handleLogin, isLoggingIn };
};
