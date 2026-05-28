import { lazy } from "react";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAssetCreate } from "../hooks/useAssetCreate"; 

const PageContainer = lazy(() => import('@/components/commons/containers/PageContainer'))
const AssetFormFields = lazy(() => import('../layouts/AssetFormFields'))

export default function AssetCreatePage() {
  const { form, onSubmit, isSubmitting, errors, canCreate, navigate } = useAssetCreate();

  return (
    <PageContainer
      // isSticky
      header={{
        title: "Add New Asset",
        description: "Register a new equipment asset to the system.",
        top: (
          <Button 
            variant="ghost" 
            size="sm" 
            className="pl-0 text-muted-foreground"
            onClick={() => navigate({ to: ".." })}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to List
          </Button>
        ),
        actions: canCreate && (
          <div className="flex gap-2">
             <Button variant="outline" onClick={() => navigate({ to: ".." })}>
               Cancel
             </Button>
             <Button form="asset-create-form" type="submit" disabled={isSubmitting}>
               {isSubmitting ? "Saving..." : "Save Asset"}
             </Button>
          </div>
        )
      }}
    >
      <form id="asset-create-form" onSubmit={onSubmit}>
        <AssetFormFields form={form} errors={errors} />
      </form>
    </PageContainer>
  );
}