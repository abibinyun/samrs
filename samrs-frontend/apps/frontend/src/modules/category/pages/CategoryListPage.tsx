import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { CategoryTable } from "../layouts/CategoryTable";
import { CategoryForm } from "../layouts/CategoryForm";
import { useGetCategoriesQuery, useDeleteCategoryMutation } from "../actions/categoryApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { Category } from "../types";
import { toast } from "sonner";

export function CategoryListPage() {
  const { data, isLoading } = useGetCategoriesQuery({});
  const [deleteCategory] = useDeleteCategoryMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | undefined>();

  const handleEdit = (category: Category) => {
    setSelectedCategory(category);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (confirm("Are you sure you want to delete this category?")) {
      try {
        await deleteCategory(id).unwrap();
        toast.success("Category deleted successfully");
      } catch (error) {
        toast.error("Failed to delete category");
      }
    }
  };

  const handleSuccess = () => {
    setIsDialogOpen(false);
    setSelectedCategory(undefined);
  };

  return (
    <>
      <PageContainer
        header={{
          title: "Categories",
          description: "Manage asset categories",
          actions: (
            <Button onClick={() => setIsDialogOpen(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Add Category
            </Button>
          ),
        }}
      >
        <CategoryTable
          data={data?.data || []}
          isLoading={isLoading}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      </PageContainer>

      <Dialog open={isDialogOpen} onOpenChange={(open) => {
        setIsDialogOpen(open);
        if (!open) setSelectedCategory(undefined);
      }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{selectedCategory ? "Edit Category" : "Add Category"}</DialogTitle>
          </DialogHeader>
          <CategoryForm
            category={selectedCategory}
            onSuccess={handleSuccess}
            onCancel={() => setIsDialogOpen(false)}
          />
        </DialogContent>
      </Dialog>
    </>
  );
}
