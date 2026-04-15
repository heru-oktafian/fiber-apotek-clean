# Project Map - fiber-apotek-clean

Dokumen ini menjelaskan peta struktur repo, alur request, lokasi route, dan status fitur saat ini.

## 1. Gambaran besar arsitektur

Project ini ditulis dengan arah **Clean Architecture**.

Prinsip sederhananya:
- `domain` berisi bentuk data dan konsep bisnis inti
- `usecase` berisi logika bisnis
- `ports` berisi kontrak interface
- `adapters` berisi implementasi detail luar seperti HTTP, Postgres, Redis, JWT
- `bootstrap` merakit semuanya menjadi aplikasi berjalan
- `cmd/api` adalah entry point aplikasi

Alur request umumnya:

1. request masuk ke **route**
2. route memanggil **handler**
3. handler memanggil **usecase**
4. usecase memanggil **repository/adapter** lewat **ports**
5. repository baca/tulis ke **database**
6. hasil dikembalikan ke handler lalu ke response HTTP

---

## 2. Struktur folder utama

### `cmd/api`
Titik masuk aplikasi.

File penting:
- `cmd/api/main.go`

Fungsi:
- menjalankan bootstrap app
- start server Fiber di port dari config/env

---

### `internal/bootstrap`
Tempat merakit dependency aplikasi.

File penting:
- `internal/bootstrap/app.go`

Fungsi:
- load `.env`
- load config
- koneksi PostgreSQL
- koneksi Redis
- inisialisasi JWT service
- inisialisasi handler/usecase/repository
- register route Fiber

Kalau mau tahu route disuntik ke mana, lihat file ini.

---

### `internal/domain`
Berisi model bisnis inti dan shape data domain.

Subfolder saat ini:
- `auth`
- `branch`
- `common`
- `member`
- `opname`
- `product`
- `purchase`
- `sale`
- `unit`
- `user`

Fungsi:
- mendefinisikan entity
- request object
- detail object
- tipe bisnis umum

Contoh:
- `internal/domain/purchase/purchase.go`
- `internal/domain/sale/sale.go`
- `internal/domain/opname/opname.go`

---

### `internal/usecase`
Berisi logika bisnis utama.

Subfolder saat ini:
- `auth`
- `product`
- `purchase`
- `sale`
- `opname`

Fungsi:
- proses login/logout/set branch
- create purchase
- create sale
- create opname dan itemnya
- mengambil detail atau item opname

Contoh:
- `internal/usecase/auth/service.go`
- `internal/usecase/purchase/service.go`
- `internal/usecase/sale/service.go`
- `internal/usecase/opname/service.go`

---

### `internal/ports`
Berisi kontrak interface antara usecase dan detail implementasi.

File penting:
- `internal/ports/ports.go`

Fungsi:
- mendefinisikan repository interface
- mendefinisikan token manager interface
- mendefinisikan blacklist interface
- mendefinisikan transaksi dan dependency umum lain

Kalau mau tahu usecase butuh method apa saja dari repo, lihat file ini.

---

### `internal/adapters`
Berisi implementasi detail di layer luar.

#### `internal/adapters/http/fiber`
Untuk lapisan HTTP dengan Fiber.

Subfolder:
- `handlers`
- `middleware`
- `presenter`
- `router`

Fungsi:
- menerima request HTTP
- parse body / params
- panggil usecase
- balikin response
- register route
- middleware auth

File penting:
- `internal/adapters/http/fiber/router/router.go`
- `internal/adapters/http/fiber/handlers/*.go`

#### `internal/adapters/persistence/postgres`
Untuk akses database PostgreSQL.

File penting:
- `internal/adapters/persistence/postgres/models.go`
- `internal/adapters/persistence/postgres/repositories.go`

Fungsi:
- mapping tabel DB
- query database
- implement repository interface
- transaksi database

#### `internal/adapters/cache/redis`
Untuk Redis.

File penting:
- `internal/adapters/cache/redis/blacklist.go`

Fungsi:
- blacklist token login/logout

#### `internal/adapters/auth/jwt`
Untuk generate dan parse JWT.

File penting:
- `internal/adapters/auth/jwt/service.go`

---

### `internal/shared`
Berisi helper umum yang dipakai lintas layer.

Subfolder:
- `apperror`
- `clock`
- `config`
- `idgen`
- `response`

Fungsi:
- error app standar
- waktu sekarang
- pembacaan config/env
- generator ID
- response helper JSON

---

### `docs`
Dokumentasi internal proyek.

File yang sudah ada:
- `docs/fiber-apotek-analysis.md`
- `docs/clean-architecture-notes.md`
- `docs/runtime-smoke-test.md`
- `docs/project-map.md`

---

### `scripts`
Script bantu operasional.

File yang sudah ada:
- `scripts/smoke_test.sh`

Fungsi:
- smoke test baseline endpoint rewrite

---

## 3. Tempat melihat route

Kalau kamu ingin tahu route ada di mana, lihat file ini:

- `internal/adapters/http/fiber/router/router.go`

Itu adalah tempat utama daftar endpoint HTTP di repo baru.

