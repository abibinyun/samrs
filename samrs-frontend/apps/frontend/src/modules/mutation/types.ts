export type MutationStatus = "ready" | "maintenance" | "broken" | (string & {});

export type Mutation = {
  id: string;
  code: string;
  name: string;
  category: string;
  room: string;
  status: MutationStatus;
  purchase_date: string;
};

export type MutationWide = Mutation & {
  serial_number: string;
  manufacturer: string;
  model: string;
  warranty_until: string;
  condition_score: number;
  last_service_date: string;
  vendor: string;
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

export type MutationDevice = {
  id: string;
  // relations (ids)
  category_id: number;
  room_id: string;
  bed_id: number | null;
  vendor_id: number;
  brand_id: number;
  model_id: number;

  // display fields
  code: string;
  name: string;
  brand: string;
  model: string;

  status: MutationStatus;
  purchase_date: string; // ISO date string
};