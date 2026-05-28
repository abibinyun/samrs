import { lazy, useMemo, useState } from "react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import { Plus, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SCOPES } from "@/constants/permissions";
import { usePermission } from "@/hooks/usePermission";
import type { AssetQrPayload } from "@/modules/asset/components/AssetQrDialog"; 
import { TableSkeleton } from "@/components/commons/data-table/TableSkeleton";
import { useGetMutationsQuery } from "@/modules/mutation/actions/mutationApi";
import { useDebounce } from "@/hooks/useDebounce";
import type { DateRange } from "react-day-picker";
import { startOfDay } from "date-fns";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { createTabKey } from "@/modules/dashboard/tabs/utils";
import { patchAssetListUI } from "@/modules/asset/actions/assetListUiSlice"; 

const AssetTable = lazy(() => import("../../asset/layouts/AssetTable"));
const PageContainer = lazy(() => import("@/components/commons/containers/PageContainer"));
const AssetQrDialog = lazy(() => import("../../asset/components/AssetQrDialog"));

function toISO(d?: Date) {
  return d ? d.toISOString() : undefined;
}
function fromISO(s?: string) {
  return s ? new Date(s) : undefined;
}

export default function MutationListPage() {
  const { hasPermission } = usePermission();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const location = useRouterState({ select: (s) => s.location });

  const { data: assets = [], isLoading, isFetching } = useGetMutationsQuery();
  const [qr, setQr] = useState<AssetQrPayload | null>(null);

  // ✅ ambil location untuk tabKey yang sama persis dengan KeepAliveOutlet
  const tabKey = useMemo(
    () => createTabKey(location.pathname, location.searchStr),
    [location.pathname, location.searchStr]
  );

  const ui = useAppSelector((s) => s.assetListUi.byTabKey[tabKey]) ?? {
    search: "",
    category: "__all",
    dateRange: undefined,
  };

  const search = ui.search;
  const category = ui.category;

  const dateRange: DateRange | undefined = useMemo(() => {
    const from = fromISO(ui.dateRange?.from);
    const to = fromISO(ui.dateRange?.to);
    return from || to ? { from, to } : undefined;
  }, [ui.dateRange?.from, ui.dateRange?.to]);

  const setSearch = (v: string) =>
    dispatch(patchAssetListUI({ tabKey, patch: { search: v } }));

  const setCategory = (v: string) =>
    dispatch(patchAssetListUI({ tabKey, patch: { category: v } }));

  const setDateRange = (r?: DateRange) =>
    dispatch(
      patchAssetListUI({
        tabKey,
        patch: {
          dateRange: r
            ? { from: toISO(r.from), to: toISO(r.to) }
            : undefined,
        },
      })
    );
  
  const setColumnVisibility = (v: any) =>
    dispatch(patchAssetListUI({ tabKey, patch: { columnVisibility: v } }));

  const setColumnFilters = (v: any) =>
    dispatch(patchAssetListUI({ tabKey, patch: { columnFilters: v } }));

  const categories = ["ready", "maintenance", "broken"];

  const normalizedAssets = useMemo(() => {
    return assets.map((a) => ({
      ...a,
      _search: `${a.name} ${a.code}`.toLowerCase(),
      _category: a.status.toString(),
      _purchase_date: startOfDay(new Date(a.purchase_date)),
    }));
  }, [assets]);

  const debouncedSearch = useDebounce(search, 400);

  const filteredAssets = useMemo(() => {
    const s = debouncedSearch.trim().toLowerCase();

    return normalizedAssets.filter((a) => {
      const matchSearch = !s || a._search.includes(s);
      const matchCategory = category === "__all" || a._category === category;

      const matchDate =
        !dateRange?.from ||
        !dateRange?.to ||
        (a._purchase_date >= startOfDay(dateRange.from) &&
          a._purchase_date <= startOfDay(dateRange.to));

      return matchSearch && matchCategory && matchDate;
    });
  }, [normalizedAssets, debouncedSearch, category, dateRange]);

  const handleAddAsset = () => {
    navigate({ to: "create" });
  };

  return (
    <PageContainer
      header={{
        title: "Mutation Aset",
        description: `Total ${assets.length.toLocaleString(
          "id-ID"
        )} aset terdaftar di sistem.`,
        actions: (
          <>
            {hasPermission(SCOPES.REPORT.EXPORT) && (
              <Button variant="outline" className="gap-2">
                <Download className="w-4 h-4" /> Export
              </Button>
            )}
            {hasPermission(SCOPES.ASSET.CREATE) && (
              <Button className="gap-2" onClick={handleAddAsset}>
                <Plus className="w-4 h-4" /> Tambah Aset
              </Button>
            )}
          </>
        )
      }}
    >
      {(isLoading || isFetching) && <TableSkeleton rows={20} />}

      {!isLoading && (
        <AssetTable
          data={filteredAssets}
          isLoading={isLoading}
          search={search}
          onSearchChange={setSearch}
          onOpenQr={setQr}
          setCategory={setCategory}
          category={category}
          categories={categories}
          dateRange={dateRange}
          setDateRange={setDateRange}
          columnVisibility={ui.columnVisibility}
          onColumnVisibilityChange={setColumnVisibility}
          columnFilters={ui.columnFilters}
          onColumnFiltersChange={setColumnFilters}
        />
      )}

      <AssetQrDialog
        open={!!qr}
        onOpenChange={(v) => !v && setQr(null)}
        asset={qr}
      />
    </PageContainer>
  );
}