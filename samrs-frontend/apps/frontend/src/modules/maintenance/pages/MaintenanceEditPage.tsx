import { useEffect } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { maintenanceSchema, type MaintenanceFormValues } from "../schemas";
import { useGetMaintenanceScheduleQuery, useUpdateMaintenanceScheduleMutation } from "../actions/maintenanceApi";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AssetSelect } from "@/modules/stock-opname/components/AssetSelect";
import { toast } from "sonner";
import { useNavigate, useParams, useRouterState } from "@tanstack/react-router";

export function MaintenanceEditPage() {
  const { scheduleId } = useParams({ strict: false });
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data, isLoading } = useGetMaintenanceScheduleQuery(Number(scheduleId));
  const [updateSchedule, { isLoading: isSaving }] = useUpdateMaintenanceScheduleMutation();

  const { register, handleSubmit, control, reset, formState: { errors } } = useForm<MaintenanceFormValues>({
    resolver: zodResolver(maintenanceSchema),
  });

  useEffect(() => {
    if (data?.data) {
      const s = data.data;
      reset({
        asset_id: s.asset_id,
        schedule_type: s.schedule_type,
        title: s.title,
        interval_days: s.interval_days,
        next_due_date: s.next_due_date,
        status: s.status,
        notes: s.notes ?? "",
      });
    }
  }, [data, reset]);

  const onSubmit = async (formData: MaintenanceFormValues) => {
    try {
      await updateSchedule({ id: Number(scheduleId), data: formData }).unwrap();
      toast.success("Jadwal berhasil diupdate");
      navigate({ to: `${basePath}/maintenance/schedules` });
    } catch {
      toast.error("Gagal mengupdate jadwal");
    }
  };

  if (isLoading) return <div className="p-8 text-center text-muted-foreground">Loading...</div>;

  const schedule = data?.data;
  if (!schedule) return <div className="p-8 text-center text-muted-foreground">Jadwal tidak ditemukan</div>;

  return (
    <PageContainer
      header={{
        title: "Edit Jadwal",
        description: `Edit jadwal: ${schedule.title}`
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label>Aset</Label>
              <Controller
                control={control}
                name="asset_id"
                render={({ field }) => (
                  <AssetSelect value={field.value} onChange={field.onChange} />
                )}
              />
              {errors.asset_id && <p className="text-sm text-red-600 mt-1">{errors.asset_id.message}</p>}
            </div>

            <div>
              <Label>Tipe</Label>
              <Controller
                control={control}
                name="schedule_type"
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="maintenance">Pemeliharaan</SelectItem>
                      <SelectItem value="calibration">Kalibrasi</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>

            <div>
              <Label htmlFor="title">Judul</Label>
              <Input id="title" {...register("title")} />
              {errors.title && <p className="text-sm text-red-600 mt-1">{errors.title.message}</p>}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="interval_days">Interval (hari)</Label>
                <Input id="interval_days" type="number" {...register("interval_days", { valueAsNumber: true })} />
                {errors.interval_days && <p className="text-sm text-red-600 mt-1">{errors.interval_days.message}</p>}
              </div>
              <div>
                <Label htmlFor="next_due_date">Jadwal Berikutnya</Label>
                <Input id="next_due_date" type="date" {...register("next_due_date")} />
                {errors.next_due_date && <p className="text-sm text-red-600 mt-1">{errors.next_due_date.message}</p>}
              </div>
            </div>

            <div>
              <Label htmlFor="notes">Catatan</Label>
              <Textarea id="notes" {...register("notes")} rows={3} />
            </div>

            <div className="flex gap-2">
              <Button type="submit" disabled={isSaving}>{isSaving ? "Menyimpan..." : "Simpan"}</Button>
              <Button type="button" variant="outline" onClick={() => navigate({ to: `${basePath}/maintenance/schedules` })}>Batal</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </PageContainer>
  );
}
