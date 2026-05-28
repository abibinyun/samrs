import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { BedTable } from "../layouts/BedTable";
import { BedForm } from "../layouts/BedForm";
import { useGetBedsQuery, useDeleteBedMutation } from "../actions/bedApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { Bed } from "../types";
import { toast } from "sonner";

export function BedListPage() {
  const { data, isLoading } = useGetBedsQuery({});
  const [deleteBed] = useDeleteBedMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedBed, setSelectedBed] = useState<Bed | undefined>();

  const handleEdit = (bed: Bed) => {
    setSelectedBed(bed);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Delete this bed?")) {
      try {
        await deleteBed(id).unwrap();
        toast.success("Bed deleted");
      } catch {
        toast.error("Failed to delete bed");
      }
    }
  };

  return (
    <>
      <PageContainer header={{ title: "Beds", description: "Manage hospital beds", actions: <Button onClick={() => setIsDialogOpen(true)}><Plus className="w-4 h-4 mr-2" />Add Bed</Button> }}>
        <BedTable data={data?.data || []} isLoading={isLoading} onEdit={handleEdit} onDelete={handleDelete} />
      </PageContainer>
      <Dialog open={isDialogOpen} onOpenChange={(open) => { setIsDialogOpen(open); if (!open) setSelectedBed(undefined); }}>
        <DialogContent>
          <DialogHeader><DialogTitle>{selectedBed ? "Edit Bed" : "Add Bed"}</DialogTitle></DialogHeader>
          <BedForm bed={selectedBed} onSuccess={() => { setIsDialogOpen(false); setSelectedBed(undefined); }} onCancel={() => setIsDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </>
  );
}
