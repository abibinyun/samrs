import { 
  LayoutDashboard, Box, Wrench, ClipboardList, FileText,
  MapPin, Tags, History, ShieldCheck, Bell, FormIcon,
  type LucideIcon, 
  FileCheck,
  Landmark,
  AlertTriangle
} from "lucide-react";
import { SCOPES } from "./permissions";

export interface SubMenuItem {
  label: string;
  to: string;
  permission: string;
}

export interface MenuItem {
  label: string;
  icon: LucideIcon;
  to?: string;
  permission?: string;
  children?: SubMenuItem[];
}

export interface MenuSection {
  group: string;
  items: MenuItem[];
}

export const MENU_ITEMS: MenuSection[] = [
  {
    group: "Dashboard",
    items: [
      { 
        label: "Ringkasan", 
        icon: LayoutDashboard, 
        to: "/", 
        permission: SCOPES.TENANT.READ
      },
      { 
        label: "Notifikasi", 
        icon: Bell, 
        to: "/notifications", 
        permission: SCOPES.NOTIFICATION.SEND 
      },
    ],
  },
  {
    group: "Manajemen Aset",
    items: [
      { 
        label: "Inventaris Aset", 
        icon: Box, 
        to: "/assets", 
        permission: SCOPES.ASSET.READ 
      },
      { 
        label: "Mutasi & Timeline", 
        icon: History, 
        to: "/assets/mutation", 
        permission: SCOPES.ASSET_MUTATION.READ 
      },
      { 
        label: "Stock Opname", 
        icon: FileCheck, 
        to: "/stock-opname", 
        permission: SCOPES.STOCK_OPNAME.READ 
      },
    ],
  },
  {
    group: "Layanan Teknis",
    items: [
      { 
        label: "Keluhan (Complaint)", 
        icon: AlertTriangle,
        to: "/complaints", 
        permission: SCOPES.COMPLAINT.READ 
      },
      { 
        label: "Pemeliharaan", 
        icon: Wrench, 
        permission: SCOPES.MAINTENANCE.READ,
        children: [
          { label: "Jadwal & Kalibrasi", to: "/maintenance/schedules", permission: SCOPES.MAINTENANCE.READ },
          { label: "Dokumen Teknis", to: "/maintenance/documents", permission: SCOPES.MAINTENANCE_DOCUMENT.READ },
        ]
      },
      { 
        label: "Forms", 
        icon: FormIcon,
        to: "/forms", 
        permission: SCOPES.COMPLAINT.READ 
      }
    ],
  },
  {
    group: "Master Data",
    items: [
      { 
        label: "Struktur Ruangan", 
        icon: MapPin, 
        children: [
          { label: "Data Ruangan", to: "/master/rooms", permission: SCOPES.ROOM.READ },
          { label: "Data Tempat Tidur", to: "/master/beds", permission: SCOPES.BED.READ },
        ]
      },
      { 
        label: "Katalog & Atribut", 
        icon: Tags, 
        children: [
          { label: "Kategori Aset", to: "/master/categories", permission: SCOPES.CATEGORY.READ },
          { label: "Vendor / Pemasok", to: "/master/vendors", permission: SCOPES.VENDOR.READ },
          { label: "Brand & Model", to: "/master/brands", permission: SCOPES.BRAND.READ },
          { label: "Status Kondisi", to: "/master/asset-status", permission: SCOPES.ASSET_STATUS.READ },
        ]
      },
    ],
  },
  {
    group: "Pusat Dokumen",
    items: [
      { 
        label: "Perpustakaan SOP", 
        icon: FileText, 
        to: "/documents", 
        permission: SCOPES.DOCUMENT.READ 
      },
      { 
        label: "Ekspor Laporan", 
        icon: ClipboardList, 
        to: "/reports", 
        permission: SCOPES.REPORT.EXPORT 
      },
    ],
  },
  {
    group: "Administrator",
    items: [
      { 
        label: "Data Tenant", 
        icon: Landmark, 
        to: "/admin/tenants", 
        permission: SCOPES.TENANT.CREATE
      },
      { 
        label: "Keamanan", 
        icon: ShieldCheck, 
        children: [
          { label: "Pengguna", to: "/settings/users", permission: SCOPES.USER.READ },
          { label: "Role & Permission", to: "/settings/roles", permission: SCOPES.ROLE.READ },
          { label: "Audit Trail", to: "/settings/audit", permission: SCOPES.AUDIT.READ },
        ]
      },
    ],
  },
];