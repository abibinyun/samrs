import { z } from "zod";

export const stockOpnameSchema = z.object({
  title: z.string().min(3, "Title must be at least 3 characters"),
  opname_at: z.string().min(1, "Date is required"),
  notes: z.string().optional(),
});

export const stockOpnameItemSchema = z.object({
  asset_id: z.string().min(1, "Asset is required"),
  condition: z.enum(["match", "missing", "damaged"]),
  note: z.string().optional(),
});

export type StockOpnameFormValues = z.infer<typeof stockOpnameSchema>;
export type StockOpnameItemFormValues = z.infer<typeof stockOpnameItemSchema>;
