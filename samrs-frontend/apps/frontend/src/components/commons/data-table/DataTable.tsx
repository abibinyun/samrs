import * as React from "react";
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  getPaginationRowModel,
  useReactTable,
  HeaderContext,
  CellContext,
  Row,
  ColumnDef,
  DisplayColumnDef,
  Table,
} from "@tanstack/react-table";
import { defaultRangeExtractor, useVirtualizer } from "@tanstack/react-virtual";
import type { Range, VirtualItem } from "@tanstack/react-virtual";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";
import { useDataTableState } from "./useDataTable";
import type { DataTableProps } from "./types";

function isVirtualItem<T>(item: Row<T> | VirtualItem): item is VirtualItem {
  return "start" in item;
}

type DataTablePaginationComponent = <TData>(props: {
  table: Table<TData>;
}) => React.JSX.Element;

const DataTablePagination = React.lazy(
  () => import("./DataTablePagination"),
) as unknown as DataTablePaginationComponent;

export default function DataTable<TData>(props: DataTableProps<TData>) {
  const {
    data,
    columns,
    enableRowSelection = false,
    enableMultiRowSelection = true,
    manualPagination,
    manualSorting,
    manualFiltering,
    pageCount,
    isLoading,
    className,
    tableClassName,
    enableColumnResizing = true,
    columnResizeMode = "onChange",
    renderToolbar,
    renderPagination,
    renderRowActions,
    renderEmpty,
    renderLoading,
  } = props;

  const s = useDataTableState(props);

  // === 1. COLUMN DEFINITIONS ===
  const selectionColumn = React.useMemo<DisplayColumnDef<TData>[]>(() => {
    if (!enableRowSelection) return [];
    return [
      {
        id: "__select",
        header: ({ table }: HeaderContext<TData, unknown>) => (
          <Checkbox
            checked={
              table.getIsAllPageRowsSelected() ||
              (table.getIsSomePageRowsSelected() && "indeterminate")
            }
            onCheckedChange={(value) =>
              table.toggleAllPageRowsSelected(!!value)
            }
            aria-label="Select all rows"
            className="border border-gray-800 dark:border-gray-500"
          />
        ),
        cell: ({ row }: CellContext<TData, unknown>) => (
          <Checkbox
            checked={row.getIsSelected()}
            onCheckedChange={(value) => row.toggleSelected(!!value)}
            aria-label="Select row"
            className="border border-gray-700 dark:border-gray-600"
          />
        ),
        enableSorting: false,
        enableHiding: false,
        enableResizing: false,
        size: 40,
      },
    ];
  }, [enableRowSelection]);

  const actionsColumn = React.useMemo<DisplayColumnDef<TData>[]>(() => {
    if (!renderRowActions) return [];
    return [
      {
        id: "__actions",
        header: () => null,
        cell: ({ row }: CellContext<TData, unknown>) => (
          <div className="flex justify-end">{renderRowActions(row)}</div>
        ),
        enableSorting: false,
        enableHiding: false,
        enableResizing: false,
        size: 60,
      },
    ];
  }, [renderRowActions]);

  const finalColumns = React.useMemo<ColumnDef<TData, unknown>[]>(
    () => [...selectionColumn, ...columns, ...actionsColumn],
    [selectionColumn, columns, actionsColumn],
  );

  // === 2. TABLE INSTANCE ===
  const table = useReactTable({
    data,
    columns: finalColumns,
    state: s.state,
    enableRowSelection,
    enableMultiRowSelection,
    onSortingChange: s.setSorting,
    onColumnFiltersChange: s.setColumnFilters,
    // onGlobalFilterChange: s.setGlobalFilter,
    onColumnVisibilityChange: s.setColumnVisibility,
    onRowSelectionChange: s.setRowSelection,
    onPaginationChange: s.setPagination,
    manualPagination,
    manualSorting,
    manualFiltering,
    pageCount,
    enableColumnResizing,
    columnResizeMode,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: manualFiltering ? undefined : getFilteredRowModel(),
    getSortedRowModel: manualSorting ? undefined : getSortedRowModel(),
    getPaginationRowModel: manualPagination
      ? undefined
      : getPaginationRowModel(),
  });

  const virtual = React.useMemo(() => props.virtual ?? {}, [props.virtual]);
  const isVirtual = !!virtual.enabled;
  const rows = table.getRowModel().rows;
  const hasRows = rows.length > 0;

  const scrollYRef = React.useRef<HTMLDivElement | null>(null);

  // === 3. VIRTUAL STICKY (TanStack Virtual style) ===
  // Non-breaking: kalau tidak diisi, stickyIndexes = []
  // Kamu bisa pakai:
  // - props.virtual.stickyRowIndexes?: number[]
  // - props.virtual.isStickyRow?: (row, index) => boolean
  const stickyIndexes = React.useMemo(() => {
    const v = virtual as {
      stickyRowIndexes?: number[];
      isStickyRow?: (row: Row<TData>, index: number) => boolean;
    };

    if (Array.isArray(v.stickyRowIndexes)) {
      return v.stickyRowIndexes.filter(
        (n): n is number => Number.isFinite(n) && n >= 0,
      );
    }

    if (typeof v.isStickyRow === "function") {
      const out: number[] = [];
      for (let i = 0; i < rows.length; i++) {
        if (v.isStickyRow(rows[i], i)) out.push(i);
      }
      return out;
    }

    return [];
  }, [rows, virtual]);

  const activeStickyIndexRef = React.useRef(0);

  const isSticky = React.useCallback(
    (index: number) => stickyIndexes.includes(index),
    [stickyIndexes],
  );

  const isActiveSticky = React.useCallback(
    (index: number) => activeStickyIndexRef.current === index,
    [],
  );

  // === 4. VIRTUALIZER ===
  const rowVirtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => (isVirtual ? scrollYRef.current : null),
    estimateSize: () => virtual.estimateRowHeight ?? 48,
    overscan: virtual.overscan ?? 10,
    measureElement: (el) => el.getBoundingClientRect().height,
    enabled: isVirtual,
    rangeExtractor: React.useCallback(
      (range: Range) => {
        if (!stickyIndexes.length) {
          return defaultRangeExtractor(range);
        }

        activeStickyIndexRef.current =
          [...stickyIndexes]
            .reverse()
            .find((index) => range.startIndex >= index) ?? 0;

        const next = new Set<number>([
          activeStickyIndexRef.current,
          ...defaultRangeExtractor(range),
        ]);

        return [...next].sort((a, b) => a - b);
      },
      [stickyIndexes],
    ),
  });

  // === 5. LAYOUT CALCULATIONS (The "Stabilizer") ===
  const visibleLeaf = table.getVisibleLeafColumns();

  // Menggunakan CSS Variables untuk performa resize yang mulus
  const gridTemplateColumns = React.useMemo(() => {
    return visibleLeaf.map((c) => `${c.getSize()}px`).join(" ");
  }, [visibleLeaf]);

  const totalTableWidth = React.useMemo(() => {
    return visibleLeaf.reduce((acc, c) => acc + c.getSize(), 0);
  }, [visibleLeaf]);

  return (
    <div className={cn("space-y-3 sticky top-40", className)}>
      {renderToolbar
        ? renderToolbar({
            table,
            resetTable: s.reset,
          })
        : null}

      <div
        className={cn(
          "rounded-lg border bg-card overflow-hidden shadow-sm",
          tableClassName,
        )}
      >
        {/* Horizontal Scroll Container */}
        <div className="overflow-auto scrollbar-thin">
          <div style={{ width: totalTableWidth, minWidth: "100%" }}>
            {/* HEADER AREA */}
            <div className="sticky top-0 border-b bg-gray-200 dark:bg-slate-950 flex flex-col">
              {table.getHeaderGroups().map((hg) => (
                <div key={hg.id} className="flex w-full">
                  {hg.headers.map((header) => (
                    <div
                      key={header.id}
                      style={{ width: header.getSize() }}
                      className="relative flex items-center p-3 text-sm font-semibold text-foreground border-r last:border-0 bg-gray-200 dark:bg-slate-950"
                    >
                      <div
                        className={cn(
                          "flex items-center gap-2 grow",
                          header.column.getCanSort() &&
                            "cursor-pointer select-none",
                        )}
                        onClick={header.column.getToggleSortingHandler()}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext(),
                            )}
                        {header.column.getCanSort() && (
                          <span className="text-muted-foreground ml-1">
                            {{ asc: "▲", desc: "▼" }[
                              header.column.getIsSorted() as string
                            ] ?? ""}
                          </span>
                        )}
                      </div>

                      {/* RESIZER */}
                      {enableColumnResizing && header.column.getCanResize() ? (
                        <div
                          onMouseDown={header.getResizeHandler()}
                          onTouchStart={header.getResizeHandler()}
                          className={cn(
                            "absolute right-0 top-0 h-full w-1 cursor-col-resize select-none touch-none",
                            "after:absolute after:right-0 after:top-0 after:h-full after:w-px after:bg-border",
                            header.column.getIsResizing() && "after:bg-primary",
                          )}
                          title="Drag to resize"
                        />
                      ) : null}
                    </div>
                  ))}
                </div>
              ))}
            </div>

            {/* BODY AREA */}
            <div
              ref={scrollYRef}
              className={cn(
                "bg-background",
                isVirtual && "overflow-x-hidden scrollbar-thin",
              )}
              style={
                isVirtual
                  ? { maxHeight: virtual.maxBodyHeight ?? 560 }
                  : undefined
              }
            >
              {isLoading ? (
                <div className="p-12 text-center text-sm text-muted-foreground">
                  {renderLoading ? renderLoading() : "Loading data..."}
                </div>
              ) : !hasRows ? (
                <div className="p-20 text-center text-sm text-muted-foreground border-b">
                  {renderEmpty ? renderEmpty() : "No results found."}
                </div>
              ) : (
                <div
                  style={
                    isVirtual
                      ? {
                          height: `${rowVirtualizer.getTotalSize()}px`,
                          width: "100%",
                          position: "relative",
                        }
                      : {}
                  }
                >
                  {(isVirtual ? rowVirtualizer.getVirtualItems() : rows).map(
                    (item: VirtualItem | Row<TData>) => {
                      const row = isVirtual
                        ? rows[(item as VirtualItem).index]
                        : (item as Row<TData>);

                      const index = isVirtual
                        ? (item as VirtualItem).index
                        : rows.indexOf(row);

                      const sticky =
                        isVirtual && stickyIndexes.length
                          ? isSticky(index)
                          : false;

                      const activeSticky =
                        sticky && isVirtual && stickyIndexes.length
                          ? isActiveSticky(index)
                          : false;

                      return (
                        <div
                          key={row.id}
                          ref={isVirtual ? rowVirtualizer.measureElement : null}
                          data-index={isVirtual ? item.index : undefined}
                          className={cn(
                            "flex w-full border-b transition-colors hover:bg-muted/50",
                            row.getIsSelected() && "bg-muted",
                            isVirtual &&
                              !activeSticky &&
                              "absolute left-0 top-0",
                          )}
                          style={
                            isVirtual && isVirtualItem(item)
                              ? {
                                  ...(sticky
                                    ? {
                                        background: "#fff",
                                        borderBottom: "1px solid #ddd",
                                        zIndex: 1,
                                      }
                                    : {}),
                                  ...(activeSticky
                                    ? {
                                        position: "sticky" as const,
                                      }
                                    : {
                                        position: "absolute" as const,
                                        transform: `translateY(${item.start}px)`,
                                      }),
                                  top: 0,
                                  left: 0,
                                  width: "100%",
                                }
                              : {}
                          }
                        >
                          <div
                            className="grid w-full"
                            style={{ gridTemplateColumns }}
                          >
                            {row.getVisibleCells().map((cell) => (
                              <div
                                key={cell.id}
                                className="p-3 text-sm flex items-center overflow-hidden border-r last:border-0"
                                style={{ width: cell.column.getSize() }}
                              >
                                {flexRender(
                                  cell.column.columnDef.cell,
                                  cell.getContext(),
                                )}
                              </div>
                            ))}
                          </div>
                        </div>
                      );
                    },
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {renderPagination ? (
        renderPagination(table)
      ) : (
        <DataTablePagination table={table} />
      )}
    </div>
  );
}
