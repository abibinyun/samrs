import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Upload, Trash2, Download, FileText } from "lucide-react";
import { useGetMaintenanceSchedulesQuery, useGetDocumentsQuery, useDeleteDocumentMutation } from "../actions/maintenanceApi";
import { DocumentUploadDialog } from "../components/DocumentUploadDialog";
import { maintenanceStatus } from "../components/MaintenanceStatus";
import { cn } from "@/lib/utils";
import { toast } from "sonner";

export function MaintenanceDocumentListPage() {
  const { data: schedulesData } = useGetMaintenanceSchedulesQuery({});
  const [selectedScheduleId, setSelectedScheduleId] = useState<string>("");
  const { data: docsData, isLoading: docsLoading } = useGetDocumentsQuery(
    Number(selectedScheduleId),
    { skip: !selectedScheduleId }
  );
  const [deleteDocument] = useDeleteDocumentMutation();
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

  const schedules = schedulesData?.data ?? [];
  const documents = docsData?.data ?? [];
  const selectedSchedule = schedules.find((s) => s.id === Number(selectedScheduleId));

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

  return (
    <PageContainer
      header={{
        title: "Dokumen Teknis",
        description: "Kelola dokumen teknis pemeliharaan dan kalibrasi",
      }}
    >
      <Card className="mb-6">
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Pilih Jadwal</CardTitle>
            {selectedScheduleId && (
              <Button size="sm" onClick={() => setUploadDialogOpen(true)}>
                <Upload className="w-4 h-4 mr-2" />
                Upload Dokumen
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <Select value={selectedScheduleId} onValueChange={setSelectedScheduleId}>
            <SelectTrigger>
              <SelectValue placeholder="Pilih jadwal pemeliharaan..." />
            </SelectTrigger>
            <SelectContent>
              {schedules.map((s) => (
                <SelectItem key={s.id} value={String(s.id)}>
                  {s.title} — {s.asset_name} ({s.asset_code})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </CardContent>
      </Card>

      {selectedScheduleId && (
        <Card>
          <CardHeader>
            <CardTitle>
              Dokumen — {selectedSchedule?.title}
              {documents.length > 0 && (
                <Badge variant="secondary" className="ml-2">{documents.length}</Badge>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {docsLoading ? (
              <div className="text-center text-muted-foreground py-8">Loading dokumen...</div>
            ) : documents.length === 0 ? (
              <div className="text-center text-muted-foreground py-8">
                <FileText className="w-12 h-12 mx-auto mb-3 opacity-50" />
                <p>Belum ada dokumen untuk jadwal ini.</p>
                <p className="text-xs mt-1">Klik "Upload Dokumen" untuk menambahkan.</p>
              </div>
            ) : (
              <div className="space-y-2">
                {documents.map((doc) => (
                  <div key={doc.id} className="flex items-center justify-between p-4 rounded-lg border hover:bg-muted/50">
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <FileText className="w-8 h-8 text-muted-foreground shrink-0" />
                      <div className="min-w-0">
                        <div className="text-sm font-medium truncate">{doc.filename}</div>
                        <div className="text-xs text-muted-foreground">
                          <Badge variant="outline" className="text-[10px] mr-2">{doc.doc_type}</Badge>
                          {formatFileSize(doc.size)}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-1 ml-4">
                      <a href={doc.file_url} target="_blank" rel="noopener noreferrer">
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
      )}

      {selectedScheduleId && (
        <DocumentUploadDialog
          open={uploadDialogOpen}
          onOpenChange={setUploadDialogOpen}
          scheduleId={Number(selectedScheduleId)}
        />
      )}
    </PageContainer>
  );
}
