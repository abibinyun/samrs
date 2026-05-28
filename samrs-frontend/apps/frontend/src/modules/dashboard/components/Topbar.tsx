import { Search, Bell, User as UserIcon, LogOut, Settings } from "lucide-react";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu"; // Shadcn
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "./ThemeToggle";
import { useLogout } from "@/modules/auth/hooks/useLogout";
import { RecentTabsBar } from "@/components/commons/tabs/RecentTabsBar";

interface TopbarProps {
  tenantName: string;
  userName?: string;
}

export const Topbar = ({ tenantName, userName }: TopbarProps) => {
  const { handleLogout } = useLogout();

  return (
    <>
      <header className="h-16 border-b bg-card px-6 flex items-center justify-between sticky top-0 z-30 shadow-sm">      
        {/* Kiri: Nama Tenant / Identitas RS */}
        <div className="flex items-center gap-4">
          <div className="flex flex-col">
            <span className="text-sm font-semibold text-primary uppercase tracking-wider leading-none">
              {tenantName}
            </span>
            <span className="text-[10px] text-muted-foreground mt-1">
              Sistem Informasi Manajemen Rumah Sakit
            </span>
          </div>
        </div>

        {/* Tengah: Global Search (Opsional - Sangat berguna untuk 100 modul) */}
        <div className="hidden lg:flex flex-1 max-w-md mx-8">
          <div className="relative w-full group">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
            <input 
              type="text" 
              placeholder="Cari modul atau data pasien... (Ctrl+K)" 
              className="w-full pl-10 pr-4 py-2 bg-muted/50 border border-transparent rounded-full text-sm focus:outline-none focus:bg-background focus:border-primary transition-all"
            />
          </div>
        </div>

        {/* Kanan: Notifikasi & Profil */}
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" className="relative">
            <Bell className="w-5 h-5 text-muted-foreground" />
            <span className="absolute top-2 right-2 w-2 h-2 bg-destructive rounded-full border-2 border-background"></span>
          </Button>

          <ThemeToggle />

          <div className="h-8 w-px bg-border mx-1"></div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="flex items-center gap-2 px-2 hover:bg-muted transition-colors">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center border border-primary/20">
                  <UserIcon className="w-4 h-4 text-primary" />
                </div>
                <div className="hidden sm:flex flex-col items-start">
                  <span className="text-sm font-medium leading-none">{userName}</span>
                  <span className="text-[10px] text-muted-foreground">Super Admin</span>
                </div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>Akun Saya</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <UserIcon className="mr-2 h-4 w-4" /> Profil
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Settings className="mr-2 h-4 w-4" /> Pengaturan
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem 
                onClick={handleLogout}
                className="text-destructive focus:bg-destructive/10 focus:text-destructive cursor-pointer"
              >
                <LogOut className="mr-2 h-4 w-4" /> Keluar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>
      <RecentTabsBar />
    </>
  );
};