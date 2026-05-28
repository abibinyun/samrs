import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Edit, Trash2 } from "lucide-react";
import type { Brand } from "../types";

export const createBrandColumns = (opts: { onEdit: (brand: Brand) => void; onDelete: (id: number) => void }): ColumnDef<Brand>[] => [
  { accessorKey: "name", header: "Name", cell: ({ row }) => <div className="font-medium">{row.original.name}</div> },
  { accessorKey: "description", header: "Description", cell: ({ row }) => <div className="max-w-md truncate">{row.original.description || "-"}</div> },
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
