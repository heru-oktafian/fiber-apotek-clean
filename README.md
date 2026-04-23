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
│   │   ├── Get - /api/user-branches/:user_id/:branch_id
│   │   └── Post - /api/user-branches
│   └── User Management/
│       ├── Get - /api/users
│       ├── Get - /api/detail-users/:id
│       ├── Post - /api/users
│       └── Put - /api/users/:id
├── Masters/
│   ├── Products/
│   │   ├── Get - /api/products
│   │   ├── Get - /api/products/excel
│   │   ├── Get - /api/products/pdf
│   │   ├── Post - /api/products
│   │   ├── Get - /api/products/:id
│   │   ├── Put - /api/products/:id
│   │   ├── Delete - /api/products/:id
│   │   ├── Get - /api/sales-products-combo
│   │   ├── Get - /api/purchase-products-combo
│   │   └── Get - /api/cmb-product-opname
│   ├── Suppliers/
│   │   ├── Get - /api/suppliers
│   │   ├── Get - /api/suppliers/excel
│   │   ├── Get - /api/suppliers/pdf
│   │   ├── Get - /api/suppliers/:id
│   │   ├── Post - /api/suppliers
│   │   ├── Put - /api/suppliers/:id
│   │   ├── Delete - /api/suppliers/:id
│   │   └── Get - /api/suppliers-combo
│   ├── Units/
│   │   ├── Get - /api/units
│   │   ├── Get - /api/units/excel
│   │   ├── Get - /api/units/pdf
│   │   ├── Get - /api/units/:id
│   │   ├── Post - /api/units
│   │   ├── Put - /api/units/:id
│   │   ├── Delete - /api/units/:id
│   │   └── Get - /api/cmb-units
│   └── Categories/
│       ├── Product Categories/
│       │   ├── Get - /api/product-categories
│       │   ├── Get - /api/product-categories/excel
│       │   ├── Get - /api/product-categories/pdf
│       │   ├── Post - /api/product-categories
│       │   ├── Get - /api/product-categories/:id
│       │   ├── Put - /api/product-categories/:id
│       │   ├── Delete - /api/product-categories/:id
│       │   └── Get - /api/product-categories-combo
│       ├── Supplier Categories/
│       │   ├── Get - /api/supplier-categories
│       │   ├── Get - /api/supplier-categories/excel
│       │   ├── Get - /api/supplier-categories/pdf
│       │   ├── Post - /api/supplier-categories
│       │   ├── Get - /api/supplier-categories/:id
│       │   ├── Put - /api/supplier-categories/:id
│       │   ├── Delete - /api/supplier-categories/:id
│       │   └── Get - /api/supplier-categories-combo
│       └── Member Categories/
│           ├── Get - /api/member-categories
│           ├── Get - /api/member-categories/excel
│           ├── Get - /api/member-categories/pdf
│           ├── Get - /api/member-categories/:id
│           ├── Post - /api/member-categories
│           ├── Put - /api/member-categories/:id
│           ├── Delete - /api/member-categories/:id
│           └── Get - /api/member-categories-combo
├── Transactions/
│   ├── Purchases/
│   │   ├── Get - /api/purchases
│   │   ├── Get - /api/purchases/excel
│   │   ├── Get - /api/purchases/pdf
│   │   ├── Get - /api/purchases/:id
│   │   ├── Post - /api/purchases
│   │   ├── Put - /api/purchases/:id
│   │   ├── Delete - /api/purchases/:id
│   │   ├── Get - /api/purchase-items/excel
│   │   ├── Get - /api/purchase-items/pdf
│   │   ├── Get - /api/purchase-items/all/:id
│   │   ├── Post - /api/purchase-items
│   │   ├── Put - /api/purchase-items/:id
│   │   └── Delete - /api/purchase-items/:id
│   ├── Sales/
│   │   ├── Get - /api/sales
│   │   ├── Get - /api/sales/excel
│   │   ├── Get - /api/sales/pdf
│   │   ├── Get - /api/sales/:id
│   │   ├── Post - /api/sales
│   │   ├── Put - /api/sales/:id
│   │   ├── Delete - /api/sales/:id
│   │   ├── Get - /api/sale-items/excel
│   │   ├── Get - /api/sale-items/pdf
│   │   ├── Get - /api/sale-items/all/:id
│   │   ├── Post - /api/sale-items
│   │   ├── Put - /api/sale-items/:id
│   │   └── Delete - /api/sale-items/:id
│   └── Duplicate Receipts/
│       ├── Get - /api/duplicate-receipts
│       ├── Get - /api/duplicate-receipts/excel
│       ├── Get - /api/duplicate-receipts/pdf
│       ├── Post - /api/duplicate-receipts
│       ├── Get - /api/duplicate-receipts/:id
│       ├── Put - /api/duplicate-receipts/:id
│       ├── Delete - /api/duplicate-receipts/:id
│       ├── Get - /api/duplicate-receipts-items/all/:id
│       ├── Get - /api/duplicate-receipts-items/excel
│       ├── Get - /api/duplicate-receipts-items/pdf
│       ├── Post - /api/duplicate-receipts-items
│       ├── Put - /api/duplicate-receipts-items/:id
│       └── Delete - /api/duplicate-receipts-items/:id
├── Audits/
│   ├── First Stocks/
│   │   ├── Get - /api/first-stocks
│   │   ├── Get - /api/first-stocks/excel
│   │   ├── Get - /api/first-stocks/pdf
│   │   ├── Post - /api/first-stocks
│   │   ├── Put - /api/first-stocks/:id
│   │   ├── Delete - /api/first-stocks/:id
│   │   ├── Get - /api/first-stock-with-items/:id
│   │   ├── Get - /api/first-stock-items/:id
│   │   ├── Get - /api/first-stock-items/excel
│   │   ├── Get - /api/first-stock-items/pdf
│   │   ├── Post - /api/first-stock-items
│   │   ├── Put - /api/first-stock-items/:id
│   │   └── Delete - /api/first-stock-items/:id
│   └── Opnames/
│       ├── Post - /api/opnames
│       ├── Get - /api/opnames/:id
│       ├── Post - /api/opname-items
│       └── Post - /api/opname-items-all
└── Finances/
    ├── Another Incomes/
    │   ├── Get - /api/another-incomes
    │   ├── Get - /api/another-incomes/excel
    │   ├── Get - /api/another-incomes/pdf
    │   ├── Post - /api/another-incomes
    │   ├── Put - /api/another-incomes/:id
    │   └── Delete - /api/another-incomes/:id
    └── Expenses/
        ├── Get - /api/expenses
        ├── Get - /api/expenses/excel
        ├── Get - /api/expenses/pdf
        ├── Post - /api/expenses
        ├── Put - /api/expenses/:id
        └── Delete - /api/expenses/:id
