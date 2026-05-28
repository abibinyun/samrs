Di bawah ini aku buat **PLAN LANJUTAN YANG TIDAK MENAMBAH FITUR**, tapi **menguatkan fondasi** supaya nanti:

* feature baru **tidak memicu bug aneh**
* AI / dev lain **tidak merusak sistem**
* kamu bisa scale dengan tenang

---

# 🧱 PHASE 2 — FOUNDATION HARDENING PLAN (NO NEW FEATURES)

> Target: **stabil, aman, bisa dikembangkan 2–3 tahun ke depan tanpa rewrite**

---

## 🎯 GOAL UTAMA PHASE INI

1. **Zero tenant leak**
2. **Predictable auth & permission**
3. **Bug mudah dilacak**
4. **Perubahan kecil tidak merusak sistem**
5. **Developer (atau AI) lain tidak “tersesat”**

---

## 1️⃣ AUTH & RBAC HARDENING (Paling penting)

### 1.1 Buat “Auth Contract” (wajib, bukan kode)

Dokumen singkat (1–2 halaman) berisi:

* definisi:

  * siapa boleh akses apa
  * dalam kondisi apa ditolak
* mapping:

  * endpoint → permission
  * endpoint → role minimum
* response standard:

  * 401 = tidak login
  * 403 = tidak punya permission
  * 404 = bukan milik tenant

📌 Ini **bukan fitur**, tapi **kontrak sistem**.

---

### 1.2 Test Matrix Auth (WAJIB)

Buat tabel test manual & automated:

| Role         | Tenant | Endpoint        | Expected |
| ------------ | ------ | --------------- | -------- |
| seller       | own    | GET /assets     | 200      |
| seller       | other  | GET /assets/:id | 404      |
| tenant admin | own    | DELETE /assets  | 403/200  |
| super admin  | any    | GET /tenants    | 200      |
| unauth       | any    | POST /assets    | 401      |

👉 Ini mencegah bug kelas “kok bisa tembus”.

---

### 1.3 Freeze Middleware Logic

* Middleware auth
* Middleware tenant
* Middleware permission

❗ **Tidak boleh diubah tanpa approval eksplisit**
Kalau perlu perubahan:

* tambah middleware baru
* jangan ubah behavior lama

---

## 2️⃣ MULTI-TENANCY SAFETY NET (Double Lock)

### 2.1 Audit Semua Repository

Checklist:

* ❌ Tidak ada query tanpa `tenant_id`
* ❌ Tidak ada `First(&model)` tanpa scope
* ❌ Tidak ada `Preload()` tanpa tenant guard

Bikin aturan:

> “Query tanpa tenant = BUG KRITIS”

---

### 2.2 Lock GORM Callback

* Pastikan:

  * callback tenant scope **tidak bisa dimatikan tanpa flag eksplisit**
  * bypass hanya untuk super admin + context flag

Tambahkan test:

* query tanpa tenant → **panic / error**

---

## 3️⃣ DATA INTEGRITY & CONSISTENCY

### 3.1 Invariant Rules (Dokumentasi + Test)

Contoh invariant:

* Asset **selalu** punya tenant
* Bed **selalu** milik room tenant yang sama
* Asset + bed **tidak boleh beda tenant**
* Soft delete **tidak boleh memecah relasi aktif**

Bikin:

* daftar invariant
* unit test khusus invariant

---

### 3.2 Transaction Boundary Review

Audit:

* create asset + audit
* mutation asset + event
* stock opname + item + event

Pastikan:

* **atomic**
* rollback aman
* tidak ada partial write

---

## 4️⃣ AUDIT TRAIL HARDENING

### 4.1 Definisikan “What must be audited”

Bukan semua perubahan harus audited.
Tetapkan:

* MUST audit: asset, role, permission, user status, tenant status
* NICE: complaint status, maintenance complete
* NO NEED: read-only

Ini mencegah audit bloat.

---

### 4.2 Audit Integrity Test

Test:

* data gagal → audit tidak tercatat
* audit ada → data pasti ada
* audit tenant A tidak bisa dibaca tenant B

---

## 5️⃣ ERROR HANDLING & OBSERVABILITY

### 5.1 Error Classification

Standarkan:

* Validation error
* Auth error
* Permission error
* Business rule violation
* System error

Pastikan:

* response konsisten
* log level benar (info/warn/error)

---

### 5.2 Structured Logging Policy

Pastikan log **selalu punya**:

* request_id
* tenant_id (jika ada)
* user_id (jika ada)

Ini bikin bug production bisa dilacak.

---

## 6️⃣ PERFORMANCE BASELINE (TANPA OPTIMASI PREMATUR)

### 6.1 Baseline Metrics

Catat:

* auth latency
* list asset (10, 50, 100 data)
* audit insert time

Bukan untuk optimasi sekarang, tapi **baseline future regression**.

---

### 6.2 Query Budget

Tetapkan:

* list endpoint max query count
* preload max depth

Kalau lewat → dianggap bug.

---

## 7️⃣ FILE & STORAGE PREPARATION (Future Proof)

### 7.1 Storage Abstraction

Pastikan:

* file storage lewat interface
* path selalu tenant-scoped
* mudah ganti ke S3/GCS

Tidak perlu implement S3 sekarang.

---

## 8️⃣ API CONTRACT FREEZE

### 8.1 Version Lock

* `/api/v1` tidak boleh breaking change
* kalau mau ubah → `/api/v2`

---

### 8.2 Response Snapshot Test

Snapshot JSON response penting:

* login
* list asset
* asset detail
* audit list

Ini mencegah “kok frontend rusak”.

---

## 9️⃣ DOCUMENTATION FOR FUTURE YOU (SANGAT PENTING)

Minimal dokumen:

1. **Auth & Permission Flow**
2. **Tenant Isolation Rules**
3. **How to Add New Module (tanpa bocor tenant)**
4. **What NOT to do**

Ini bukan dokumentasi marketing, tapi **tameng dari kesalahan**.

---

## 10️⃣ EXIT CRITERIA (PHASE SELESAI KALAU…)

Phase ini selesai jika:

* Semua auth & tenant test hijau
* Tidak ada query tanpa tenant
* Tidak ada endpoint tanpa permission
* Bug kelas login loop **tidak mungkin terjadi**
* Kamu berani bilang:

  > “kalau ada dev baru masuk, sistem ini aman”

---

# 🧠 Prinsip Emas Selanjutnya

> **Lebih baik menolak fitur baru, daripada menambal bug lama**

Kamu sudah melewati fase “build fast”.
Sekarang kamu masuk fase **“build correct”**.

---