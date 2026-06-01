import type { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import { Badge } from "@/components/ui/badge";
import type { AssetMutation } from "../types";

export function createMutationColumns(): ColumnDef<AssetMutation>[] {
  return [
    {
      accessorKey: "asset.code",
      header: "Kode Aset",
      cell: ({ row }) => (
        <span className="font-mono text-sm">{row.original.asset?.code || "-"}</span>
      ),
    },
    {
      accessorKey: "asset.name",
      header: "Nama Aset",
      cell: ({ row }) => (
        <span className="max-w-[200px] truncate block">{row.original.asset?.name || "-"}</span>
      ),
    },
    {
      accessorKey: "from_room.name",
      header: "Dari Ruangan",
      cell: ({ row }) => (
        <span>{row.original.from_room?.name || "-"}</span>
      ),
    },
    {
      accessorKey: "to_room.name",
      header: "Ke Ruangan",
      cell: ({ row }) => (
        <Badge variant="default">{row.original.to_room?.name || "-"}</Badge>
      ),
    },
    {
      accessorKey: "reason",
      header: "Alasan",
      cell: ({ row }) => (
        <span className="max-w-[150px] truncate block text-muted-foreground">
          {row.original.reason || "-"}
        </span>
      ),
    },
    {
      accessorKey: "mover.username",
      header: "Dipindahkan Oleh",
      cell: ({ row }) => (
        <span className="text-sm">{row.original.mover?.username || "-"}</span>
      ),
    },
    {
      accessorKey: "created_at",
      header: "Tanggal",
      cell: ({ row }) => (
        <span className="text-sm text-muted-foreground">
          {format(new Date(row.original.created_at), "dd MMM yyyy HH:mm")}
        </span>
      ),
    },
  ];
}
