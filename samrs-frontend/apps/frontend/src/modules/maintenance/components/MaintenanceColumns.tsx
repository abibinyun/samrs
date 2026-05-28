import type { ColumnDef } from "@tanstack/react-table";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MoreVertical, CheckCircle, Calendar } from "lucide-react";
import { cn } from "@/lib/utils";
import type { MaintenanceSchedule } from "../types";
import { maintenanceStatus } from "./MaintenanceStatus";
import { format } from "date-fns";
import { id as idLocale } from "date-fns/locale";

export function createMaintenanceColumns(opts: {
  onComplete: (schedule: MaintenanceSchedule) => void;
  onUpdate: (schedule: MaintenanceSchedule) => void;
  onDelete: (id: number) => void;
}): ColumnDef<MaintenanceSchedule, unknown>[] {
  return [
    {
      accessorKey: "title",
      header: "Jadwal",
      cell: ({ row }) => {
        const m = row.original;
        return (
          <div>
            <div className="font-semibold text-sm">{m.title}</div>
            <div className="text-xs text-muted-foreground mt-0.5">
              {m.schedule_type === "maintenance" ? "Pemeliharaan" : "Kalibrasi"}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "asset_name",
      header: "Aset",
      cell: ({ row }) => {
        const m = row.original;
        return (
          <div>
            <div className="text-sm">{m.asset_name}</div>
            <div className="text-xs text-muted-foreground font-mono mt-0.5">
              {m.asset_code}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "next_due_date",
      header: "Jadwal Berikutnya",
      cell: ({ getValue }) => {
        const date = new Date(getValue<string>());
        return (
          <div className="flex items-center gap-1.5 text-sm">
            <Calendar className="w-3.5 h-3.5 text-muted-foreground" />
            {format(date, "dd MMM yyyy", { locale: idLocale })}
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ getValue }) => {
        const s = maintenanceStatus(getValue<string>());
        return (
          <span className={cn("px-2.5 py-1 rounded-full text-[10px] font-bold border uppercase tracking-wide", s.className)}>
            {s.label}
          </span>
        );
      },
    },
    {
      id: "actions",
      header: () => <div className="text-right">Aksi</div>,
      enableSorting: false,
      cell: ({ row }) => {
        const m = row.original;
        return (
          <div className="flex items-center justify-end gap-1">
            {m.status !== "completed" && (
              <Button variant="ghost" size="icon" onClick={() => opts.onComplete(m)} title="Tandai Selesai">
                <CheckCircle className="w-4 h-4" />
              </Button>
            )}
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" title="Menu">
                  <MoreVertical className="w-4 h-4 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => opts.onUpdate(m)}>Edit</DropdownMenuItem>
                <DropdownMenuItem onClick={() => opts.onDelete(m.id)} className="text-red-600">Hapus</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        );
      },
    },
  ];
}
