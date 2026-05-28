import DataTable from "@/components/commons/data-table/DataTable";
import { createMaintenanceColumns } from "../components/MaintenanceColumns";
import type { MaintenanceSchedule } from "../types";

type MaintenanceTableProps = {
  data: MaintenanceSchedule[];
  isLoading?: boolean;
  onComplete: (schedule: MaintenanceSchedule) => void;
  onUpdate: (schedule: MaintenanceSchedule) => void;
  onDelete: (id: number) => void;
};

export function MaintenanceTable({ data, isLoading, onComplete, onUpdate, onDelete }: MaintenanceTableProps) {
  const columns = createMaintenanceColumns({ onComplete, onUpdate, onDelete });
  return <DataTable columns={columns} data={data} isLoading={isLoading} enableRowSelection={false} />;
}
