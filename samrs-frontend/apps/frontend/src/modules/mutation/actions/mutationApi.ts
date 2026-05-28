// src/store/api/assetApi.ts
import { createApi } from "@reduxjs/toolkit/query/react";
import { runMutationWorker } from "../mock/mocks"; 
// import type { AssetDevice } from "../types"; 

export const mutationApi = createApi({
  reducerPath: "mutationApi",
  baseQuery: async () => ({ data: {} }), // dummy
  tagTypes: ["Mutations"],
  endpoints: (builder) => ({
    // getAssets: builder.query<AssetDevice[], void>({
    getMutations: builder.query<any[], void>({
      async queryFn() {
        try {
          // simulasi network
          await new Promise(r => setTimeout(r, 1500));

          const data = await runMutationWorker();

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
      providesTags: ["Mutations"],
    }),
  }),
});

export const {
  useGetMutationsQuery,
  usePrefetch,
} = mutationApi;