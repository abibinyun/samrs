# Repository Guidelines (Safe Repository Pattern)

Tujuan: menjaga query selalu bersih, deterministic, dan aman saat ada middleware/hook.

## Aturan Utama

1) Selalu mulai query dari session baru + model eksplisit:

```
db := r.db.Session(&gorm.Session{NewDB: true}).Model(&domain.X{})
```

2) Selalu gunakan tenant scoping helper:

```
query := withTenant(db, tenantID)
```

3) Hindari `Save()` untuk update.
   Gunakan `Updates(map)` + `WHERE id = ?`:

```
updates := map[string]interface{}{
  "field": value,
  "updated_at": time.Now(),
}

return withTenant(db, entity.TenantID).
  Set("audit_record_id", entity.ID).
  Where("id = ?", entity.ID).
  Updates(updates).
  Error
```

## Kenapa

- Mencegah statement leak antar query (TX + middleware + hook).
- Menghindari update massal tanpa `WHERE`.
- Aman untuk audit hook dan tenant scope guard.

## Checklist Saat Menambah Module

- [ ] Semua query pakai `Session(NewDB:true)` + `Model`.
- [ ] Semua update pakai `Updates(map)` + `Where("id = ?")`.
- [ ] `withTenant(...)` dipakai untuk semua query tenant-scoped.
- [ ] Seeder dan permission repository juga mengikuti pola ini agar konsisten.
