import { z } from "zod";

export const maintenanceSchema = z.object({
  asset_id: z.string().min(1, "Aset wajib dipilih"),
  schedule_type: z.enum(["maintenance", "calibration"]),
  title: z.string().min(3, "Judul minimal 3 karakter"),
  interval_days: z.number().min(1, "Interval minimal 1 hari"),
  next_due_date: z.string().min(1, "Jadwal berikutnya wajib diisi"),
  status: z.enum(["scheduled", "due", "completed"]),
  notes: z.string().optional(),
});

export type MaintenanceFormValues = z.infer<typeof maintenanceSchema>;
