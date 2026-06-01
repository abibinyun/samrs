import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { assetSchema, type AssetFormValues } from "../schemas";
import { useUpdateAssetMutation } from "../actions/assetApiNew";
import type { Asset } from "../types";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import AssetFormFields from "./AssetFormFields";

type AssetUpdateFormProps = {
  asset: Asset;
  onCancel: () => void;
};

export function AssetUpdateForm({ asset, onCancel }: AssetUpdateFormProps) {
  const [updateAsset, { isLoading }] = useUpdateAssetMutation();
  
  const form = useForm<AssetFormValues>({
    resolver: zodResolver(assetSchema),
    defaultValues: {
      category_id: asset.category_id,
      room_id: asset.room_id || "",
      bed_id: asset.bed_id || 0,
      vendor_id: asset.vendor_id || 0,
      brand_id: asset.brand_id || 0,
      model_id: asset.model_id || 0,
      code: asset.code,
      name: asset.name,
      brand: asset.brand || "",
      model: asset.model || "",
      status: asset.status,
      purchase_date: asset.purchase_date ? asset.purchase_date.split("T")[0] : new Date().toISOString().split("T")[0],
    },
  });

  const onSubmit = async (data: AssetFormValues) => {
    try {
      await updateAsset({ id: asset.id, data }).unwrap();
      toast.success("Asset updated successfully");
      onCancel();
    } catch (error: any) {
      const message = error?.data?.message || "Failed to update asset";
      toast.error(message);
    }
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
      <AssetFormFields form={form} errors={form.formState.errors} />
      
      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>
          {isLoading ? "Updating..." : "Update Asset"}
        </Button>
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </form>
  );
}
