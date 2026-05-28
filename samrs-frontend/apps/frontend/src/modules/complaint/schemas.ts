import { z } from "zod";

export const complaintSchema = z.object({
  asset_id: z.string().min(1, "Asset is required"),
  title: z.string().min(3, "Title must be at least 3 characters"),
  description: z.string().min(10, "Description must be at least 10 characters"),
});

export const complaintUpdateSchema = z.object({
  status: z.enum(["new", "in_progress", "resolved", "closed"]),
  resolution_note: z.string().optional(),
});

export type ComplaintFormValues = z.infer<typeof complaintSchema>;
export type ComplaintUpdateValues = z.infer<typeof complaintUpdateSchema>;
