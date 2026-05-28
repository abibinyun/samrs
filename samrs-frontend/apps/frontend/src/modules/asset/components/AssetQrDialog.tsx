import { Download } from "lucide-react";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

export type AssetQrPayload = { id: string; name: string };

type Props = {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  asset: AssetQrPayload | null;
};

export default function AssetQrDialog({ open, onOpenChange, asset }: Props) {
  const qrUrl = asset
    ? `${import.meta.env.VITE_API_URL}/api/v1/assets/${asset.id}/qr?size=256`
    : "";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>{asset?.name ?? "QR Code"}</DialogTitle>
          <DialogDescription>Scan QR untuk detail aset publik</DialogDescription>
        </DialogHeader>

        {asset ? (
          <div className="text-center space-y-4">
            <div className="bg-white p-4 border-2 border-dashed rounded-xl inline-block">
              <img src={qrUrl} alt="Asset QR" className="w-48 h-48 mx-auto" />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <Button variant="outline" onClick={() => onOpenChange(false)}>
                Tutup
              </Button>
              <Button className="gap-2" onClick={() => window.print()}>
                <Download className="w-4 h-4" />
                Cetak
              </Button>
            </div>
          </div>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
