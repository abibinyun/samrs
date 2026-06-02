import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { maintenanceSchema, type MaintenanceFormValues } from "../schemas";
import { useCreateMaintenanceScheduleMutation } from "../actions/maintenanceApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import { AssetSelect } from "@/modules/stock-opname/components/AssetSelect";

export function MaintenanceForm() {
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const [createSchedule, { isLoading }] = useCreateMaintenanceScheduleMutation();
  const { register, handleSubmit, control, formState: { errors }, setValue } = useForm<MaintenanceFormValues>({
    resolver: zodResolver(maintenanceSchema),
    defaultValues: { schedule_type: "maintenance", status: "scheduled" },
  });

  const onSubmit = async (data: MaintenanceFormValues) => {
    try {
      await createSchedule(data).unwrap();
      toast.success("Jadwal pemeliharaan berhasil dibuat");
      navigate({ to: `${basePath}/maintenance/schedules` });
    } catch {
      toast.error("Gagal membuat jadwal");
    }
  };

  return (
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
        <Label htmlFor="schedule_type">Tipe</Label>
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
        <Input id="title" {...register("title")} placeholder="Judul jadwal" />
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
        <Textarea id="notes" {...register("notes")} rows={3} placeholder="Catatan tambahan..." />
      </div>

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Menyimpan..." : "Simpan"}</Button>
        <Button type="button" variant="outline" onClick={() => navigate({ to: `${basePath}/maintenance/schedules` })}>Batal</Button>
      </div>
    </form>
  );
}
