import { z } from "zod";
export const bedSchema = z.object({ room_id: z.string().min(1, "Room is required"), name: z.string().min(1, "Name is required"), is_available: z.boolean() });
export type BedFormValues = z.infer<typeof bedSchema>;
