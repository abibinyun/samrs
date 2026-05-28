import type { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Eye, Lock } from "lucide-react";
import { cn } from "@/lib/utils";
import type { StockOpnameSession } from "../types";
import { stockOpnameStatus } from "./StockOpnameStatus";
import { format } from "date-fns";
import { id as idLocale } from "date-fns/locale";

export function createStockOpnameColumns(opts: {
  onView: (session: StockOpnameSession) => void;
  onClose: (id: number) => void;
}): ColumnDef<StockOpnameSession, unknown>[] {
  return [
    {
      accessorKey: "title",
      header: "Sesi Stock Opname",
      cell: ({ row }) => {
        const s = row.original;
        return (
          <div>
            <div className="font-semibold text-sm">{s.title}</div>
            <div className="text-xs text-muted-foreground mt-0.5">
              {format(new Date(s.opname_at), "dd MMM yyyy", { locale: idLocale })}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ getValue }) => {
        const s = stockOpnameStatus(getValue<string>());
        return <span className={cn("px-2.5 py-1 rounded-full text-[10px] font-bold border uppercase tracking-wide", s.className)}>{s.label}</span>;
      },
    },
    {
      id: "actions",
      header: () => <div className="text-right">Aksi</div>,
      cell: ({ row }) => {
        const s = row.original;
        return (
          <div className="flex items-center justify-end gap-1">
            <Button variant="ghost" size="icon" onClick={() => opts.onView(s)} title="Lihat Items">
              <Eye className="w-4 h-4" />
            </Button>
            {s.status === "open" && (
              <Button variant="ghost" size="icon" onClick={() => opts.onClose(s.id)} title="Tutup Sesi">
                <Lock className="w-4 h-4" />
              </Button>
            )}
          </div>
        );
      },
    },
  ];
}
