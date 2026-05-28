import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { BrandTable } from "../layouts/BrandTable";
import { BrandForm } from "../layouts/BrandForm";
import { useGetBrandsQuery, useDeleteBrandMutation } from "../actions/brandApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { Brand } from "../types";
import { toast } from "sonner";

export function BrandListPage() {
  const { data, isLoading } = useGetBrandsQuery({});
  const [deleteBrand] = useDeleteBrandMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedBrand, setSelectedBrand] = useState<Brand | undefined>();

  const handleEdit = (brand: Brand) => {
    setSelectedBrand(brand);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Delete this brand?")) {
      try {
        await deleteBrand(id).unwrap();
        toast.success("Brand deleted");
      } catch {
        toast.error("Failed to delete brand");
      }
    }
  };

  return (
    <>
      <PageContainer header={{ title: "Brands", description: "Manage asset brands", actions: <Button onClick={() => setIsDialogOpen(true)}><Plus className="w-4 h-4 mr-2" />Add Brand</Button> }}>
        <BrandTable data={data?.data || []} isLoading={isLoading} onEdit={handleEdit} onDelete={handleDelete} />
      </PageContainer>
      <Dialog open={isDialogOpen} onOpenChange={(open) => { setIsDialogOpen(open); if (!open) setSelectedBrand(undefined); }}>
        <DialogContent>
          <DialogHeader><DialogTitle>{selectedBrand ? "Edit Brand" : "Add Brand"}</DialogTitle></DialogHeader>
          <BrandForm brand={selectedBrand} onSuccess={() => { setIsDialogOpen(false); setSelectedBrand(undefined); }} onCancel={() => setIsDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </>
  );
}
