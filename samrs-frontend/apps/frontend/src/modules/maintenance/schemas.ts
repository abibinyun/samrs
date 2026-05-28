import { z } from "zod";

export const maintenanceSchema = z.object({
  asset_id: z.string().min(1, "Asset is required"),
  schedule_type: z.enum(["maintenance", "calibration"]),
  title: z.string().min(3, "Title must be at least 3 characters"),
  interval_days: z.number().min(1, "Interval must be at least 1 day"),
  next_due_date: z.string().min(1, "Next due date is required"),
  status: z.enum(["scheduled", "in_progress", "completed", "overdue"]),
  notes: z.string().optional(),
});

export type MaintenanceFormValues = z.infer<typeof maintenanceSchema>;
