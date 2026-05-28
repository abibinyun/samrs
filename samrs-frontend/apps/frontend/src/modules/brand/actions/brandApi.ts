import { api } from "@/store/api";
import type { Brand, BrandFormInput } from "../types";

export const brandApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getBrands: builder.query<{ data: Brand[]; total: number }, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({ url: "/brands", params }),
      providesTags: ["MasterData"],
    }),
    createBrand: builder.mutation<{ data: Brand }, BrandFormInput>({
      query: (body) => ({ url: "/brands", method: "POST", body }),
      invalidatesTags: ["MasterData"],
    }),
    updateBrand: builder.mutation<{ data: Brand }, { id: number; data: BrandFormInput }>({
      query: ({ id, data }) => ({ url: `/brands/${id}`, method: "PATCH", body: data }),
      invalidatesTags: ["MasterData"],
    }),
    deleteBrand: builder.mutation<void, number>({
      query: (id) => ({ url: `/brands/${id}`, method: "DELETE" }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const { useGetBrandsQuery, useCreateBrandMutation, useUpdateBrandMutation, useDeleteBrandMutation } = brandApi;
