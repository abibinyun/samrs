export type Room = {
  id: string;
  tenant_id: string;
  name: string;
  code: string;
  location?: string;
  created_at: string;
  updated_at: string;
};

export type RoomFormInput = {
  name: string;
  code: string;
  location?: string;
};