### Route yang saat ini sudah terdaftar
- `GET /health`
- `POST /api/login`
- `POST /api/set_branch`
- `POST /api/logout`
- `POST /api/products`
- `GET /api/sales-products-combo`
- `GET /api/purchase-products-combo`
- `GET /api/cmb-product-opname`
- `POST /api/purchases`
- `POST /api/sales`
- `POST /api/opnames`
- `GET /api/opnames/:id`
- `POST /api/opname-items`
- `POST /api/opname-items-all`

---

## 4. Cara trace satu route

Kalau ingin tahu satu endpoint kerjanya lewat file mana, pakai pola ini:

1. lihat route di:
   - `internal/adapters/http/fiber/router/router.go`
2. lihat handler di:
   - `internal/adapters/http/fiber/handlers/*.go`
3. lihat logic bisnis di:
   - `internal/usecase/*/service.go`
4. lihat query DB di:
   - `internal/adapters/persistence/postgres/repositories.go`
5. lihat shape request/response/domain di:
   - `internal/domain/*/*.go`

### Contoh: `POST /api/sales`
- route:
  - `internal/adapters/http/fiber/router/router.go`
- handler:
  - `internal/adapters/http/fiber/handlers/sale_handler.go`
- usecase:
  - `internal/usecase/sale/service.go`
- repository:
  - `internal/adapters/persistence/postgres/repositories.go`
- domain:
  - `internal/domain/sale/sale.go`

### Contoh: `POST /api/opname-items`
- route:
  - `internal/adapters/http/fiber/router/router.go`
- handler:
  - `internal/adapters/http/fiber/handlers/opname_handler.go`
- usecase:
  - `internal/usecase/opname/service.go`
- repository:
  - `internal/adapters/persistence/postgres/repositories.go`
- domain:
  - `internal/domain/opname/opname.go`

---

## 5. Status fitur saat ini

### Sudah ada baseline
- auth
  - login
  - set branch
  - logout
- product
  - create product
  - sales combo
  - purchase combo
  - opname combo
- purchase
  - create purchase
- sale
  - create sale
- opname
  - create header
  - create/update item by upsert
  - detail opname
  - ambil semua item berdasarkan `opname_id`

### Sudah ada tapi belum full parity
- purchase
  - belum list/detail/update/delete
- sale
  - belum list/detail/update/delete
- opname
  - belum full parity legacy
  - belum delete/update item explicit
  - branch guard belum lengkap
  - response contract belum 1:1
  - belum report sync seperti legacy

### Belum ada sama sekali
- suppliers
- first_stocks
- duplicate_receipts
- sale_returns
- buy_returns / purchase_returns
- another_incomes
- expenses
- master data lain di luar product baseline
- dashboard/report/export

---

## 6. Catatan bisnis penting

### First Stock
- dipakai untuk memasukkan stok produk awal sebelum sistem digunakan
- harus memengaruhi `products.stock`
- tetapi secara perhitungan bisnis dianggap **nol**, bukan pembelian dan bukan pengeluaran
- artinya domain ini tidak boleh disamakan mentah dengan purchase biasa

### Another Incomes
- dipakai untuk mencatat pendapatan tambahan di luar jual beli utama
- termasuk domain pencatatan finansial operasional

### Expenses
- dipakai untuk mencatat pengeluaran apotek
- bisa pengeluaran terencana maupun aksidental
- termasuk domain pencatatan finansial operasional

### Export Excel dan PDF
- fitur export tidak boleh dilupakan
- setiap domain/fungsi utama nantinya perlu dipikirkan export Excel dan PDF
- implementasinya sebaiknya diposisikan sebagai adapter/layer luar, bukan dicampur ke inti usecase

---

## 7. Urutan implementasi yang direkomendasikan

Urutan yang paling masuk akal untuk kelanjutan project:

1. `suppliers`
2. `first_stocks`
3. `expenses`
4. `another_incomes`
5. `buy_returns / purchase_returns`
6. `sale_returns`
7. `duplicate_receipts`
8. export excel/pdf per domain
9. report/dashboard/export lanjutan

Alasannya:
- `suppliers` menopang purchase
- `first_stocks` menyentuh stok inti dan punya aturan bisnis khusus
- `expenses` dan `another_incomes` penting untuk pencatatan operasional finansial
- `buy_returns` dan `sale_returns` dekat dengan domain transaksi inti
- `duplicate_receipts` cenderung fitur turunan, bukan prioritas pertama

---

## 7. File penting yang sering dilihat

Kalau bingung mulai dari mana, biasanya lihat file-file ini dulu:

- `cmd/api/main.go`
- `internal/bootstrap/app.go`
- `internal/adapters/http/fiber/router/router.go`
- `internal/adapters/http/fiber/handlers/`
- `internal/usecase/`
- `internal/adapters/persistence/postgres/repositories.go`
- `internal/ports/ports.go`

---

## 8. Ringkasan super singkat

Kalau disederhanakan:

- route ada di:
  - `internal/adapters/http/fiber/router/router.go`
- handler ada di:
  - `internal/adapters/http/fiber/handlers`
- logika bisnis ada di:
  - `internal/usecase`
- kontrak/interface ada di:
  - `internal/ports/ports.go`
- query database ada di:
  - `internal/adapters/persistence/postgres/repositories.go`
- model/domain ada di:
  - `internal/domain`
- wiring aplikasi ada di:
  - `internal/bootstrap/app.go`

Kalau mau menelusuri fitur baru, paling aman mulai dari **router → handler → usecase → repository → domain**.
