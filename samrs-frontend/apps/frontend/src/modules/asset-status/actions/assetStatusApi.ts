import { api } from "@/store/api";
import type { AssetStatus, AssetStatusFormInput } from "../types";

export const assetStatusApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getAssetStatuses: builder.query<{ data: AssetStatus[]; total: number }, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({ url: "/asset-statuses", params }),
      providesTags: ["MasterData"],
    }),
    createAssetStatus: builder.mutation<{ data: AssetStatus }, AssetStatusFormInput>({
      query: (body) => ({ url: "/asset-statuses", method: "POST", body }),
      invalidatesTags: ["MasterData"],
    }),
    updateAssetStatus: builder.mutation<{ data: AssetStatus }, { id: number; data: AssetStatusFormInput }>({
      query: ({ id, data }) => ({ url: `/asset-statuses/${id}`, method: "PATCH", body: data }),
      invalidatesTags: ["MasterData"],
    }),
    deleteAssetStatus: builder.mutation<void, number>({
      query: (id) => ({ url: `/asset-statuses/${id}`, method: "DELETE" }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const { useGetAssetStatusesQuery, useCreateAssetStatusMutation, useUpdateAssetStatusMutation, useDeleteAssetStatusMutation } = assetStatusApi;
