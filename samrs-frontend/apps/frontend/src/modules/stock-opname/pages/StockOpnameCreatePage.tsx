import PageContainer from "@/components/commons/containers/PageContainer";
import { StockOpnameForm } from "../layouts/StockOpnameForm";
import { useCreateStockOpnameSessionMutation } from "../actions/stockOpnameApi";
import { useRouterState } from "@tanstack/react-router";
import type { StockOpnameFormData } from "../types";
import { toast } from "sonner";

export function StockOpnameCreatePage() {
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const [createSession, { isLoading }] = useCreateStockOpnameSessionMutation();

  const handleSubmit = async (data: StockOpnameFormData) => {
    try {
      await createSession(data).unwrap();
      toast.success("Sesi berhasil dibuat");
      window.location.href = `${basePath}/stock-opname`;
    } catch {
      toast.error("Gagal membuat sesi");
    }
  };

  return (
    <PageContainer
      header={{
        title: "Buat Sesi Stock Opname",
        description: "Isi data untuk membuat sesi baru",
      }}
    >
      <div className="bg-card border border-border rounded-xl p-6">
        <StockOpnameForm onSubmit={handleSubmit} isLoading={isLoading} />
      </div>
    </PageContainer>
  );
}
