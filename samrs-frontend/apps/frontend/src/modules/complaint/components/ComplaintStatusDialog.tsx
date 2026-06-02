import { useEffect } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { complaintUpdateSchema, type ComplaintUpdateValues } from "../schemas";
import { useUpdateComplaintMutation, useGetUsersQuery } from "../actions/complaintApi";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";
import type { Complaint, ComplaintStatus } from "../types";

type ComplaintStatusDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  complaint: Complaint | null;
};

export function ComplaintStatusDialog({ open, onOpenChange, complaint }: ComplaintStatusDialogProps) {
  const [updateComplaint, { isLoading }] = useUpdateComplaintMutation();
  const { data: usersData } = useGetUsersQuery();
  const users = usersData?.data ?? [];

  const { register, handleSubmit, control, watch, reset, formState: { errors } } = useForm<ComplaintUpdateValues>({
    resolver: zodResolver(complaintUpdateSchema),
    defaultValues: {
      status: complaint?.status ?? "open",
      assigned_to: complaint?.assigned_to ?? null,
      resolution_note: complaint?.resolution_note ?? "",
    },
  });

  const currentStatus = watch("status");

  useEffect(() => {
    if (complaint) {
      reset({
        status: complaint.status,
        assigned_to: complaint.assigned_to ?? null,
        resolution_note: complaint.resolution_note ?? "",
      });
    }
  }, [complaint, reset]);

  const onSubmit = async (data: ComplaintUpdateValues) => {
    if (!complaint) return;
    try {
      await updateComplaint({ id: complaint.id, data }).unwrap();
      toast.success("Keluhan berhasil diupdate");
      onOpenChange(false);
    } catch {
      toast.error("Gagal mengupdate keluhan");
    }
  };

  if (!complaint) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Update Keluhan</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <Label>Status</Label>
            <Controller
              control={control}
              name="status"
              render={({ field }) => (
                <Select value={field.value} onValueChange={field.onChange}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="open">Open</SelectItem>
                    <SelectItem value="in_progress">In Progress</SelectItem>
                    <SelectItem value="done">Done</SelectItem>
                  </SelectContent>
                </Select>
              )}
            />
            {errors.status && (
              <p className="text-sm text-red-600 mt-1">{errors.status.message}</p>
            )}
          </div>

          <div>
            <Label>Ditugaskan Ke</Label>
            <Controller
              control={control}
              name="assigned_to"
              render={({ field }) => (
                <Select
                  value={field.value ?? "__none__"}
                  onValueChange={(v) => field.onChange(v === "__none__" ? null : v)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Pilih teknisi..." />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">Belum ditugaskan</SelectItem>
                    {users.map((user) => (
                      <SelectItem key={user.id} value={user.id}>
                        {user.username}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
          </div>

          {currentStatus === "done" && (
            <div>
              <Label htmlFor="resolution_note">Catatan Resolusi</Label>
              <Textarea
                id="resolution_note"
                {...register("resolution_note")}
                rows={3}
                placeholder="Jelaskan cara menyelesaikan keluhan..."
              />
              {errors.resolution_note && (
                <p className="text-sm text-red-600 mt-1">{errors.resolution_note.message}</p>
              )}
            </div>
          )}

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Batal
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Menyimpan..." : "Simpan"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
