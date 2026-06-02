import { useState, useCallback } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { ComplaintTable } from "../layouts/ComplaintTable";
import { useGetComplaintsQuery, useDeleteComplaintMutation } from "../actions/complaintApi";
import { ComplaintStatusDialog } from "../components/ComplaintStatusDialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Plus, Search } from "lucide-react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import type { Complaint } from "../types";
import { toast } from "sonner";

export function ComplaintListPage() {
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const perPage = 10;

  const { data, isLoading } = useGetComplaintsQuery({
    page,
    per_page: perPage,
    search: search || undefined,
    status: statusFilter === "all" ? undefined : statusFilter,
  });

  const [deleteComplaint] = useDeleteComplaintMutation();
  const [statusDialogOpen, setStatusDialogOpen] = useState(false);
  const [selectedComplaint, setSelectedComplaint] = useState<Complaint | null>(null);

  const totalPages = data?.meta ? Math.ceil(data.meta.total / perPage) : 1;

  const handleView = useCallback((complaint: Complaint) => {
    navigate({ to: `${basePath}/complaints/${complaint.id}` });
  }, [navigate, basePath]);

  const handleUpdateStatus = useCallback((complaint: Complaint) => {
    setSelectedComplaint(complaint);
    setStatusDialogOpen(true);
  }, []);

  const handleDelete = async (id: string) => {
    if (confirm("Yakin ingin menghapus keluhan ini?")) {
      try {
        await deleteComplaint(id).unwrap();
        toast.success("Keluhan berhasil dihapus");
      } catch {
        toast.error("Gagal menghapus keluhan");
      }
    }
  };

  return (
    <PageContainer
      header={{
        title: "Keluhan (Complaint)",
        description: "Kelola keluhan dan laporan kerusakan aset",
        actions: (
          <Button onClick={() => navigate({ to: `${basePath}/complaints/create` })}>
            <Plus className="w-4 h-4 mr-2" />
            Buat Keluhan
          </Button>
        )
      }}
    >
      <div className="flex gap-3 mb-4">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            placeholder="Cari keluhan..."
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            className="pl-9"
          />
        </div>
        <Select value={statusFilter} onValueChange={(v) => { setStatusFilter(v); setPage(1); }}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder="Semua Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Semua Status</SelectItem>
            <SelectItem value="open">Open</SelectItem>
            <SelectItem value="in_progress">In Progress</SelectItem>
            <SelectItem value="done">Done</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <ComplaintTable
        data={data?.data || []}
        isLoading={isLoading}
        onView={handleView}
        onUpdateStatus={handleUpdateStatus}
        onDelete={handleDelete}
      />

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <div className="text-sm text-muted-foreground">
            {data?.meta?.total ?? 0} keluhan
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              Sebelumnya
            </Button>
            <span className="flex items-center px-3 text-sm">
              {page} / {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Selanjutnya
            </Button>
          </div>
        </div>
      )}

      <ComplaintStatusDialog
        open={statusDialogOpen}
        onOpenChange={setStatusDialogOpen}
        complaint={selectedComplaint}
      />
    </PageContainer>
  );
}
