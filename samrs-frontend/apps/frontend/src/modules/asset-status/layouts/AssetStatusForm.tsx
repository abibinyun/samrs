import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { assetStatusSchema, type AssetStatusFormValues } from "../schemas";
import { useCreateAssetStatusMutation, useUpdateAssetStatusMutation } from "../actions/assetStatusApi";
import type { AssetStatus } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

type AssetStatusFormProps = { assetStatus?: AssetStatus; onSuccess: () => void; onCancel: () => void };

export function AssetStatusForm({ assetStatus, onSuccess, onCancel }: AssetStatusFormProps) {
  const [createAssetStatus, { isLoading: isCreating }] = useCreateAssetStatusMutation();
  const [updateAssetStatus, { isLoading: isUpdating }] = useUpdateAssetStatusMutation();
  const isLoading = isCreating || isUpdating;

  const { register, handleSubmit, formState: { errors } } = useForm<AssetStatusFormValues>({
    resolver: zodResolver(assetStatusSchema),
    defaultValues: assetStatus || undefined,
  });

  const onSubmit = async (data: AssetStatusFormValues) => {
    try {
      if (assetStatus) {
        await updateAssetStatus({ id: assetStatus.id, data }).unwrap();
        toast.success("Asset status updated");
      } else {
        await createAssetStatus(data).unwrap();
        toast.success("Asset status created");
      }
      onSuccess();
    } catch {
      toast.error("Failed to save asset status");
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
        <Button type="submit" disabled={isLoading}>{isLoading ? "Saving..." : assetStatus ? "Update" : "Create"}</Button>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}