```

---

### Catatan implementasi terbaru

Beberapa milestone parity yang sudah lebih matang di repo ini sekarang mencakup:
- transaksi `purchases` dan `sales` berikut item CRUD dan export baseline
- `buy returns` dan `sale returns` berikut combo sumber transaksi dan export baseline
- `first stocks` berikut header/item flow dan export baseline
- `duplicate receipts` sebagai **sale-like transaction** berbasis resep dokter / kopi resep, dengan implementasi saat ini mencakup header CRUD + item CRUD + export baseline:
  - `GET /api/duplicate-receipts`
  - `GET /api/duplicate-receipts/excel`
  - `GET /api/duplicate-receipts/pdf`
  - `POST /api/duplicate-receipts`
  - `GET /api/duplicate-receipts/:id`
  - `PUT /api/duplicate-receipts/:id`
  - `DELETE /api/duplicate-receipts/:id`
  - `GET /api/duplicate-receipts-items/all/:id`
  - `GET /api/duplicate-receipts-items/excel`
  - `GET /api/duplicate-receipts-items/pdf`
  - `POST /api/duplicate-receipts-items`
  - `PUT /api/duplicate-receipts-items/:id`
  - `DELETE /api/duplicate-receipts-items/:id`

Untuk detail request/response dan status next batch, lihat `docs/api-implemented-endpoints.md` dan `docs/implemented-vs-next-target.md`.

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
- `docs/implemented-vs-next-target.md`
  - ringkasan endpoint yang sudah hidup vs target domain berikutnya

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
