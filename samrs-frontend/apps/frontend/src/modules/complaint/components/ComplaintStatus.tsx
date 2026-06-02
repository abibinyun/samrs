import type { ComplaintStatus } from "../types";

export const complaintStatus = (status: ComplaintStatus) => {
  switch (status) {
    case "open":
      return { label: "Open", className: "bg-blue-100 text-blue-700 border-blue-200" };
    case "in_progress":
      return { label: "In Progress", className: "bg-amber-100 text-amber-700 border-amber-200" };
    case "done":
      return { label: "Done", className: "bg-green-100 text-green-700 border-green-200" };
    default:
      return { label: status ?? "unknown", className: "bg-gray-100 text-gray-700 border-gray-200" };
  }
};
