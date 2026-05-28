import { useParams, useNavigate } from "@tanstack/react-router";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useGetAssetByIdQuery, useDeleteAssetMutation } from "../actions/assetApiNew";
import { ArrowLeft, Edit, Trash2, QrCode } from "lucide-react";
import { format } from "date-fns";
import { toast } from "sonner";
import { useState } from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export function AssetDetailPage() {
  const { id } = useParams({ strict: false });
  const navigate = useNavigate();
  const { data, isLoading, error } = useGetAssetByIdQuery(id!);
  const [deleteAsset, { isLoading: isDeleting }] = useDeleteAssetMutation();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  const handleDelete = async () => {
    try {
      await deleteAsset(id).unwrap();
      toast.success("Asset deleted successfully");
      navigate({ to: "/assets" });
    } catch (error: any) {
      const message = error?.data?.message || "Failed to delete asset";
      toast.error(message);
    }
  };

  if (isLoading) {
    return (
      <PageContainer>
        <div className="flex items-center justify-center h-64">
          <p>Loading...</p>
        </div>
      </PageContainer>
    );
  }

  if (error || !data) {
    return (
      <PageContainer>
        <div className="flex items-center justify-center h-64">
          <p className="text-red-600">Failed to load asset</p>
        </div>
      </PageContainer>
    );
  }

  const asset = data.data;

  const getStatusBadge = (status: string) => {
    const variants: Record<string, "default" | "secondary" | "destructive"> = {
      ready: "default",
      maintenance: "secondary",
      broken: "destructive",
    };
    return <Badge variant={variants[status] || "default"}>{status}</Badge>;
  };

  return (
    <PageContainer
      header={{
        title: asset.name,
        description: `Asset Code: ${asset.code}`,
        actions: (
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={() => navigate({ to: "/assets" })}>
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back
            </Button>
            <Button variant="outline" size="sm" onClick={() => toast.info("QR Code feature coming soon")}>
              <QrCode className="w-4 h-4 mr-2" />
              QR Code
            </Button>
            <Button variant="outline" size="sm" onClick={() => navigate({ to: `/assets/${id}/edit` })}>
              <Edit className="w-4 h-4 mr-2" />
              Edit
            </Button>
            <Button variant="destructive" size="sm" onClick={() => setShowDeleteDialog(true)} disabled={isDeleting}>
              <Trash2 className="w-4 h-4 mr-2" />
              Delete
            </Button>
          </div>
        ),
      }}
    >
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Code</p>
              <p className="font-medium">{asset.code}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Name</p>
              <p className="font-medium">{asset.name}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Status</p>
              {getStatusBadge(asset.status)}
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Category</p>
              <p className="font-medium">{asset.category?.name || "-"}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Location & Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Room</p>
              <p className="font-medium">{asset.room?.name || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Bed</p>
              <p className="font-medium">{asset.bed?.name || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Brand</p>
              <p className="font-medium">{asset.brand || asset.brand_master?.name || "-"}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Model</p>
              <p className="font-medium">{asset.model || asset.model_master?.name || "-"}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Purchase Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Purchase Date</p>
              <p className="font-medium">
                {asset.purchase_date ? format(new Date(asset.purchase_date), "dd MMM yyyy") : "-"}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Vendor</p>
              <p className="font-medium">{asset.vendor?.name || "-"}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm text-muted-foreground">Created At</p>
              <p className="font-medium">{format(new Date(asset.created_at), "dd MMM yyyy HH:mm")}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Last Updated</p>
              <p className="font-medium">{format(new Date(asset.updated_at), "dd MMM yyyy HH:mm")}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete the asset "{asset.name}" (Code: {asset.code}). This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </PageContainer>
  );
}
