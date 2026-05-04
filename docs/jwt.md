# JWT Tools

Tool untuk **decode** dan **validasi** JWT token berbasis HMAC dalam satu endpoint.

---

## Endpoints

| Method | Path | Deskripsi |
|--------|------|-----------|
| `POST` | `/tools/jwt/decode-validation` | Parse header & claims, dan opsional verifikasi signature HMAC |

> Semua request wajib menyertakan header `request-id` (alphanumeric, maks 50 karakter).

---

## POST /tools/jwt/decode-validation

Memecah JWT menjadi bagian **header** dan **claims**.  
Jika `secret` disertakan, signature HMAC (HS256 / HS384 / HS512) juga diverifikasi dan hasilnya tercermin di field `valid`.

> ⚠️ Verifikasi signature hanya mendukung algoritma HMAC. Token dengan algoritma RSA/ECDSA akan ditolak saat validasi.

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
| `token` | string | ✅ | JWT token yang akan di-decode |
| `secret` | string | ❌ | HMAC secret key untuk verifikasi signature. Jika tidak diisi, `valid` selalu `false` |

### Response

**200 OK — Decode berhasil, signature valid**

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
    },
    "valid": true
  }
}
```

**200 OK — Decode berhasil, signature tidak valid atau secret tidak diisi**

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
    },
    "valid": false
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
