import * as React from "react";
// import type { AssetDevice } from "../types";
import { createAssetDeviceColumns } from "@/modules/asset/components/AssetColumns"; 
import type { DateRange } from "react-day-picker";
import { resetAssetListUIToDefault } from "../actions/mutationListUiSlice";
import { useAppDispatch } from "@/store/hooks";
import { useTabKey } from "@/hooks/useTabKey";

const DataTableFiltersBar = React.lazy(
  () => import("../../../components/commons/data-table/DataTableFilterBar")
);
const PersistedDataTable = React.lazy(
  () => import("@/components/commons/data-table/PersistedDataTable")
);

type Props = {
  data: any[];
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
  onOpenQr: (asset: Pick<any, "id" | "name">) => void;
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

  const columns = React.useMemo(() => createAssetDeviceColumns(), [onOpenQr]);

  const handleResetExternal = () => {
    dispatch(resetAssetListUIToDefault({ tabKey }));
  };

  return (
    <PersistedDataTable
      data={data}
      columns={columns as any}
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
          categories={categories.map((c) => ({ value: c, label: c.toUpperCase() }))}
          dateRange={dateRange}
          onDateRangeChange={setDateRange}
          onResetFilters={resetTable}
        />
      )}
    />
  );
}
