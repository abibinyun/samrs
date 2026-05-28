import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import { AssetDevice } from "../types";

interface AssetUIState {
  search: string;
  category: string;
  selectedId: string | null;
  data: AssetDevice[];
}

const initialState: AssetUIState = {
  search: "",
  category: "__all",
  selectedId: null,
  data: [],
};

export const assetSlice = createSlice({
  name: "assetUI",
  initialState,
  reducers: {
    setSearch(state, action: PayloadAction<string>) {
      state.search = action.payload;
    },
    setCategory(state, action: PayloadAction<string>) {
      state.category = action.payload;
    },
    selectAsset(state, action: PayloadAction<string | null>) {
      state.selectedId = action.payload;
    },
    setData(state, action: PayloadAction<AssetDevice[]>) {
      state.data = action.payload;
    },
  },
});

export const { setData, setSearch, setCategory, selectAsset } =
  assetSlice.actions;

export default assetSlice.reducer;
