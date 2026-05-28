import { createTabKey } from "@/modules/dashboard/tabs/utils";
import { useRouterState } from "@tanstack/react-router";
import React from "react";

export function useTabKey() {
  const location = useRouterState({ select: (s) => s.location });
  return React.useMemo(
    () => createTabKey(location.pathname, location.searchStr),
    [location.pathname, location.searchStr]
  );
}
