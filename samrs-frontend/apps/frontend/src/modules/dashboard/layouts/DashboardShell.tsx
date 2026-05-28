import { Sidebar } from "../components/Sidebar";
import { Topbar } from "../components/Topbar";
import { Suspense } from "react";
import { Loading } from "@/routes/__root";
import { KeepAliveOutlet } from "../tabs/layouts/KeepAliveOutlet";
// import { RecentTabsBar } from "@/components/commons/tabs/RecentTabsBar";
import { useAppSelector } from "@/store/hooks";

export default function DashboardShell() {
  const { user } = useAppSelector((state) => state.auth);

  return (
    <div className="flex min-h-screen bg-background text-foreground font-sans">
      <Sidebar />

      <div className="flex-1 flex flex-col min-w-0">
        <Topbar 
          tenantName={user?.tenant_slug || "Memuat..."} 
          userName={user?.name}
        />

        <main className="flex-1 p-4 md:p-8 custom-scrollbar">
          <div className="max-w-7xl mx-auto container">
            <Suspense fallback={<Loading />}>
              <KeepAliveOutlet />
            </Suspense>
          </div>
        </main>
      </div>
    </div>
  );
};