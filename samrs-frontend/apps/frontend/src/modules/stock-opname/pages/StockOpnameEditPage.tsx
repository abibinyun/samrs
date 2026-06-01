import PageContainer from "@/components/commons/containers/PageContainer";
import { StockOpnameForm } from "../layouts/StockOpnameForm";
import { useGetStockOpnameSessionQuery, useUpdateStockOpnameSessionMutation } from "../actions/stockOpnameApi";
import { useRouterState, useParams } from "@tanstack/react-router";
import type { StockOpnameFormData } from "../types";
import { toast } from "sonner";

export function StockOpnameEditPage() {
  const { sessionId } = useParams({ strict: false });
  const id = Number(sessionId);
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data: session, isLoading: isSessionLoading } = useGetStockOpnameSessionQuery(id);
  const [updateSession, { isLoading }] = useUpdateStockOpnameSessionMutation();

  const handleSubmit = async (data: StockOpnameFormData) => {
    try {
      await updateSession({ id, data }).unwrap();
      toast.success("Sesi berhasil diupdate");
      window.location.href = `${basePath}/stock-opname`;
    } catch {
      toast.error("Gagal mengupdate sesi");
    }
  };

  return (
    <PageContainer
      header={{
        title: "Edit Sesi Stock Opname",
        description: `Editing session: ${session?.data?.title || ""}`,
      }}
    >
      <div className="bg-card border border-border rounded-xl p-6">
        {isSessionLoading ? (
          <div className="text-muted-foreground">Loading...</div>
        ) : (
          <StockOpnameForm
            onSubmit={handleSubmit}
            isLoading={isLoading}
            defaultValues={session?.data}
          />
        )}
      </div>
    </PageContainer>
  );
}
