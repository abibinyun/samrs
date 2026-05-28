# Rencana Teknis Multi‑Tenant (Backend Go)

Tujuan: backend siap untuk **subdomain tenant**, **custom domain**, **on‑prem**, dan **SaaS JWT** tanpa rewrite besar, tetap backward compatible.

## Ringkasan Target
- Tenant bisa di‑resolve dari **Host → Header → JWT → DEFAULT_TENANT_ID**.
- Tenant di‑inject ke **gin.Context** dan **request.Context()**.
- Repository bisa memakai **WithTenant()** dari context tanpa harus menerima tenant_id eksplisit.
- Mode **single‑tenant** via config (`SINGLE_TENANT=true`).

---

## A. Middleware TenantResolver (core)

### Urutan resolusi (prioritas)
1) **Host** (subdomain/custom domain)
2) **Header** (mis. `X-Tenant-ID` atau `X-Tenant-Slug`)
3) **JWT** (tenant_id dari claims)
4) **DEFAULT_TENANT_ID** (fallback on‑prem / single‑tenant)

### Contoh implementasi (ringkas)
```go
// middleware/tenant_resolver.go
func TenantResolver(tenantRepo repository.TenantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, source, err := resolveTenantFromRequest(c, tenantRepo)
		if err != nil {
			util.ErrorResponseFromErr(c, "Tenant tidak ditemukan", util.ErrUnauthorized("tenant resolve failed"))
			c.Abort()
			return
		}

		// Inject ke gin.Context
		c.Set("tenant_id", tenantID.String())
		c.Set("tenant_id_uuid", tenantID)
		c.Set("tenant_source", source) // host|header|jwt|default

		// Inject ke request.Context
		ctx := context.WithValue(c.Request.Context(), httputil.ContextTenantIDKey{}, tenantID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
```

### Resolver (ringkas)
```go
func resolveTenantFromRequest(c *gin.Context, tenantRepo repository.TenantRepository) (uuid.UUID, string, error) {
	// 1) Host (subdomain/custom domain)
	if host := parseHost(c); host != "" {
		if tenant, err := tenantRepo.FindByDomainOrSubdomain(host); err == nil {
			return tenant.ID, "host", nil
		}
	}
	// 2) Header
	if raw := strings.TrimSpace(c.GetHeader("X-Tenant-ID")); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			return id, "header", nil
		}
	}
	if slug := strings.TrimSpace(c.GetHeader("X-Tenant-Slug")); slug != "" {
		if tenant, err := tenantRepo.FindBySlug(slug); err == nil {
			return tenant.ID, "header", nil
		}
	}
	// 3) JWT (compat dengan current)
	if raw, ok := c.Get("tenant_id"); ok {
		if id, err := uuid.Parse(fmt.Sprint(raw)); err == nil {
			return id, "jwt", nil
		}
	}
	// 4) DEFAULT_TENANT_ID (on‑prem)
	if raw := strings.TrimSpace(os.Getenv("DEFAULT_TENANT_ID")); raw != "" {
		if id, err := uuid.Parse(raw); err == nil {
			return id, "default", nil
		}
	}
	return uuid.Nil, "", errors.New("tenant not resolved")
}
```

Catatan: middleware ini **dipasang sebelum** `TenantScope()` (atau menggantikan `TenantScope` sepenuhnya).

---

## B. Context Key untuk Tenant

```go
// httputil/context_keys.go
type ContextTenantIDKey struct{}

func TenantIDFromRequest(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(ContextTenantIDKey{})
	id, ok := v.(uuid.UUID)
	return id, ok && id != uuid.Nil
}
```

Handler tetap bisa pakai `TenantIDFromContext(c)` yang sudah ada untuk backward compatibility.

---

## C. Repository Wrapper (tanpa tenant_id eksplisit)

### Target
Tambahkan wrapper agar repository **bisa menerima context** dan apply tenant scope otomatis.

### Contoh
```go
// repository/base/with_tenant.go
func WithTenantFromContext(db *gorm.DB, ctx context.Context) (*gorm.DB, error) {
	tenantID, ok := httputil.TenantIDFromRequest(ctx)
	if !ok {
		return nil, util.ErrUnauthorized("tenant id missing")
	}
	return WithTenant(db, tenantID), nil
}
```

### Contoh repository (baru)
```go
// repository/room_repository.go
func (r *roomRepository) FindAll(ctx context.Context, filter RoomFilter) ([]domain.Room, int64, error) {
	db := repobase.NewDB(r.db).Model(&domain.Room{})
	query, err := repobase.WithTenantFromContext(db, ctx)
	if err != nil {
		return nil, 0, err
	}
	// lanjutkan query...
}
```

**Backward compatible:** method lama `FindAllByTenant(tenantID, filter)` tetap dipertahankan sementara.

---

## D. Contoh Handler + Usecase (context‑based)

