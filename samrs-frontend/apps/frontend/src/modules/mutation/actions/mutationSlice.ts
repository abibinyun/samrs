// src/modules/assets/store/assetSlice.ts
import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface MutationUIState {
  search: string;
  category: string;
  selectedId: string | null;
  data: any[];
}

const initialState: MutationUIState = {
  search: "",
  category: "__all",
  selectedId: null,
  data: [],
};

export const mutationSlice = createSlice({
  name: "mutationUI",
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
    setData(state, action: PayloadAction<any | null>) {
      state.data = action.payload;
    },
  },
});

export const {
  setData,
  setSearch,
  setCategory,
  selectAsset,
} = mutationSlice.actions;

export default mutationSlice.reducer;
