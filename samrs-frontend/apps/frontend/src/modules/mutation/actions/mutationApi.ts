import { api } from "@/store/api";
import type { AssetMutation, MutationFormInput, MutationListResponse, MutationApiResponse } from "../types";

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
      transformResponse: (response: MutationApiResponse<AssetMutation[]>) => ({
        data: response.data,
        total: (response as any).meta?.total || 0,
      }),
    }),

    getMutationById: builder.query<MutationApiResponse<AssetMutation>, number>({
      query: (id) => `/api/v1/asset-mutations/${id}`,
      providesTags: (_result, _error, id) => [{ type: "Assets", id: `mutation-${id}` }],
    }),

    createMutation: builder.mutation<MutationApiResponse<{ mutation: AssetMutation }>, MutationFormInput>({
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
