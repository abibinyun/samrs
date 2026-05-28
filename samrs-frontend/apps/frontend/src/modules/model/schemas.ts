import { z } from "zod";
export const modelSchema = z.object({ brand_id: z.number().min(1, "Brand is required"), name: z.string().min(1, "Name is required"), description: z.string().optional() });
export type ModelFormValues = z.infer<typeof modelSchema>;
