import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { stockOpnameSchema, type StockOpnameFormValues } from "../schemas";
import { useCreateStockOpnameSessionMutation } from "../actions/stockOpnameApi";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";

export function StockOpnameForm() {
  const navigate = useNavigate();
  const [createSession, { isLoading }] = useCreateStockOpnameSessionMutation();
  const { register, handleSubmit, formState: { errors } } = useForm<StockOpnameFormValues>({
    resolver: zodResolver(stockOpnameSchema),
  });

  const onSubmit = async (data: StockOpnameFormValues) => {
    try {
      await createSession(data).unwrap();
      toast.success("Sesi stock opname berhasil dibuat");
      navigate({ to: "/assets/stock-opname" });
    } catch (error) {
      toast.error("Gagal membuat sesi");
    }
  };

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
        <Button type="button" variant="outline" onClick={() => navigate({ to: "/assets/stock-opname" })}>Batal</Button>
      </div>
    </form>
  );
}
