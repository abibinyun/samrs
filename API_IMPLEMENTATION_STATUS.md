# SAMRS API Implementation Status

> Last updated: 2026-06-02

## Legend

| Status | Keterangan |
|--------|-----------|
| ✅ DONE | RTK Query API + UI sudah ada |
| 🔶 PARTIAL | API ada tapi belum lengkap / UI belum final |
| ❌ TODO | Belum ada implementasi di frontend |

---

## 1. AUTH

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/auth/login` | ✅ DONE | |
| GET | `/api/v1/auth/me` | ✅ DONE | |
| GET | `/ping` | ❌ TODO | Health check, low priority |

---

## 2. ROOM

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/rooms` | ✅ DONE | |
| GET | `/api/v1/rooms` | ✅ DONE | |
| GET | `/api/v1/rooms/:id` | ✅ DONE | |
| PATCH | `/api/v1/rooms/:id` | ✅ DONE | |
| DELETE | `/api/v1/rooms/:id` | ✅ DONE | |

---

## 3. BED

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/beds` | ✅ DONE | |
| GET | `/api/v1/beds` | ✅ DONE | |
| GET | `/api/v1/beds/:id` | ✅ DONE | |
| PATCH | `/api/v1/beds/:id` | ✅ DONE | |
| DELETE | `/api/v1/beds/:id` | ✅ DONE | |

---

## 4. CATEGORY

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/categories` | ✅ DONE | |
| GET | `/api/v1/categories` | ✅ DONE | |
| GET | `/api/v1/categories/:id` | ✅ DONE | |
| PATCH | `/api/v1/categories/:id` | ✅ DONE | |
| DELETE | `/api/v1/categories/:id` | ✅ DONE | |

---

## 5. VENDOR

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/vendors` | ❌ TODO | |
| GET | `/api/v1/vendors` | ❌ TODO | |
| GET | `/api/v1/vendors/:id` | ❌ TODO | |
| PATCH | `/api/v1/vendors/:id` | ❌ TODO | |
| DELETE | `/api/v1/vendors/:id` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `vendorApi.ts` (RTK Query)
- [ ] Buat `VendorListPage.tsx`
- [ ] Buat `VendorForm.tsx`
- [ ] Daftarkan route di `routes/index.tsx`
- [ ] Tambahkan sidebar menu

---

## 6. BRAND

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/brands` | ✅ DONE | |
| GET | `/api/v1/brands` | ✅ DONE | |
| GET | `/api/v1/brands/:id` | ✅ DONE | |
| PATCH | `/api/v1/brands/:id` | ✅ DONE | |
| DELETE | `/api/v1/brands/:id` | ✅ DONE | |

---

## 7. MODEL

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/models` | ✅ DONE | |
| GET | `/api/v1/models` | ✅ DONE | |
| GET | `/api/v1/models/:id` | ✅ DONE | |
| PATCH | `/api/v1/models/:id` | ✅ DONE | |
| DELETE | `/api/v1/models/:id` | ✅ DONE | |

---

## 8. ASSET STATUS

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/asset-statuses` | ✅ DONE | |
| GET | `/api/v1/asset-statuses` | ✅ DONE | |
| GET | `/api/v1/asset-statuses/:id` | ✅ DONE | |
| PATCH | `/api/v1/asset-statuses/:id` | ✅ DONE | |
| DELETE | `/api/v1/asset-statuses/:id` | ✅ DONE | |

---

## 9. ASSET

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/assets` | ✅ DONE | |
| GET | `/api/v1/assets` | ✅ DONE | |
| GET | `/api/v1/assets/:id` | ✅ DONE | |
| GET | `/api/v1/assets/:id/timeline` | ❌ TODO | Asset history/audit trail |
| GET | `/api/v1/assets/:id/qr` | 🔶 PARTIAL | QR dialog ada, tapi fetch langsung |
| PATCH | `/api/v1/assets/:id` | ✅ DONE | |
| DELETE | `/api/v1/assets/:id` | ✅ DONE | |

**Frontend tasks:**
- [ ] Tambah endpoint `getAssetTimeline` di `assetApi.ts`
- [ ] Buat `AssetTimelinePage.tsx` atau tab di detail

---

## 10. COMPLAINT

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/complaints` | ✅ DONE | |
| GET | `/api/v1/complaints` | ✅ DONE | |
| GET | `/api/v1/complaints/:id` | ✅ DONE | |
| PATCH | `/api/v1/complaints/:id` | ✅ DONE | |
| DELETE | `/api/v1/complaints/:id` | ✅ DONE | |

