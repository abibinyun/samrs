import DataTable from "@/components/commons/data-table/DataTable";
import { createAssetStatusColumns } from "../components/AssetStatusColumns";
import type { AssetStatus } from "../types";

type AssetStatusTableProps = { data: AssetStatus[]; isLoading: boolean; onEdit: (status: AssetStatus) => void; onDelete: (id: number) => void };

export function AssetStatusTable({ data, isLoading, onEdit, onDelete }: AssetStatusTableProps) {
  return <DataTable columns={createAssetStatusColumns({ onEdit, onDelete })} data={data} isLoading={isLoading} />;
}
