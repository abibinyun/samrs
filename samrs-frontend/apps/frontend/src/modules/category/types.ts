export type Category = {
  id: number;
  tenant_id: string;
  name: string;
  slug: string;
  description?: string;
  created_at: string;
  updated_at: string;
};

export type CategoryFormInput = {
  name: string;
  description?: string;
};
