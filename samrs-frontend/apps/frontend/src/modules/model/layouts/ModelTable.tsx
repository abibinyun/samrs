import DataTable from "@/components/commons/data-table/DataTable";
import { createModelColumns } from "../components/ModelColumns";
import type { Model } from "../types";

type ModelTableProps = { data: Model[]; isLoading: boolean; onEdit: (model: Model) => void; onDelete: (id: number) => void };

export function ModelTable({ data, isLoading, onEdit, onDelete }: ModelTableProps) {
  return <DataTable columns={createModelColumns({ onEdit, onDelete })} data={data} isLoading={isLoading} />;
}
