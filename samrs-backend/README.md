# SAMRS Project

SAMRS (Sistem Aset Manajemen Rumah Sakit) adalah backend multi-tenant berbasis Go untuk manajemen aset rumah sakit. Proyek ini fokus pada arsitektur clean, RBAC, tenant scoping, audit trail, dan modul inti operasional untuk MVP.

## Ringkasan Fitur (Saat Ini)
- Multi-tenant SaaS (isolasi data per tenant).
- Auth JWT + RBAC (permission per endpoint).
- Master data: Room, Bed, Category, Asset.
- Relasi aset ke room/bed, mutasi aset, dan stock opname.
- Asset lifecycle events + timeline.
- Complaint / work order dasar.
- Maintenance & calibration schedule + completion.
- Dokumen maintenance (upload/list/download/delete).
- Document management global (SOP/manual/sertifikat, multi-file).
- Reporting export CSV/XLSX/PDF + field selection.
- Audit trail (GORM hook + transactional).
- Hardening: tenant scope guard, structured logging, rate limit, metrics.
- Swagger docs + API test collection.

## Arsitektur

```
samrs-backend/
├── cmd/api                 # entrypoint server
├── internal/
│   ├── config              # DB, hooks, guards
│   ├── domain              # entity & model
│   ├── repository          # data access
│   ├── usecase             # business logic
│   └── delivery/http       # handler + middleware
└── pkg                     # util & helpers
```

## Prasyarat
- Go (lihat `go.mod`)
- PostgreSQL

## Konfigurasi Env
Contoh variabel penting:
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- `JWT_SECRET`
- `PUBLIC_ASSET_BASE_URL` (opsional untuk QR/public asset)
- `PUBLIC_FILE_BASE_URL` (opsional untuk file download)
- `DOCUMENTS_DIR` (opsional, default `./uploads/documents`)
- `MAINTENANCE_DOCS_DIR` (opsional, default `./uploads/maintenance`)

## Menjalankan Server
```
cd samrs-backend
go run cmd/api/main.go
```

## Dokumentasi API
- Swagger spec: `Dokumentasi/swagger.yaml`
- Swagger UI: `/swagger-ui`
- Redoc: `/swagger`
- Spec endpoint: `/swagger.yaml`
- Koleksi test: `samrs-backend/api_test.http`

## Monitoring (Prometheus + Grafana)
File tersedia di `monitoring/`.

Jalankan:
```
cd monitoring
METRICS_TOKEN=your-token docker compose -f docker-compose.monitoring.yml up -d
```

Prometheus: http://localhost:9090  
Grafana: http://localhost:3000 (default login: admin/admin)

Catatan:
- Untuk Linux, ganti target di `monitoring/prometheus.yml` menjadi `172.17.0.1:8080` atau gunakan `--network=host`.

## Logging
Logging berbentuk JSON dan dikirim ke stdout (untuk aggregator) **dan** ke file dengan rotation.

Env yang tersedia:
- `LOG_FILE` (default `./logs/samrs.log`, set ke `-` untuk disable file logging)
- `LOG_MAX_SIZE_MB` (default 100)
- `LOG_MAX_BACKUPS` (default 7)
- `LOG_MAX_AGE_DAYS` (default 7)
- `LOG_COMPRESS` (default true)

## Export Report
Endpoint export saat ini:
- `/api/v1/reports/assets/export`
- `/api/v1/reports/complaints/export`
- `/api/v1/reports/maintenance-schedules/export`

Query penting:
- `format=csv|xlsx|pdf`
- `fields=code,name,status` (whitelist per modul)
- `columns=code:Kode,name:Nama` (override header)

## Catatan Multi‑Tenant & RBAC
- Semua query harus tenant-scoped.
- Endpoint diproteksi oleh permission slug.
- Super admin memiliki hak lintas tenant.

## Pengujian
```
cd samrs-backend
go test ./...
```

## Status Proyek
- Status MVP: `docs/PROJECT_STATUS.md`
- Phase hardening: `docs/PROJECT_STATUS_v2.md`

## Lisensi
TBD (tentukan sesuai kebutuhan sebelum publikasi).
