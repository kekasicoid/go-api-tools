# Instagram Media Downloader

Tool untuk **mengunduh media** (foto, video, reel, carousel) dari URL post Instagram publik — tanpa memerlukan API key.

---

## Endpoints

| Method | Path | Deskripsi |
|--------|------|-----------|
| `POST` | `/tools/instagram/download` | Ekstrak URL download langsung dari post Instagram |

> Semua request wajib menyertakan header `request-id` (alphanumeric, maks 50 karakter).

---

## POST /tools/instagram/download

Mengambil URL media yang bisa di-download langsung dari sebuah post Instagram publik.  
Mendukung **foto**, **video**, **reel**, dan **carousel** (album multi-foto/video).

> ⚠️ **Batasan:**
> - Hanya mendukung konten **publik**. Konten akun private tidak bisa diakses tanpa session cookie login.
> - **Stories** membutuhkan cookie session akun Instagram yang sudah login — tidak didukung tanpa auth.
> - Instagram dapat mengubah struktur halaman sewaktu-waktu yang dapat mempengaruhi hasil ekstraksi.

> ⚠️ **ToS Instagram:** Men-scrape Instagram tanpa izin dapat melanggar [Instagram Terms of Use](https://help.instagram.com/581066165581870). Gunakan tool ini hanya untuk konten milik sendiri atau keperluan edukasi/riset. Untuk keperluan produksi, gunakan [Instagram Graph API](https://developers.facebook.com/docs/instagram-api) resmi.

---

### Request

**Header**

| Key | Tipe | Wajib | Keterangan |
|-----|------|-------|------------|
| `request-id` | string | ✅ | ID unik request (alphanumeric, maks 50 karakter) |

**Body** (`application/json`)

```json
{
  "url": "https://www.instagram.com/arditya.kekasi/p/C2uaPTYShvn/"
}
```

| Field | Tipe | Wajib | Keterangan |
|-------|------|-------|------------|
| `url` | string | ✅ | URL post Instagram (harus mengandung domain `instagram.com`) |

**Format URL yang didukung:**

| Format | Tipe Konten |
|--------|-------------|
| `https://www.instagram.com/p/{shortcode}/` | Post foto / video |
| `https://www.instagram.com/{username}/p/{shortcode}/` | Post foto / video (dengan username) |
| `https://www.instagram.com/reel/{shortcode}/` | Reel |
| `https://www.instagram.com/tv/{shortcode}/` | IGTV |

---

### Response

**200 OK — Media berhasil diekstrak**

```json
{
  "response_code": "200",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "success",
  "response_data": {
    "post_type": "post",
    "items": [
      {
        "media_type": "photo",
        "media_url": "https://scontent.cdninstagram.com/v/t51.2885-15/...",
        "thumb_url": "https://scontent.cdninstagram.com/v/t51.2885-15/..."
      }
    ]
  }
}
```

| Field | Tipe | Keterangan |
|-------|------|------------|
| `post_type` | string | Tipe post: `post`, `reel`, `tv`, atau `story` |
| `items` | array | Daftar media dalam post (lebih dari satu untuk carousel) |
| `items[].media_type` | string | `photo` atau `video` |
| `items[].media_url` | string | URL langsung ke file media (foto/video) |
| `items[].thumb_url` | string | URL thumbnail (opsional, mungkin kosong untuk beberapa tipe) |

**Contoh — Carousel (album):**

```json
{
  "response_code": "200",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "success",
  "response_data": {
    "post_type": "post",
    "items": [
      {
        "media_type": "photo",
        "media_url": "https://scontent.cdninstagram.com/v/t51.2885-15/slide1...",
        "thumb_url": "https://scontent.cdninstagram.com/v/t51.2885-15/slide1..."
      },
      {
        "media_type": "video",
        "media_url": "https://scontent.cdninstagram.com/v/t50.2886-16/slide2...",
        "thumb_url": "https://scontent.cdninstagram.com/v/t51.2885-15/slide2_thumb..."
      }
    ]
  }
}
```

**400 Bad Request** — URL tidak valid atau bukan URL Instagram

```json
{
  "response_code": "400",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "url must be an Instagram URL (instagram.com)",
  "response_data": null
}
```

**500 Internal Server Error** — post private, tidak ada media ditemukan, atau struktur halaman berubah

```json
{
  "response_code": "500",
  "response_refnum": "",
  "response_id": "abc-123",
  "response_desc": "no downloadable media found; the post may be private or Instagram's page structure has changed",
  "response_data": null
}
```

---

## Caching

Response di-cache di **Redis** selama **30 menit** berdasarkan URL. Jika Redis tidak dikonfigurasi, setiap request akan langsung mengambil dari Instagram.

---

## Struktur Kode

```
internal/
  domain/instagram.go              # Interface InstagramDownloader + domain types
  model/instagram.go               # Request & response models
  usecase/instagram_usecase.go     # Business logic + Redis caching
  delivery/http/
    instagram_handler.go           # Handler
    router.go                      # Route registration

pkg/
  instagramutil/scraper.go         # Implementasi HTML scraper
```

---

## Lihat Juga

- [Swagger UI](http://localhost:8080/swagger/index.html) — dokumentasi interaktif semua endpoint
- [readme.md](../readme.md) — panduan menjalankan server
