import * as React from "react";
import type { Table } from "@tanstack/react-table";
import type { DateRange } from "react-day-picker";
import { format, startOfDay, subDays, subMonths, subYears, addYears } from "date-fns";
import { Search, CalendarIcon, X } from "lucide-react";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

type CategoryOption = { value: string; label: string };

type Props<TData> = {
  table: Table<TData>;

  search?: string;
  onSearchChange?: (v: string) => void;
  searchPlaceholder?: string;

  dateRange?: DateRange;
  onDateRangeChange?: (range?: DateRange) => void;
  datePresets?: { label: string; range: DateRange }[];

  category?: string;
  onCategoryChange?: (v: string) => void;
  categories?: CategoryOption[];
  categoryAllLabel?: string;
  categoryAllValue?: string;

  enableColumns?: boolean;
  columnsLabel?: string;

  onResetFilters?: () => void;
  resetLabel?: string;
};

const today = startOfDay(new Date());

const defaultPresets: { label: string; range: DateRange }[] = [
  { label: "Sekarang", range: { from: today, to: today } },
  { label: "1 Minggu", range: { from: subDays(today, 7), to: today } },
  { label: "1 Bulan", range: { from: subMonths(today, 1), to: today } },
  { label: "3 Bulan", range: { from: subMonths(today, 3), to: today } },
  { label: "1 Tahun", range: { from: subYears(today, 1), to: today } },
];

const isSameRange = (a?: DateRange, b?: DateRange) =>
  a?.from?.getTime() === b?.from?.getTime() && a?.to?.getTime() === b?.to?.getTime();

