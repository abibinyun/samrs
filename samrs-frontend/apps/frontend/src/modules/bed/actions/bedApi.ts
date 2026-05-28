import { api } from "@/store/api";
import type { Bed, BedFormInput } from "../types";

export const bedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getBeds: builder.query<{ data: Bed[]; total: number }, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({ url: "/beds", params }),
      providesTags: ["MasterData"],
    }),
    createBed: builder.mutation<{ data: Bed }, BedFormInput>({
      query: (body) => ({ url: "/beds", method: "POST", body }),
      invalidatesTags: ["MasterData"],
    }),
    updateBed: builder.mutation<{ data: Bed }, { id: number; data: BedFormInput }>({
      query: ({ id, data }) => ({ url: `/beds/${id}`, method: "PATCH", body: data }),
      invalidatesTags: ["MasterData"],
    }),
    deleteBed: builder.mutation<void, number>({
      query: (id) => ({ url: `/beds/${id}`, method: "DELETE" }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const { useGetBedsQuery, useCreateBedMutation, useUpdateBedMutation, useDeleteBedMutation } = bedApi;
