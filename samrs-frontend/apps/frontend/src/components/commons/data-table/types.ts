import type {
  ColumnDef,
  SortingState,
  ColumnFiltersState,
  VisibilityState,
  PaginationState,
  RowSelectionState,
  Table,
  Row,
} from "@tanstack/react-table";
import type { ReactNode } from "react";

export type DataTableState = {
  sorting: SortingState;
  columnFilters: ColumnFiltersState;
  columnVisibility: VisibilityState;
  rowSelection: RowSelectionState;
  pagination: PaginationState;
};

export type DataTableCallbacks = Partial<{
  onSortingChange: (next: SortingState) => void;
  onColumnFiltersChange: (next: ColumnFiltersState) => void;
  onGlobalFilterChange: (next: string) => void;
  onColumnVisibilityChange: (next: VisibilityState) => void;
  onRowSelectionChange: (next: RowSelectionState) => void;
  onPaginationChange: (next: PaginationState) => void;
}>;

export type DataTableVirtualOptions = {
  enabled?: boolean;
  estimateRowHeight?: number;
  overscan?: number;
  maxBodyHeight?: number;
};

export type DataTableProps<TData> = {
  data: TData[];
  columns: ColumnDef<TData, unknown>[];

  // states (optional controlled)
  state?: Partial<DataTableState>;
  defaultState?: Partial<DataTableState>;
  callbacks?: DataTableCallbacks;

  // features
  enableRowSelection?: boolean;
  enableMultiRowSelection?: boolean;
  enableGlobalFilter?: boolean;

  // server-side mode flags (kalau nanti kamu butuh)
  manualPagination?: boolean;
  manualSorting?: boolean;
  manualFiltering?: boolean;
  pageCount?: number;

  // ui
  isLoading?: boolean;
  className?: string;
  tableClassName?: string;

  // options
  enableColumnResizing?: boolean;
  columnResizeMode?: "onChange" | "onEnd";

  // virtual scroll
  virtual?: DataTableVirtualOptions;

  globalSearchTextFn?: (row: TData) => string;

  // slots
  renderToolbar?: (props: {
    table: Table<TData>;
    resetTable: () => void;
  }) => React.ReactNode;
  renderPagination?: (table: Table<TData>) => ReactNode;
  renderRowActions?: (row: Row<TData>) => ReactNode;
  renderEmpty?: () => ReactNode;
  renderLoading?: () => ReactNode;

  isHeaderSticky?: boolean;
};
