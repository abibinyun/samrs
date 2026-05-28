# Backend Testing & Verification Plan

## 📊 Status Overview

Berdasarkan analisa codebase dan `api_test.http`, backend SAMRS memiliki:

- **Total Files**: 875 files
- **Lines of Code**: ~23,279 LOC
- **Test Coverage**: Unit tests tersedia untuk usecase layer
- **API Endpoints**: ~100+ endpoints
- **Modules**: 15+ domain modules

---

## 🎯 Tujuan Testing

1. **Verifikasi semua endpoint berfungsi** dengan benar
2. **Validasi business logic** sesuai requirement
3. **Cek error handling** dan edge cases
4. **Pastikan RBAC & permission** working
5. **Test multi-tenant isolation**
6. **Validasi audit trail** logging

---

## 📋 Testing Scope

### ✅ Sudah Ada (dari `api_test.http`)

1. **Authentication** - Login, JWT, Me endpoint
2. **Room Management** - CRUD rooms
3. **Bed Management** - CRUD beds dengan room relation
4. **Category Management** - CRUD categories
5. **Vendor/Brand/Model/Status** - Master data CRUD
6. **Asset Inventory** - CRUD assets dengan relations
7. **Asset Timeline** - Event tracking
8. **Asset QR Code** - Generate & public access
9. **Complaint/Work Order** - CRUD complaints
10. **Maintenance Schedule** - CRUD schedules
11. **Maintenance Documents** - Upload/download files
12. **Asset Mutation** - Move assets between locations
13. **Stock Opname** - Session & items management
14. **Report Export** - CSV/XLSX/PDF exports
15. **Notifications** - Send notifications (noop)
16. **RBAC** - Roles, permissions, users
17. **Audit Trail** - View audit logs
18. **Admin Tenants** - Multi-tenant management (Super Admin)
19. **Document Management** - Global document storage

### ⚠️ Perlu Diverifikasi

1. **Permission enforcement** - Apakah RBAC benar-benar block unauthorized access?
2. **Tenant isolation** - Apakah data tenant A tidak bisa diakses tenant B?
3. **Super Admin bypass** - Apakah `is_system` role bypass semua permission?
4. **File upload limits** - Max file size, allowed types
5. **Pagination** - Edge cases (page 0, negative, huge number)
6. **Sorting & filtering** - All combinations working?
7. **Concurrent operations** - Race conditions?
8. **Database constraints** - Unique, foreign key violations
9. **Audit trail completeness** - Semua operasi tercatat?
10. **Error messages** - User-friendly & informative?

### ❌ Belum Ada / TODO (dari comment di api_test.http)

1. **Detail Location tracking** - Building/Floor/Room/Bed hierarchy
2. **Grouping Category** - Category grouping/nesting
3. **Complaint Status CRUD** - Dynamic status management dengan priority
4. **Complaint permission** - Update status dengan permission check
5. **Notification Handler** - Real notification (email/SMS/push) untuk staff
6. **Custom Form Dynamic** - Form builder untuk schedule/monitoring

---

## 🗂️ Testing Structure

```
docs/
└── backend-testing/
    ├── 01-test-plan.md                    # This file
    ├── 02-api-endpoints-checklist.md      # Checklist semua endpoint
    ├── 03-rbac-permission-matrix.md       # Permission testing matrix
    ├── 04-tenant-isolation-tests.md       # Multi-tenant testing
    ├── 05-edge-cases-scenarios.md         # Edge cases & error scenarios
    ├── 06-performance-tests.md            # Load & performance testing
    ├── 07-security-audit.md               # Security checklist
    ├── 08-test-results/                   # Hasil testing
    │   ├── test-run-YYYY-MM-DD.md
    │   └── issues-found.md
    └── 09-automation/                     # Test automation scripts
        ├── run-all-tests.sh
        └── test-data-generator.go
```

---

## 📝 Next Steps

1. **Create detailed checklists** untuk setiap module
2. **Setup test environment** (Docker compose sudah ada)
3. **Run manual tests** menggunakan `api_test.http`
4. **Document results** dengan screenshot/logs
5. **Fix issues** yang ditemukan
6. **Automate tests** (optional - menggunakan Go test atau Postman/Newman)

---

## 🚀 Execution Plan

### Phase 1: Core Functionality (Priority: HIGH)
- [ ] Authentication & Authorization
- [ ] Asset Management (CRUD)
- [ ] Room & Bed Management
- [ ] RBAC & Permissions

### Phase 2: Business Logic (Priority: HIGH)
- [ ] Complaint Management
- [ ] Maintenance Scheduling
- [ ] Asset Mutation
- [ ] Stock Opname

### Phase 3: Supporting Features (Priority: MEDIUM)
- [ ] Document Management
- [ ] Report Export
- [ ] Audit Trail
- [ ] Notifications

### Phase 4: Admin Features (Priority: MEDIUM)
- [ ] Tenant Management (Super Admin)
- [ ] User Management
- [ ] Master Data Management

### Phase 5: Advanced Testing (Priority: LOW)
- [ ] Performance Testing
- [ ] Security Audit
- [ ] Edge Cases
- [ ] Concurrent Operations

---

## 📊 Success Criteria

- ✅ **100% endpoint coverage** - Semua endpoint tested
- ✅ **Zero critical bugs** - No blocking issues
- ✅ **RBAC working** - Permission enforcement correct
- ✅ **Tenant isolation** - No data leakage
- ✅ **Audit trail complete** - All operations logged
- ✅ **Error handling** - Proper error messages
- ✅ **Documentation updated** - All findings documented

---

**Created**: 2026-02-28  
**Status**: Planning Phase  
**Next**: Create detailed endpoint checklist
