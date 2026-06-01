import { configureStore } from "@reduxjs/toolkit";

import { api } from "./api";
import { assetApi } from "@/modules/asset/actions/assetApi";
import { mutationApi } from "@/modules/mutation/actions/mutationApi";

import authReducer from "@/modules/auth/actions/slice";
import assetReducer from "@/modules/asset/actions/assetSlice";
import tableStateReducer from "@/modules/dashboard/tabs/actions/tableStateSlice";
import assetListUiReducer from "@/modules/asset/actions/assetListUiSlice";
import mutationReducer from "@/modules/mutation/actions/mutationSlice";
import { formTemplatesApi } from "@/modules/forms/actions/formTemplatesApi";

export const store = configureStore({
  reducer: {
    auth: authReducer,

    // RTK Query APIs
    [api.reducerPath]: api.reducer,
    [assetApi.reducerPath]: assetApi.reducer,
    [mutationApi.reducerPath]: mutationApi.reducer,
    [formTemplatesApi.reducerPath]: formTemplatesApi.reducer,

    // UI slice
    assets: assetReducer,
    tableState: tableStateReducer,
    assetListUi: assetListUiReducer,
    mutations: mutationReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }).concat(
      api.middleware,
      assetApi.middleware,
      mutationApi.middleware,
      formTemplatesApi.middleware
    ),
});

// Expose store to window for debugging
if (typeof window !== 'undefined') {
  (window as any).store = store;
}

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
