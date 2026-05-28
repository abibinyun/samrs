import * as React from "react";
import { useDispatch } from "react-redux";
import { useRouter } from "@tanstack/react-router";
import { closeTab, setActive } from "@/modules/dashboard/tabs/actions/tabSlice";
import { useAppSelector } from "@/store/hooks";
import { cn } from "@/lib/utils"; // kalau kamu punya helper cn (shadcn). kalau belum, hapus cn dan pakai string biasa.

type RecentTab = {
  key: string;
  title: string;
  to: string;
};

export function RecentTabsBar() {
  const router = useRouter();
  const dispatch = useDispatch();

  const { tabs, activeKey } = useAppSelector(
    (s) => s.tabs as { tabs: RecentTab[]; activeKey: string | null },
  );

  const listRef = React.useRef<HTMLDivElement | null>(null);
  const tabRefs = React.useRef<Record<string, HTMLButtonElement | null>>({});

  const activeIndex = React.useMemo(() => {
    const i = tabs.findIndex((t) => t.key === activeKey);
    return i >= 0 ? i : 0;
  }, [tabs, activeKey]);

  const activate = React.useCallback(
    (tab: RecentTab) => {
      dispatch(setActive(tab.key));
      router.navigate({ to: tab.to });
    },
    [dispatch, router],
  );

  const close = React.useCallback(
    (key: string) => {
      dispatch(closeTab(key));
    },
    [dispatch],
  );

  const location = router.state.location;

  React.useEffect(() => {
    if (!activeKey) return;

    const activeTab = tabs.find((t) => t.key === activeKey);
    if (!activeTab) return;

    const currentUrl = location.pathname + (location.searchStr ?? "");
    if (currentUrl === activeTab.to) return;

    router.navigate({ to: activeTab.to });
  }, [activeKey, tabs, location.pathname, location.searchStr, router]);

  const onKeyDown = (e: React.KeyboardEvent<HTMLDivElement>) => {
    if (!tabs.length) return;

    const moveFocus = (nextIndex: number) => {
      const tab = tabs[Math.max(0, Math.min(tabs.length - 1, nextIndex))];
      tabRefs.current[tab.key]?.focus();
    };

    switch (e.key) {
      case "ArrowRight":
        e.preventDefault();
        moveFocus(activeIndex + 1);
        break;
      case "ArrowLeft":
        e.preventDefault();
        moveFocus(activeIndex - 1);
        break;
      case "Home":
        e.preventDefault();
        moveFocus(0);
        break;
      case "End":
        e.preventDefault();
        moveFocus(tabs.length - 1);
        break;
      case "Delete":
      case "Backspace":
        if (activeKey) {
          e.preventDefault();
          close(activeKey);
        }
        break;
    }
  };

  if (!tabs.length) return null;

  return (
    <div className="border-b bg-card sticky top-16 z-10">
      <div
        ref={listRef}
        className="flex gap-2 px-4 py-2"
        role="tablist"
        aria-label="Recent tabs"
        onKeyDown={onKeyDown}
      >
        {tabs.map((t) => {
          const isActive = t.key === activeKey;

          return (
            <div
              key={t.key}
              className={cn(
                "inline-flex items-center gap-2 rounded-md border px-2 py-1 shrink-0",
                isActive
                  ? "bg-muted/60 border-primary/40"
                  : "bg-background hover:bg-muted/40",
              )}
            >
              <button
                ref={(node) => {
                  tabRefs.current[t.key] = node;
                }}
                role="tab"
                aria-selected={isActive}
                tabIndex={isActive ? 0 : -1}
                title={t.title}
                className={cn(
                  "max-w-45 truncate text-sm outline-none",
                  isActive
                    ? "font-semibold text-foreground"
                    : "opacity-80 hover:opacity-100",
                )}
                onClick={() => activate(t)}
              >
                {t.title}
              </button>

              <button
                type="button"
                title="Close tab"
                aria-label={`Close ${t.title}`}
                className="h-6 w-6 rounded flex items-center justify-center leading-none opacity-60 hover:opacity-100 hover:bg-muted/60"
                onClick={() => close(t.key)}
              >
                ×
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
