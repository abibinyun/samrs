import type { StockOpnameStatus, StockOpnameCondition } from "../types";

export const stockOpnameStatus = (status: StockOpnameStatus) => {
  switch (status) {
    case "open": return { label: "Open", className: "bg-blue-100 text-blue-700 border-blue-200" };
    case "closed": return { label: "Closed", className: "bg-gray-100 text-gray-700 border-gray-200" };
    default: return { label: status ?? "unknown", className: "bg-gray-100 text-gray-700 border-gray-200" };
  }
};

export const stockOpnameCondition = (condition: StockOpnameCondition) => {
  switch (condition) {
    case "match": return { label: "Match", className: "bg-green-100 text-green-700 border-green-200" };
    case "missing": return { label: "Missing", className: "bg-red-100 text-red-700 border-red-200" };
    case "excess": return { label: "Excess", className: "bg-amber-100 text-amber-700 border-amber-200" };
    case "damaged": return { label: "Damaged", className: "bg-orange-100 text-orange-700 border-orange-200" };
    default: return { label: condition ?? "unknown", className: "bg-gray-100 text-gray-700 border-gray-200" };
  }
};
