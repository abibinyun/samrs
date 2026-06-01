import { useNavigate } from "@tanstack/react-router";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useGetAssetByIdQuery } from "../actions/assetApiNew";
import { AssetUpdateForm } from "../layouts/AssetUpdateForm";

function parsePath() {
  const parts = window.location.pathname.split("/").filter(Boolean);
  const idx = parts.indexOf("assets");
  return {
    tenant: parts[0] || "",
    id: idx >= 0 ? parts[idx + 1] : undefined,
  };
}

export function AssetUpdatePage() {
  const navigate = useNavigate();
  const { tenant, id } = parsePath();
  const { data: asset, isLoading, error, refetch } = useGetAssetByIdQuery(id!, { skip: !id });

  if (isLoading) {
    return (
      <PageContainer>
        <div className="flex items-center justify-center h-64">
          <p>Loading...</p>
        </div>
      </PageContainer>
    );
  }

  if (error || !asset) {
    return (
      <PageContainer>
        <div className="flex flex-col items-center justify-center h-64 gap-4">
          <p className="text-destructive">Failed to load asset</p>
          <Button variant="outline" onClick={() => refetch()}>Retry</Button>
        </div>
      </PageContainer>
    );
  }

  return (
    <PageContainer
      header={{
        title: "Edit Asset",
        description: `Update asset: ${asset.name}`,
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <AssetUpdateForm asset={asset} onCancel={() => navigate({ to: `/${tenant}/assets/${id}` })} />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