---

## 11. MAINTENANCE

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/maintenance-schedules` | ✅ DONE | |
| GET | `/api/v1/maintenance-schedules` | ✅ DONE | |
| GET | `/api/v1/maintenance-schedules/:id` | ✅ DONE | |
| PATCH | `/api/v1/maintenance-schedules/:id` | ✅ DONE | |
| PATCH | `/api/v1/maintenance-schedules/:id/complete` | ✅ DONE | |
| DELETE | `/api/v1/maintenance-schedules/:id` | ✅ DONE | |
| POST | `/api/v1/maintenance-schedules/:id/documents` | ❌ TODO | Upload document |
| GET | `/api/v1/maintenance-schedules/:id/documents` | ❌ TODO | List documents |
| GET | `/api/v1/maintenance-documents/:id/download` | ❌ TODO | Download document |
| DELETE | `/api/v1/maintenance-documents/:id` | ❌ TODO | Delete document |

**Frontend tasks:**
- [ ] Tambah endpoints dokumen di `maintenanceApi.ts`
- [ ] Buat `MaintenanceDocuments.tsx` (upload + list + download + delete)

---

## 12. MUTATION

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/asset-mutations` | ❌ TODO | **Masih mock data** |
| GET | `/api/v1/asset-mutations` | ❌ TODO | **Masih mock data** |
| GET | `/api/v1/asset-mutations/:id` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `mutationApi.ts` (RTK Query real, bukan mock)
- [ ] Ganti `runMutationWorker()` dengan real API calls
- [ ] Buat `MutationForm.tsx` (form create mutation)

---

## 13. STOCK OPNAME

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/stock-opnames` | ✅ DONE | |
| GET | `/api/v1/stock-opnames` | ✅ DONE | |
| GET | `/api/v1/stock-opnames/:id` | ✅ DONE | |
| PATCH | `/api/v1/stock-opnames/:id` | ❌ TODO | Update session |
| PATCH | `/api/v1/stock-opnames/:id/close` | ✅ DONE | |
| POST | `/api/v1/stock-opnames/:id/items` | ✅ DONE | |
| GET | `/api/v1/stock-opnames/:id/items` | ✅ DONE | |

**Frontend tasks:**
- [ ] Tambah `updateStockOpname` endpoint

---

## 14. DOCUMENT

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/documents` | ❌ TODO | Create document |
| GET | `/api/v1/documents` | ❌ TODO | List documents |
| GET | `/api/v1/documents/:id` | ❌ TODO | Detail document |
| PATCH | `/api/v1/documents/:id` | ❌ TODO | Update document |
| DELETE | `/api/v1/documents/:id` | ❌ TODO | Delete document |
| POST | `/api/v1/documents/:id/files` | ❌ TODO | Upload file |
| GET | `/api/v1/documents/:id/files` | ❌ TODO | List files |
| GET | `/api/v1/document-files/:id/download` | ❌ TODO | Download file |
| DELETE | `/api/v1/document-files/:id` | ❌ TODO | Delete file |

**Frontend tasks:**
- [ ] Buat `documentApi.ts` (RTK Query)
- [ ] Buat `DocumentListPage.tsx`
- [ ] Buat `DocumentDetailPage.tsx`
- [ ] Buat `DocumentForm.tsx`
- [ ] Buat `DocumentFileUpload.tsx`

---

