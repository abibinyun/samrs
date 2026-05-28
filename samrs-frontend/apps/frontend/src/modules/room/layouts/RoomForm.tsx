import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { roomSchema, type RoomFormValues } from "../schemas";
import { useCreateRoomMutation, useUpdateRoomMutation } from "../actions/roomApi";
import type { Room } from "../types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

type RoomFormProps = { room?: Room; onSuccess: () => void; onCancel: () => void };

export function RoomForm({ room, onSuccess, onCancel }: RoomFormProps) {
  const [createRoom, { isLoading: isCreating }] = useCreateRoomMutation();
  const [updateRoom, { isLoading: isUpdating }] = useUpdateRoomMutation();
  const isLoading = isCreating || isUpdating;

  const { register, handleSubmit, formState: { errors } } = useForm<RoomFormValues>({
    resolver: zodResolver(roomSchema),
    defaultValues: room || undefined,
  });

  const onSubmit = async (data: RoomFormValues) => {
    try {
      if (room) {
        await updateRoom({ id: room.id, data }).unwrap();
        toast.success("Room updated");
      } else {
        await createRoom(data).unwrap();
        toast.success("Room created");
      }
      onSuccess();
    } catch {
      toast.error("Failed to save room");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <Label htmlFor="code">Room Code</Label>
        <Input id="code" {...register("code")} placeholder="e.g. R001, VIP-01" />
        {errors.code && <p className="text-sm text-red-600 mt-1">{errors.code.message}</p>}
      </div>
      <div>
        <Label htmlFor="name">Room Name</Label>
        <Input id="name" {...register("name")} />
        {errors.name && <p className="text-sm text-red-600 mt-1">{errors.name.message}</p>}
      </div>
      <div>
        <Label htmlFor="location">Location</Label>
        <Input id="location" {...register("location")} placeholder="e.g. Building A Floor 2" />
      </div>
      <div className="flex gap-2">
        <Button type="submit" disabled={isLoading}>{isLoading ? "Saving..." : room ? "Update" : "Create"}</Button>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}
