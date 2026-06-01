import { type UseFormReturn } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { type AssetFormValues } from "../schemas";
import { useGetRoomsQuery } from "@/modules/room/actions/roomApi";
import { useGetCategoriesQuery } from "@/modules/category/actions/categoryApi";

type Props = {
  form: UseFormReturn<AssetFormValues>;
  errors: Record<string, { message?: string }>;
};

export default function AssetFormFields({ form, errors }: Props) {
  const { data: roomsData, isLoading: isLoadingRooms, isError: isErrorRooms } = useGetRoomsQuery({ limit: 100 });
  const { data: categoriesData, isLoading: isLoadingCategories, isError: isErrorCategories } = useGetCategoriesQuery({ limit: 100 });
  
  const rooms = roomsData?.data || [];
  const categories = categoriesData?.data || [];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 bg-card p-6 border rounded-xl shadow-sm">
      {/* General Information */}
      <div className="col-span-full border-b pb-2 mb-2">
        <h3 className="font-semibold text-lg">General Information</h3>
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">Asset Name</Label>
        <Input id="name" {...form.register("name")} placeholder="e.g. Ventilator Philips V60" />
        {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="code">Asset Code</Label>
        <Input id="code" {...form.register("code")} placeholder="VENT-003-4" />
        {errors.code && <p className="text-xs text-destructive">{errors.code.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="brand">Brand Name</Label>
        <Input id="brand" {...form.register("brand")} placeholder="e.g. Philips" />
        {errors.brand && <p className="text-xs text-destructive">{errors.brand.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="model">Model Name</Label>
        <Input id="model" {...form.register("model")} placeholder="e.g. V60" />
        {errors.model && <p className="text-xs text-destructive">{errors.model.message}</p>}
      </div>

      {/* Logistics Information */}
      <div className="col-span-full border-b pb-2 mt-4 mb-2">
        <h3 className="font-semibold text-lg">Placement & Logistics</h3>
      </div>

      <div className="space-y-2">
        <Label>Room</Label>
        <Select onValueChange={(v) => form.setValue("room_id", v, { shouldValidate: true })} disabled={isLoadingRooms || isErrorRooms}>
          <SelectTrigger>
            <SelectValue placeholder={isLoadingRooms ? "Loading..." : isErrorRooms ? "Failed to load rooms" : "Select Room"} />
          </SelectTrigger>
          <SelectContent>
            {rooms.map((room) => (
              <SelectItem key={room.id} value={room.id}>
                {room.name} ({room.code})
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.room_id && <p className="text-xs text-destructive">{errors.room_id.message}</p>}
      </div>

      <div className="space-y-2">
        <Label>Category</Label>
        <Select 
          onValueChange={(v) => form.setValue("category_id", Number(v), { shouldValidate: true })}
          disabled={isLoadingCategories || isErrorCategories}
        >
          <SelectTrigger>
            <SelectValue placeholder={isLoadingCategories ? "Loading..." : isErrorCategories ? "Failed to load categories" : "Select Category"} />
          </SelectTrigger>
          <SelectContent>
            {categories.map((cat) => (
              <SelectItem key={cat.id} value={String(cat.id)}>
                {cat.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.category_id && <p className="text-xs text-destructive">{errors.category_id.message}</p>}
      </div>

      <div className="space-y-2">
        <Label>Status</Label>
        <Select defaultValue="ready" onValueChange={(v: "ready" | "maintenance" | "broken") => form.setValue("status", v)}>
          <SelectTrigger>
            <SelectValue placeholder="Select Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="ready">Ready</SelectItem>
            <SelectItem value="maintenance">Maintenance</SelectItem>
            <SelectItem value="broken">Broken</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label>Purchase Date</Label>
        <Input type="date" {...form.register("purchase_date")} />
      </div>

      <div className="space-y-2">
        <Label>Bed ID</Label>
        <Input type="number" {...form.register("bed_id", { valueAsNumber: true })} placeholder="1" />
      </div>
    </div>
  );
}