import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Pencil } from "lucide-react";
import { useNavigate, useParams, useRouterState } from "@tanstack/react-router";
import { useGetComplaintQuery } from "../actions/complaintApi";
import { complaintStatus } from "../components/ComplaintStatus";
import { ComplaintStatusDialog } from "../components/ComplaintStatusDialog";
import type { Complaint } from "../types";

export function ComplaintDetailPage() {
  const { complaintId } = useParams({ strict: false });
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data, isLoading } = useGetComplaintQuery(complaintId!);
  const [statusDialogOpen, setStatusDialogOpen] = useState(false);
  const [selectedComplaint, setSelectedComplaint] = useState<Complaint | null>(null);

  const handleOpenStatusDialog = () => {
    if (data?.data) {
      setSelectedComplaint(data.data);
      setStatusDialogOpen(true);
    }
  };

  if (isLoading) return <div className="p-8 text-center text-muted-foreground">Loading...</div>;

  const complaint = data?.data;
  if (!complaint) return <div className="p-8 text-center text-muted-foreground">Keluhan tidak ditemukan</div>;

  const statusInfo = complaintStatus(complaint.status);

  return (
    <PageContainer
      header={{
        title: complaint.title,
        description: `Detail keluhan #${complaint.id.slice(0, 8)}`,
        actions: (
          <Button variant="outline" onClick={() => navigate({ to: `${basePath}/complaints` })}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
        ),
      }}
    >
      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Deskripsi</CardTitle>
                <Badge variant="outline" className={statusInfo.className}>{statusInfo.label}</Badge>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm whitespace-pre-wrap">{complaint.description}</p>
            </CardContent>
          </Card>

          {complaint.resolution_note && (
            <Card>
              <CardHeader>
                <CardTitle>Catatan Resolusi</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm whitespace-pre-wrap">{complaint.resolution_note}</p>
              </CardContent>
            </Card>
          )}
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Informasi</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div>
                <div className="text-sm text-muted-foreground">Aset</div>
                <div className="text-sm font-medium">{complaint.asset?.name || "-"}</div>
                <div className="text-xs text-muted-foreground font-mono">{complaint.asset?.code || "-"}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Dilaporkan Oleh</div>
                <div className="text-sm font-medium">{complaint.reporter?.username || "-"}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Ditugaskan Ke</div>
                <div className="flex items-center gap-2">
                  <div className="text-sm font-medium">
                    {complaint.assignee?.username || "Belum ditugaskan"}
                  </div>
                  <button
                    type="button"
                    onClick={handleOpenStatusDialog}
                    className="text-muted-foreground hover:text-primary transition-colors"
                    title="Ubah penugasan"
                  >
                    <Pencil className="w-3 h-3" />
                  </button>
                </div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Dibuat</div>
                <div className="text-sm">{new Date(complaint.created_at).toLocaleString("id-ID")}</div>
              </div>
              <div>
                <div className="text-sm text-muted-foreground">Diupdate</div>
                <div className="text-sm">{new Date(complaint.updated_at).toLocaleString("id-ID")}</div>
              </div>
            </CardContent>
          </Card>

          <Button className="w-full" onClick={handleOpenStatusDialog}>
            Update Status
          </Button>
        </div>
      </div>

      <ComplaintStatusDialog
        open={statusDialogOpen}
        onOpenChange={setStatusDialogOpen}
        complaint={selectedComplaint}
      />
    </PageContainer>
  );
}
