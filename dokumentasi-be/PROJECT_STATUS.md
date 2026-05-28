# Project Status Summary (SAMRS)

## Ringkasan Proyek
SAMRS adalah sistem manajemen aset rumah sakit berbasis SaaS multi-tenant. Target akhirnya adalah platform cloud yang dipakai banyak rumah sakit dengan isolasi data per tenant. Backend menggunakan Go (Gin), database PostgreSQL, dan JWT untuk autentikasi. Arsitektur diarahkan ke Clean Architecture agar kode mudah dirawat dan dikembangkan.

## Fokus MVP
- Struktur proyek yang rapi dan terpisah per layer (delivery, usecase, repository, domain).
- Autentikasi dan otorisasi (JWT + RBAC).
- Core module master data rumah sakit.
- Audit trail untuk setiap perubahan data penting.
- Multi-tenancy: semua data terisolasi per tenant.
- Guard tenant scope di middleware + validasi relasi antar resource menggunakan query tenant-scoped.

## Catatan Teknis
- Saat ini response sudah distandarkan via helper `pkg/util/response.go`.
- JWT menggunakan HS256 dengan secret dari `.env`.
- Struktur folder sudah sesuai rancangan Clean Architecture di dokumentasi.

## Konsep Filter & Pagination (Disepakati)
- **Query param global** untuk semua modul: `page`, `per_page`, `sort_by`, `sort_dir`, `search`, `date_from`, `date_to`.
- **Default** filter tanggal memakai `created_at`; modul tertentu bisa punya field tanggal tambahan (mis. `purchase_date` di asset).
- **Search** hanya pada kolom whitelist per modul (bukan full scan semua kolom).
- **Pagination MVP**: offset (`page`, `per_page`, maksimum `per_page <= 100`).
- **Future-proof**: reserve parameter `cursor` untuk migrasi ke cursor pagination tanpa breaking change.
- **Index minimum** untuk performa: `(tenant_id, created_at)`, `(tenant_id, code)`, `(tenant_id, status)`, `(tenant_id, foreign_id)` sesuai relasi.
- **Tenant scoping** wajib di semua query (tidak boleh query tanpa `tenant_id`).

## Spesifikasi Filter per Modul (Draft Final)
- **Room**: `search` (name, code, location), `date_from`, `date_to`, `sort_by` (created_at, name, code).
- **Bed**: `search` (code), `room_id`, `status`, `date_from`, `date_to`, `sort_by` (created_at, code, status).
- **Category**: `search` (name, slug), `date_from`, `date_to`, `sort_by` (created_at, name).
- **Asset**: `search` (code, name, brand, model), `status`, `category_id`, `room_id`, `purchase_date` range, `date_from`, `date_to`, `sort_by` (created_at, code, status).

## Yang Sudah Dibuat
- Dokumentasi MVP: BRD, FSD, TDD (Part 1 & 2), SAD, Deployment Plan, Roadmap.
- Struktur backend (Clean Architecture):
  - `internal/domain`: Tenant, Role, User, Room, Bed, Category, Asset.
  - `internal/repository`: akses data per entity.
  - `internal/usecase`: logika bisnis per entity.
  - `internal/delivery/http`: handler Gin + middleware.
- Auth JWT:
  - Login menghasilkan JWT (`user_id`, `tenant_id`, `role_id`).
  - Middleware untuk validasi token dan inject claim ke context.
- Guard tenant scope di middleware dan helper repository `withTenant` untuk memastikan query tenant-scoped.
- Guard tenant scope global via GORM callbacks (query/update/delete) untuk model bertenant, dengan bypass eksplisit saat dibutuhkan.
- RBAC lengkap:
  - Model permissions + role_permissions.
  - Endpoint list permissions.
  - Role management per tenant (create/list/update/delete).
  - Assign/revoke permission ke role.
  - User management per tenant (create/list/pagination/update profile/set role/set status/reset password/delete).
  - Reset password user tenant.
  - List permissions by role.
  - List permissions by role sudah tenant-scoped (no cross-tenant access).
  - Endpoint `GET /roles/:id` tenant-scoped untuk detail role.
  - Endpoint `/auth/me` mengembalikan permission list.
  - Enforcement permission per endpoint (core modules + RBAC endpoints).
  - Super admin dapat set/unset flag tenant admin pada role.
- Super admin:
  - Endpoint list semua tenant + role mapping (`GET /api/v1/admin/tenants`).
  - Endpoint detail tenant termasuk roles + counts (`GET /api/v1/admin/tenants/:id`).
  - Detail tenant juga menyertakan list users (paged), list roles (paged), dan counts audit trail.
  - Endpoint create/update/status tenant (super admin only) dengan metadata tenant (`address`, `type`) dan status `active|suspended`.
- AutoMigrate GORM dan seeder admin + tenant default.
- Modul master data lengkap:
  - Room, Bed, Category, Asset: CRUD lengkap, validasi unik per tenant, dan soft delete.
- Relasi aset:
  - Asset bisa ditempatkan di Room atau Bed (bed adalah sub-lokasi dari room).
  - Jika `bed_id` diisi, sistem mengikat asset ke bed dan menyesuaikan `room_id` bila perlu.
