export type StockOpnameStatus = "draft" | "closed";
export type StockOpnameCondition = "match" | "missing" | "damaged";

export type StockOpnameSession = {
  id: number;
  tenant_id: string;
  title: string;
  opname_at: string;
  status: StockOpnameStatus;
  notes?: string;
  created_at: string;
  updated_at: string;
};

export type StockOpnameItem = {
  id: number;
  tenant_id: string;
  session_id: number;
  asset_id: string;
  condition: StockOpnameCondition;
  note?: string;
  checked_by: string;
  created_at: string;
  asset?: {
    id: string;
    code: string;
    name: string;
    status: string;
  };
  checker?: {
    id: string;
    username: string;
  };
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
