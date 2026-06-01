export type AssetMutation = {
  id: number;
  tenant_id: string;
  asset_id: string;
  from_room_id: string | null;
  from_bed_id: number | null;
  to_room_id: string | null;
  to_bed_id: number | null;
  reason: string;
  moved_by: string;
  created_at: string;

  // Preloaded relations
  asset?: {
    id: string;
    code: string;
    name: string;
    status: string;
    category?: { id: number; name: string };
  };
  mover?: {
    id: string;
    username: string;
  };
  from_room?: {
    id: string;
    code: string;
    name: string;
  };
  to_room?: {
    id: string;
    code: string;
    name: string;
  };
  from_bed?: {
    id: number;
    code: string;
    name: string;
  };
  to_bed?: {
    id: number;
    code: string;
    name: string;
  };
};

export type MutationFormInput = {
  asset_id: string;
  room_id?: string;
  bed_id?: number;
  reason?: string;
};

export type MutationListResponse = {
  data: AssetMutation[];
  total: number;
};

export type MutationApiResponse<T = unknown> = {
  success: boolean;
  message: string;
  data: T;
};
