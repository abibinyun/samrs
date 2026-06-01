import PageContainer from "@/components/commons/containers/PageContainer";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ArrowLeft, Plus, Trash2 } from "lucide-react";
import { useNavigate, useRouterState, useParams } from "@tanstack/react-router";
import {
  useGetStockOpnameSessionQuery,
  useGetStockOpnameItemsQuery,
  useAddStockOpnameItemMutation,
  useDeleteStockOpnameItemMutation,
  useCloseStockOpnameSessionMutation,
} from "../actions/stockOpnameApi";
import { stockOpnameStatus, stockOpnameCondition } from "../components/StockOpnameStatus";
import { toast } from "sonner";
import { useState } from "react";
import type { StockOpnameCondition } from "../types";

export function StockOpnameDetailPage() {
  const { sessionId } = useParams({ strict: false });
  const id = Number(sessionId);
  const navigate = useNavigate();
  const basePath = useRouterState({
    select: (s) => {
      const parts = s.location.pathname.split("/").filter(Boolean);
      return `/${parts[0]}`;
    },
  });
  const { data: session, isLoading } = useGetStockOpnameSessionQuery(id);
  const { data: items, isLoading: itemsLoading } = useGetStockOpnameItemsQuery(id);
  const [addItem] = useAddStockOpnameItemMutation();
  const [deleteItem] = useDeleteStockOpnameItemMutation();
  const [closeSession] = useCloseStockOpnameSessionMutation();

  const [assetId, setAssetId] = useState("");
  const [condition, setCondition] = useState<StockOpnameCondition>("match");
  const [note, setNote] = useState("");

  const handleAddItem = async () => {
    if (!assetId) return;
    try {
      await addItem({ sessionId: id, data: { asset_id: assetId, condition, note: note || undefined } }).unwrap();
      setAssetId("");
      setCondition("match");
      setNote("");
    } catch {
      toast.error("Gagal menambah item");
    }
  };

  const handleDeleteItem = async (itemId: number) => {
    if (confirm("Hapus item ini?")) {
      try {
        await deleteItem({ sessionId: id, itemId }).unwrap();
      } catch {
        toast.error("Gagal menghapus item");
      }
    }
  };

  const handleClose = async () => {
    if (confirm("Tutup sesi ini? Setelah ditutup, tidak bisa diubah lagi.")) {
      try {
        await closeSession(id).unwrap();
      } catch {
        toast.error("Gagal menutup sesi");
      }
    }
  };

  if (isLoading) return <div className="p-8 text-center text-muted-foreground">Loading...</div>;

  const data = session?.data;
  if (!data) return <div className="p-8 text-center text-muted-foreground">Sesi tidak ditemukan</div>;

  const statusInfo = stockOpnameStatus(data.status);

  return (
    <PageContainer
      header={{
        title: data.title,
        description: `Detail sesi stock opname #${data.id}`,
        actions: (
          <Button variant="outline" onClick={() => navigate({ to: `${basePath}/stock-opname` })}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Kembali
          </Button>
        ),
      }}
    >
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-muted-foreground mb-1">Status</div>
              <Badge variant="outline" className={statusInfo.className}>{statusInfo.label}</Badge>
            </div>
            <div className="flex gap-2">
              {data.status === "draft" && (
                <>
                  <Button variant="outline" size="sm" onClick={() => navigate({ to: `${basePath}/stock-opname/${data.id}/edit` })}>Edit</Button>
                  <Button size="sm" onClick={handleClose}>Tutup Sesi</Button>
                </>
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm font-medium text-muted-foreground">Tanggal</div>
              <div>{new Date(data.opname_at).toLocaleDateString("id-ID")}</div>
            </div>
            <div>
              <div className="text-sm font-medium text-muted-foreground">Dibuat</div>
              <div>{new Date(data.created_at).toLocaleString("id-ID")}</div>
            </div>
          </div>
          {data.notes && (
            <div>
              <div className="text-sm font-medium text-muted-foreground">Catatan</div>
              <div>{data.notes}</div>
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Daftar Item</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {data.status === "draft" && (
            <div className="flex gap-2 items-end p-4 bg-muted/50 rounded-lg">
              <Input placeholder="Asset ID" value={assetId} onChange={(e) => setAssetId(e.target.value)} className="w-48" />
              <Select value={condition} onValueChange={(v) => setCondition(v as StockOpnameCondition)}>
                <SelectTrigger className="w-36"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="match">Sesuai</SelectItem>
                  <SelectItem value="missing">Hilang</SelectItem>
                  <SelectItem value="damaged">Rusak</SelectItem>
                </SelectContent>
              </Select>
              <Input placeholder="Catatan" value={note} onChange={(e) => setNote(e.target.value)} className="flex-1" />
              <Button onClick={handleAddItem} disabled={!assetId}>
                <Plus className="w-4 h-4 mr-1" /> Tambah
              </Button>
            </div>
          )}

          {itemsLoading ? (
            <div className="text-center text-muted-foreground py-4">Loading...</div>
          ) : !items?.data?.length ? (
            <div className="text-center text-muted-foreground py-8">Belum ada item</div>
          ) : (
            <div className="border rounded-lg">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b bg-muted/50">
                    <th className="p-2 text-left">Kode</th>
                    <th className="p-2 text-left">Nama Aset</th>
                    <th className="p-2 text-left">Kondisi</th>
                    <th className="p-2 text-left">Catatan</th>
                    <th className="p-2 text-left">Diperiksa Oleh</th>
                    {data.status === "draft" && <th className="p-2 w-10"></th>}
                  </tr>
                </thead>
                <tbody>
                  {items.data.map((item) => {
                    const cond = stockOpnameCondition(item.condition);
                    return (
                      <tr key={item.id} className="border-b last:border-b-0">
                        <td className="p-2">{item.asset?.code || "-"}</td>
                        <td className="p-2">{item.asset?.name || "-"}</td>
                        <td className="p-2">
                          <Badge variant="outline" className={cond.className}>{cond.label}</Badge>
                        </td>
                        <td className="p-2">{item.note || "-"}</td>
                        <td className="p-2">{item.checker?.username || "-"}</td>
                        {data.status === "draft" && (
                          <td className="p-2">
                            <Button variant="ghost" size="sm" onClick={() => handleDeleteItem(item.id)}>
                              <Trash2 className="w-3 h-3 text-destructive" />
                            </Button>
                          </td>
                        )}
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </PageContainer>
  );
}

// toast is imported from sonner
