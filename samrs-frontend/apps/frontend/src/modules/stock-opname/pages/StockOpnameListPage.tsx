import PageContainer from "@/components/commons/containers/PageContainer";
import { StockOpnameTable } from "../layouts/StockOpnameTable";
import { useGetStockOpnameSessionsQuery, useCloseStockOpnameSessionMutation } from "../actions/stockOpnameApi";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import type { StockOpnameSession } from "../types";
import { toast } from "sonner";

export function StockOpnameListPage() {
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data, isLoading } = useGetStockOpnameSessionsQuery();
  const [closeSession] = useCloseStockOpnameSessionMutation();

  const handleView = (session: StockOpnameSession) => {
    navigate({ to: `${basePath}/stock-opname/${session.id}` });
  };

  const handleClose = async (id: number) => {
    if (confirm("Yakin ingin menutup sesi ini?")) {
      try {
        await closeSession(id).unwrap();
        toast.success("Sesi ditutup");
      } catch {
        toast.error("Gagal menutup sesi");
      }
    }
  };

  return (
    <PageContainer
      header={{
        title: "Stock Opname",
        description: "Kelola sesi stock opname aset",
        actions: (
          <Button onClick={() => navigate({ to: `${basePath}/stock-opname/create` })}>
            <Plus className="w-4 h-4 mr-2" />
            Buat Sesi
          </Button>
        )
      }}
    >
      <StockOpnameTable data={data?.data || []} isLoading={isLoading} onView={handleView} onClose={handleClose} />
    </PageContainer>
  );
}
