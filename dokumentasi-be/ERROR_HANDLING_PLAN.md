# Error Handling Plan

## Tujuan
- Menstandarkan format response error agar stabil untuk frontend/mobile.
- Memastikan status code konsisten sesuai kategori error.
- Mengurangi leak error internal (DB/ORM) ke klien.
- Meningkatkan observability lewat logging error yang terstruktur.
- Mengurangi potensi panic dari hook/middleware.

## Rencana Pengerjaan (Planned)
1) Standarkan payload error (errors sebagai object dengan code/detail).
2) Buat helper mapping error → HTTP status (validation/not found/conflict/internal).
3) Terapkan error type/sentinel di usecase agar handler bisa mapping konsisten.
4) Normalisasi error di middleware (RBAC/Auth) agar status code tepat.
5) Guard audit hook untuk map/slice update dan fallback record_id.
6) Tambahkan logging error terstruktur + request id (tanpa bocor detail ke response).
7) Audit handler yang masih pakai err.Error() langsung → ganti pesan aman.
8) Update dokumentasi API untuk error format yang konsisten.

## Log Pekerjaan (Done)
- [x] Standarkan payload error agar `errors` konsisten (detail/details).
- [x] Tambah helper mapping error → HTTP status (belum diterapkan di handler).
- [x] Terapkan error type/sentinel di usecase (validation/not_found/conflict/forbidden/unauthorized).
- [x] Normalisasi error di middleware (RBAC/Auth/Tenant scope/admin).
- [x] Guard audit hook untuk map/slice agar tidak panic saat update batch/map.
- [x] Tambahkan request ID + logging error terstruktur di middleware logger.
- [x] Audit handler dan ganti error response agar mapping status konsisten + detail aman.
- [x] Update dokumentasi format error response + request id.
- [x] Rapikan pesan error agar tidak bocor detail internal (message stabil).
- [x] Audit status code: validasi upload dokumen (400) + role not found di permission list (404).
