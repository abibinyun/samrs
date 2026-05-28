# Project Status v2 (Phase 2 Foundation Hardening)

Dokumen ini adalah status kerja Phase 2 untuk memperkuat fondasi sistem SAMRS tanpa menambah fitur baru. Fokus utamanya: keamanan multi-tenant, konsistensi data, hardening auth/RBAC, dan observability.

## Prinsip Umum
- Tidak menambah fitur baru (hanya hardening).
- Zero tenant leak adalah target utama.
- Perubahan di middleware/guard harus terkontrol dan terdokumentasi.
- Semua perubahan wajib dicatat di bagian log eksekusi.

## Status Ringkas
- Target fase: Stabil, aman, dan siap dikembangkan 2–3 tahun ke depan tanpa rewrite.
- Fokus utama: Auth/RBAC, tenant isolation, audit trail, error handling, observability.

## Daftar Akan Dikerjakan (Backlog Phase 2)
Catat item yang akan dikerjakan berikutnya di sini.
- [ ] (TBD)  
- [ ] (TBD)  
- [ ] (TBD)

## Checklist Eksekusi (Disiapkan untuk diisi)
Gunakan format: `[ ]` untuk belum, `[~]` untuk sedang, `[x]` untuk selesai.

### 1) Auth & RBAC Hardening
- [ ] Auth Contract (dokumen mapping endpoint→permission + response standard).
- [ ] Test matrix auth (role/tenant/endpoint).
- [ ] Freeze middleware logic (catatan perubahan + approval rule).

### 2) Multi‑Tenancy Safety Net
- [ ] Audit semua repository (no tenantless query).
- [ ] GORM tenant scope guard tidak bisa dimatikan tanpa flag eksplisit.
- [ ] Test: query tanpa tenant harus error/panic.

### 3) Data Integrity & Consistency
- [ ] Daftar invariant (asset/bed/room/tenant relation).
- [ ] Unit test invariant.
- [ ] Review transaction boundary (asset, mutation, stock opname, audit).

### 4) Audit Trail Hardening
- [ ] Definisi “what must be audited”.
- [ ] Audit integrity test (data ↔ audit, cross‑tenant).

### 5) Error Handling & Observability
- [ ] Klasifikasi error + standar response.
- [ ] Logging policy (request_id, tenant_id, user_id).

### 6) Performance Baseline
- [ ] Baseline metrics (auth, list assets, audit insert).
- [ ] Query budget (max query count, preload depth).

### 7) File & Storage Preparation
- [ ] Storage abstraction (interface + tenant‑scoped path).

### 8) API Contract Freeze
- [ ] Version lock (`/api/v1` no breaking change).
- [ ] Response snapshot test (login, list asset, asset detail, audit list).

### 9) Documentation for Future Dev
- [ ] Auth & Permission Flow.
- [ ] Tenant Isolation Rules.
- [ ] How to Add New Module (tenant safe).
- [ ] What NOT to do.

### 10) Exit Criteria
- [ ] All auth & tenant tests green.
- [ ] No tenantless query found.
- [ ] No endpoint without permission.
- [ ] Ready for new dev onboarding.

## Log Eksekusi (Catatan Perubahan)
Catat setiap item yang dikerjakan agar mudah audit dan tracking.
- (TBD)  
- (TBD)  
- (TBD)
