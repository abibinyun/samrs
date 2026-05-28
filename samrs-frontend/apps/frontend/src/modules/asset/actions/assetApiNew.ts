import { api } from "@/store/api";
import type { Asset, AssetFormInput } from "../types";

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

type AssetListResponse = {
  data: Asset[];
  total: number;
};

export const assetApiNew = api.injectEndpoints({
  endpoints: (builder) => ({
    getAssets: builder.query<AssetListResponse, { page?: number; limit?: number; search?: string; status?: string; category_id?: number }>({
      query: (params) => ({
        url: "/api/v1/assets",
        params,
      }),
      providesTags: ["Assets"],
      transformResponse: (response: ApiResponse<Asset[]>) => ({
        data: response.data,
        total: response.meta?.total || 0,
      }),
    }),
    
    getAssetById: builder.query<ApiResponse<Asset>, string>({
      query: (id) => `/api/v1/assets/${id}`,
      providesTags: (_result, _error, id) => [{ type: "Assets", id }],
    }),
    
    createAsset: builder.mutation<ApiResponse<Asset>, AssetFormInput>({
      query: (body) => ({
        url: "/api/v1/assets",
        method: "POST",
        body,
      }),
      invalidatesTags: ["Assets"],
    }),
    
    updateAsset: builder.mutation<ApiResponse<Asset>, { id: string; data: Partial<AssetFormInput> }>({
      query: ({ id, data }) => ({
        url: `/api/v1/assets/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: (_result, _error, { id }) => ["Assets", { type: "Assets", id }],
    }),
    
    deleteAsset: builder.mutation<ApiResponse<null>, string>({
      query: (id) => ({
        url: `/api/v1/assets/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Assets"],
    }),
  }),
});

export const {
  useGetAssetsQuery,
  useGetAssetByIdQuery,
  useCreateAssetMutation,
  useUpdateAssetMutation,
  useDeleteAssetMutation,
} = assetApiNew;
