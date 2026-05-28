import type { ColumnDef } from "@tanstack/react-table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { MapPin, Tag, QrCode, MoreVertical } from "lucide-react";
import { cn } from "@/lib/utils";
import type { Asset, AssetWide, AssetDevice } from "../types";
import { assetStatus } from "./AssetStatus";
import { useNavigate } from "@tanstack/react-router";

export function createAssetColumns(opts: {
  onOpenQr: (asset: Pick<Asset, "id" | "name">) => void;
}): ColumnDef<Asset, unknown>[] {
  return [
    {
      accessorKey: "name",
      header: "Info Aset",
      cell: ({ row }) => {
        const a = row.original;
        return (
          <div>
            <div className="font-semibold text-sm">{a.name}</div>
            <div className="text-xs text-muted-foreground font-mono mt-0.5">
              {a.code}
            </div>
          </div>
        );
      },
    },
    {
      id: "locationCategory",
      header: "Lokasi & Kategori",
      cell: ({ row }) => {
        const a = row.original;
        return (
          <div>
            <div className="flex items-center gap-1.5 text-sm">
              <MapPin className="w-3.5 h-3.5 text-muted-foreground" />
              {a.room?.name || "-"}
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-1">
              <Tag className="w-3.5 h-3.5" />
              {a.category?.name || "-"}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const s = assetStatus(row.original.status);
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
      id: "actions",
      header: () => <div className="text-right">Aksi</div>,
      enableSorting: false,
      cell: ({ row }) => {
        const a = row.original;
        const navigate = useNavigate();
        
        return (
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => opts.onOpenQr({ id: a.id, name: a.name })}
              title="Lihat QR Code"
            >
              <QrCode className="w-4 h-4" />
            </Button>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" title="Menu">
                  <MoreVertical className="w-4 h-4 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => navigate({ to: `./$id`, params: { id: a.id } })}>
                  Lihat Detail
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => navigate({ to: `./$id/edit`, params: { id: a.id } })}>
                  Edit
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        );
      },
    },
  ];
}

export function createAssetWideColumns(): ColumnDef<AssetWide, unknown>[] {
  return [
    { accessorKey: "code", header: "Kode" },
    { accessorKey: "name", header: "Nama Aset" },
    { accessorKey: "category", header: "Kategori" },
    { accessorKey: "room", header: "Ruangan" },

    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const s = assetStatus(row.original.status);
        return (
          <span
            className={cn(
              "px-2 py-1 rounded-full text-[10px] font-bold border uppercase",
              s.className,
            )}
          >
            {s.label}
          </span>
        );
      },
    },

    { accessorKey: "purchase_date", header: "Tanggal Beli" },
    { accessorKey: "serial_number", header: "Serial Number" },
    { accessorKey: "manufacturer", header: "Pabrikan" },
    { accessorKey: "model", header: "Model" },
    { accessorKey: "vendor", header: "Vendor" },
    { accessorKey: "price", header: "Harga (Rp)" },
    { accessorKey: "depreciation_rate", header: "Depresiasi (%)" },
    { accessorKey: "ownership", header: "Kepemilikan" },
    { accessorKey: "risk_level", header: "Risiko" },
    { accessorKey: "maintenance_cycle", header: "Siklus Maintenance" },
    { accessorKey: "last_service_date", header: "Service Terakhir" },
    {
      accessorKey: "calibration_required",
      header: "Kalibrasi",
      cell: ({ getValue }) => (getValue<boolean>() ? "Ya" : "Tidak"),
    },
    { accessorKey: "power_rating", header: "Daya" },
    { accessorKey: "weight_kg", header: "Berat (kg)" },
    {
      accessorKey: "notes",
      header: "Catatan",
      cell: ({ getValue }) => (
        <div className="max-w-70 truncate">{getValue<string>()}</div>
      ),
    },
  ];
}

export function createAssetDeviceColumns(): ColumnDef<AssetDevice, unknown>[] {
  return [
    {
      accessorKey: "code",
      header: "Kode",
      cell: ({ getValue }) => (
        <span className="font-mono text-sm">{getValue<string>()}</span>
      ),
    },
    {
      accessorKey: "name",
      header: "Nama Aset",
      cell: ({ getValue }) => (
        <div className="font-medium truncate">{getValue<string>()}</div>
      ),
    },
    {
      accessorKey: "brand",
      header: "Brand",
    },
    {
      accessorKey: "model",
      header: "Model",
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const s = assetStatus(row.original.status);
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
      accessorKey: "purchase_date",
      header: "Tanggal Beli",
      cell: ({ getValue }) => {
        const date = new Date(getValue<string>());
        return (
          <span className="text-sm">
            {date.toLocaleDateString("id-ID", {
              day: "2-digit",
              month: "short",
              year: "numeric",
            })}
          </span>
        );
      },
    },
  ];
}
