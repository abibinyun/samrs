// src/modules/assets/workers/assets.worker.ts
import { generateMockAssets } from "./assets.mock";

self.onmessage = (e) => {
  const count = e.data ?? 10_000;
  const data = generateMockAssets(count);
  self.postMessage(data);
};
