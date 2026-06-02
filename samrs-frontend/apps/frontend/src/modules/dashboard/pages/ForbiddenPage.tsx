import { ShieldOff, ArrowLeft } from "lucide-react";
import { Link, useParams } from "@tanstack/react-router";

export default function ForbiddenPage() {
  const { tenantId } = useParams({ from: "/$tenantId" });

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center p-6">
      <div className="w-16 h-16 bg-destructive/10 rounded-full flex items-center justify-center mb-6">
        <ShieldOff className="w-8 h-8 text-destructive" />
      </div>
      <h1 className="text-2xl font-bold mb-2">Akses Ditolak</h1>
      <p className="text-muted-foreground mb-6 max-w-md">
        Anda tidak memiliki izin untuk mengakses halaman ini. 
        Hubungi administrator jika Anda membutuhkan akses.
      </p>
      <Link
        to={`/$tenantId`}
        params={{ tenantId }}
        className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors"
      >
        <ArrowLeft className="w-4 h-4" />
        Kembali ke Dashboard
      </Link>
    </div>
  );
}