### Handler
```go
func (h *RoomHandler) List(c *gin.Context) {
	filter := parseFilter(c)
	items, total, err := h.usecase.ListRooms(c.Request.Context(), filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengambil data", err)
		return
	}
	util.SuccessResponseWithMeta(c, "Data ditemukan", items, meta(total, filter))
}
```

### Usecase
```go
func (u *roomUsecase) ListRooms(ctx context.Context, filter RoomFilter) ([]domain.Room, int64, error) {
	return u.roomRepo.FindAll(ctx, filter)
}
```

---

## E. Proposal Implementasi Bertahap (<= 5 Langkah)

1) **Tambahkan TenantResolver middleware** dan context key (tanpa hapus JWT flow).  
   - Pasang sebelum `TenantScope` / replace `TenantScope` jika aman.

2) **Tambah DEFAULT_TENANT_ID + SINGLE_TENANT** (config).  
   - Jika `SINGLE_TENANT=true`, skip JWT requirement untuk tenant.

3) **Tambah repository wrapper context‑based** (`WithTenantFromContext`).  
   - Buat versi baru method repo yang menerima `context.Context`.

4) **Migrasi handler/usecase bertahap** ke context‑based API.  
   - Modul baru langsung pakai versi context, modul lama tetap pakai tenantID param.

5) **Deprecate tenantID explicit** setelah semua modul migrasi.  
   - Hapus method lama ketika semua endpoint stabil.

Semua langkah di atas **backward compatible**.

---

## F. Contoh Alur Request Nyata (Verifikasi)

### 1) Subdomain tenant (tenant1.samrs.com)
1. Request: `GET https://tenant1.samrs.com/api/v1/rooms`
2. TenantResolver → parse Host `tenant1.samrs.com` → `tenant_id` dari DB.
3. Inject ke gin.Context + request.Context.
4. Usecase → repo → query tenant‑scoped.

### 2) Custom domain (rs‑persahabatan.com)
1. Request: `GET https://rs-persahabatan.com/api/v1/rooms`
2. TenantResolver → domain mapping (`tenant_domains` table).
3. Inject tenant id → proceed.

### 3) On‑prem (no DNS, no JWT)
1. Config: `SINGLE_TENANT=true`, `DEFAULT_TENANT_ID=<uuid>`
2. Request: `GET http://localhost:8080/api/v1/rooms`
3. TenantResolver → fallback DEFAULT_TENANT_ID → proceed tanpa JWT.

### 4) SaaS default (JWT)
1. Request: `GET /api/v1/rooms` with Authorization: Bearer
2. Auth middleware sets `tenant_id` from JWT (current behavior).
3. TenantResolver uses JWT (fallback #3) → proceed.

---

## G. Minimal Perubahan Data Model (Opsional)
- `tenant_domains` table: `tenant_id`, `domain`, `is_primary`.
- `tenants` table: add `subdomain` or use `slug` as subdomain.

---

## H. Checklist Implementasi
- [ ] TenantResolver middleware tersedia & terpasang.
- [ ] Context tenant key tersedia (request.Context).
- [ ] SINGLE_TENANT + DEFAULT_TENANT_ID config.
- [ ] Repo wrapper context‑based + test.
- [ ] Minimal 1 modul migrasi sebagai contoh (Room/Asset).

## Log Implementasi
- [x] TenantResolver + context key ditambahkan dan dipasang di routes protected.
- [x] Env `ROOT_DOMAIN`, `SINGLE_TENANT`, `DEFAULT_TENANT_ID` ditambahkan ke `.env`/`.env.example` dan README.
- [x] Repo wrapper `WithTenantFromContext` + migrasi contoh modul Room (handler/usecase/repo).
- [x] Migrasi context-based untuk Asset list (handler/usecase/repo).
- [x] Migrasi context-based untuk User list (handler/usecase/repo).
- [x] Migrasi context-based untuk Category list (handler/usecase/repo).
- [x] Migrasi context-based untuk Bed list (handler/usecase/repo).
- [x] Migrasi context-based untuk Vendor list (handler/usecase/repo).
- [x] Migrasi context-based untuk Brand list (handler/usecase/repo).
- [x] Migrasi context-based untuk Model list (handler/usecase/repo).
- [x] Migrasi context-based untuk Asset Status list (handler/usecase/repo).
- [x] Unit test TenantResolver (subdomain/header/custom domain/jwt/default).
- [x] Migrasi context-based untuk Complaint list (handler/usecase/repo).
- [x] Migrasi context-based untuk Maintenance Schedule list (handler/usecase/repo).
- [x] Migrasi context-based untuk Stock Opname list (handler/usecase/repo).
- [x] Migrasi context-based untuk Audit Trail list (handler/usecase/repo).
- [x] Migrasi context-based untuk Document list (handler/usecase/repo).
- [x] Migrasi context-based untuk Asset Mutation list (handler/usecase/repo).
- [x] Migrasi context-based untuk Role list (handler/usecase/repo).
