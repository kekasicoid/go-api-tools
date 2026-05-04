# JWT Tools

Kumpulan tools untuk **decode** dan **validasi** JWT token berbasis HMAC.

---

## Endpoints

| Method | Path | Deskripsi |
|--------|------|-----------|
| `POST` | `/tools/jwt/decode` | Parse header & claims tanpa verifikasi signature |
| `POST` | `/tools/jwt/validate` | Verifikasi signature HMAC + kembalikan claims |

> Semua request wajib menyertakan header `request-id` (alphanumeric, maks 50 karakter).

---

## POST /tools/jwt/decode

Memecah JWT menjadi bagian **header** dan **claims** tanpa memverifikasi signature.  
Berguna untuk inspeksi isi token secara cepat.

### Request

**Header**

| Key | Tipe | Wajib | Keterangan |
|-----|------|-------|------------|
| `request-id` | string | ✅ | ID unik request |

**Body** (`application/json`)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `token` | string | ✅ | JWT token yang akan di-decode |

### Response

**200 OK**

```json
{
  "response_code": "200",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "success",
  "response_data": {
    "header": {
      "alg": "HS256",
      "typ": "JWT"
    },
    "claims": {
      "sub": "1234567890",
      "name": "John Doe",
      "iat": 1516239022
    }
  }
}
```

**400 Bad Request** — token kosong atau format JWT tidak valid

```json
{
  "response_code": "400",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "token is required"
}
```

---

## POST /tools/jwt/validate

Memverifikasi **signature HMAC** (HS256 / HS384 / HS512) JWT menggunakan secret yang disediakan.  
Mengembalikan claims jika token valid.

> ⚠️ Hanya mendukung algoritma HMAC. Token dengan algoritma RSA/ECDSA akan ditolak.

### Request

**Header**

| Key | Tipe | Wajib | Keterangan |
|-----|------|-------|------------|
| `request-id` | string | ✅ | ID unik request |

**Body** (`application/json`)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "secret": "my-secret-key"
}
```

| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `token` | string | ✅ | JWT token yang akan divalidasi |
| `secret` | string | ✅ | HMAC secret key |

### Response

**200 OK — Token valid**

```json
{
  "response_code": "200",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "token is valid",
  "response_data": {
    "valid": true,
    "claims": {
      "sub": "1234567890",
      "name": "John Doe",
      "iat": 1516239022
    }
  }
}
```

**200 OK — Token tidak valid / signature salah**

```json
{
  "response_code": "200",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "token is not valid",
  "response_data": {
    "valid": false
  }
}
```

**400 Bad Request** — token atau secret kosong

```json
{
  "response_code": "400",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "secret is required"
}
```

---

## Perbedaan Decode vs Validate

| Aspek | `/jwt/decode` | `/jwt/validate` |
|-------|---------------|-----------------|
| Verifikasi signature | ❌ Tidak | ✅ Ya |
| Perlu secret | ❌ Tidak | ✅ Ya |
| Cocok untuk | Inspeksi isi token | Autentikasi token |

---

## Struktur Kode

```
internal/
  domain/jwt.go              # Interface JWTDecoder
  model/jwt.go               # Request & response models
  usecase/jwt_usecase.go     # Business logic
  delivery/http/
    jwt_handler.go           # Handler (terpisah dari handler.go)
    router.go                # Route registration

pkg/
  jwtutil/decoder.go         # Implementasi decode & validate
```

---

## Lihat Juga

- [Swagger UI](http://localhost:8080/swagger/index.html) — dokumentasi interaktif semua endpoint
- [readme.md](../readme.md) — panduan menjalankan server
