import { api } from "@/store/api";
import type { AssetMutation, MutationFormInput } from "../types";

type ApiResponse<T> = {
  success: boolean;
  message: string;
  data: T;
  meta?: {
    total: number;
    page: number;
    per_page: number;
  };
};

type MutationListResponse = {
  data: AssetMutation[];
  total: number;
};

export const mutationApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMutations: builder.query<MutationListResponse, {
      page?: number;
      limit?: number;
      search?: string;
      asset_id?: string;
      date_from?: string;
      date_to?: string;
    }>({
      query: (params) => ({
        url: "/api/v1/asset-mutations",
        params,
      }),
      providesTags: ["Assets"],
      transformResponse: (response: ApiResponse<AssetMutation[]>) => ({
        data: response.data,
        total: response.meta?.total || 0,
      }),
    }),

    getMutationById: builder.query<ApiResponse<AssetMutation>, number>({
      query: (id) => `/api/v1/asset-mutations/${id}`,
      providesTags: (_result, _error, id) => [{ type: "Assets", id: `mutation-${id}` }],
    }),

    createMutation: builder.mutation<ApiResponse<{ mutation: AssetMutation }>, MutationFormInput>({
      query: (body) => ({
        url: "/api/v1/asset-mutations",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Assets"],
    }),
  }),
});

export const {
  useGetMutationsQuery,
  useGetMutationByIdQuery,
  useCreateMutationMutation,
} = mutationApi;
