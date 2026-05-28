import DataTable from "@/components/commons/data-table/DataTable";
import { createRoomColumns } from "../components/RoomColumns";
import type { Room } from "../types";

type RoomTableProps = { data: Room[]; isLoading: boolean; onEdit: (room: Room) => void; onDelete: (id: string) => void };

export function RoomTable({ data, isLoading, onEdit, onDelete }: RoomTableProps) {
  return <DataTable columns={createRoomColumns({ onEdit, onDelete })} data={data} isLoading={isLoading} />;
}
