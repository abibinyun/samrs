import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { complaintSchema, type ComplaintFormValues } from "../schemas";
import { useCreateComplaintMutation } from "../actions/complaintApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";

export function ComplaintForm() {
  const navigate = useNavigate();
  const [createComplaint, { isLoading }] = useCreateComplaintMutation();
  
  const { register, handleSubmit, formState: { errors } } = useForm<ComplaintFormValues>({
    resolver: zodResolver(complaintSchema),
  });

  const onSubmit = async (data: ComplaintFormValues) => {
    try {
      await createComplaint(data).unwrap();
      toast.success("Keluhan berhasil dibuat");
      navigate({ to: "/complaints" });
    } catch (error) {
      toast.error("Gagal membuat keluhan");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="asset_id">Asset ID</Label>
        <Input id="asset_id" {...register("asset_id")} />
        {errors.asset_id && (
          <p className="text-sm text-red-600 mt-1">{errors.asset_id.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="title">Judul</Label>
        <Input id="title" {...register("title")} />
        {errors.title && (
          <p className="text-sm text-red-600 mt-1">{errors.title.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="description">Deskripsi</Label>
        <Textarea id="description" {...register("description")} rows={5} />
        {errors.description && (
          <p className="text-sm text-red-600 mt-1">{errors.description.message}</p>
        )}
      </div>

      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>
          {isLoading ? "Menyimpan..." : "Simpan"}
        </Button>
        <Button
          type="button"
          variant="outline"
          onClick={() => navigate({ to: "/complaints" })}
        >
          Batal
        </Button>
      </div>
    </form>
  );
}
