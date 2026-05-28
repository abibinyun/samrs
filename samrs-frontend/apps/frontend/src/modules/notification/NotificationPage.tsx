// src/modules/notifications/NotificationPage.tsx
import { Bell, Wrench, AlertTriangle, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";

// Dummy data berdasarkan skenario API (Maintenance & Complaint)
const NOTIFICATIONS = [
  {
    id: "1",
    title: "Pengingat Maintenance",
    message: "Ventilator Philips V60 (VENT-001) jatuh tempo pemeliharaan hari ini.",
    type: "maintenance",
    time: "5 menit yang lalu",
    isRead: false,
    icon: Wrench,
    color: "text-amber-600",
    bg: "bg-amber-100",
  },
  {
    id: "2",
    title: "Keluhan Baru",
    message: "Laporan kerusakan baru: Ventilator mati total di Ruang Melati 200.",
    type: "complaint",
    time: "1 jam yang lalu",
    isRead: false,
    icon: AlertTriangle,
    color: "text-red-600",
    bg: "bg-red-100",
  },
  {
    id: "3",
    title: "Stock Opname Selesai",
    message: "Sesi Stock Opname Q1 telah ditutup oleh Admin.",
    type: "system",
    time: "Yesterday",
    isRead: true,
    icon: CheckCircle2,
    color: "text-green-600",
    bg: "bg-green-100",
  },
];

export const NotificationPage = () => {
  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Notifikasi</h1>
          <p className="text-muted-foreground text-sm">
            Kelola pengingat dan pemberitahuan sistem Anda.
          </p>
        </div>
        <button className="text-xs text-primary font-medium hover:underline">
          Tandai semua terbaca
        </button>
      </div>

      <div className="space-y-4">
        {NOTIFICATIONS.map((item) => (
          <div
            key={item.id}
            className={cn(
              "flex gap-4 p-4 rounded-xl border transition-all hover:bg-muted/50",
              !item.isRead ? "bg-primary/5 border-primary/20" : "bg-card"
            )}
          >
            <div className={cn("w-10 h-10 rounded-full flex items-center justify-center shrink-0", item.bg)}>
              <item.icon className={cn("w-5 h-5", item.color)} />
            </div>
            
            <div className="flex-1">
              <div className="flex items-center justify-between">
                <h4 className={cn("text-sm font-semibold", !item.isRead ? "text-primary" : "text-foreground")}>
                  {item.title}
                </h4>
                <span className="text-[10px] text-muted-foreground font-medium">
                  {item.time}
                </span>
              </div>
              <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                {item.message}
              </p>
              
              {!item.isRead && (
                <div className="mt-3 flex gap-2">
                  <button className="text-[11px] bg-primary text-white px-3 py-1 rounded-md hover:bg-primary/90 transition-colors">
                    Cek Detail
                  </button>
                  <button className="text-[11px] border px-3 py-1 rounded-md hover:bg-background transition-colors">
                    Abaikan
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {NOTIFICATIONS.length === 0 && (
        <div className="text-center py-20 border rounded-xl border-dashed">
          <Bell className="w-12 h-12 text-muted-foreground mx-auto opacity-20" />
          <p className="mt-4 text-muted-foreground font-medium">Tidak ada notifikasi baru</p>
        </div>
      )}
    </div>
  );
};