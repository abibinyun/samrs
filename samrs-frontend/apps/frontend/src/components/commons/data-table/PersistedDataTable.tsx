import * as React from "react";
import type {
  ColumnDef,
  ColumnFiltersState,
  Table,
} from "@tanstack/react-table";
import DataTable from "./DataTable";
import { Virtualizer } from "@tanstack/react-virtual";

type PersistedDataTableProps<TData> = {
  data: TData[];
  columns: ColumnDef<TData, unknown>[];
  isLoading?: boolean;
  columnVisibility?: Record<string, boolean>;
  onColumnVisibilityChange?: (visibility: Record<string, boolean>) => void;
  columnFilters?: ColumnFiltersState;
  onColumnFiltersChange?: (filters: ColumnFiltersState) => void;
  enableRowSelection?: boolean;
  virtual?: {
    enabled?: boolean;
    estimateRowHeight?: number;
    maxBodyHeight?: number;
    rowVirtualizer?: Virtualizer<HTMLDivElement, HTMLDivElement>;
    overscan?: number;
  };
  renderToolbar?: (props: {
    table: Table<TData>;
    resetTable: () => void;
  }) => React.ReactNode;
  onResetExternal?: () => void;
};

export default function PersistedDataTable<TData>({
  data,
  columns,
  isLoading,
  columnVisibility,
  onColumnVisibilityChange,
  columnFilters,
  onColumnFiltersChange,
  enableRowSelection = true,
  virtual = {
    enabled: true,
    estimateRowHeight: 56,
    maxBodyHeight: 640,
    overscan: 20,
  },
  renderToolbar,
  onResetExternal,
}: PersistedDataTableProps<TData>) {
  return (
    <DataTable
      data={data}
      columns={columns}
      isLoading={isLoading}
      enableRowSelection={enableRowSelection}
      virtual={virtual}
      defaultState={{
        columnVisibility: columnVisibility ?? {},
        columnFilters: columnFilters ?? [],
      }}
      callbacks={{
        onColumnVisibilityChange,
        onColumnFiltersChange,
      }}
      renderToolbar={({ table, resetTable }) => {
        const resetAll = () => {
          resetTable();
          table.resetColumnVisibility();
          table.resetColumnFilters();
          table.resetSorting?.();
          table.resetPagination?.();
          onResetExternal?.();
        };
        return renderToolbar
          ? renderToolbar({ table, resetTable: resetAll })
          : null;
      }}
    />
  );
}
