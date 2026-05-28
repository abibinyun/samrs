import type { Table } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";

type Props<TData> = {
  table: Table<TData>;
};

export default function DataTablePagination<TData>({ table }: Props<TData>) {
  const { pageIndex, pageSize } = table.getState().pagination;
  const pageCount = table.getPageCount();

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      {/* INFO */}
      <div className="text-sm text-muted-foreground text-center sm:text-left">
        {table.getFilteredSelectedRowModel().rows.length} selected /{" "}
        {table.getFilteredRowModel().rows.length} rows
      </div>

      {/* CONTROLS */}
      <div className="flex flex-wrap items-center justify-center gap-2 sm:justify-end">
        {/* PAGE SIZE (hidden on very small screens if needed) */}
        <select
          className="h-9 rounded-md border border-input bg-background px-2 text-sm"
          value={pageSize}
          onChange={(e) => table.setPageSize(Number(e.target.value))}
        >
          {[10, 20, 50, 100].map((size) => (
            <option key={size} value={size}>
              {size} / page
            </option>
          ))}
        </select>

        {/* PAGE INFO */}
        <div className="text-sm text-muted-foreground whitespace-nowrap">
          {pageIndex + 1} / {pageCount || 1}
        </div>

        {/* DESKTOP BUTTONS */}
        <div className="hidden sm:flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => table.setPageIndex(0)}
            disabled={!table.getCanPreviousPage()}
          >
            « First
          </Button>
          <Button
            variant="outline"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            ‹ Prev
          </Button>
          <Button
            variant="outline"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Next ›
          </Button>
          <Button
            variant="outline"
            onClick={() => table.setPageIndex(Math.max(0, pageCount - 1))}
            disabled={!table.getCanNextPage()}
          >
            Last »
          </Button>
        </div>

        {/* MOBILE BUTTONS */}
        <div className="flex sm:hidden items-center gap-2">
          <Button
            size="icon"
            variant="outline"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
            aria-label="Previous page"
          >
            ‹
          </Button>
          <Button
            size="icon"
            variant="outline"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
            aria-label="Next page"
          >
            ›
          </Button>
        </div>
      </div>
    </div>
  );
}

