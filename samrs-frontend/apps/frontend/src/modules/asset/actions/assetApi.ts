// src/store/api/assetApi.ts
import { createApi } from "@reduxjs/toolkit/query/react";
import { runAssetWorker } from "../mock/mocks";
import type { AssetDevice } from "../types";

export const assetApi = createApi({
  reducerPath: "assetApi",
  baseQuery: async () => ({ data: {} }),
  tagTypes: ["Assets"],
  endpoints: (builder) => ({
    getAssets: builder.query<AssetDevice[], void>({
      async queryFn() {
        try {
          // simulasi network
          await new Promise((r) => setTimeout(r, 1500));

          const data = await runAssetWorker();

          return { data };
        } catch (error) {
          return {
            error: {
              status: "CUSTOM_ERROR",
              error: String(error),
            },
          };
        }
      },
      providesTags: ["Assets"],
    }),
  }),
});

export const { useGetAssetsQuery, usePrefetch } = assetApi;
