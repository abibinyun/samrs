# SAMRS — Bug Fix Notes & Gotchas

> Dokumentasi bug yang sulit di-fix agar tidak terulang. Setiap bagian berisi: **gejala**, **akar masalah**, **solusi**, dan **cara deteksi**.

---

## 1. `lazyRouteComponent` — Halaman Stuck di Loading Spinner

### Gejala
- URL sudah benar (`/rs-pusat/assets/{id}`)
- Halaman menampilkan spinner tak terbatas dari `__root.tsx` Suspense
- **Tidak ada request masuk ke backend** (dicek di backend logs)
- Tidak ada error di browser console maupun Vite logs

### Akar Masalah
`lazyRouteComponent()` dari TanStack Router melakukan lazy import. Jika import gagal secara **silent** (runtime error saat module evaluation, circular dependency, atau Vite HMR cache corrupt), Suspense tidak menampilkan error — hanya stuck di fallback spinner.

### Solusi
Import component secara langsung (eager) di route definition:

```tsx
// ❌ Bisa gagal silent
import { AssetDetailPage } from "./pages/AssetDetailPage";

createRoute({
  component: lazyRouteComponent(() => import("./pages/AssetDetailPage")),
  ...
})

// ✅ Aman
import { AssetDetailPage } from "./pages/AssetDetailPage";

createRoute({
  component: AssetDetailPage,
  ...
})
```

### Cara Deteksi
1. Buka Network tab di DevTools — jika **tidak ada request** ke API, component tidak render
2. Cek `docker logs samrs_backend --tail 20` — jika tidak ada request dari browser (hanya curl), masalah di frontend
3. Restart container frontend: `docker compose restart frontend`

### Kapan pakai `lazyRouteComponent`
Hanya untuk route yang **sudah terbukti works** (misal list page). Untuk route baru/debugging, pakai import langsung dulu.

---

## 2. Vite Proxy Tidak Work di Docker

### Gejala
- Request dari browser ke `/api/v1/assets/{id}` tidak sampai ke backend
- `curl http://localhost:5173/api/v1/assets` dari host juga gagal
- Tapi `curl http://localhost:8090/api/v1/assets` (langsung ke backend) works

### Akar Masalah
Vite dev server berjalan di dalam Docker container. `vite.config.ts` punya proxy config:
```ts
proxy: {
  "/api": {
    target: process.env.VITE_PROXY_TARGET || "http://localhost:3000",
    // ...
  }
}
```
Default `http://localhost:3000` — ini localhost **di dalam container**, bukan di host. BFF container bisa diakses lewat `http://bff:3000` via Docker network.

### Solusi
1. Set env var di `docker-compose.dev.yml`:
```yaml
environment:
  VITE_PROXY_TARGET: http://bff:3000
```
2. Di `vite.config.ts` gunakan env var:
```ts
target: process.env.VITE_PROXY_TARGET || "http://localhost:3000"
```
3. **Rebuild/restart container** setelah perubahan (env var dibaca saat Vite startup, bukan hot-reload)

### Cara Deteksi
```bash
# Cek env var di container
docker exec samrs_frontend printenv VITE_PROXY_TARGET

# Test proxy dari host
curl -s http://localhost:5173/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"password123"}'
```

---

## 3. Radix Select Tidak Menampilkan Value Saat Edit

### Gejala
- Form edit asset: dropdown Room, Category, Status **kosong** meski data sudah ada
- Input text (name, code) terisi dengan benar
- Data tersimpan di form state (bisa di-submit), tapi UI tidak menampilkan selection

### Akar Masalah
Radix UI `Select` adalah **uncontrolled component** kecuali diberi prop `value`. Tanpa `value`, Select tidak tahu `SelectItem` mana yang harus ditampilkan.

```tsx
// ❌ Tidak akan menampilkan selection yang sudah ada
<Select onValueChange={...}>

// ✅ Controlled — tampilkan value dari form state
<Select value={form.watch("room_id") || ""} onValueChange={...}>
```

### Tambahkan: Type Mismatch (number vs string)
`category_id` di form state bertipe `number`, tapi `SelectItem` menggunakan `String(cat.id)`. Radix Select hanya cocokkan **string**.

```tsx
// ❌ Number vs String tidak cocok
<Select value={form.watch("category_id")} ...>
<SelectItem value={cat.id}>  // cat.id is number

// ✅ Konversi ke string
<Select value={form.watch("category_id") ? String(form.watch("category_id")) : ""} ...>
<SelectItem value={String(cat.id)}>
```

### Pattern untuk Setiap Select di Edit Form
```tsx
<Select
  value={form.watch("field_name") || ""}   // controlled + fallback string
  onValueChange={(v) => form.setValue("field_name", v)}
>
  <SelectTrigger>
    <SelectValue placeholder="Select..." />
  </SelectTrigger>
  <SelectContent>
    {items.map((item) => (
      <SelectItem key={item.id} value={String(item.id)}>
        {item.name}
      </SelectItem>
    ))}
  </SelectContent>
</Select>
```

