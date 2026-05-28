import DataTable from "@/components/commons/data-table/DataTable";
import { createCategoryColumns } from "../components/CategoryColumns";
import type { Category } from "../types";

type CategoryTableProps = {
  data: Category[];
  isLoading: boolean;
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
};

export function CategoryTable({ data, isLoading, onEdit, onDelete }: CategoryTableProps) {
  const columns = createCategoryColumns({ onEdit, onDelete });

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
    />
  );
}
