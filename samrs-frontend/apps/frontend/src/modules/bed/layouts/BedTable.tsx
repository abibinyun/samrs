import DataTable from "@/components/commons/data-table/DataTable";
import { createBedColumns } from "../components/BedColumns";
import type { Bed } from "../types";

type BedTableProps = { data: Bed[]; isLoading: boolean; onEdit: (bed: Bed) => void; onDelete: (id: number) => void };

export function BedTable({ data, isLoading, onEdit, onDelete }: BedTableProps) {
  return <DataTable columns={createBedColumns({ onEdit, onDelete })} data={data} isLoading={isLoading} />;
}
