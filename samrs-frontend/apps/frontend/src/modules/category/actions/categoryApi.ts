import { api } from "@/store/api";
import type { Category, CategoryFormInput } from "../types";

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

type CategoryListResponse = {
  data: Category[];
  total: number;
};

export const categoryApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getCategories: builder.query<CategoryListResponse, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({
        url: "/api/v1/categories",
        params,
      }),
      providesTags: ["MasterData"],
      transformResponse: (response: ApiResponse<Category[]>) => ({
        data: response.data,
        total: response.meta?.total || 0,
      }),
    }),
    
    getCategoryById: builder.query<{ data: Category }, number>({
      query: (id) => `/api/v1/categories/${id}`,
      providesTags: (_result, _error, id) => [{ type: "MasterData", id }],
    }),
    
    createCategory: builder.mutation<{ data: Category }, CategoryFormInput>({
      query: (body) => ({
        url: "/api/v1/categories",
        method: "POST",
        body,
      }),
      invalidatesTags: ["MasterData"],
    }),
    
    updateCategory: builder.mutation<{ data: Category }, { id: number; data: CategoryFormInput }>({
      query: ({ id, data }) => ({
        url: `/api/v1/categories/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: ["MasterData"],
    }),
    
    deleteCategory: builder.mutation<void, number>({
      query: (id) => ({
        url: `/api/v1/categories/${id}`,
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
