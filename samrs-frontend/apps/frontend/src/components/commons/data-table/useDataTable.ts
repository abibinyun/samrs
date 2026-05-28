import { useMemo, useState, useCallback } from "react";
import type { DataTableProps, DataTableState } from "./types";
import type {
  SortingState,
  ColumnFiltersState,
  VisibilityState,
  PaginationState,
  RowSelectionState,
  Updater,
} from "@tanstack/react-table";

const defaultPagination: PaginationState = { pageIndex: 0, pageSize: 10 };

function resolveUpdater<T>(updater: Updater<T>, prev: T): T {
  return typeof updater === "function"
    ? (updater as (old: T) => T)(prev)
    : updater;
}

function coalesce<T>(controlled: T | undefined, internal: T): T {
  return controlled ?? internal;
}

export function useDataTableState<TData>(props: DataTableProps<TData>) {
  const { state: controlledState, defaultState, callbacks } = props;

  // === INTERNAL STATE (TABLE UI ONLY) ===
  const [sorting, setSorting] = useState<SortingState>(
    defaultState?.sorting ?? []
  );
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(
    defaultState?.columnFilters ?? []
  );
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>(
    defaultState?.columnVisibility ?? {}
  );
  const [rowSelection, setRowSelection] = useState<RowSelectionState>(
    defaultState?.rowSelection ?? {}
  );
  const [pagination, setPagination] = useState<PaginationState>(
    defaultState?.pagination ?? defaultPagination
  );

  // === INITIAL SNAPSHOT (UNTUK RESET) ===
  const initialState = useMemo(
    () => ({
      sorting: defaultState?.sorting ?? [],
      columnFilters: defaultState?.columnFilters ?? [],
      columnVisibility: defaultState?.columnVisibility ?? {},
      rowSelection: defaultState?.rowSelection ?? {},
      pagination: defaultState?.pagination ?? defaultPagination,
    }),
    [defaultState]
  );

  // === MERGED STATE (CONTROLLED + INTERNAL) ===
  const state: DataTableState = useMemo(
    () => ({
      sorting: coalesce(controlledState?.sorting, sorting),
      columnFilters: coalesce(controlledState?.columnFilters, columnFilters),
      columnVisibility: coalesce(
        controlledState?.columnVisibility,
        columnVisibility
      ),
      rowSelection: coalesce(controlledState?.rowSelection, rowSelection),
      pagination: coalesce(controlledState?.pagination, pagination),
    }),
    [
      controlledState,
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      pagination,
    ]
  );

  // === RESET ===
  const reset = useCallback(() => {
    callbacks?.onSortingChange?.(initialState.sorting);
    callbacks?.onColumnFiltersChange?.(initialState.columnFilters);
    callbacks?.onColumnVisibilityChange?.(initialState.columnVisibility);
    callbacks?.onRowSelectionChange?.(initialState.rowSelection);
    callbacks?.onPaginationChange?.(initialState.pagination);

    if (controlledState?.sorting === undefined) setSorting(initialState.sorting);
    if (controlledState?.columnFilters === undefined)
      setColumnFilters(initialState.columnFilters);
    if (controlledState?.columnVisibility === undefined)
      setColumnVisibility(initialState.columnVisibility);
    if (controlledState?.rowSelection === undefined)
      setRowSelection(initialState.rowSelection);
    if (controlledState?.pagination === undefined)
      setPagination(initialState.pagination);
  }, [callbacks, controlledState, initialState]);

  // === SETTERS ===
  return {
    state,

    setSorting: (updater: Updater<SortingState>) => {
      const next = resolveUpdater(updater, state.sorting);
      callbacks?.onSortingChange?.(next);
      if (controlledState?.sorting === undefined) setSorting(next);
    },

    setColumnFilters: (updater: Updater<ColumnFiltersState>) => {
      const next = resolveUpdater(updater, state.columnFilters);
      callbacks?.onColumnFiltersChange?.(next);
      if (controlledState?.columnFilters === undefined) setColumnFilters(next);
    },

    setColumnVisibility: (updater: Updater<VisibilityState>) => {
      const next = resolveUpdater(updater, state.columnVisibility);
      callbacks?.onColumnVisibilityChange?.(next);
      if (controlledState?.columnVisibility === undefined)
        setColumnVisibility(next);
    },

    setRowSelection: (updater: Updater<RowSelectionState>) => {
      const next = resolveUpdater(updater, state.rowSelection);
      callbacks?.onRowSelectionChange?.(next);
      if (controlledState?.rowSelection === undefined) setRowSelection(next);
    },

    setPagination: (updater: Updater<PaginationState>) => {
      const next = resolveUpdater(updater, state.pagination);
      callbacks?.onPaginationChange?.(next);
      if (controlledState?.pagination === undefined) setPagination(next);
    },

    reset,
  };
}
