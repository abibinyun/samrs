import { api } from "@/store/api";
import type { Model, ModelFormInput } from "../types";

export const modelApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getModels: builder.query<{ data: Model[]; total: number }, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({ url: "/models", params }),
      providesTags: ["MasterData"],
    }),
    createModel: builder.mutation<{ data: Model }, ModelFormInput>({
      query: (body) => ({ url: "/models", method: "POST", body }),
      invalidatesTags: ["MasterData"],
    }),
    updateModel: builder.mutation<{ data: Model }, { id: number; data: ModelFormInput }>({
      query: ({ id, data }) => ({ url: `/models/${id}`, method: "PATCH", body: data }),
      invalidatesTags: ["MasterData"],
    }),
    deleteModel: builder.mutation<void, number>({
      query: (id) => ({ url: `/models/${id}`, method: "DELETE" }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const { useGetModelsQuery, useCreateModelMutation, useUpdateModelMutation, useDeleteModelMutation } = modelApi;
