export type AssetStatus = "ready" | "maintenance" | "broken";

export type Asset = {
  id: string;
  tenant_id: string;
  category_id: number;
  room_id?: string | null;
  bed_id?: number | null;
  vendor_id?: number | null;
  brand_id?: number | null;
  model_id?: number | null;
  code: string;
  name: string;
  brand?: string;
  model?: string;
  status: AssetStatus;
  purchase_date?: string | null;
  created_at: string;
  updated_at: string;
  category?: {
    id: number;
    name: string;
  };
  room?: {
    id: string;
    name: string;
  };
  bed?: {
    id: number;
    code?: string;
    name?: string;
  };
  vendor?: {
    id: number;
    name: string;
  };
  brand_master?: {
    id: number;
    name: string;
  };
  model_master?: {
    id: number;
    name: string;
  };
};

export type AssetFormInput = {
  category_id: number;
  room_id?: string;
  bed_id?: number;
  vendor_id?: number;
  brand_id?: number;
  model_id?: number;
  code: string;
  name: string;
  brand?: string;
  model?: string;
  status: AssetStatus;
  purchase_date?: string;
};

// Legacy types - keep for backward compatibility
export type AssetWide = Asset & {
  serial_number: string;
  manufacturer: string;
  warranty_until: string;
  condition_score: number;
  last_service_date: string;
  price: number;
  depreciation_rate: number;
  ownership: string;
  risk_level: string;
  maintenance_cycle: string;
  calibration_required: boolean;
  power_rating: string;
  weight_kg: number;
  notes: string;
};

export type AssetDevice = {
  id: string;
  category_id: number;
  room_id: string;
  bed_id: number | null;
  vendor_id: number;
  brand_id: number;
  model_id: number;
  code: string;
  name: string;
  brand: string;
  model: string;
  status: AssetStatus;
  purchase_date: string;
};
