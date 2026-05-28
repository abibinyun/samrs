# API Error Format

Dokumen ini menjelaskan format response error yang konsisten di API SAMRS.

## Format Response
Semua error mengikuti struktur berikut:

```json
{
  "success": false,
  "message": "pesan ringkas untuk user",
  "errors": {
    "detail": "detail teknis yang aman untuk ditampilkan"
  }
}
```

Catatan:
- `message` adalah pesan utama yang konsisten.
- `errors.detail` berisi detail yang aman (tidak membocorkan SQL/stack trace).
- Untuk error list, `errors.details` dapat berupa array.

## Request ID
Setiap response menyertakan header:

```
X-Request-ID: <uuid>
```

Gunakan nilai ini saat debugging (log server dan response akan sama).

## Mapping Status Code (Ringkas)
- `400` validation/format request tidak valid
- `401` unauthorized (token/identitas tidak valid)
- `403` forbidden (akses ditolak)
- `404` not found
- `409` conflict (kode/slug sudah digunakan)
- `500` internal error (default)

## Contoh

### Validation Error (400)
```json
{
  "success": false,
  "message": "Input tidak valid",
  "errors": {
    "detail": "kode aset wajib diisi"
  }
}
```

### Not Found (404)
```json
{
  "success": false,
  "message": "Aset tidak ditemukan",
  "errors": {
    "detail": "aset tidak ditemukan"
  }
}
```

### Conflict (409)
```json
{
  "success": false,
  "message": "Gagal membuat ruangan",
  "errors": {
    "detail": "kode ruangan sudah digunakan"
  }
}
```