- Partial unique index:
  - Unique kode/slug hanya untuk data aktif (`deleted_at IS NULL`) agar soft delete tidak mengunci kode lama.
  - Dibuat otomatis saat startup lewat `ensurePartialUniqueIndexes`.
- Filter & pagination:
  - Semua modul mendukung `page`, `per_page`, `sort_by`, `sort_dir`, `search`, `date_from`, `date_to`.
  - Asset mendukung `purchase_from`, `purchase_to`.
  - Response list mengembalikan `meta` (total, page, per_page).
- Validasi status enum:
  - Bed: `available`, `occupied`, `maintenance`.
  - Asset: `ready`, `broken`, `maintenance`.
- Swagger:
  - Spesifikasi OpenAPI disimpan di `Dokumentasi/swagger.yaml`.
  - UI docs tersedia di `/swagger` (Redoc) dan `/swagger-ui` (Swagger UI), spec di `/swagger.yaml`.
- QR code & public asset:
  - Endpoint generate QR code asset (PNG) untuk link publik.
  - Endpoint public read-only asset via tenant slug + asset code (untuk scan QR).
  - Base URL QR bisa di-override lewat env `PUBLIC_ASSET_BASE_URL`.
- Complaint / work order dasar:
  - CRUD complaint tenant-scoped dengan status workflow `open`, `in_progress`, `done`.
  - Filter list complaint (status, asset_id, assigned_to, reported_by, search, date range).
- Master data tambahan:
  - Vendor (CRUD, tenant-scoped).
  - Merek/brand aset (CRUD, tenant-scoped).
  - Tipe/model alat (CRUD, tenant-scoped, terkait merek).
  - Status aset terstandar (CRUD, tenant-scoped).
  - Asset mendukung referensi vendor/brand/model opsional.
- Asset lifecycle events:
  - Event otomatis saat asset create/update/delete.
  - Event terkait complaint (reported, in_progress, done).
  - Endpoint timeline asset (`GET /api/v1/assets/:id/timeline`).
- Maintenance & calibration:
  - CRUD jadwal maintenance/calibration per aset.
  - Endpoint complete schedule yang otomatis menulis asset event.
  - Filter jadwal berdasarkan tipe, status, asset_id, created_at, dan due date.
- Dokumen/sertifikat maintenance:
  - Upload dokumen (sertifikat/laporan) ke jadwal maintenance/calibration.
  - Endpoint list, download, dan delete dokumen.
  - Storage file lokal dengan base URL `PUBLIC_FILE_BASE_URL` (opsional).
- Document management global:
  - Dokumen SOP/manual/sertifikat tenant-scoped dengan metadata (judul, tipe, deskripsi).
  - Upload multi-file per dokumen + tambah file ke dokumen yang sudah ada.
  - Endpoint list/detail, list file, download file, dan delete dokumen/file.
  - Storage file lokal dengan base URL `PUBLIC_FILE_BASE_URL` (opsional).
- Asset mutation:
  - Endpoint mutasi aset untuk memindahkan asset antar ruangan/bed.
  - Event asset dicatat pada timeline saat mutasi.
- Stock opname:
  - Session opname (draft/closed) + item opname per asset.
  - Event asset dicatat saat item opname dibuat.
- Reporting & export:
  - Export CSV/XLSX/PDF untuk assets, complaints, dan maintenance schedules (via `format` param).
  - Support field selection via `fields` dan header override via `columns`.
  - Catatan implementasi modul baru: buat whitelist field + mapping value (mirip `resolveAssetExportFields`/`buildAssetRows`) lalu panggil `writeReportExport` agar otomatis support `format`, `fields`, dan `columns`.
  - Filter mengikuti query list masing-masing modul.
- Notifikasi (handler only):
  - Endpoint send notification dengan channel `email`, `whatsapp`, `generic` (noop handler).
  - Struktur request siap untuk integrasi provider eksternal.
- Audit trail:
  - Tabel `audit_trails` dibuat dan logging otomatis via GORM hooks untuk create/update/delete pada model bertenant (termasuk Room/Bed/Category/Asset, RBAC/User).
  - Audit transactional (data + audit dalam satu transaksi) melalui transaksi handler.
  - Endpoint audit trail list dengan filter/pagination (`GET /api/v1/audit-trails`).
  - Endpoint detail audit (`GET /api/v1/audit-trails/:id`).
  - Audit endpoints dibatasi ke tenant admin flag (atau system role).
  - Swagger audit ditambah contoh response.
  - Audit manual untuk tenant management (create/update/status) agar tercatat per tenant target.
- Helper response format standar JSON.
- `api_test.http` untuk pengujian endpoint (auth, room, bed, category, asset).
- Unit tests dasar untuk auth/RBAC serta guard tenant scope dan helper audit hook.
- Integration test audit hook (DB) tersedia via env `SAMRS_TEST_DSN`.
- Hardening: rate limit middleware, structured logging (Zap), dan endpoint metrics Prometheus `/metrics`.

## Yang Belum Dibuat (Sesuai Roadmap MVP)
- Mobile app (Android) untuk scan QR + komplain + update work order.
