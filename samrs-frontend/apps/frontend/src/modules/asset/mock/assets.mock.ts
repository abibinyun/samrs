// src/modules/assets/mocks/assets.mock.ts
import type { AssetDevice, AssetStatus } from "../types";

/* ================= CONSTANT SOURCES ================= */

const STATUSES: AssetStatus[] = ["ready", "maintenance", "broken"];

const BRANDS = [
  { id: 1, name: "Philips" },
  { id: 2, name: "GE Healthcare" },
  { id: 3, name: "Siemens" },
  { id: 4, name: "Mindray" },
];

const MODELS = [
  { id: 1, brand_id: 1, name: "V60" },
  { id: 2, brand_id: 1, name: "IntelliVue MX800" },
  { id: 3, brand_id: 2, name: "LOGIQ E10" },
  { id: 4, brand_id: 3, name: "ACUSON Sequoia" },
  { id: 5, brand_id: 4, name: "SV300" },
];

const ROOMS = [
  { id: "IGD-1", name: "IGD Utama" },
  { id: "ICU-1", name: "ICU Barat" },
  { id: "ICU-2", name: "ICU Timur" },
  { id: "OK-1", name: "Kamar Operasi 1" },
];

/* ================= HELPERS ================= */

const randomFrom = <T>(arr: T[]) =>
  arr[Math.floor(Math.random() * arr.length)];

const randomDate = (startYear = 2020) => {
  const year = startYear + Math.floor(Math.random() * 5);
  const month = String(Math.floor(Math.random() * 12) + 1).padStart(2, "0");
  const day = String(Math.floor(Math.random() * 28) + 1).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

/* ================= GENERATOR (PENTING) ================= */

export function generateMockAssets(
  count = 100_000
): AssetDevice[] {
  return Array.from({ length: count }, (_, i) => {
    const brand = randomFrom(BRANDS);
    const model = randomFrom(
      MODELS.filter((m) => m.brand_id === brand.id)
    );
    const room = randomFrom(ROOMS);

    return {
      id: String(i + 1),

      // relations (ids)
      category_id: i % 2 === 0 ? 1 : 2, // 1, 2, 3 , dst number
      room_id: room.id,
      bed_id: i % 4 === 0 ? null : (i % 20) + 1,
      vendor_id: (i % 5) + 1,
      brand_id: brand.id,
      model_id: model.id,

      // display fields
      code: `ASSET-${String(i + 1).padStart(5, "0")}`,
      name: `Medical Device ${i + 1}`, // Adul, RoZak, LamPu, Meja,
      brand: brand.name,
      model: model.name,

      status: STATUSES[i % STATUSES.length],
      purchase_date: randomDate(2020),

      // _search: `${a.name} ${a.code}`.toLowerCase(), // adul 123
      // _category: a.category_id.toString(), // "1" string
    };
  });
}
