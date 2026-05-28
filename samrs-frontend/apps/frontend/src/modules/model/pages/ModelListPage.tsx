import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { ModelTable } from "../layouts/ModelTable";
import { ModelForm } from "../layouts/ModelForm";
import { useGetModelsQuery, useDeleteModelMutation } from "../actions/modelApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { Model } from "../types";
import { toast } from "sonner";

export function ModelListPage() {
  const { data, isLoading } = useGetModelsQuery({});
  const [deleteModel] = useDeleteModelMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedModel, setSelectedModel] = useState<Model | undefined>();

  const handleEdit = (model: Model) => {
    setSelectedModel(model);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Delete this model?")) {
      try {
        await deleteModel(id).unwrap();
        toast.success("Model deleted");
      } catch {
        toast.error("Failed to delete model");
      }
    }
  };

  return (
    <>
      <PageContainer header={{ title: "Models", description: "Manage asset models", actions: <Button onClick={() => setIsDialogOpen(true)}><Plus className="w-4 h-4 mr-2" />Add Model</Button> }}>
        <ModelTable data={data?.data || []} isLoading={isLoading} onEdit={handleEdit} onDelete={handleDelete} />
      </PageContainer>
      <Dialog open={isDialogOpen} onOpenChange={(open) => { setIsDialogOpen(open); if (!open) setSelectedModel(undefined); }}>
        <DialogContent>
          <DialogHeader><DialogTitle>{selectedModel ? "Edit Model" : "Add Model"}</DialogTitle></DialogHeader>
          <ModelForm model={selectedModel} onSuccess={() => { setIsDialogOpen(false); setSelectedModel(undefined); }} onCancel={() => setIsDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </>
  );
}
