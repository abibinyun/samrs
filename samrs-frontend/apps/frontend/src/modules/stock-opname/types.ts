export type StockOpnameStatus = "open" | "closed" | (string & {});
export type StockOpnameCondition = "match" | "missing" | "excess" | "damaged" | (string & {});

export type StockOpnameSession = {
  id: number;
  title: string;
  opname_at: string;
  status: StockOpnameStatus;
  notes?: string;
  created_at: string;
  updated_at: string;
};

export type StockOpnameItem = {
  id: number;
  session_id: number;
  asset_id: string;
  asset_code: string;
  asset_name: string;
  condition: StockOpnameCondition;
  note?: string;
  created_at: string;
};

export type StockOpnameFormData = {
  title: string;
  opname_at: string;
  notes?: string;
};

export type StockOpnameItemFormData = {
  asset_id: string;
  condition: StockOpnameCondition;
  note?: string;
};
