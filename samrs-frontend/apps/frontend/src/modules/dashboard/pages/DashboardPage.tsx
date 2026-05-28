import { SCOPES } from "@/constants/permissions";
import { usePermission } from "@/hooks/usePermission";
import { 
  Box, Wrench, AlertTriangle, CheckCircle2, 
  ArrowRight, 
  FileText
} from "lucide-react";
import { Link } from "@tanstack/react-router";

export default function DashboardPage() {
  const { hasPermission } = usePermission();

  const stats = [
    { 
      label: "Total Aset", 
      value: "1,284", 
      icon: Box, 
      color: "text-blue-600", 
      permission: SCOPES.ASSET.READ,
      link: "/assets"
    },
    { 
      label: "Maintenance Pending", 
      value: "12", 
      icon: Wrench, 
      color: "text-amber-600", 
      permission: SCOPES.MAINTENANCE.READ,
      link: "/maintenance"
    },
    { 
      label: "Keluhan Aktif", 
      value: "5", 
      icon: AlertTriangle, 
      color: "text-red-600", 
      permission: SCOPES.COMPLAINT.READ,
      link: "/complaints"
    },
  ];

  return (
    <div className="p-6 space-y-8">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Selamat Datang di SAMRS</h1>
        <p className="text-muted-foreground">Ringkasan operasional aset Anda hari ini.</p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-3">
        {stats.map((stat) => (
          hasPermission(stat.permission) && (
            <div key={stat.label} className="bg-card p-6 rounded-xl border shadow-sm">
              <div className="flex items-center justify-between">
                <stat.icon className={`w-8 h-8 ${stat.color}`} />
                <span className="text-2xl font-bold">{stat.value}</span>
              </div>
              <div className="mt-4 flex items-center justify-between">
                <p className="text-sm font-medium text-muted-foreground">{stat.label}</p>
                <Link to={`/$tenantId${stat.link}`} className="text-xs flex items-center gap-1 hover:underline">
                  Lihat Detail <ArrowRight className="w-3 h-3" />
                </Link>
              </div>
            </div>
          )
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {/* Widget: Shortcut Cepat berdasarkan API Section 5.1 & 5.5 */}
        <div className="bg-card p-6 rounded-xl border">
          <h3 className="font-semibold mb-4">Aksi Cepat</h3>
          <div className="space-y-2">
            {hasPermission(SCOPES.ASSET.CREATE) && (
              <button className="w-full text-left px-4 py-3 rounded-lg border hover:bg-muted transition-colors flex items-center gap-3">
                <Box className="w-4 h-4" /> Tambah Inventaris Baru
              </button>
            )}
            {hasPermission(SCOPES.COMPLAINT.CREATE) && (
              <button className="w-full text-left px-4 py-3 rounded-lg border hover:bg-muted transition-colors flex items-center gap-3">
                <AlertTriangle className="w-4 h-4" /> Laporkan Kerusakan Aset
              </button>
            )}
            {hasPermission(SCOPES.REPORT.EXPORT) && (
              <button className="w-full text-left px-4 py-3 rounded-lg border hover:bg-muted transition-colors flex items-center gap-3">
                <FileText className="w-4 h-4" /> Export Laporan Bulanan
              </button>
            )}
          </div>
        </div>

        {/* Widget: Status Sistem (Section 1.2 API) */}
        <div className="bg-card p-6 rounded-xl border flex flex-col justify-center items-center text-center">
          <div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
             <CheckCircle2 className="w-6 h-6 text-green-600" />
          </div>
          <h3 className="font-semibold">Konektivitas Cloud Aktif</h3>
          <p className="text-xs text-muted-foreground mt-1">Tenant ID Verified: rs-pusat</p>
          <div className="mt-4 py-1 px-3 bg-muted rounded-full text-[10px] font-mono">
             Latensi Server: 24ms
          </div>
        </div>
      </div>
    </div>
  );
};