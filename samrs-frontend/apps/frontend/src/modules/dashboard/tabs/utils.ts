export const createTabKey = (pathname: string, searchStr?: string) =>
  `${pathname}${searchStr || ""}`;
