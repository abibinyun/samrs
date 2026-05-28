import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { modelSchema, type ModelFormValues } from "../schemas";
import { useCreateModelMutation, useUpdateModelMutation } from "../actions/modelApi";
import { useGetBrandsQuery } from "@/modules/brand/actions/brandApi";
import type { Model } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { toast } from "sonner";

type ModelFormProps = { model?: Model; onSuccess: () => void; onCancel: () => void };

export function ModelForm({ model, onSuccess, onCancel }: ModelFormProps) {
  const [createModel, { isLoading: isCreating }] = useCreateModelMutation();
  const [updateModel, { isLoading: isUpdating }] = useUpdateModelMutation();
  const { data: brandsData } = useGetBrandsQuery({});
  const isLoading = isCreating || isUpdating;

  const { register, handleSubmit, formState: { errors }, setValue } = useForm<ModelFormValues>({
    resolver: zodResolver(modelSchema),
    defaultValues: model ? { brand_id: model.brand_id, name: model.name, description: model.description } : undefined,
  });

  const onSubmit = async (data: ModelFormValues) => {
    try {
      if (model) {
        await updateModel({ id: model.id, data }).unwrap();
        toast.success("Model updated");
      } else {
        await createModel(data).unwrap();
        toast.success("Model created");
      }
      onSuccess();
    } catch {
      toast.error("Failed to save model");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="brand_id">Brand</Label>
        <Select onValueChange={(v) => setValue("brand_id", parseInt(v))} defaultValue={model?.brand_id?.toString()}>
          <SelectTrigger><SelectValue placeholder="Select brand" /></SelectTrigger>
          <SelectContent>
            {brandsData?.data.map((brand) => <SelectItem key={brand.id} value={brand.id.toString()}>{brand.name}</SelectItem>)}
          </SelectContent>
        </Select>
        {errors.brand_id && <p className="text-sm text-red-600 mt-1">{errors.brand_id.message}</p>}
      </div>
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
        <Button type="submit" disabled={isLoading}>{isLoading ? "Saving..." : model ? "Update" : "Create"}</Button>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}
