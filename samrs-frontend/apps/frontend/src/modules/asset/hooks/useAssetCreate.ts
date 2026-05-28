import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { assetSchema, type AssetFormValues } from "../schemas";
import { usePermission } from "@/hooks/usePermission";
import { SCOPES } from "@/constants";
import { useCreateAssetMutation } from "../actions/assetApiNew";

export function useAssetCreate() {
  const navigate = useNavigate();
  const { hasPermission } = usePermission();
  const [createAsset] = useCreateAssetMutation();

  const form = useForm<AssetFormValues>({
    resolver: zodResolver(assetSchema) as any,
    defaultValues: {
      bed_id: 0,
      brand_id: 0,
      model_id: 0,
      status: "ready",
      name: "",
      code: "",
      brand: "",
      model: "",
      room_id: "",
      purchase_date: new Date().toISOString().split("T")[0],
    },
  });

  const onSubmit: SubmitHandler<AssetFormValues> = async (data) => {
    try {
      await createAsset(data).unwrap();
      toast.success("Asset created successfully");
      navigate({ to: "/assets" });
    } catch (error: any) {
      const message = error?.data?.message || "Failed to create asset";
      toast.error(message);
    }
  };

  const onErrors = (err: any) => {
    console.log("Form Errors:", err);
    toast.error("Please check the form for errors");
  };

  return {
    form,
    onSubmit: form.handleSubmit(onSubmit, onErrors),
    isSubmitting: form.formState.isSubmitting,
    errors: form.formState.errors,
    canCreate: hasPermission(SCOPES.ASSET.CREATE),
    navigate,
  };
}