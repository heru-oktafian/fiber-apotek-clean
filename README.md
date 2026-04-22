# fiber-apotek-clean

Rewrite backend **API Apotek** berbasis Go dengan pendekatan **Clean Architecture**.

Project ini dibuat sebagai penerus bertahap dari repo legacy `fiber-apotek`, dengan tujuan:
- menjaga kontrak endpoint penting tetap tersedia
- merapikan struktur kode agar lebih maintainable
- mengurangi coupling ke framework/ORM di layer bisnis
- mendorong pengembangan domain per domain, bukan big-bang rewrite

## Status proyek

Project ini masih dalam fase **incremental rewrite**.

Sudah tersedia dan sudah pernah divalidasi sebagian:
- auth flow utama
- menus
- profile
- list branches user
- branches read
- user branches read
- user management read
- product create + combo endpoints
- purchase create baseline
- sale create baseline
- opname baseline

Belum semua domain legacy selesai. Lihat bagian **Struktur API Saat Ini** dan **Roadmap / Progress** di bawah.

---

## Teknologi yang digunakan

### Bahasa & runtime
- **Go 1.24**

### Web framework
- **Fiber v2**

### Database & persistence
- **PostgreSQL**
- **GORM**

### Authentication & security
- **JWT** (`github.com/golang-jwt/jwt/v5`)
- **bcrypt** (`golang.org/x/crypto`)

### Cache / token blacklist
- **Redis** (`github.com/redis/go-redis/v9`)

### Configuration
- **godotenv** untuk load `.env`

### Validation
- **go-playground/validator**

---

## Dependensi utama

Dependensi utama dari `go.mod`:
- `github.com/gofiber/fiber/v2`
- `gorm.io/gorm`
- `gorm.io/driver/postgres`
- `github.com/golang-jwt/jwt/v5`
- `github.com/redis/go-redis/v9`
- `github.com/joho/godotenv`
- `golang.org/x/crypto`
- `github.com/go-playground/validator/v10`

Untuk install dependency:

```bash
go mod tidy
```

---

## Kebutuhan untuk menjalankan repo

Agar repo ini bisa jalan penuh, minimal butuh:

### 1. Go
- Go `1.24.x`

### 2. PostgreSQL
- database dev yang sesuai dengan schema legacy API Apotek

### 3. Redis
- dipakai untuk token blacklist / auth support

### 4. File environment
Buat `.env` di root project.

Contoh variabel penting:

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=aptdev_db

REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASS=
REDIS_DB=0

