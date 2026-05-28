export type Bed = { id: number; tenant_id: string; room_id: string; name: string; is_available: boolean; created_at: string; updated_at: string; room?: { id: string; name: string } };
export type BedFormInput = { room_id: string; name: string; is_available: boolean };
