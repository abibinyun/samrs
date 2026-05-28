import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

export type RecentTab = {
  key: string;      // unik per halaman (path + search)
  title: string;    // label tab
  to: string;       // path untuk navigate
};

type TabsState = {
  activeKey: string | null;
  tabs: RecentTab[];
  max: number;
};

const initialState: TabsState = {
  activeKey: null,
  tabs: [],
  max: 8, // bebas (mis. 8 tab terakhir)
};

const tabsSlice = createSlice({
  name: "tabs",
  initialState,
  reducers: {
    upsertTab(state, action: PayloadAction<RecentTab>) {
      const tab = action.payload;

      // remove existing
      state.tabs = state.tabs.filter((t) => t.key !== tab.key);
      // push to end = paling recent
      state.tabs.push(tab);

      // limit
      if (state.tabs.length > state.max) state.tabs.shift();

      state.activeKey = tab.key;
    },
    setActive(state, action: PayloadAction<string>) {
      state.activeKey = action.payload;
    },
    closeTab(state, action: PayloadAction<string>) {
      const key = action.payload;
      state.tabs = state.tabs.filter((t) => t.key !== key);
      if (state.activeKey === key) {
        state.activeKey = state.tabs[state.tabs.length - 1]?.key ?? null;
      }
    },
    resetTabs(state) {
      state.tabs = [];
      state.activeKey = null;
    },
  },
});

export const { upsertTab, setActive, closeTab, resetTabs } = tabsSlice.actions;
export default tabsSlice.reducer;
