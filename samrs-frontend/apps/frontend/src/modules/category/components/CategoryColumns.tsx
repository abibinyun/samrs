import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Edit, Trash2 } from "lucide-react";
import type { Category } from "../types";
import { format } from "date-fns";

export const createCategoryColumns = (opts: {
  onEdit: (category: Category) => void;
  onDelete: (id: number) => void;
}): ColumnDef<Category>[] => [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => <div className="font-medium">{row.original.name}</div>,
  },
  {
    accessorKey: "slug",
    header: "Slug",
    cell: ({ row }) => <div className="text-sm text-muted-foreground font-mono">{row.original.slug}</div>,
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) => <div className="max-w-md truncate">{row.original.description || "-"}</div>,
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => <div className="text-sm">{format(new Date(row.original.created_at), "dd MMM yyyy")}</div>,
  },
  {
    id: "actions",
    header: "Actions",
    cell: ({ row }) => (
      <div className="flex gap-2">
        <Button variant="ghost" size="sm" onClick={() => opts.onEdit(row.original)}>
          <Edit className="w-4 h-4" />
        </Button>
        <Button variant="ghost" size="sm" onClick={() => opts.onDelete(row.original.id)}>
          <Trash2 className="w-4 h-4 text-red-600" />
        </Button>
      </div>
    ),
  },
];
