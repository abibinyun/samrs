import { useState, useRef } from "react";
import { useUploadDocumentMutation } from "../actions/maintenanceApi";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";
import { Upload, FileText } from "lucide-react";
import { cn } from "@/lib/utils";

type DocumentUploadDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  scheduleId: number;
};

export function DocumentUploadDialog({ open, onOpenChange, scheduleId }: DocumentUploadDialogProps) {
  const [uploadDocument, { isLoading }] = useUploadDocumentMutation();
  const [docType, setDocType] = useState("other");
  const [file, setFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0];
    if (selected) {
      if (selected.size > 10 * 1024 * 1024) {
        toast.error("Ukuran file maksimal 10MB");
        return;
      }
      setFile(selected);
    }
  };

  const handleSubmit = async () => {
    if (!file) {
      toast.error("Pilih file terlebih dahulu");
      return;
    }
    try {
      await uploadDocument({ scheduleId, docType, file }).unwrap();
      toast.success("Dokumen berhasil diunggah");
      setFile(null);
      setDocType("other");
      onOpenChange(false);
    } catch {
      toast.error("Gagal mengunggah dokumen");
    }
  };

  const handleClose = () => {
    setFile(null);
    setDocType("other");
    onOpenChange(false);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Upload Dokumen</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div>
            <Label>Tipe Dokumen</Label>
            <Select value={docType} onValueChange={setDocType}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="certificate">Sertifikat</SelectItem>
                <SelectItem value="report">Laporan</SelectItem>
                <SelectItem value="other">Lainnya</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Label>File</Label>
            <input
              ref={fileInputRef}
              type="file"
              className="hidden"
              accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png"
              onChange={handleFileChange}
            />
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className={cn(
                "flex w-full items-center justify-center gap-2 rounded-md border-2 border-dashed p-6 text-sm transition-colors hover:border-primary/50 hover:bg-muted/50",
                file ? "border-primary" : "border-muted-foreground/25"
              )}
            >
              {file ? (
                <>
                  <FileText className="w-5 h-5 text-primary" />
                  <div className="text-left">
                    <div className="font-medium">{file.name}</div>
                    <div className="text-xs text-muted-foreground">{formatFileSize(file.size)}</div>
                  </div>
                </>
              ) : (
                <>
                  <Upload className="w-5 h-5 text-muted-foreground" />
                  <span>Klik untuk memilih file (PDF, DOC, XLS, JPG, PNG)</span>
                </>
              )}
            </button>
            <p className="text-xs text-muted-foreground mt-1">Maksimal 10MB</p>
          </div>
        </div>
        <DialogFooter>
          <Button type="button" variant="outline" onClick={handleClose}>Batal</Button>
          <Button onClick={handleSubmit} disabled={!file || isLoading}>
            {isLoading ? "Mengunggah..." : "Upload"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
