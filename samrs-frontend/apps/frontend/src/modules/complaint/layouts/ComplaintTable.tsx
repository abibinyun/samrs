import DataTable from "@/components/commons/data-table/DataTable";
import { createComplaintColumns } from "../components/ComplaintColumns";
import type { Complaint } from "../types";

type ComplaintTableProps = {
  data: Complaint[];
  isLoading?: boolean;
  onView: (complaint: Complaint) => void;
  onUpdateStatus: (complaint: Complaint) => void;
  onDelete: (id: string) => void;
};

export function ComplaintTable({
  data,
  isLoading,
  onView,
  onUpdateStatus,
  onDelete,
}: ComplaintTableProps) {
  const columns = createComplaintColumns({ onView, onUpdateStatus, onDelete });

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      enableRowSelection={false}
    />
  );
}
