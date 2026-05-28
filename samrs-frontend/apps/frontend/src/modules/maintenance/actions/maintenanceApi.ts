import { api } from "@/store/api";
import type { MaintenanceSchedule, MaintenanceFormData } from "../types";

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
  }),
});

export const {
  useGetMaintenanceSchedulesQuery,
  useGetMaintenanceScheduleQuery,
  useCreateMaintenanceScheduleMutation,
  useUpdateMaintenanceScheduleMutation,
  useCompleteMaintenanceScheduleMutation,
  useDeleteMaintenanceScheduleMutation,
} = maintenanceApi;
