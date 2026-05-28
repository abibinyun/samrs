import { z } from "zod";
export const assetStatusSchema = z.object({ name: z.string().min(1, "Name is required"), description: z.string().optional() });
export type AssetStatusFormValues = z.infer<typeof assetStatusSchema>;
