import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { bedSchema, type BedFormValues } from "../schemas";
import { useCreateBedMutation, useUpdateBedMutation } from "../actions/bedApi";
import type { Bed } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";

type BedFormProps = { bed?: Bed; onSuccess: () => void; onCancel: () => void };

export function BedForm({ bed, onSuccess, onCancel }: BedFormProps) {
  const [createBed, { isLoading: isCreating }] = useCreateBedMutation();
  const [updateBed, { isLoading: isUpdating }] = useUpdateBedMutation();
  const isLoading = isCreating || isUpdating;

  const { register, handleSubmit, formState: { errors }, setValue, watch } = useForm<BedFormValues>({
    resolver: zodResolver(bedSchema),
    defaultValues: bed ? { room_id: bed.room_id, name: bed.name, is_available: bed.is_available } : { is_available: true },
  });

  const onSubmit = async (data: BedFormValues) => {
    try {
      if (bed) {
        await updateBed({ id: bed.id, data }).unwrap();
        toast.success("Bed updated");
      } else {
        await createBed(data).unwrap();
        toast.success("Bed created");
      }
      onSuccess();
    } catch {
      toast.error("Failed to save bed");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="room_id">Room ID</Label>
        <Input id="room_id" {...register("room_id")} />
        {errors.room_id && <p className="text-sm text-red-600 mt-1">{errors.room_id.message}</p>}
      </div>
      <div>
        <Label htmlFor="name">Bed Name</Label>
        <Input id="name" {...register("name")} />
        {errors.name && <p className="text-sm text-red-600 mt-1">{errors.name.message}</p>}
      </div>
      <div className="flex items-center space-x-2">
        <Switch id="is_available" checked={watch("is_available")} onCheckedChange={(checked) => setValue("is_available", checked)} />
        <Label htmlFor="is_available">Available</Label>
      </div>
      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Saving..." : bed ? "Update" : "Create"}</Button>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}
