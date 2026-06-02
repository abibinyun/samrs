export type ComplaintStatus = "open" | "in_progress" | "done";

export type Complaint = {
  id: string;
  tenant_id: string;
  asset_id: string;
  reported_by: string;
  assigned_to?: string | null;
  title: string;
  description: string;
  resolution_note?: string;
  status: ComplaintStatus;
  created_at: string;
  updated_at: string;
  asset?: {
    id: string;
    code: string;
    name: string;
  };
  reporter?: {
    id: string;
    username: string;
  };
  assignee?: {
    id: string;
    username: string;
  };
};

export type ComplaintFormData = {
  asset_id: string;
  title: string;
  description: string;
};

export type ComplaintUpdateData = {
  status?: ComplaintStatus;
  assigned_to?: string | null;
  resolution_note?: string;
};
