import { api } from "@/store/api";
import type { Category, CategoryFormInput } from "../types";

export const categoryApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getCategories: builder.query<{ data: Category[]; total: number }, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({
        url: "/categories",
        params,
      }),
      providesTags: ["MasterData"],
    }),
    
    getCategoryById: builder.query<{ data: Category }, number>({
      query: (id) => `/categories/${id}`,
      providesTags: (_result, _error, id) => [{ type: "MasterData", id }],
    }),
    
    createCategory: builder.mutation<{ data: Category }, CategoryFormInput>({
      query: (body) => ({
        url: "/categories",
        method: "POST",
        body,
      }),
      invalidatesTags: ["MasterData"],
    }),
    
    updateCategory: builder.mutation<{ data: Category }, { id: number; data: CategoryFormInput }>({
      query: ({ id, data }) => ({
        url: `/categories/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: ["MasterData"],
    }),
    
    deleteCategory: builder.mutation<void, number>({
      query: (id) => ({
        url: `/categories/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const {
  useGetCategoriesQuery,
  useGetCategoryByIdQuery,
  useCreateCategoryMutation,
  useUpdateCategoryMutation,
  useDeleteCategoryMutation,
} = categoryApi;
