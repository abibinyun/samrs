import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { maintenanceSchema, type MaintenanceFormValues } from "../schemas";
import { useCreateMaintenanceScheduleMutation } from "../actions/maintenanceApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";

export function MaintenanceForm() {
  const navigate = useNavigate();
  const [createSchedule, { isLoading }] = useCreateMaintenanceScheduleMutation();
  const { register, handleSubmit, formState: { errors }, setValue } = useForm<MaintenanceFormValues>({
    resolver: zodResolver(maintenanceSchema),
    defaultValues: { schedule_type: "maintenance", status: "scheduled" },
  });

  const onSubmit = async (data: MaintenanceFormValues) => {
    try {
      await createSchedule(data).unwrap();
      toast.success("Jadwal pemeliharaan berhasil dibuat");
      navigate({ to: "/maintenance/schedules" });
    } catch (error) {
      toast.error("Gagal membuat jadwal");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="asset_id">Asset ID</Label>
        <Input id="asset_id" {...register("asset_id")} />
        {errors.asset_id && <p className="text-sm text-red-600 mt-1">{errors.asset_id.message}</p>}
      </div>

      <div>
        <Label htmlFor="schedule_type">Tipe</Label>
        <Select onValueChange={(v) => setValue("schedule_type", v as any)} defaultValue="maintenance">
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="maintenance">Pemeliharaan</SelectItem>
            <SelectItem value="calibration">Kalibrasi</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div>
        <Label htmlFor="title">Judul</Label>
        <Input id="title" {...register("title")} />
        {errors.title && <p className="text-sm text-red-600 mt-1">{errors.title.message}</p>}
      </div>

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

      <div>
        <Label htmlFor="notes">Catatan</Label>
        <Textarea id="notes" {...register("notes")} rows={3} />
      </div>

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Menyimpan..." : "Simpan"}</Button>
        <Button type="button" variant="outline" onClick={() => navigate({ to: "/maintenance/schedules" })}>Batal</Button>
      </div>
    </form>
  );
}
