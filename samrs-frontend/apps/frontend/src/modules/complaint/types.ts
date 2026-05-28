export type ComplaintStatus = "new" | "in_progress" | "resolved" | "closed" | (string & {});

export type Complaint = {
  id: string;
  asset_id: string;
  asset_code: string;
  asset_name: string;
  title: string;
  description: string;
  status: ComplaintStatus;
  resolution_note?: string;
  created_at: string;
  updated_at: string;
};

export type ComplaintFormData = {
  asset_id: string;
  title: string;
  description: string;
};

export type ComplaintUpdateData = {
  status: ComplaintStatus;
  resolution_note?: string;
};
