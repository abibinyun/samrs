import * as React from "react";
import { Outlet, useRouterState } from "@tanstack/react-router";
import { useDispatch } from "react-redux";
import { upsertTab } from "../actions/tabSlice";
import { useAppSelector } from "@/store/hooks";

// helper: bikin key unik
const createKey = (pathname: string, searchStr: string) => `${pathname}${searchStr || ""}`;

// helper title sederhana (nanti bisa pakai map dari Sidebar menu)
const defaultTitleFromPath = (pathname: string) => {
  const parts = pathname.split("/").filter(Boolean);
  return parts[parts.length - 1] ?? "dashboard";
};

type CachedView = { key: string; element: React.ReactNode };

export function KeepAliveOutlet() {
  const dispatch = useDispatch();

  const location = useRouterState({ select: (s) => s.location });
  const pathname = location.pathname;
  const searchStr = location.searchStr ?? "";
  const key = React.useMemo(() => createKey(pathname, searchStr), [pathname, searchStr]);

  const tabsState = useAppSelector((s: any) => s.tabs);
  const activeKey = tabsState.activeKey ?? key;

  // cache element per tab key (di memori)
  const cacheRef = React.useRef<Map<string, CachedView>>(new Map());

  // update tab list + cache saat route berubah
  React.useEffect(() => {
    // simpan element untuk key ini jika belum ada
    if (!cacheRef.current.has(key)) {
      cacheRef.current.set(key, { key, element: <Outlet key={key} /> });
    }

    // upsert tab ke redux (recent)
    dispatch(
      upsertTab({
        key,
        title: defaultTitleFromPath(pathname), // nanti bisa diubah jadi label menu beneran
        to: pathname + searchStr,
      })
    );
  }, [key, pathname, searchStr, dispatch]);

  // render semua cached views, hide yang tidak aktif
  const cached = Array.from(cacheRef.current.values());

  return (
    <div className="relative">
      {cached.map((view) => (
        <div key={view.key} style={{ display: view.key === activeKey ? "block" : "none" }}>
          {view.element}
        </div>
      ))}
    </div>
  );
}