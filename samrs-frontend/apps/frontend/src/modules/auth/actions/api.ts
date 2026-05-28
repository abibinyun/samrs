import { api } from "@/store/api";
import { setCredentials } from "./slice";
import type { ILoginResponse, ILoginRequest, IMeResponse } from "../types";
import { API_ENDPOINTS } from "@/constants/endpoints";

export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<ILoginResponse, ILoginRequest>({
      query: (credentials) => ({
        url: API_ENDPOINTS.AUTH.LOGIN,
        method: "POST",
        body: credentials,
      }),
      async onQueryStarted(_, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;

          // update authSlice
          dispatch(
            setCredentials({
              token: data.data.token
            })
          );
        } catch {
          // error tetap ditangani RTK Query
        }
      },
    }),

    me: builder.query<IMeResponse, void>({
      query: () => API_ENDPOINTS.AUTH.ME,
      providesTags: ["User"],
    }),
  }),
  overrideExisting: false,
});

export const {
  useLoginMutation,
  useMeQuery,
  useLazyMeQuery
} = authApi;
