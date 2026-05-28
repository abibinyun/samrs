import { useState } from "react";
import PageContainer from "@/components/commons/containers/PageContainer";
import { RoomTable } from "../layouts/RoomTable";
import { RoomForm } from "../layouts/RoomForm";
import { useGetRoomsQuery, useDeleteRoomMutation } from "../actions/roomApi";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import type { Room } from "../types";
import { toast } from "sonner";

export function RoomListPage() {
  const { data, isLoading } = useGetRoomsQuery({});
  const [deleteRoom] = useDeleteRoomMutation();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedRoom, setSelectedRoom] = useState<Room | undefined>();

  const handleEdit = (room: Room) => {
    setSelectedRoom(room);
    setIsDialogOpen(true);
  };

  const handleDelete = async (id: string) => {
    if (confirm("Delete this room?")) {
      try {
        await deleteRoom(id).unwrap();
        toast.success("Room deleted");
      } catch {
        toast.error("Failed to delete room");
      }
    }
  };

  return (
    <>
      <PageContainer header={{ title: "Rooms", description: "Manage hospital rooms", actions: <Button onClick={() => setIsDialogOpen(true)}><Plus className="w-4 h-4 mr-2" />Add Room</Button> }}>
        <RoomTable data={data?.data || []} isLoading={isLoading} onEdit={handleEdit} onDelete={handleDelete} />
      </PageContainer>
      <Dialog open={isDialogOpen} onOpenChange={(open) => { setIsDialogOpen(open); if (!open) setSelectedRoom(undefined); }}>
        <DialogContent>
          <DialogHeader><DialogTitle>{selectedRoom ? "Edit Room" : "Add Room"}</DialogTitle></DialogHeader>
          <RoomForm room={selectedRoom} onSuccess={() => { setIsDialogOpen(false); setSelectedRoom(undefined); }} onCancel={() => setIsDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </>
  );
}
