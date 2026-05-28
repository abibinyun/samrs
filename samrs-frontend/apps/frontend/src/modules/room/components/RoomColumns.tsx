import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Edit, Trash2 } from "lucide-react";
import type { Room } from "../types";

export const createRoomColumns = (opts: { onEdit: (room: Room) => void; onDelete: (id: string) => void }): ColumnDef<Room>[] => [
  { accessorKey: "code", header: "Code", cell: ({ row }) => <div className="font-mono font-medium">{row.original.code}</div> },
  { accessorKey: "name", header: "Name", cell: ({ row }) => <div className="font-medium">{row.original.name}</div> },
  { accessorKey: "location", header: "Location", cell: ({ row }) => row.original.location || "-" },
  {
    id: "actions",
    header: "Actions",
    cell: ({ row }) => (
      <div className="flex gap-2">
        <Button variant="ghost" size="sm" onClick={() => opts.onEdit(row.original)}><Edit className="w-4 h-4" /></Button>
        <Button variant="ghost" size="sm" onClick={() => opts.onDelete(row.original.id)}><Trash2 className="w-4 h-4 text-red-600" /></Button>
      </div>
    ),
  },
];
