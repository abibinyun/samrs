import { api } from "@/store/api";
import type { StockOpnameSession, StockOpnameItem, StockOpnameFormData, StockOpnameItemFormData } from "../types";

type SessionListResponse = { success: boolean; data: StockOpnameSession[]; };
type SessionResponse = { success: boolean; data: StockOpnameSession; };
type ItemListResponse = { success: boolean; data: StockOpnameItem[]; };

export const stockOpnameApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getStockOpnameSessions: builder.query<SessionListResponse, void>({
      query: () => "/api/v1/stock-opnames",
      providesTags: ["StockOpname"],
    }),
    getStockOpnameSession: builder.query<SessionResponse, number>({
      query: (id) => `/api/v1/stock-opnames/${id}`,
      providesTags: (_r, _e, id) => [{ type: "StockOpname", id }],
    }),
    createStockOpnameSession: builder.mutation<SessionResponse, StockOpnameFormData>({
      query: (data) => ({ url: "/api/v1/stock-opnames", method: "POST", body: data }),
      invalidatesTags: ["StockOpname"],
    }),
    closeStockOpnameSession: builder.mutation<SessionResponse, number>({
      query: (id) => ({ url: `/api/v1/stock-opnames/${id}/close`, method: "PATCH" }),
      invalidatesTags: (_r, _e, id) => ["StockOpname", { type: "StockOpname", id }],
    }),
    getStockOpnameItems: builder.query<ItemListResponse, number>({
      query: (sessionId) => `/api/v1/stock-opnames/${sessionId}/items`,
      providesTags: (_r, _e, id) => [{ type: "StockOpname", id: `items-${id}` }],
    }),
    addStockOpnameItem: builder.mutation<{ success: boolean }, { sessionId: number; data: StockOpnameItemFormData }>({
      query: ({ sessionId, data }) => ({ url: `/api/v1/stock-opnames/${sessionId}/items`, method: "POST", body: data }),
      invalidatesTags: (_r, _e, { sessionId }) => [{ type: "StockOpname", id: `items-${sessionId}` }],
    }),
  }),
});

export const {
  useGetStockOpnameSessionsQuery,
  useGetStockOpnameSessionQuery,
  useCreateStockOpnameSessionMutation,
  useCloseStockOpnameSessionMutation,
  useGetStockOpnameItemsQuery,
  useAddStockOpnameItemMutation,
} = stockOpnameApi;
