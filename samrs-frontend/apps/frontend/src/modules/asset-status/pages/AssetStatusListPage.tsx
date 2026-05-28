import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { AssetStatusTable } from "../layouts/AssetStatusTable";
import { AssetStatusForm } from "../layouts/AssetStatusForm";
import { useGetAssetStatusesQuery, useDeleteAssetStatusMutation } from "../actions/assetStatusApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { AssetStatus } from "../types";
import { toast } from "sonner";

export function AssetStatusListPage() {
  const { data, isLoading } = useGetAssetStatusesQuery({});
  const [deleteAssetStatus] = useDeleteAssetStatusMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedAssetStatus, setSelectedAssetStatus] = useState<AssetStatus | undefined>();

  const handleEdit = (assetStatus: AssetStatus) => {
    setSelectedAssetStatus(assetStatus);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Delete this asset status?")) {
      try {
        await deleteAssetStatus(id).unwrap();
        toast.success("Asset status deleted");
      } catch {
        toast.error("Failed to delete asset status");
      }
    }
  };

  return (
    <>
      <PageContainer header={{ title: "Asset Statuses", description: "Manage asset statuses", actions: <Button onClick={() => setIsDialogOpen(true)}><Plus className="w-4 h-4 mr-2" />Add Status</Button> }}>
        <AssetStatusTable data={data?.data || []} isLoading={isLoading} onEdit={handleEdit} onDelete={handleDelete} />
      </PageContainer>
      <Dialog open={isDialogOpen} onOpenChange={(open) => { setIsDialogOpen(open); if (!open) setSelectedAssetStatus(undefined); }}>
        <DialogContent>
          <DialogHeader><DialogTitle>{selectedAssetStatus ? "Edit Asset Status" : "Add Asset Status"}</DialogTitle></DialogHeader>
          <AssetStatusForm assetStatus={selectedAssetStatus} onSuccess={() => { setIsDialogOpen(false); setSelectedAssetStatus(undefined); }} onCancel={() => setIsDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </>
  );
}
