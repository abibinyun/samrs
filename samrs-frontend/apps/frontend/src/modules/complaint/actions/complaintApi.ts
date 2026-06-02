import { api } from "@/store/api";
import { API_ENDPOINTS } from "@/constants/endpoints";
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

type UserOption = {
  id: string;
  username: string;
};

type UserListResponse = {
  success: boolean;
  data: UserOption[];
  meta?: {
    total: number;
    page: number;
    per_page: number;
  };
};

export const complaintApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getUsers: builder.query<UserListResponse, void>({
      query: () => ({
        url: API_ENDPOINTS.USERS.BASE,
        params: { per_page: 100 },
      }),
    }),

    getComplaints: builder.query<ComplaintListResponse, {
      page?: number;
      per_page?: number;
      search?: string;
      status?: string;
      asset_id?: string;
    }>({
      query: (params) => ({
        url: API_ENDPOINTS.COMPLAINTS.BASE,
        params,
      }),
      providesTags: ["Complaints"],
    }),

    getComplaint: builder.query<ComplaintResponse, string>({
      query: (id) => API_ENDPOINTS.COMPLAINTS.DETAIL(id),
      providesTags: (_result, _error, id) => [{ type: "Complaints", id }],
    }),

    createComplaint: builder.mutation<ComplaintResponse, ComplaintFormData>({
      query: (data) => ({
        url: API_ENDPOINTS.COMPLAINTS.BASE,
        method: "POST",
        body: data,
      }),
      invalidatesTags: ["Complaints"],
    }),

    updateComplaint: builder.mutation<ComplaintResponse, { id: string; data: ComplaintUpdateData }>({
      query: ({ id, data }) => ({
        url: API_ENDPOINTS.COMPLAINTS.DETAIL(id),
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
        url: API_ENDPOINTS.COMPLAINTS.DETAIL(id),
        method: "DELETE",
      }),
      invalidatesTags: ["Complaints"],
    }),
  }),
});

export const {
  useGetUsersQuery,
  useGetComplaintsQuery,
  useGetComplaintQuery,
  useCreateComplaintMutation,
  useUpdateComplaintMutation,
  useDeleteComplaintMutation,
} = complaintApi;
