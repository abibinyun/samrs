import { z } from "zod";

export const assetSchema = z.object({
  category_id: z.number().min(1, "Category is required"),
  room_id: z.string().optional(),
  bed_id: z.number().optional(),
  vendor_id: z.number().optional(),
  brand_id: z.number().optional(),
  model_id: z.number().optional(),
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  brand: z.string().optional(),
  model: z.string().optional(),
  status: z.enum(["ready", "maintenance", "broken"]),
  purchase_date: z.string().optional(),
});

export type AssetFormValues = z.infer<typeof assetSchema>;