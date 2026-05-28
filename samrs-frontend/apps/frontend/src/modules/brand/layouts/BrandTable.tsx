import DataTable from "@/components/commons/data-table/DataTable";
import { createBrandColumns } from "../components/BrandColumns";
import type { Brand } from "../types";

type BrandTableProps = { data: Brand[]; isLoading: boolean; onEdit: (brand: Brand) => void; onDelete: (id: number) => void };

export function BrandTable({ data, isLoading, onEdit, onDelete }: BrandTableProps) {
  return <DataTable columns={createBrandColumns({ onEdit, onDelete })} data={data} isLoading={isLoading} />;
}
