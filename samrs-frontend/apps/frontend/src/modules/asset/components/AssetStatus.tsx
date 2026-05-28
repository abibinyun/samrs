import type { AssetStatus } from "../types";

export const assetStatus = (status: AssetStatus) => {
  switch (status) {
    case "ready":
      return { label: "ready", className: "bg-green-100 text-green-700 border-green-200" };
    case "maintenance":
      return { label: "maintenance", className: "bg-amber-100 text-amber-700 border-amber-200" };
    case "broken":
      return { label: "broken", className: "bg-red-100 text-red-700 border-red-200" };
    default:
      return { label: status ?? "unknown", className: "bg-gray-100 text-gray-700 border-gray-200" };
  }
};
