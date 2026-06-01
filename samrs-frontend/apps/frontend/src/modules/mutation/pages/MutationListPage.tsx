import { Suspense, lazy, useMemo, useState } from "react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { SCOPES } from "@/constants/permissions";
import { usePermission } from "@/hooks/usePermission";
import { useGetMutationsQuery } from "../actions/mutationApi";
import { createMutationColumns } from "../components/MutationColumns";
import { useDebounce } from "@/hooks/useDebounce";
import PageContainer from "@/components/commons/containers/PageContainer";
import MutationTable from "../components/MutationTable";

export default function MutationListPage() {
  const { hasPermission } = usePermission();
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });

  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const debouncedSearch = useDebounce(search, 400);

  const { data, isLoading, isFetching, error } = useGetMutationsQuery({
    page,
    limit: 20,
    search: debouncedSearch || undefined,
  });

  const mutations = data?.data || [];
  const total = data?.total || 0;
  const totalPages = Math.ceil(total / 20);

  const columns = useMemo(
    () =>
      createMutationColumns({
        onViewDetail: (id) => navigate({ to: `${basePath}/assets/mutation/${id}` }),
      }),
    [navigate, basePath],
  );

  return (
    <PageContainer
      header={{
        title: "Mutasi & Timeline",
        description: `Total ${total.toLocaleString("id-ID")} mutasi tercatat.`,
        actions: (
          <>
            {hasPermission(SCOPES.ASSET_MUTATION.CREATE) && (
              <Button className="gap-2" onClick={() => navigate({ to: `${basePath}/assets/mutation/create` })}>
                <Plus className="w-4 h-4" /> Tambah Mutasi
              </Button>
            )}
          </>
        ),
      }}
    >
      <div className="mb-4">
        <Input
          placeholder="Cari alasan, kode aset..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
          className="max-w-sm"
        />
      </div>

      {error && (
        <div className="flex flex-col items-center justify-center h-64 gap-4">
          <p className="text-destructive">Gagal memuat data mutasi</p>
          <Button variant="outline" onClick={() => window.location.reload()}>Retry</Button>
        </div>
      )}

      {!error && (
        <Suspense
          fallback={
            <div className="space-y-2">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="h-12 bg-muted animate-pulse rounded" />
              ))}
            </div>
          }
        >
          <MutationTable
            data={mutations}
            columns={columns}
            isLoading={isLoading || isFetching}
          />
        </Suspense>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <p className="text-sm text-muted-foreground">
            Halaman {page} dari {totalPages}
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </PageContainer>
  );
}
