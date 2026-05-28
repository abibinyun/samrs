import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { FormTemplate } from "../types";

type CreateTemplateDto = {
  name: string;
  description?: string;
  schema: any;
};

type UpdateTemplateDto = Partial<Pick<FormTemplate, "name" | "description" | "schema">>;

export const formTemplatesApi = createApi({
  reducerPath: "formTemplatesApi",
  baseQuery: fetchBaseQuery({ baseUrl: "/api" }),
  tagTypes: ["FormTemplates"],
  endpoints: (b) => ({
    listTemplates: b.query<FormTemplate[], void>({
      query: () => `/form-templates`,
      providesTags: ["FormTemplates"],
    }),

    getTemplate: b.query<FormTemplate, string>({
      query: (id) => `/form-templates/${id}`,
      providesTags: ["FormTemplates"],
    }),

    createTemplate: b.mutation<FormTemplate, CreateTemplateDto>({
      query: (body) => ({
        url: `/form-templates`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["FormTemplates"],
    }),

    updateTemplate: b.mutation<FormTemplate, { id: string; patch: UpdateTemplateDto }>({
      query: ({ id, patch }) => ({
        url: `/form-templates/${id}`,
        method: "PATCH",
        body: patch,
      }),
      invalidatesTags: ["FormTemplates"],
    }),

    publishTemplate: b.mutation<FormTemplate, { id: string }>({
      query: ({ id }) => ({
        url: `/form-templates/${id}/publish`,
        method: "POST",
      }),
      invalidatesTags: ["FormTemplates"],
    }),

    deleteTemplate: b.mutation<{ ok: true }, { id: string }>({
      query: ({ id }) => ({
        url: `/form-templates/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["FormTemplates"],
    }),
  }),
});

export const {
  useListTemplatesQuery,
  useGetTemplateQuery,
  useCreateTemplateMutation,
  useUpdateTemplateMutation,
  usePublishTemplateMutation,
  useDeleteTemplateMutation,
} = formTemplatesApi;