---

## 4. `<input type="date">` Tidak Terisi dari API Response

### Gejala
- `purchase_date` dari API: `"2024-01-15T00:00:00Z"` atau `"2024-01-15"`
- `<Input type="date">` tetap kosong meski value sudah di-set

### Akar Masalah
HTML `<input type="date">` hanya menerima format `YYYY-MM-DD`. Jika value adalah ISO datetime (`2024-01-15T00:00:00Z`), browser mengabaikannya.

### Solusi
Normalisasi di `defaultValues`:
```ts
purchase_date: asset.purchase_date
  ? asset.purchase_date.split("T")[0]   // "2024-01-15T00:00:00Z" → "2024-01-15"
  : new Date().toISOString().split("T")[0],
```

---

## 5. Navigasi dari TanStack Table Cell Salah Path

### Gejala
- Klik "Detail" di table column → navigasi ke `/rs-pusat` (root tenant) bukan `/rs-pusat/assets/{id}`
- URL benar di code, tapi browser tidak navigate ke sana

### Akar Masalah
Relative path (`./$id`) di `navigate()` dari dalam TanStack Table cell callback tidak resolve dengan benar karena context render-nya berbeda dari route context.

### Solusi
Gunakan **absolute path** dengan tenant prefix:
```tsx
// ❌ Relative path — tidak resolve dengan benar dari table cell
navigate({ to: `./${a.id}` });

// ✅ Absolute path dengan tenant
navigate({ to: `/${tenant}/assets/${a.id}` });
```

Untuk mendapatkan tenant, pakai `useRouterState`:
```tsx
const basePath = useRouterState({
  select: (s) => {
    const parts = s.location.pathname.split("/").filter(Boolean);
    return `/${parts[0]}`;
  },
});
navigate({ to: `${basePath}/assets/${a.id}` });
```

---

## 6. API Response Bungkus Ganda (Double-Wrap)

### Gejala
- `data?.data` mengembalikan `undefined`
- List page tidak menampilkan data meski API return 200

### Akar Masalah
Backend mengembalikan:
```json
{ "success": true, "data": [...], "meta": { "total": 5 } }
```
RTK Query `transformResponse` harus unwrap sekali, lalu component akses `.data` lagi untuk array.

### Solusi
```ts
// RTK Query
transformResponse: (response: ApiResponse<Asset[]>) => ({
  data: response.data,
  total: response.meta?.total || 0,
}),

// Component
const rooms = roomsData?.data || [];  // roomsData sudah di-unwrap oleh transformResponse
```

---

## 7. `useParams({ strict: false })` Tidak Resolve Params

### Gejala
- `useParams()` dari TanStack Router mengembalikan `{}` (kosong)
- `useParams({ from: "/$tenantId/assets/$id" })` throws error

### Akar Masalah
`useParams` di TanStack Router v1 hanya mengembalikan params dari **route yang match**. Jika route menggunakan `lazyRouteComponent` dan komponen belum mount, params mungkin belum tersedia. Selain itu, API `from` parameter belum tentu kompatibel dengan semua versi.

### Solusi
Fallback ke `window.location.pathname`:
```tsx
function parsePath() {
  const parts = window.location.pathname.split("/").filter(Boolean);
  const idx = parts.indexOf("assets");
  return {
    tenant: parts[0] || "",
    id: idx >= 0 ? parts[idx + 1] : undefined,
  };
}

export function AssetDetailPage() {
  const { tenant, id } = parsePath();
  // ...
}
```

---

## 8. BFF Path Handling (Dev vs Prod)

### Gejala
- Dev mode: request ke `/api/v1/assets` → BFF terima `/v1/assets` (Vite strip `/api`)
- Prod mode: request ke `/api/v1/assets` → BFF terima `/api/v1/assets` (nginx forward langsung)

### Akar Masalah
Vite proxy punya config `rewrite: (p) => p.replace(/^\/api/, "")`. Ini menghapus prefix `/api` sebelum forward ke BFF.

### Solusi
BFF harus handle kedua kasus:
```ts
const apiPath = pathname.startsWith('/api') ? pathname : `/api${pathname}`;
const targetUrl = `${GO_SERVICE_URL}${apiPath}`;
```

---

## Quick Reference — Checklist Debugging Frontend

| Gejala | Cek Pertama |
|---|---|
| Halaman stuck loading | Cek backend logs — ada request dari browser? |
| Select/datepicker kosong di edit | Cek apakah ada `value` prop di Select |
| API return 200 tapi data kosong | Cek `transformResponse` — double-wrap? |
| Navigasi salah setelah submit | Cek path — pakai absolute dengan tenant prefix? |
| Proxy error di Docker | Cek `VITE_PROXY_TARGET` env var |
| Component tidak render | Coba ganti `lazyRouteComponent` → direct import |
