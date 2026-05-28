import PageContainer from "@/components/commons/containers/PageContainer";
import { Card, CardContent } from "@/components/ui/card";
import { MaintenanceForm } from "../layouts/MaintenanceForm";

export function MaintenanceCreatePage() {
  return (
    <PageContainer
      header={{
        title: "Buat Jadwal Pemeliharaan",
        description: "Jadwalkan pemeliharaan atau kalibrasi aset"
      }}
    >
      <Card>
        <CardContent className="pt-6">
          <MaintenanceForm />
        </CardContent>
      </Card>
    </PageContainer>
  );
}
