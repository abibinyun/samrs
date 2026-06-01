import { configureStore } from "@reduxjs/toolkit";

import { api } from "./api";
import { assetApi } from "@/modules/asset/actions/assetApi";
import { formTemplatesApi } from "@/modules/forms/actions/formTemplatesApi";

import authReducer from "@/modules/auth/actions/slice";
import assetReducer from "@/modules/asset/actions/assetSlice";
import tableStateReducer from "@/modules/dashboard/tabs/actions/tableStateSlice";
import assetListUiReducer from "@/modules/asset/actions/assetListUiSlice";

export const store = configureStore({
  reducer: {
    auth: authReducer,

    // RTK Query APIs
    [api.reducerPath]: api.reducer,
    [assetApi.reducerPath]: assetApi.reducer,
    [formTemplatesApi.reducerPath]: formTemplatesApi.reducer,

    // UI slice
    assets: assetReducer,
    tableState: tableStateReducer,
    assetListUi: assetListUiReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }).concat(
      api.middleware,
      assetApi.middleware,
      formTemplatesApi.middleware
    ),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
