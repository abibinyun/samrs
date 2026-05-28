export type MaintenanceStatus = "scheduled" | "in_progress" | "completed" | "overdue" | (string & {});
export type MaintenanceType = "maintenance" | "calibration" | (string & {});

export type MaintenanceSchedule = {
  id: number;
  asset_id: string;
  asset_code: string;
  asset_name: string;
  schedule_type: MaintenanceType;
  title: string;
  interval_days: number;
  next_due_date: string;
  last_done_at?: string;
  status: MaintenanceStatus;
  notes?: string;
  created_at: string;
  updated_at: string;
};

export type MaintenanceFormData = {
  asset_id: string;
  schedule_type: MaintenanceType;
  title: string;
  interval_days: number;
  next_due_date: string;
  status: MaintenanceStatus;
  notes?: string;
};
