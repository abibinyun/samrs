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
});

/**
 * Base query with auto logout on 401
 */
const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError,
  {},
  FetchBaseQueryMeta
> = async (args, api, extraOptions) => {
  const result = await rawBaseQuery(args, api, extraOptions);

  if (result.error?.status === 401) {
    api.dispatch(clearCredentials());
  }

  return result;
};

export const api = createApi({
  reducerPath: "api",
  baseQuery: baseQueryWithReauth,
  tagTypes: ["User", "Auth", "MasterData", "Complaints", "Maintenance", "StockOpname", "Assets"],
  endpoints: () => ({}),
});
