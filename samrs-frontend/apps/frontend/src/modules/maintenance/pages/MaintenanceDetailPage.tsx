import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Upload, Trash2, Download, Calendar } from "lucide-react";
import { useNavigate, useParams, useRouterState } from "@tanstack/react-router";
import { useGetMaintenanceScheduleQuery, useGetDocumentsQuery, useDeleteDocumentMutation } from "../actions/maintenanceApi";
import { maintenanceStatus } from "../components/MaintenanceStatus";
import { DocumentUploadDialog } from "../components/DocumentUploadDialog";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { id as idLocale } from "date-fns/locale";
import { toast } from "sonner";

export function MaintenanceDetailPage() {
  const { scheduleId } = useParams({ strict: false });
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data, isLoading } = useGetMaintenanceScheduleQuery(Number(scheduleId));
  const { data: docsData, isLoading: docsLoading } = useGetDocumentsQuery(Number(scheduleId));
  const [deleteDocument] = useDeleteDocumentMutation();
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

  const handleDeleteDoc = async (id: number) => {
    if (confirm("Yakin ingin menghapus dokumen ini?")) {
      try {
        await deleteDocument(id).unwrap();
        toast.success("Dokumen berhasil dihapus");
      } catch {
        toast.error("Gagal menghapus dokumen");
      }
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  if (isLoading) return <div className="p-8 text-center text-muted-foreground">Loading...</div>;

  const schedule = data?.data;
  if (!schedule) return <div className="p-8 text-center text-muted-foreground">Jadwal tidak ditemukan</div>;

  const statusInfo = maintenanceStatus(schedule.status);
  const documents = docsData?.data ?? [];

  return (
    <PageContainer
      header={{
        title: schedule.title,
        description: `${schedule.schedule_type === "maintenance" ? "Pemeliharaan" : "Kalibrasi"} — ${schedule.asset_code}`,
        actions: (
          <Button variant="outline" onClick={() => navigate({ to: `${basePath}/maintenance/schedules` })}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
        ),
      }}
    >
      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2 space-y-6">
          {/* Schedule Info */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Informasi Jadwal</CardTitle>
                <Badge variant="outline" className={statusInfo.className}>{statusInfo.label}</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <div className="text-sm text-muted-foreground">Aset</div>
                  <div className="text-sm font-medium">{schedule.asset_name}</div>
                  <div className="text-xs text-muted-foreground font-mono">{schedule.asset_code}</div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">Tipe</div>
                  <div className="text-sm font-medium">
                    {schedule.schedule_type === "maintenance" ? "Pemeliharaan" : "Kalibrasi"}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground">Interval</div>
                  <div className="text-sm font-medium">{schedule.interval_days} hari</div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground flex items-center gap-1">
                    <Calendar className="w-3.5 h-3.5" /> Jadwal Berikutnya
                  </div>
                  <div className="text-sm font-medium">
                    {format(new Date(schedule.next_due_date), "dd MMMM yyyy", { locale: idLocale })}
                  </div>
                </div>
              </div>
              {schedule.notes && (
                <div>
                  <div className="text-sm text-muted-foreground">Catatan</div>
                  <p className="text-sm whitespace-pre-wrap">{schedule.notes}</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Documents Section */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Dokumen Teknis</CardTitle>
                <Button size="sm" onClick={() => setUploadDialogOpen(true)}>
                  <Upload className="w-4 h-4 mr-2" />
                  Upload
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {docsLoading ? (
                <div className="text-center text-muted-foreground py-4">Loading dokumen...</div>
              ) : documents.length === 0 ? (
                <div className="text-center text-muted-foreground py-8">
                  Belum ada dokumen. Klik "Upload" untuk menambahkan.
                </div>
              ) : (
                <div className="space-y-2">
                  {documents.map((doc) => (
                    <div key={doc.id} className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50">
                      <div className="flex-1 min-w-0">
                        <div className="text-sm font-medium truncate">{doc.filename}</div>
                        <div className="text-xs text-muted-foreground">
                          {doc.doc_type} • {formatFileSize(doc.size)}
                        </div>
                      </div>
                      <div className="flex items-center gap-1 ml-4">
                        <a
                          href={doc.file_url}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          <Button variant="ghost" size="icon" title="Download">
                            <Download className="w-4 h-4" />
                          </Button>
                        </a>
                        <Button variant="ghost" size="icon" onClick={() => handleDeleteDoc(doc.id)} title="Hapus">
                          <Trash2 className="w-4 h-4 text-red-500" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Riwayat</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3 text-sm">
              <div>
                <div className="text-muted-foreground">Dibuat</div>
                <div>{format(new Date(schedule.created_at), "dd MMM yyyy HH:mm", { locale: idLocale })}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Diupdate</div>
                <div>{format(new Date(schedule.updated_at), "dd MMM yyyy HH:mm", { locale: idLocale })}</div>
              </div>
              {schedule.last_done_at && (
                <div>
                  <div className="text-muted-foreground">Terakhir Selesai</div>
                  <div>{format(new Date(schedule.last_done_at), "dd MMM yyyy HH:mm", { locale: idLocale })}</div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      <DocumentUploadDialog
        open={uploadDialogOpen}
        onOpenChange={setUploadDialogOpen}
        scheduleId={Number(scheduleId)}
      />
    </PageContainer>
  );
}
