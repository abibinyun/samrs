import { api } from "@/store/api";
import type { Room, RoomFormInput } from "../types";

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

type RoomListResponse = {
  data: Room[];
  total: number;
};

export const roomApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getRooms: builder.query<RoomListResponse, { page?: number; limit?: number; search?: string }>({
      query: (params) => ({ url: "/api/v1/rooms", params }),
      providesTags: ["MasterData"],
      transformResponse: (response: ApiResponse<Room[]>) => ({
        data: response.data,
        total: response.meta?.total || 0,
      }),
    }),
    createRoom: builder.mutation<{ data: Room }, RoomFormInput>({
      query: (body) => ({ url: "/api/v1/rooms", method: "POST", body }),
      invalidatesTags: ["MasterData"],
    }),
    updateRoom: builder.mutation<{ data: Room }, { id: string; data: RoomFormInput }>({
      query: ({ id, data }) => ({ url: `/api/v1/rooms/${id}`, method: "PATCH", body: data }),
      invalidatesTags: ["MasterData"],
    }),
    deleteRoom: builder.mutation<void, string>({
      query: (id) => ({ url: `/api/v1/rooms/${id}`, method: "DELETE" }),
      invalidatesTags: ["MasterData"],
    }),
  }),
});

export const { useGetRoomsQuery, useCreateRoomMutation, useUpdateRoomMutation, useDeleteRoomMutation } = roomApi;
