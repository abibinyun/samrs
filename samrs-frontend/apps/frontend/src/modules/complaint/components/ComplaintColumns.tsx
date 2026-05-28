import type { ColumnDef } from "@tanstack/react-table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MoreVertical, Eye } from "lucide-react";
import { cn } from "@/lib/utils";
import type { Complaint } from "../types";
import { complaintStatus } from "./ComplaintStatus";
import { formatDistanceToNow } from "date-fns";
import { id as idLocale } from "date-fns/locale";

export function createComplaintColumns(opts: {
  onView: (complaint: Complaint) => void;
  onUpdateStatus: (complaint: Complaint) => void;
  onDelete: (id: string) => void;
}): ColumnDef<Complaint, unknown>[] {
  return [
    {
      accessorKey: "title",
      header: "Keluhan",
      cell: ({ row }) => {
        const c = row.original;
        return (
          <div>
            <div className="font-semibold text-sm">{c.title}</div>
            <div className="text-xs text-muted-foreground mt-0.5 line-clamp-1">
              {c.description}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "asset_name",
      header: "Aset",
      cell: ({ row }) => {
        const c = row.original;
        return (
          <div>
            <div className="text-sm">{c.asset_name}</div>
            <div className="text-xs text-muted-foreground font-mono mt-0.5">
              {c.asset_code}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ getValue }) => {
        const s = complaintStatus(getValue<string>());
        return (
          <span
            className={cn(
              "px-2.5 py-1 rounded-full text-[10px] font-bold border uppercase tracking-wide",
              s.className,
            )}
          >
            {s.label}
          </span>
        );
      },
    },
    {
      accessorKey: "created_at",
      header: "Dibuat",
      cell: ({ getValue }) => {
        const date = new Date(getValue<string>());
        return (
          <div className="text-xs text-muted-foreground">
            {formatDistanceToNow(date, { addSuffix: true, locale: idLocale })}
          </div>
        );
      },
    },
    {
      id: "actions",
      header: () => <div className="text-right">Aksi</div>,
      enableSorting: false,
      cell: ({ row }) => {
        const c = row.original;
        return (
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => opts.onView(c)}
              title="Lihat Detail"
            >
              <Eye className="w-4 h-4" />
            </Button>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" title="Menu">
                  <MoreVertical className="w-4 h-4 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => opts.onUpdateStatus(c)}>
                  Update Status
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => opts.onDelete(c.id)}
                  className="text-red-600"
                >
                  Hapus
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        );
      },
    },
  ];
}
