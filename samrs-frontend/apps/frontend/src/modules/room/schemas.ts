import { z } from "zod";

export const roomSchema = z.object({
  name: z.string().min(1, "Name is required"),
  code: z.string().min(1, "Code is required"),
  location: z.string().optional(),
});

export type RoomFormValues = z.infer<typeof roomSchema>;
