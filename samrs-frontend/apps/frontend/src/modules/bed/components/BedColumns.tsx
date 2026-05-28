import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Edit, Trash2 } from "lucide-react";
import type { Bed } from "../types";

export const createBedColumns = (opts: { onEdit: (bed: Bed) => void; onDelete: (id: number) => void }): ColumnDef<Bed>[] => [
  { accessorKey: "name", header: "Name", cell: ({ row }) => <div className="font-medium">{row.original.name}</div> },
  { accessorKey: "room", header: "Room", cell: ({ row }) => row.original.room?.name || "-" },
  { accessorKey: "is_available", header: "Status", cell: ({ row }) => <Badge variant={row.original.is_available ? "default" : "secondary"}>{row.original.is_available ? "Available" : "Occupied"}</Badge> },
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
