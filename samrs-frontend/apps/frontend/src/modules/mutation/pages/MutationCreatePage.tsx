import { useState } from "react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { useCreateMutationMutation } from "../actions/mutationApi";
import { useGetAssetsQuery } from "@/modules/asset/actions/assetApiNew";
import { useGetRoomsQuery } from "@/modules/room/actions/roomApi";
import { SCOPES } from "@/constants/permissions";
import { usePermission } from "@/hooks/usePermission";
import PageContainer from "@/components/commons/containers/PageContainer";

export default function MutationCreatePage() {
  const navigate = useNavigate();
  const { hasPermission } = usePermission();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const [createMutation, { isLoading }] = useCreateMutationMutation();

  const { data: assetsData, isLoading: isLoadingAssets } = useGetAssetsQuery({ limit: 100 });
  const { data: roomsData, isLoading: isLoadingRooms } = useGetRoomsQuery({ limit: 100 });

  const [form, setForm] = useState({
    asset_id: "",
    room_id: "",
    bed_id: "",
    reason: "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const assets = assetsData?.data || [];
  const rooms = roomsData?.data || [];

  const validate = () => {
    const newErrors: Record<string, string> = {};
    if (!form.asset_id) newErrors.asset_id = "Asset harus dipilih";
    if (!form.room_id) newErrors.room_id = "Ruangan tujuan harus dipilih";
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    try {
      await createMutation({
        asset_id: form.asset_id,
        room_id: form.room_id || undefined,
        bed_id: form.bed_id ? Number(form.bed_id) : undefined,
        reason: form.reason || undefined,
      }).unwrap();

      toast.success("Mutasi aset berhasil");
      navigate({ to: `${basePath}/assets/mutation` });
    } catch (error: any) {
      const message = error?.data?.message || "Gagal melakukan mutasi";
      toast.error(message);
    }
  };

  if (!hasPermission(SCOPES.ASSET_MUTATION.CREATE)) {
    return (
      <PageContainer>
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Anda tidak memiliki akses untuk membuat mutasi</p>
        </div>
      </PageContainer>
    );
  }

  return (
    <PageContainer
      header={{
        title: "Tambah Mutasi",
        description: "Pindahkan aset ke ruangan lain.",
        top: (
          <Button
            variant="ghost"
            size="sm"
            className="pl-0 text-muted-foreground"
            onClick={() => navigate({ to: `${basePath}/assets/mutation` })}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Kembali
          </Button>
        ),
        actions: (
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => navigate({ to: `${basePath}/assets/mutation` })}>
              Batal
            </Button>
            <Button form="mutation-form" type="submit" disabled={isLoading}>
              {isLoading ? "Menyimpan..." : "Simpan Mutasi"}
            </Button>
          </div>
        ),
      }}
    >
      <form id="mutation-form" onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Detail Mutasi</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="asset_id">Aset *</Label>
              <Select
                value={form.asset_id}
                onValueChange={(v) => setForm({ ...form, asset_id: v })}
                disabled={isLoadingAssets}
              >
                <SelectTrigger>
                  <SelectValue placeholder={isLoadingAssets ? "Memuat aset..." : "Pilih aset"} />
                </SelectTrigger>
                <SelectContent>
                  {assets.map((asset) => (
                    <SelectItem key={asset.id} value={asset.id}>
                      {asset.code} - {asset.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.asset_id && <p className="text-xs text-destructive">{errors.asset_id}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="room_id">Ruangan Tujuan *</Label>
              <Select
                value={form.room_id}
                onValueChange={(v) => setForm({ ...form, room_id: v })}
                disabled={isLoadingRooms}
              >
                <SelectTrigger>
                  <SelectValue placeholder={isLoadingRooms ? "Memuat ruangan..." : "Pilih ruangan tujuan"} />
                </SelectTrigger>
                <SelectContent>
                  {rooms.map((room) => (
                    <SelectItem key={room.id} value={room.id}>
                      {room.code} - {room.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {errors.room_id && <p className="text-xs text-destructive">{errors.room_id}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="bed_id">Tempat Tidur (Opsional)</Label>
              <Input
                id="bed_id"
                type="number"
                value={form.bed_id}
                onChange={(e) => setForm({ ...form, bed_id: e.target.value })}
                placeholder="Nomor tempat tidur"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="reason">Alasan Mutasi</Label>
              <Textarea
                id="reason"
                value={form.reason}
                onChange={(e) => setForm({ ...form, reason: e.target.value })}
                placeholder="Contoh: Pindah ke instalasi radiologi"
                rows={3}
              />
            </div>
          </CardContent>
        </Card>
      </form>
    </PageContainer>
  );
}
