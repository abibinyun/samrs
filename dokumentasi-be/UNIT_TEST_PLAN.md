# Rencana Unit Test (SAMRS Backend)

Dokumen ini berisi daftar test yang akan dibuat.  
Setelah test selesai dikerjakan, pindahkan dari **Planned** ke **Done**.

---

## Planned

### RBAC

### User Management

### Room/Bed/Category/Asset (Core)

### Asset Events / Timeline

### Maintenance & Documents

### Complaint / Work Order

### Stock Opname

### Reports

### Public & Misc

---

## Done

### Auth & Security
- [x] Auth login: valid credentials → token issued.
- [x] Auth login: invalid credentials → 401 + message.
- [x] Auth me: token valid → user + permissions.
- [x] Auth me: token invalid/expired → 401.

### Tenant Scoping & Guard
- [x] Tenant scope guard: query tanpa tenant → error.
- [x] Tenant scope guard: query dengan tenant → ok.
- [x] Skip tenant scope: khusus path tertentu → ok.

### RBAC
- [x] Role create/update/delete: happy path.
- [x] Permission assign/revoke: role permissions updated.
- [x] List permissions by role: tenant-scoped.
- [x] Tenant admin flag: set/unset (super admin only).
- [x] Access control: endpoint protected by permission → 403 jika tidak punya.

### User Management
- [x] Create user: role valid + tenant scoped.
- [x] Update user profile.
- [x] Update user role.
- [x] Activate/Deactivate user.
- [x] Reset password.

### Room/Bed/Category/Asset (Core)
- [x] Room create/update/delete.
- [x] Bed create/update/delete (room scoped).
- [x] Category create/update/delete.
- [x] Asset create/update/delete.
- [x] Asset status validation (registered vs default).

### Asset Events / Timeline
- [x] Create asset event on asset create.
- [x] List timeline with filters (event_type, date).

### Maintenance & Documents
- [x] Maintenance schedule create/update/complete.
- [x] Maintenance document upload/list/delete.
- [x] Document create/update/delete.
- [x] Document file add/list/delete.

### Complaint / Work Order
- [x] Complaint create/update/delete.
- [x] Complaint status transitions.

### Stock Opname
- [x] Stock opname session create/close.
- [x] Stock opname item add/list.

### Reports
- [x] Export assets: csv/xlsx/pdf (basic response).
- [x] Export complaints: csv/xlsx/pdf.
- [x] Export maintenance schedules: csv/xlsx/pdf.

### Public & Misc
- [x] Public asset endpoint by tenant slug + asset code.
- [x] Notifications send (noop handler).
- [x] Audit trail list/detail.
