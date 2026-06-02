import { z } from "zod";

export const complaintSchema = z.object({
  asset_id: z.string().min(1, "Aset wajib dipilih"),
  title: z.string().min(3, "Judul minimal 3 karakter"),
  description: z.string().min(10, "Deskripsi minimal 10 karakter"),
});

export const complaintUpdateSchema = z.object({
  status: z.enum(["open", "in_progress", "done"]),
  assigned_to: z.string().optional().nullable(),
  resolution_note: z.string().optional(),
});

export type ComplaintFormValues = z.infer<typeof complaintSchema>;
export type ComplaintUpdateValues = z.infer<typeof complaintUpdateSchema>;
