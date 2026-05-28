import {
  createApi,
  fetchBaseQuery,
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
  type FetchBaseQueryMeta,
  type QueryReturnValue,
} from "@reduxjs/toolkit/query/react";
import type { RootState } from "./index";
import {
  setCredentials,
  clearCredentials,
} from "@/modules/auth/actions/slice";

/**
 * Base query
 */
const rawBaseQuery = fetchBaseQuery({
  baseUrl: import.meta.env.VITE_API_URL,
  credentials: "include",
  prepareHeaders: (headers, { getState }) => {
    const state = getState() as RootState;

    if (state.auth.token) {
      headers.set("Authorization", `Bearer ${state.auth.token}`);
    }

    if (state.auth.user?.tenant_slug) {
      headers.set("X-Tenant-ID", state.auth.user.tenant_slug);
    }

    return headers;
  },

  // prepareHeaders: (headers, { getState, endpoint }) => {
  //   const state = getState() as RootState;

  //   // 1. Ambil token dari Redux
  //   const token = state.auth.token;

  //   // 2. Jangan pasang Authorization header untuk endpoint login / refresh
  //   if (token && endpoint !== "login" && endpoint !== "refresh") {
  //     headers.set("Authorization", `Bearer ${token}`);
  //   }

  //   // tenant header optional
  //   if (state.auth.user?.tenant_slug) {
  //     headers.set("X-Tenant-ID", state.auth.user.tenant_slug);
  //   }

  //   return headers;
  // },

});

/**
 * Refresh mutex
 */
let isRefreshing = false;
let refreshPromise:
  | Promise<
      QueryReturnValue<
        unknown,
        FetchBaseQueryError,
        FetchBaseQueryMeta
      >
    >
  | null = null;

/**
 * Base query with auto refresh
 */
const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError,
  {},
  FetchBaseQueryMeta
> = async (args, api, extraOptions) => {
  let result = await rawBaseQuery(args, api, extraOptions);

  if (result.error?.status !== 401) {
    return result;
  }

  // 🔒 single refresh request
  if (!isRefreshing) {
    isRefreshing = true;

    refreshPromise = (async () => {
      try {
        return await rawBaseQuery(
          { url: "/auth/refresh", method: "POST" },
          api,
          extraOptions
        );
      } finally {
        isRefreshing = false;
      }
    })();
  }

  const refreshResult = await refreshPromise;

  if (refreshResult?.data && typeof refreshResult.data === "object") {
    const token = (refreshResult.data as { token?: string }).token;

    if (token) {
      api.dispatch(setCredentials({ token }));

      // 🔁 retry original request
      result = await rawBaseQuery(args, api, extraOptions);
      return result;
    }
  }

  api.dispatch(clearCredentials());
  return result;
};

export const api = createApi({
  reducerPath: "api",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["User", "Auth", "MasterData", "Complaints", "Maintenance", "StockOpname", "Assets"],
  endpoints: () => ({}),
});
