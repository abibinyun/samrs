import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { stockOpnameSchema, type StockOpnameFormValues } from "../schemas";
import type { StockOpnameFormData } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { useNavigate } from "@tanstack/react-router";

type Props = {
  onSubmit: (data: StockOpnameFormData) => Promise<void>;
  isLoading?: boolean;
  defaultValues?: StockOpnameFormData;
};

export function StockOpnameForm({ onSubmit, isLoading, defaultValues }: Props) {
  const navigate = useNavigate();
  const { register, handleSubmit, formState: { errors } } = useForm<StockOpnameFormValues>({
    resolver: zodResolver(stockOpnameSchema),
    defaultValues,
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="title">Judul Sesi</Label>
        <Input id="title" {...register("title")} />
        {errors.title && <p className="text-sm text-red-600 mt-1">{errors.title.message}</p>}
      </div>
      <div>
        <Label htmlFor="opname_at">Tanggal</Label>
        <Input id="opname_at" type="date" {...register("opname_at")} />
        {errors.opname_at && <p className="text-sm text-red-600 mt-1">{errors.opname_at.message}</p>}
      </div>
      <div>
        <Label htmlFor="notes">Catatan</Label>
        <Textarea id="notes" {...register("notes")} rows={3} />
      </div>
      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Menyimpan..." : "Simpan"}</Button>
        <Button type="button" variant="outline" onClick={() => navigate({ to: "/stock-opname" })}>Batal</Button>
      </div>
    </form>
  );
}
