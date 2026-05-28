# 📘 SAMRS CLOUD

**Sistem Aset Manajemen Rumah Sakit – Full Feature Specification**

---

## 1️⃣ Ringkasan Proyek

**SAMRS Cloud** adalah platform **SaaS multi-tenant** untuk manajemen aset rumah sakit (medis & non-medis) yang mendukung **inventarisasi, pelacakan, pemeliharaan, perbaikan, kalibrasi, mutasi aset**, serta **audit dan kepatuhan regulasi**, dengan akses **web & mobile (Android)**.

Target pengguna:

* IPSRS
* Teknisi
* Unit pelayanan
* Manajemen RS
* Auditor / Akreditasi

---

## 2️⃣ Arsitektur Sistem (Final)

### 2.1 Arsitektur Teknis

```
Web App / Android App
        │
        ▼
 REST API (Gin, JWT, RBAC)
        │
        ▼
 Business Layer (Usecase)
        │
        ▼
 Repository (Tenant-scoped)
        │
        ▼
 PostgreSQL  ←→ Object Storage (S3)
```

### 2.2 Prinsip Arsitektur

* Clean Architecture
* Multi-tenant (hard isolation)
* Tenant guard di:

  * Middleware
  * Repository
  * GORM callback
* Audit trail transactional
* Cloud-native & stateless

---

## 3️⃣ Modul Sistem (FULL FEATURE)

---

## A. Core Platform (WAJIB)

### A1. Tenant Management

* Create / update / deactivate tenant
* Tenant metadata (nama RS, alamat, tipe)
* Tenant status (active / suspended)
* Statistik tenant (user, asset, audit count)

**Akses:** Super Admin

---

### A2. Authentication & Authorization

* Login JWT
* Refresh token (future)
* Role-Based Access Control (RBAC)
* Permission-based enforcement
* Tenant admin flag
* `/auth/me` → profile + permission list

---

### A3. User Management

* CRUD user per tenant
* Assign role
* Reset password
* Set status (active / inactive)
* Pagination & filter

---

## B. Master Data (MVP)

### B1. Room

* Code, name, location
* Soft delete
* Tenant-scoped
* Filter + pagination

### B2. Bed

* Sub-lokasi dari room
* Status: `available | occupied | maintenance`
* Relasi ke asset

### B3. Category

* Hierarki kategori (future)
* Slug unik per tenant
* Soft delete

---

## C. Asset Management (INTI)

### C1. Asset Master

* Code, name, brand, model
* Category
* Purchase date & value
* Status: `ready | broken | maintenance`
* Placement: room / bed
* Soft delete

---

### C2. Asset Lifecycle Event (KRITIS)

**(INI YANG MELENGKAPI AUDIT TRAIL)**

```
asset_events
- created
- placed
- moved
- broken
- repaired
- maintenance
- calibrated
- disposed
```

* Timeline aset
* Digunakan untuk akreditasi
* Tampil di detail asset

---

### C3. QR Code & Tracking

* Generate QR code (PNG/SVG)
* QR berisi URL + asset_id
* Public read-only endpoint
* Scan via Android app

---

## D. Operasional IPSRS

### D1. Complaint / Work Order (MVP+)

* Report kerusakan (mobile/web)
* Status workflow:

  * open
  * in_progress
  * done
  * cancelled
* Assign teknisi
* Catatan perbaikan
* Relasi ke asset
* Auto asset_event

---

### D2. Maintenance Schedule

* Jadwal rutin (bulanan / tahunan)
* Checklist digital
* Reminder (future)
* Riwayat maintenance

---

### D3. Calibration

* Jadwal kalibrasi
* Sertifikat upload
* Expiry tracking
* History per asset

---

## E. Asset Movement & Inventory

### E1. Asset Mutation

* Pindah room / bed
* Approval (optional)
* Riwayat mutasi

---

### E2. Stock Opname

* Periodik
* Checklist aset
* Missing / found status
* Audit evidence

---

## F. Sparepart Management (Phase 2)

* Master sparepart
* Stok masuk/keluar
* Pemakaian di work order
* Minimum stock alert

---

## G. Dokumentasi & Kepatuhan

### G1. Document Management

* SOP
* Manual alat
* Sertifikat kalibrasi
* Upload multi-file
* Versioning (future)

---

### G2. Audit Trail (SUDAH ADA – DIPERLUAS)

* Create / update / delete
* User, waktu, data before/after
* Tenant scoped
* Filter & pagination
* Export (future)

---

## H. Reporting & Analytics

* Asset by category
* Asset by status
* Work order SLA
* Maintenance compliance
* Audit summary
* Export PDF / Excel

---

## I. Security & Hardening

* JWT HS256
* Rate limiting
* Structured logging (Zap)
* Metrics Prometheus
* Soft delete everywhere
* Tenant isolation enforced

---

## J. Mobile Application (Android)

### Fitur Mobile MVP:

* Login
* Scan QR
* View asset
* Report complaint
* Update work order

---

## 4️⃣ MVP FINAL (Disepakati)

### 🎯 MVP Modules:

✅ Auth + RBAC
✅ Tenant + User
✅ Room / Bed / Category
✅ Asset + QR
✅ Complaint / Work Order
✅ Asset Event
✅ Audit Trail

🚫 Belum MVP:

* Depresiasi
* Sparepart detail
* Digital signature
* Notification canggih

---

## 5️⃣ Roadmap Implementasi

### Phase 1 (Pilot RS)

* MVP lengkap
* 1–2 RS
* Feedback lapangan

### Phase 2

* Maintenance & calibration
* Reporting
* Mobile enhancement

### Phase 3

* Depresiasi
* Sparepart
* Integrasi SIMRS

---