## 15. NOTIFICATION

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/notifications/send` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `notificationApi.ts` (RTK Query)
- [ ] Implementasi di halaman yang relevan

---

## 16. RBAC (ROLES & PERMISSIONS)

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/roles` | ❌ TODO | |
| GET | `/api/v1/roles` | ❌ TODO | |
| GET | `/api/v1/roles/:id` | ❌ TODO | |
| PATCH | `/api/v1/roles/:id` | ❌ TODO | |
| PATCH | `/api/v1/roles/:id/tenant-admin` | ❌ TODO | SuperAdmin only |
| DELETE | `/api/v1/roles/:id` | ❌ TODO | |
| POST | `/api/v1/roles/:id/permissions` | ❌ TODO | Assign permissions |
| DELETE | `/api/v1/roles/:id/permissions` | ❌ TODO | Revoke permissions |
| GET | `/api/v1/permissions` | ❌ TODO | List all permissions |
| GET | `/api/v1/roles/:id/permissions` | ❌ TODO | Get role permissions |

**Frontend tasks:**
- [ ] Buat `rbacApi.ts` (RTK Query)
- [ ] Buat `RoleListPage.tsx`
- [ ] Buat `RoleForm.tsx` (dengan permission checkboxes)
- [ ] Buat `PermissionListPage.tsx`

---

## 17. USER

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| POST | `/api/v1/users` | ❌ TODO | |
| GET | `/api/v1/users` | ❌ TODO | |
| PATCH | `/api/v1/users/:id` | ❌ TODO | |
| PATCH | `/api/v1/users/:id/role` | ❌ TODO | |
| PATCH | `/api/v1/users/:id/status` | ❌ TODO | |
| PATCH | `/api/v1/users/:id/password` | ❌ TODO | |
| DELETE | `/api/v1/users/:id` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `userApi.ts` (RTK Query)
- [ ] Buat `UserListPage.tsx`
- [ ] Buat `UserForm.tsx`
- [ ] Buat `UserDetailPage.tsx`

---

## 18. AUDIT TRAIL

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| GET | `/api/v1/audit-trails` | ❌ TODO | |
| GET | `/api/v1/audit-trails/:id` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `auditApi.ts` (RTK Query)
- [ ] Buat `AuditTrailListPage.tsx`

---

## 19. ADMIN TENANT (SuperAdmin)

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| GET | `/api/v1/admin/tenants` | ❌ TODO | |
| POST | `/api/v1/admin/tenants` | ❌ TODO | |
| GET | `/api/v1/admin/tenants/:id` | ❌ TODO | |
| PATCH | `/api/v1/admin/tenants/:id` | ❌ TODO | |
| PATCH | `/api/v1/admin/tenants/:id/status` | ❌ TODO | |

**Frontend tasks:**
- [ ] Buat `adminTenantApi.ts` (RTK Query)
- [ ] Buat `TenantListPage.tsx`
- [ ] Buat `TenantForm.tsx`

---

## 20. REPORT

| Method | Endpoint | Status | Keterangan |
|--------|----------|--------|------------|
| GET | `/api/v1/reports/assets/export` | ❌ TODO | CSV/XLSX/PDF |
| GET | `/api/v1/reports/complaints/export` | ❌ TODO | CSV/XLSX/PDF |
| GET | `/api/v1/reports/maintenance-schedules/export` | ❌ TODO | CSV/XLSX/PDF |

**Frontend tasks:**
- [ ] Buat `reportApi.ts` (RTK Query)
- [ ] Buat export dialog (pilih format + download)

---

## RECOMMENDED PRIORITY

### Phase 1 - Quick Wins (mirror CRUD pattern)
1. **Vendor** - Copy pattern dari Category/Brand
2. **Mutation** - Ganti mock dengan real API

### Phase 2 - Admin Essentials
3. **User Management** - CRUD users
4. **RBAC (Roles)** - Role + permission management

### Phase 3 - Features
5. **Document Management** - Upload/download files
6. **Maintenance Documents** - Upload ke schedule
7. **Asset Timeline** - History perubahan asset

### Phase 4 - Reporting & Admin
8. **Report Export** - CSV/XLSX/PDF download
9. **Audit Trail** - Log aktivitas
10. **Admin Tenant** - SuperAdmin panel
11. **Notification** - Send notifications

---

## STATS

```
Total backend routes  : ~95
Implemented           : ~50 (53%)
Partial               :  ~6 (6%)
Not implemented       : ~39 (41%)
```
