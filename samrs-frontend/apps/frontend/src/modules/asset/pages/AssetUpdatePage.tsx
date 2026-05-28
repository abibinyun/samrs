import { useParams, useNavigate } from "@tanstack/react-router";
import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { useGetAssetByIdQuery } from "../actions/assetApiNew";
import { AssetUpdateForm } from "../layouts/AssetUpdateForm";

export function AssetUpdatePage() {
  const { id } = useParams({ strict: false });
  const navigate = useNavigate();
  const { data, isLoading, error } = useGetAssetByIdQuery(id!);

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

  return (
    <PageContainer
      header={{
        title: "Edit Asset",
        description: `Update asset: ${data.data.name}`,
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <AssetUpdateForm asset={data.data} onCancel={() => navigate({ to: `/assets/${id}` })} />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
