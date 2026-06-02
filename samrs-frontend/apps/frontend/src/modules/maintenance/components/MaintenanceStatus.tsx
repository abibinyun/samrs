import type { MaintenanceStatus } from "../types";

export const maintenanceStatus = (status: MaintenanceStatus) => {
  switch (status) {
    case "scheduled":
      return { label: "Scheduled", className: "bg-blue-100 text-blue-700 border-blue-200" };
    case "due":
      return { label: "Due", className: "bg-amber-100 text-amber-700 border-amber-200" };
    case "completed":
      return { label: "Completed", className: "bg-green-100 text-green-700 border-green-200" };
    default:
      return { label: status ?? "unknown", className: "bg-gray-100 text-gray-700 border-gray-200" };
  }
};