SERVER_PORT=1113
JWT_SECRET_KEY=your-secret
APPNAME=Rest API Apotek | v2.0
APP_TIMEZONE=Asia/Jakarta
```

### Catatan penting env
Project ini membaca secret JWT dengan fallback berikut:
- `JWT_SECRET`
- jika kosong, pakai `JWT_SECRET_KEY`

Jadi env lama dan env baru sama-sama bisa didukung.

---

## Cara menjalankan aplikasi

### Jalankan langsung

```bash
cd fiber-apotek-clean
go run ./cmd/api
```

Secara default server listen di:

```text
http://0.0.0.0:1113
```

### Build

```bash
go build ./...
```

---

## Struktur arsitektur

Project ini mengikuti pola **Clean Architecture**.

### Alur request
1. request masuk ke **router**
2. router memanggil **handler**
3. handler memanggil **usecase**
4. usecase menggunakan **ports/interface**
5. adapter repository menjalankan query ke DB / Redis / JWT
6. response dikembalikan ke client

### Struktur folder utama

```text
cmd/api                         # entry point aplikasi
internal/bootstrap              # perakitan dependency app
internal/domain                 # entity dan bentuk data domain
internal/usecase                # logika bisnis
internal/ports                  # kontrak interface
internal/adapters/http/fiber    # handler, middleware, router HTTP
internal/adapters/persistence   # implementasi akses DB
internal/adapters/cache         # Redis adapter
internal/adapters/auth          # JWT adapter
internal/shared                 # helper umum (config, response, idgen, error)
docs                            # dokumentasi proyek
scripts                         # script pengujian / utilitas
menus.json                      # sumber menu sementara (JSON-based)
```

---

## Struktur API saat ini (implemented)

README ini sekarang mengikuti **implemented API** yang benar-benar sudah hidup di repo saat ini.
Untuk detail header, body, dan catatan request, lihat juga:
- `docs/api-implemented-endpoints.md`

### Diagram struktur API implemented

```text
Api Apotek (Implemented)
├── Health/
│   └── Get - /health
├── Sys/
│   ├── Auth/
│   │   ├── Post - /api/login
│   │   ├── Get - /api/list_branches
│   │   ├── Post - /api/set_branch
│   │   ├── Get - /api/profile
│   │   ├── Get - /api/menus
│   │   └── Post - /api/logout
│   ├── Branches/
│   │   ├── Get - /api/branches
│   │   ├── Get - /api/branches/:id
│   │   ├── Post - /api/branches
│   │   └── Delete - /api/branches/:id
│   ├── User Branches/
│   │   ├── Get - /api/user-branches
│   │   └── Get - /api/user-branches/:user_id/:branch_id
│   └── User Management/
│       ├── Get - /api/users
│       ├── Get - /api/detail-users/:id
│       ├── Post - /api/users
│       └── Put - /api/users/:id
├── Masters/
│   └── Products/
│       ├── Post - /api/products
│       ├── Get - /api/sales-products-combo
│       ├── Get - /api/purchase-products-combo
│       └── Get - /api/cmb-product-opname
├── Transactions/
│   ├── Purchases/
│   │   └── Post - /api/purchases
│   └── Sales/
│       └── Post - /api/sales
└── Audit & Finances/
    └── Opnames/
        ├── Post - /api/opnames
        ├── Get - /api/opnames/:id
        ├── Post - /api/opname-items
        └── Post - /api/opname-items-all
```

---

## Dokumentasi yang tersedia

Di folder `docs/` sudah ada beberapa dokumen penting:

- `docs/fiber-apotek-analysis.md`
  - analisis repo legacy
- `docs/clean-architecture-notes.md`
  - catatan pendekatan arsitektur
- `docs/project-map.md`
  - peta struktur repo
- `docs/feature-matrix.md`
  - status fitur rewrite vs legacy
- `docs/runtime-smoke-test.md`
  - panduan validasi runtime
- `docs/api-implemented-endpoints.md`
  - daftar endpoint yang sudah implemented, lengkap dengan header/body penting

---

## Roadmap / progress singkat

### Sudah lumayan stabil
- Auth flow utama
- Menus
- Branches read
- User branches read
- User management read
- Product combo endpoints
- Purchase create baseline
- Sale create baseline
- Opname baseline

### Sedang / akan dilanjutkan
1. Sys > User Management full
2. Masters > Suppliers
3. Masters > Units / Categories / Products
4. Transactions full CRUD
5. First Stocks
6. Defectas
7. Reports / Dashboard / Export

---

## Known gaps

- belum semua domain legacy selesai
- banyak endpoint masih baru baseline, belum full CRUD
- export Excel/PDF belum dibangun untuk sebagian besar domain
- reports/dashboard belum dibangun
- menus masih JSON-based, belum DB-based
- beberapa nilai data runtime bisa berbeda tergantung isi DB dev

---

## Catatan pengembangan

Project ini dikerjakan dengan prinsip:
- incremental rewrite
- compatibility-first untuk endpoint penting
- clean architecture sebagai arah jangka panjang
- perbaikan kualitas internal tanpa harus mewarisi semua kelemahan legacy

Jadi repo ini memang bukan salinan 1:1 repo lama, tapi backend baru yang dibangun bertahap agar bisa mencapai parity fitur secara sehat.

---

## Repository notes

Jika ingin lanjut pengembangan, titik penting yang paling sering disentuh:
- route: `internal/adapters/http/fiber/router/router.go`
- handler: `internal/adapters/http/fiber/handlers/`
- usecase: `internal/usecase/`
- repository: `internal/adapters/persistence/postgres/repositories.go`
- config: `internal/shared/config/config.go`

---

## License / internal usage

Belum ada lisensi publik khusus di repo ini. Gunakan sesuai kebutuhan internal proyek API Apotek.
