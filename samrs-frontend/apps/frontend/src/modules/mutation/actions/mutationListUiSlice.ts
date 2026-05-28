import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { ColumnFiltersState, VisibilityState } from "@tanstack/react-table";

type DateRangeISO = { from?: string; to?: string };

export type MutationListUI = {
  search: string;
  category: string;
  dateRange?: DateRangeISO;

  columnVisibility: VisibilityState;
  columnFilters: ColumnFiltersState;
};

type MutationListUIState = {
  byTabKey: Record<string, MutationListUI>;
};

const defaultUI: MutationListUI = {
  search: "",
  category: "__all",
  dateRange: undefined,
  columnVisibility: {},
  columnFilters: [], 
};

const initialState: MutationListUIState = { byTabKey: {} };

const slice = createSlice({
  name: "mutationListUI",
  initialState,
  reducers: {
    patchAssetListUI(
      state,
      action: PayloadAction<{ tabKey: string; patch: Partial<MutationListUI> }>
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
