import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface User {
  id: string;
  tenant_slug: string;
  role: string;
  name?: string;
  is_system?: boolean; // Deprecated: use role.is_system instead
}

interface AuthState {
  token: string | null;
  user: User | null;
  permissions: string[];
}

const initialState: AuthState = {
  token:
    typeof window !== "undefined"
      ? localStorage.getItem("samrs_token")
      : null,
  user:
    typeof window !== "undefined" && localStorage.getItem("samrs_user")
      ? JSON.parse(localStorage.getItem("samrs_user")!)
      : null,
  permissions:
    typeof window !== "undefined" && localStorage.getItem("samrs_permissions")
      ? JSON.parse(localStorage.getItem("samrs_permissions")!)
      : [],
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{
        token: string;
        user?: User;
        permissions?: string[];
      }>
    ) => {
      const { token, user, permissions } = action.payload;

      state.token = token;
      if (user) state.user = user;
      if (permissions) state.permissions = permissions;

      localStorage.setItem("samrs_token", token);
      if (user) localStorage.setItem("samrs_user", JSON.stringify(user));
      if (permissions) localStorage.setItem("samrs_permissions", JSON.stringify(permissions));
    },

    clearCredentials: (state) => {
      state.token = null;
      state.user = null;
      state.permissions = [];
      localStorage.removeItem("samrs_token");
      localStorage.removeItem("samrs_user");
      localStorage.removeItem("samrs_permissions");
    },
  },
});

export const { setCredentials, clearCredentials } = authSlice.actions;
export default authSlice.reducer;
