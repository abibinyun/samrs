import { api } from "@/store/api";
import type { Complaint, ComplaintFormData, ComplaintUpdateData } from "../types";

type ComplaintListResponse = {
  success: boolean;
  data: Complaint[];
  meta?: {
    total: number;
    page: number;
    per_page: number;
  };
};

type ComplaintResponse = {
  success: boolean;
  data: Complaint;
};

export const complaintApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getComplaints: builder.query<ComplaintListResponse, {
      page?: number;
      per_page?: number;
      search?: string;
      status?: string;
      asset_id?: string;
    }>({
      query: (params) => ({
        url: "/api/v1/complaints",
        params,
      }),
      providesTags: ["Complaints"],
    }),

    getComplaint: builder.query<ComplaintResponse, string>({
      query: (id) => `/api/v1/complaints/${id}`,
      providesTags: (_result, _error, id) => [{ type: "Complaints", id }],
    }),

    createComplaint: builder.mutation<ComplaintResponse, ComplaintFormData>({
      query: (data) => ({
        url: "/api/v1/complaints",
        method: "POST",
        body: data,
      }),
      invalidatesTags: ["Complaints"],
    }),

    updateComplaint: builder.mutation<ComplaintResponse, { id: string; data: ComplaintUpdateData }>({
      query: ({ id, data }) => ({
        url: `/api/v1/complaints/${id}`,
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: (_result, _error, { id }) => [
        "Complaints",
        { type: "Complaints", id },
      ],
    }),

    deleteComplaint: builder.mutation<{ success: boolean }, string>({
      query: (id) => ({
        url: `/api/v1/complaints/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["Complaints"],
    }),
  }),
});

export const {
  useGetComplaintsQuery,
  useGetComplaintQuery,
  useCreateComplaintMutation,
  useUpdateComplaintMutation,
  useDeleteComplaintMutation,
} = complaintApi;
