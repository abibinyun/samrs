export type Model = { id: number; tenant_id: string; brand_id: number; name: string; description?: string; created_at: string; updated_at: string; brand?: { id: number; name: string } };
export type ModelFormInput = { brand_id: number; name: string; description?: string };
