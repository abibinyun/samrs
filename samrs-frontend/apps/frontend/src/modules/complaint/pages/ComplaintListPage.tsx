import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { ComplaintTable } from "../layouts/ComplaintTable";
import { useGetComplaintsQuery, useDeleteComplaintMutation } from "../actions/complaintApi";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import type { Complaint } from "../types";
import { toast } from "sonner";

export function ComplaintListPage() {
  const navigate = useNavigate();
  const { data, isLoading } = useGetComplaintsQuery({});
  const [deleteComplaint] = useDeleteComplaintMutation();
  const [selectedComplaint, setSelectedComplaint] = useState<Complaint | null>(null);

  const handleView = (complaint: Complaint) => {
    setSelectedComplaint(complaint);
    // TODO: Open detail dialog
  };

  const handleUpdateStatus = (complaint: Complaint) => {
    setSelectedComplaint(complaint);
    // TODO: Open update status dialog
  };

  const handleDelete = async (id: string) => {
    if (confirm("Yakin ingin menghapus keluhan ini?")) {
      try {
        await deleteComplaint(id).unwrap();
        toast.success("Keluhan berhasil dihapus");
      } catch (error) {
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
          <Button onClick={() => navigate({ to: "/complaints/create" })}>
            <Plus className="w-4 h-4 mr-2" />
            Buat Keluhan
          </Button>
        )
      }}
    >
      <ComplaintTable
        data={data?.data || []}
        isLoading={isLoading}
        onView={handleView}
        onUpdateStatus={handleUpdateStatus}
        onDelete={handleDelete}
      />
    </PageContainer>
  );
}
