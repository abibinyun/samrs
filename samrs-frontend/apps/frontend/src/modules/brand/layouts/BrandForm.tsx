import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { brandSchema, type BrandFormValues } from "../schemas";
import { useCreateBrandMutation, useUpdateBrandMutation } from "../actions/brandApi";
import type { Brand } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

type BrandFormProps = { brand?: Brand; onSuccess: () => void; onCancel: () => void };

export function BrandForm({ brand, onSuccess, onCancel }: BrandFormProps) {
  const [createBrand, { isLoading: isCreating }] = useCreateBrandMutation();
  const [updateBrand, { isLoading: isUpdating }] = useUpdateBrandMutation();
  const isLoading = isCreating || isUpdating;

  const { register, handleSubmit, formState: { errors } } = useForm<BrandFormValues>({
    resolver: zodResolver(brandSchema),
    defaultValues: brand || undefined,
  });

  const onSubmit = async (data: BrandFormValues) => {
    try {
      if (brand) {
        await updateBrand({ id: brand.id, data }).unwrap();
        toast.success("Brand updated");
      } else {
        await createBrand(data).unwrap();
        toast.success("Brand created");
      }
      onSuccess();
    } catch {
      toast.error("Failed to save brand");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="name">Name</Label>
        <Input id="name" {...register("name")} />
        {errors.name && <p className="text-sm text-red-600 mt-1">{errors.name.message}</p>}
      </div>
      <div>
        <Label htmlFor="description">Description</Label>
        <Textarea id="description" {...register("description")} rows={3} />
      </div>
      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Saving..." : brand ? "Update" : "Create"}</Button>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}
