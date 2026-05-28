import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { ColumnFiltersState, VisibilityState } from "@tanstack/react-table";

type DateRangeISO = { from?: string; to?: string };

export type AssetListUI = {
  search: string;
  category: string;
  dateRange?: DateRangeISO;

  columnVisibility: VisibilityState;
  columnFilters: ColumnFiltersState;
};

type AssetListUIState = {
  byTabKey: Record<string, AssetListUI>;
};

const defaultUI: AssetListUI = {
  search: "",
  category: "__all",
  dateRange: undefined,
  columnVisibility: {},
  columnFilters: [], 
};

const initialState: AssetListUIState = { byTabKey: {} };

const slice = createSlice({
  name: "assetListUI",
  initialState,
  reducers: {
    patchAssetListUI(
      state,
      action: PayloadAction<{ tabKey: string; patch: Partial<AssetListUI> }>
    ) {
      const { tabKey, patch } = action.payload;

      state.byTabKey[tabKey] = {
        ...(state.byTabKey[tabKey] ?? defaultUI),
        ...patch,
        columnVisibility: patch.columnVisibility ?? state.byTabKey[tabKey]?.columnVisibility ?? {},
        columnFilters: patch.columnFilters ?? state.byTabKey[tabKey]?.columnFilters ?? [],
      };
    },

    resetAssetListUIToDefault(state, action: PayloadAction<{ tabKey: string }>) {
      state.byTabKey[action.payload.tabKey] = { ...defaultUI };
    },

    removeAssetListUIForTab(state, action: PayloadAction<{ tabKey: string }>) {
      delete state.byTabKey[action.payload.tabKey];
    },
  },
});

export const {
  patchAssetListUI,
  resetAssetListUIToDefault,
  removeAssetListUIForTab,
} = slice.actions;

export { defaultUI };
export default slice.reducer;