export default function DataTableFiltersBar<TData>({
  table,

  search,
  onSearchChange,
  searchPlaceholder = "Cari...",

  dateRange,
  onDateRangeChange,
  datePresets = defaultPresets,

  category,
  onCategoryChange,
  categories,
  categoryAllLabel = "Semua",
  categoryAllValue = "__all",

  enableColumns = true,
  columnsLabel = "Columns",

  onResetFilters,
  resetLabel = "Reset",
}: Props<TData>) {
  const [dateOpen, setDateOpen] = React.useState(false);

  // columns menu local buffer
  const [columnsOpen, setColumnsOpen] = React.useState(false);
  const [tempVisibility, setTempVisibility] = React.useState<Record<string, boolean>>({});

  const state = table.getState();

  const hasSearch = !!search && search.trim().length > 0;
  const hasCategory = !!category && category !== categoryAllValue;
  const hasDateRange = !!dateRange?.from || !!dateRange?.to;
  const hasSorting = (state.sorting?.length ?? 0) > 0;
  const hasRealColumnFilters = (state.columnFilters?.length ?? 0) > 0;

  const hasHiddenColumns = table
    .getAllLeafColumns()
    .filter((c) => c.getCanHide?.())
    .some((c) => !c.getIsVisible());

  const isFiltered =
    hasSearch ||
    hasCategory ||
    hasDateRange ||
    hasSorting ||
    hasRealColumnFilters ||
    hasHiddenColumns;

  const handleColumnsOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      const current = table.getState().columnVisibility ?? {};
      const initial: Record<string, boolean> = {};
      table.getAllColumns().forEach((col) => {
        if (col.getCanHide()) initial[col.id] = current[col.id] ?? true;
      });
      setTempVisibility(initial);
    }
    setColumnsOpen(isOpen);
  };

  const applyColumns = () => {
    table.setColumnVisibility(tempVisibility);
    setColumnsOpen(false);
  };

  const handleSelectRange = (range?: DateRange) => {
    if (!onDateRangeChange) return;

    if (!range?.from || !range?.to) {
      onDateRangeChange(range);
      return;
    }

    const diffYears = range.to.getFullYear() - range.from.getFullYear();
    if (diffYears > 1 || (diffYears === 1 && range.to.getMonth() > range.from.getMonth())) {
      const maxTo = addYears(range.from, 1);
      onDateRangeChange({ from: range.from, to: range.to > maxTo ? maxTo : range.to });
    } else {
      onDateRangeChange(range);
    }
  };

  const handleReset = () => {
    // reset table internal
    table.resetColumnVisibility();
    table.resetColumnFilters();
    table.resetSorting?.();
    table.resetPagination?.();

    // reset external filters (redux/page)
    onResetFilters?.();
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-12 gap-3 bg-card p-3 rounded-xl border shadow-sm">
      {/* Search */}
      {typeof search === "string" && onSearchChange ? (
        <div className="relative md:col-span-5">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder={searchPlaceholder}
            className="pl-9"
          />
        </div>
      ) : null}

      {/* Date range */}
      {onDateRangeChange ? (
        <div className="md:col-span-3">
          <Popover open={dateOpen} onOpenChange={setDateOpen}>
            <PopoverTrigger asChild>
              <Button variant="outline" className="w-full justify-start text-left font-normal">
                {!(dateRange?.from && dateRange?.to) && (
                  <CalendarIcon className="mr-2 h-4 w-4 shrink-0" />
                )}
                <span className="truncate flex-1">
                  {dateRange?.from ? (
                    dateRange.to ? (
                      `${format(dateRange.from, "dd MMM yyyy")} – ${format(dateRange.to, "dd MMM yyyy")}`
                    ) : (
                      format(dateRange.from, "dd MMM yyyy")
                    )
                  ) : (
                    "Pilih rentang tanggal"
                  )}
                </span>
                {dateRange?.from && (
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onDateRangeChange(undefined);
                      setDateOpen(false);
                    }}
                    className="ml-2 p-1 rounded hover:bg-muted/20"
                  >
                    <X className="h-4 w-4 text-muted-foreground" />
                  </button>
                )}
              </Button>
            </PopoverTrigger>

            <PopoverContent className="w-auto p-3 space-y-3" align="start">
              <div className="flex flex-wrap gap-2">
                {datePresets.map((preset) => (
                  <Button
                    key={preset.label}
                    size="sm"
                    variant={isSameRange(dateRange, preset.range) ? "default" : "secondary"}
                    onClick={() => {
                      handleSelectRange(preset.range);
                      setDateOpen(false);
                    }}
                  >
                    {preset.label}
                  </Button>
                ))}
              </div>

              <Calendar
                autoFocus
                mode="range"
                selected={dateRange}
                onSelect={handleSelectRange}
                defaultMonth={dateRange?.from}
                numberOfMonths={2}
                captionLayout="dropdown"
                startMonth={new Date(1990, 0)}
                endMonth={new Date(2026, 11)}
              />
            </PopoverContent>
          </Popover>
        </div>
      ) : null}

      {/* Category */}
      {typeof category === "string" && onCategoryChange ? (
        <div className="md:col-span-2">
          <Select value={category} onValueChange={onCategoryChange}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder={categoryAllLabel} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={categoryAllValue}>{categoryAllLabel}</SelectItem>
              {categories?.map((c) => (
                <SelectItem key={c.value} value={c.value}>
                  {c.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      ) : null}

      {/* Columns */}
      {enableColumns ? (
        <div className="md:col-span-2">
          <DropdownMenu open={columnsOpen} onOpenChange={handleColumnsOpenChange}>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="w-full">
                {columnsLabel}
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent align="end">
              {table
                .getAllColumns()
                .filter((col) => col.getCanHide())
                .map((col) => (
                  <DropdownMenuCheckboxItem
                    key={col.id}
                    className="capitalize"
                    checked={tempVisibility[col.id] ?? true}
                    onCheckedChange={(v) => {
                      setTempVisibility((prev) => ({ ...prev, [col.id]: !!v }));
                    }}
                    onSelect={(e) => e.preventDefault()}
                  >
                    {col.id}
                  </DropdownMenuCheckboxItem>
                ))}

              <DropdownMenuSeparator />

              <div className="p-2 space-y-2">
                <Button size="sm" className="w-full h-8" onClick={applyColumns}>
                  Apply Changes
                </Button>

                <Button
                  size="sm"
                  variant="secondary"
                  className="w-full h-8"
                  onClick={() => {
                    const allVisible: Record<string, boolean> = {};
                    table
                      .getAllColumns()
                      .filter((c) => c.getCanHide())
                      .forEach((c) => (allVisible[c.id] = true));
                    table.setColumnVisibility(allVisible);
                    setTempVisibility(allVisible);
                  }}
                >
                  Show All
                </Button>
              </div>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      ) : null}

      {/* Reset */}
      <div className="flex flex-1 items-center gap-2">
        {isFiltered ? (
          <Button variant="secondary" onClick={handleReset}>
            {resetLabel}
          </Button>
        ) : null}
      </div>
    </div>
  );
}
