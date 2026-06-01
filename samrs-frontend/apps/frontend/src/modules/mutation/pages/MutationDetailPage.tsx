import { useNavigate, useParams } from "@tanstack/react-router";
import { useGetMutationByIdQuery } from "../actions/mutationApi";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ArrowLeft } from "lucide-react";
import { format } from "date-fns";

function parsePath() {
  const parts = window.location.pathname.split("/").filter(Boolean);
  const idx = parts.indexOf("mutation");
  return {
    tenant: parts[0] || "",
    id: idx >= 0 ? parts[idx + 1] : undefined,
  };
}

export default function MutationDetailPage() {
  const navigate = useNavigate();
  const { tenant, id } = parsePath();
  const mutationId = id ? Number(id) : undefined;
  const { data, isLoading, error } = useGetMutationByIdQuery(mutationId!, { skip: !mutationId });

  if (isLoading) {
    return (
      <PageContainer>
        <div className="flex items-center justify-center h-64">
          <p>Loading...</p>
        </div>
      </PageContainer>
    );
  }

  if (error || !data?.data) {
    return (
      <PageContainer>
        <div className="flex flex-col items-center justify-center h-64 gap-4">
          <p className="text-destructive">Gagal memuat data mutasi</p>
          <Button variant="outline" onClick={() => navigate({ to: `/${tenant}/assets/mutation` })}>
            Kembali
          </Button>
        </div>
      </PageContainer>
    );
  }

  const m = data.data;

  return (
    <PageContainer
      header={{
        title: "Detail Mutasi",
        description: `Mutasi #${m.id}`,
        actions: (
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate({ to: `/${tenant}/assets/mutation` })}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
        ),
      }}
    >
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Informasi Aset</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Kode Aset</p>
              <p className="font-mono font-medium">{m.asset?.code || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Nama Aset</p>
              <p className="font-medium">{m.asset?.name || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Status</p>
              <Badge variant={m.asset?.status === "ready" ? "default" : m.asset?.status === "broken" ? "destructive" : "secondary"}>
                {m.asset?.status || "-"}
              </Badge>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Kategori</p>
              <p className="font-medium">{m.asset?.category?.name || "-"}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Detail Mutasi</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Dari Ruangan</p>
              <p className="font-medium">{m.from_room?.name || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Ke Ruangan</p>
              <Badge variant="default">{m.to_room?.name || "-"}</Badge>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Alasan</p>
              <p className="font-medium">{m.reason || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Dipindahkan Oleh</p>
              <p className="font-medium">{m.mover?.username || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Tanggal Mutasi</p>
              <p className="font-medium">{format(new Date(m.created_at), "dd MMMM yyyy HH:mm")}</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </PageContainer>
  );
}
