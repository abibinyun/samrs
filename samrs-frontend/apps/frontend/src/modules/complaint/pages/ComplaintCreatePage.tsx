import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { ComplaintForm } from "../layouts/ComplaintForm";

export function ComplaintCreatePage() {
  return (
    <PageContainer
      header={{
        title: "Buat Keluhan Baru",
        description: "Laporkan kerusakan atau masalah pada aset"
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <ComplaintForm />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
