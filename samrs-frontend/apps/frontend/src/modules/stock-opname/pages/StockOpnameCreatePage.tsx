import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { StockOpnameForm } from "../layouts/StockOpnameForm";

export function StockOpnameCreatePage() {
  return (
    <PageContainer 
      header={{
        title: "Buat Sesi Stock Opname",
        description: "Buat sesi baru untuk stock opname aset"
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <StockOpnameForm />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
