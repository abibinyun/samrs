import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { DateRange } from "react-day-picker";

type TableState = {
  // filter UI
  search?: string;
  category?: string;
  dateRange?: DateRange;

  // table state
  globalFilter?: string;
  columnFilters?: any[];
  sorting?: any[];
  pagination?: { pageIndex: number; pageSize: number };
  columnVisibility?: Record<string, boolean>;
  rowSelection?: Record<string, boolean>;
};

type State = {
  byKey: Record<string, TableState>;
};

const initialState: State = { byKey: {} };

const tableStateSlice = createSlice({
  name: "tableState",
  initialState,
  reducers: {
    setTableState(
      state,
      action: PayloadAction<{ key: string; patch: Partial<TableState> }>
    ) {
      const { key, patch } = action.payload;
      state.byKey[key] = { ...state.byKey[key], ...patch };
    },
    resetTableState(state, action: PayloadAction<{ key: string }>) {
      delete state.byKey[action.payload.key];
    },
  },
});

export const { setTableState, resetTableState } = tableStateSlice.actions;
export default tableStateSlice.reducer;
