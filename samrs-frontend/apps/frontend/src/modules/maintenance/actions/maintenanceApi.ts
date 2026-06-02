import { api } from "@/store/api";
import type { MaintenanceSchedule, MaintenanceDocument, MaintenanceFormData } from "../types";

type MaintenanceListResponse = {
  success: boolean;
  data: MaintenanceSchedule[];
  meta?: {
    total: number;
    page: number;
    per_page: number;
  };
};

type MaintenanceResponse = {
  success: boolean;
  data: MaintenanceSchedule;
};

type DocumentListResponse = {
  success: boolean;
  data: MaintenanceDocument[];
};

type DocumentResponse = {
  success: boolean;
  data: MaintenanceDocument;
};

export const maintenanceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getMaintenanceSchedules: builder.query<MaintenanceListResponse, {
      page?: number;
      per_page?: number;
      search?: string;
      status?: string;
      asset_id?: string;
    }>({
      query: (params) => ({
        url: "/api/v1/maintenance-schedules",
        params,
      }),
      providesTags: ["Maintenance"],
    }),

    getMaintenanceSchedule: builder.query<MaintenanceResponse, number>({
      query: (id) => `/api/v1/maintenance-schedules/${id}`,
      providesTags: (_result, _error, id) => [{ type: "Maintenance", id }],
    }),

    createMaintenanceSchedule: builder.mutation<MaintenanceResponse, MaintenanceFormData>({
      query: (data) => ({
        url: "/api/v1/maintenance-schedules",
        method: "POST",
        body: data,
      }),
      invalidatesTags: ["Maintenance"],
    }),

    updateMaintenanceSchedule: builder.mutation<MaintenanceResponse, { id: number; data: Partial<MaintenanceFormData> }>({
      query: ({ id, data }) => ({
        url: `/api/v1/maintenance-schedules/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: (_result, _error, { id }) => [
        "Maintenance",
        { type: "Maintenance", id },
      ],
    }),

    completeMaintenanceSchedule: builder.mutation<MaintenanceResponse, { id: number; note: string }>({
      query: ({ id, note }) => ({
        url: `/api/v1/maintenance-schedules/${id}/complete`,
        method: "PATCH",
        body: { note },
      }),
      invalidatesTags: (_result, _error, { id }) => [
        "Maintenance",
        { type: "Maintenance", id },
      ],
    }),

    deleteMaintenanceSchedule: builder.mutation<{ success: boolean }, number>({
      query: (id) => ({
        url: `/api/v1/maintenance-schedules/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Maintenance"],
    }),

    // Document endpoints
    getDocuments: builder.query<DocumentListResponse, number>({
      query: (scheduleId) => ({
        url: `/api/v1/maintenance-schedules/${scheduleId}/documents`,
      }),
      providesTags: ["MaintenanceDocuments"],
    }),

    uploadDocument: builder.mutation<DocumentResponse, { scheduleId: number; docType: string; file: File }>({
      query: ({ scheduleId, docType, file }) => {
        const formData = new FormData();
        formData.append("file", file);
        formData.append("doc_type", docType);
        return {
          url: `/api/v1/maintenance-schedules/${scheduleId}/documents`,
          method: "POST",
          body: formData,
        };
      },
      invalidatesTags: ["MaintenanceDocuments"],
    }),

    deleteDocument: builder.mutation<{ success: boolean }, number>({
      query: (id) => ({
        url: `/api/v1/maintenance-documents/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["MaintenanceDocuments"],
    }),
  }),
});

export const {
  useGetMaintenanceSchedulesQuery,
  useGetMaintenanceScheduleQuery,
  useCreateMaintenanceScheduleMutation,
  useUpdateMaintenanceScheduleMutation,
  useCompleteMaintenanceScheduleMutation,
  useDeleteMaintenanceScheduleMutation,
  useGetDocumentsQuery,
  useUploadDocumentMutation,
  useDeleteDocumentMutation,
} = maintenanceApi;
