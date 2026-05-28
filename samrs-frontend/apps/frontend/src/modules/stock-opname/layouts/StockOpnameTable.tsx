import DataTable from "@/components/commons/data-table/DataTable";
import { createStockOpnameColumns } from "../components/StockOpnameColumns";
import type { StockOpnameSession } from "../types";

type Props = {
  data: StockOpnameSession[];
  isLoading?: boolean;
  onView: (session: StockOpnameSession) => void;
  onClose: (id: number) => void;
};

export function StockOpnameTable({ data, isLoading, onView, onClose }: Props) {
  const columns = createStockOpnameColumns({ onView, onClose });
  return <DataTable columns={columns} data={data} isLoading={isLoading} enableRowSelection={false} />;
}
