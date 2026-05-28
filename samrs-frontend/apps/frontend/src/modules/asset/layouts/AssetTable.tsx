import * as React from "react";
import type { Asset } from "../types";
import { createAssetColumns } from "../components/AssetColumns";
import type { DateRange } from "react-day-picker";
import { resetAssetListUIToDefault } from "../actions/assetListUiSlice";
import { useAppDispatch } from "@/store/hooks";
import { useTabKey } from "@/hooks/useTabKey";
import PersistedDataTable from "@/components/commons/data-table/PersistedDataTable";
import DataTableFiltersBar from "@/components/commons/data-table/DataTableFilterBar";

type Props = {
  data: Asset[];
  isLoading?: boolean;
  search: string;
  onSearchChange: (v: string) => void;
  category: string;
  setCategory: (v: string) => void;
  dateRange?: DateRange;
  setDateRange: (v?: DateRange) => void;
  categories: string[];
  columnVisibility?: any;
  onColumnVisibilityChange?: (v: any) => void;
  columnFilters?: any;
  onColumnFiltersChange?: (v: any) => void;
  onOpenQr: (asset: Pick<Asset, "id" | "name">) => void;
};

export default function AssetTable({
  data,
  isLoading,
  search,
  onSearchChange,
  onOpenQr,
  setCategory,
  category,
  categories,
  dateRange,
  setDateRange,
  columnVisibility,
  onColumnVisibilityChange,
  columnFilters,
  onColumnFiltersChange,
}: Props) {
  const dispatch = useAppDispatch();
  const tabKey = useTabKey();

  const columns = React.useMemo(() => createAssetColumns({ onOpenQr }), [onOpenQr]);

  const handleResetExternal = () => {
    dispatch(resetAssetListUIToDefault({ tabKey }));
  };

  return (
    <PersistedDataTable
      data={data}
      columns={columns}
      isLoading={isLoading}
      columnVisibility={columnVisibility}
      onColumnVisibilityChange={onColumnVisibilityChange}
      columnFilters={columnFilters}
      onColumnFiltersChange={onColumnFiltersChange}
      onResetExternal={handleResetExternal}
      renderToolbar={({ table, resetTable }) => (
        <DataTableFiltersBar
          table={table}
          search={search}
          onSearchChange={onSearchChange}
          searchPlaceholder="Cari kode atau nama aset..."
          category={category}
          onCategoryChange={setCategory}
          categories={categories.map((c) => ({
            value: c,
            label: c.toUpperCase(),
          }))}
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
          onResetFilters={resetTable}
        />
      )}
    />
  );
}
