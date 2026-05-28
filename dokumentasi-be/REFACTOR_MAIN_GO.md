# Refactor Plan: `cmd/api/main.go`

Dokumen ini mencatat rencana refactor `main.go` agar maintainable, modular, future‑proof, dan scale‑proof. Tujuan utamanya: memisahkan concern, membuat wiring/route modular, dan memudahkan test/operasional.

## Target Hasil
- `main.go` tipis (hanya config → init → run).
- DI & routing modular.
- Seeding terkontrol (manual/flag).
- Konfigurasi terpusat dan konsisten.
- Mudah ditambah modul baru tanpa edit file besar.

## Daftar Rencana (Akan Dikerjakan)
Gunakan status `[ ]` belum, `[~]` sedang, `[x]` selesai.

### P0 (Prioritas Tinggi)
- [x] Pindahkan seeder ke package/command terpisah (atau env gated).
- [x] Buat `internal/app` sebagai App container (DI + router).
- [x] Pisahkan routing per modul (`routes` per modul).
- [x] Graceful shutdown untuk server.
- [x] Konsisten logger (gunakan satu logger).
- [x] Config via env untuk port, seed, limiter, default tenant, admin password.

### P1 (Penting)
- [x] Helper middleware RBAC agar routing lebih deklaratif.
- [x] Const/wrapper untuk `skip_tenant_scope`.
- [x] Proteksi `/metrics` (port terpisah atau auth).

### P2 (Nice to Have)
- [ ] Snapshot response test (auth/list/detail).
- [ ] Test matrix auth (role × endpoint).
- [ ] Seeder locking (advisory lock) jika tetap seed‑on‑start.

## Log Eksekusi (Yang Sudah Dikerjakan)
Tulis item yang sudah selesai di sini.
- [x] Env gating untuk seeder (`SEED_ON_START`) + port dari `APP_PORT` di `samrs-backend/cmd/api/main.go`.
- [x] Seeder dipindahkan ke `samrs-backend/internal/seed` + command `samrs-backend/cmd/seed`.
- [x] App container di `samrs-backend/internal/app/app.go`, `main.go` hanya init logger + router.
- [x] Routing dipisah per modul di `samrs-backend/internal/app/routes_*.go`.
- [x] Graceful shutdown + close DB pool di `samrs-backend/cmd/api/main.go`.
- [x] Helper RBAC `requirePerm` untuk routes di `samrs-backend/internal/app/routes_helpers.go`.
- [x] Config env untuk seeder + rate limiter di `samrs-backend/internal/seed/seed.go` dan routes login/api.
- [x] Proteksi `/metrics` dengan `METRICS_TOKEN` di `samrs-backend/internal/app/routes_public.go`.
- [x] Const/wrapper `WithSkipTenantScope` di `samrs-backend/internal/config/tenant_scope_guard.go`.
- [x] Std log diarahkan ke Zap di `samrs-backend/cmd/api/main.go` dan `samrs-backend/cmd/seed/main.go`.
