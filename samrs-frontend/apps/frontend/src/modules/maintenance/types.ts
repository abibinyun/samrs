export type MaintenanceStatus = "scheduled" | "due" | "completed" | (string & {});
export type MaintenanceType = "maintenance" | "calibration" | (string & {});
export type MaintenanceDocType = "certificate" | "report" | "other" | (string & {});

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

export type MaintenanceDocument = {
  id: number;
  tenant_id: string;
  schedule_id: number;
  asset_id: string;
  uploaded_by: string;
  doc_type: MaintenanceDocType;
  filename: string;
  file_path: string;
  mime_type: string;
  size: number;
  file_url?: string;
  created_at: string;
  uploader?: {
    id: string;
    username: string;
  };
  schedule?: {
    id: number;
    title: string;
  };
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

export type MaintenanceDocumentFormData = {
  schedule_id: number;
  doc_type: MaintenanceDocType;
  file: File;
};
