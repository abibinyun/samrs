import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { MaintenanceTable } from "../layouts/MaintenanceTable";
import { useGetMaintenanceSchedulesQuery, useDeleteMaintenanceScheduleMutation, useCompleteMaintenanceScheduleMutation } from "../actions/maintenanceApi";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import type { MaintenanceSchedule } from "../types";
import { toast } from "sonner";

export function MaintenanceListPage() {
  const navigate = useNavigate();
  const { data, isLoading } = useGetMaintenanceSchedulesQuery({});
  const [deleteSchedule] = useDeleteMaintenanceScheduleMutation();
  const [completeSchedule] = useCompleteMaintenanceScheduleMutation();

  const handleComplete = async (schedule: MaintenanceSchedule) => {
    const note = prompt("Catatan penyelesaian:");
    if (note) {
      try {
        await completeSchedule({ id: schedule.id, note }).unwrap();
        toast.success("Jadwal ditandai selesai");
      } catch (error) {
        toast.error("Gagal menandai selesai");
      }
    }
  };

  const handleUpdate = (schedule: MaintenanceSchedule) => {
    // TODO: Open update dialog
    toast.info("Fitur edit akan segera hadir");
  };

  const handleDelete = async (id: number) => {
    if (confirm("Yakin ingin menghapus jadwal ini?")) {
      try {
        await deleteSchedule(id).unwrap();
        toast.success("Jadwal berhasil dihapus");
      } catch (error) {
        toast.error("Gagal menghapus jadwal");
      }
    }
  };

  return (
    <PageContainer
      header={{
        title: "Jadwal Pemeliharaan",
        description: "Kelola jadwal pemeliharaan dan kalibrasi aset",
        actions: (
          <Button onClick={() => navigate({ to: "/maintenance/schedules/create" })}>
            <Plus className="w-4 h-4 mr-2" />
            Buat Jadwal
          </Button>
        )
      }}
    >
      <MaintenanceTable
        data={data?.data || []}
        isLoading={isLoading}
        onComplete={handleComplete}
        onUpdate={handleUpdate}
        onDelete={handleDelete}
      />
    </PageContainer>
  );
}
