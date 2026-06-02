import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { complaintSchema, type ComplaintFormValues } from "../schemas";
import { useCreateComplaintMutation } from "../actions/complaintApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { AssetSelect } from "@/modules/stock-opname/components/AssetSelect";
import { toast } from "sonner";
import { useNavigate, useRouterState } from "@tanstack/react-router";

export function ComplaintForm() {
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const [createComplaint, { isLoading }] = useCreateComplaintMutation();

  const { register, handleSubmit, control, formState: { errors } } = useForm<ComplaintFormValues>({
    resolver: zodResolver(complaintSchema),
  });

  const onSubmit = async (data: ComplaintFormValues) => {
    try {
      await createComplaint(data).unwrap();
      toast.success("Keluhan berhasil dibuat");
      navigate({ to: `${basePath}/complaints` });
    } catch {
      toast.error("Gagal membuat keluhan");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label>Aset</Label>
        <Controller
          control={control}
          name="asset_id"
          render={({ field }) => (
            <AssetSelect value={field.value} onChange={field.onChange} />
          )}
        />
        {errors.asset_id && (
          <p className="text-sm text-red-600 mt-1">{errors.asset_id.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="title">Judul</Label>
        <Input id="title" {...register("title")} placeholder="Judul keluhan" />
        {errors.title && (
          <p className="text-sm text-red-600 mt-1">{errors.title.message}</p>
        )}
      </div>

      <div>
        <Label htmlFor="description">Deskripsi</Label>
        <Textarea id="description" {...register("description")} rows={5} placeholder="Jelaskan masalah yang terjadi..." />
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
          onClick={() => navigate({ to: `${basePath}/complaints` })}
        >
          Batal
        </Button>
      </div>
    </form>
  );
}
