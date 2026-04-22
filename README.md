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

## Struktur API target

Secara bisnis, struktur besar API yang sedang dibangun adalah:

### 1. Sys
- Auth
- Branches
- User Management
- Membership
- Defectas

### 2. Masters
- Product Categories
- Units
- Products
- Unit Conversions
- Supplier Categories
- Suppliers

### 3. Transactions
- Purchases
- Sales
- Duplicate Receipts / Kopi Resep
- Buy Returns
- Sale Returns
- Expenses
- Another Incomes

### 4. Audit & Finances
- First Stocks
- Opnames

### 5. Reports
- Neraca saldo
- Profit bulanan
- Asset
- export laporan

### 6. Complements
- Dashboard
- Menus
- Test

### Diagram struktur API target

```text
Api Apotek
├── Sys/
│   ├── Auth/
│   │    ├── Post - Login
│   │    ├── Get - Branches List
│   │    ├── Post - Set Branches
│   │    └── Post - Logout
│   ├── Branches/
│   │    ├── Get - All Branches
│   │    ├── Get - Branch by BranchID
│   │    ├── Post - Create Branch
│   │    └── Delete - Branch by BranchID
│   ├── User Management/
│   │    ├── Get - All Users
│   │    ├── Get - Detail User & Branches
│   │    ├── Post - Create User
│   │    ├── Put - Update User by UserID
│   │    └── Post - Adding Branch by UserID
│   ├── Membership/
│   │    ├── Member Category/
│   │    │    ├── Get - All Member Category
│   │    │    ├── Post - Create Member Category
│   │    │    ├── Get - Member Category by MemberCategoryID
│   │    │    ├── Put - Update Category by MemberCategoryID
│   │    │    └── Delete - Delete Member Category by MemberCategoryID
│   │    └── Members/
│   │         ├── Get - Combobox Member Categories
│   │         ├── Get - All Members
│   │         ├── Post - Create Member
│   │         ├── Put - Update Member
│   │         └── Delete - Delete Member
│   └── Defectas/
│        ├── Get - Combobox Products (Purchase Price)
│        ├── Get - All Defectas
│        ├── Post - Create Defecta
│        ├── Put - Update Defecta by DefectaID
│        ├── Delete - Delete Defecta by DefectaID
│        ├── Get - Get All Defecta Items by DefectaID
│        ├── Post - Create Defecta Item by DefectaItemID
│        ├── Put - Update Defecta Item by DefectaItemID
│        ├── Delete - Delete Defecta Item by DefectaItemID
│        └── Get - Download Defectas Excel by DefectaID
├── Masters/
│        ├── Product Categories/
│        │    ├── Get - All Product Categories
│        │    ├── Post - Create Product Category
│        │    ├── Get - Product Category by CategoryProductID
│        │    ├── Put - Update Product Category by CategoryProductID
│        │    ├── Delete - Delete Product Category by CategoryProductID
│        │    ├── Get - Download Product Categories PDF
│        │    └── Get - Download Product Categories Excel
│        ├── Units/
│        │    ├── Get - All Units
│        │    ├── Post - Create Unit
│        │    ├── Get - Unit by UnitID
│        │    ├── Put - Update Unit by UnitID
│        │    ├── Delete - Delete Unit by UnitID
│        │    ├── Get - Download Units PDF
│        │    └── Get - Download Units Excel
│        ├── Products/
│        │    ├── Get - Combobox Categories
│        │    ├── Get - Combobox Units
│        │    ├── Get - All Products
│        │    ├── Post - Create Product
│        │    ├── Get - Product by ProductID
│        │    ├── Put - Update Product by ProductID
│        │    ├── Delete - Delete Product by ProductID
│        │    ├── Get - Download Products Label
│        │    ├── Get - Download Products PDF
│        │    └── Get - Download Products Excel
│        ├── Unit Conversions/
│        │    ├── Get - Combobox Products
│        │    ├── Get - Combobox Units
│        │    ├── Get - All Unit Conversions
│        │    ├── Post - Create Unit Conversion
│        │    ├── Get - Unit Conversion by UnitConversionID
│        │    ├── Put - Update Unit Conversion by UnitConversionID
│        │    ├── Delete - Delete Unit Conversion by UnitConversionID
│        │    ├── Get - Download Unit Conversions PDF
│        │    └── Get - Download Unit Conversions Excel
│        ├── Supplier Categories/
│        │    ├── Get - All Supplier Categories
│        │    ├── Post - Create Supplier Category
│        │    ├── Get - Supplier Category by CategorySupplierID
│        │    ├── Put - Update Supplier Category by CategorySupplierID
│        │    ├── Delete - Delete Supplier Category by CategorySupplierID
│        │    ├── Get - Download Supplier Categories PDF
│        │    └── Get - Download Supplier Categories Excel
│        └── Suppliers/
│             ├── Get - Combobox Supplier Categories
│             ├── Get - All Suppliers
│             ├── Post - Create Supplier
│             ├── Get - Supplier by SupplierID
│             ├── Put - Update Supplier by SupplierID
│             ├── Delete - Delete Supplier by SupplierID
│             ├── Get - Download Suppliers PDF
│             └── Get - Download Suppliers Excel
├── Transactions/
│        ├── Purchases/
│        │    ├── Get - Combobox Product (Purchase Price)
│        │    ├── Get - Combobox Suppliers
│        │    ├── Get - All Purchases
│        │    ├── Post - Create Purchase
│        │    ├── Put - Update Purchase by PurchaseID
│        │    ├── Delete - Delete Purchase by PurchaseID
│        │    ├── Get - All Purchase Items by PurchaseID
│        │    ├── Post - Create Purchase Item
│        │    ├── Put - Update Purchase Item by PurchaseItemID
│        │    ├── Delete - Delete Purchase Item by PurchaseItemID
│        │    ├── Get - Cetak / Print Struk Pembelian
│        │    ├── Get - Download Purchase - PDF
│        │    ├── Get - Download Purchase - Excel
│        │    ├── Get - Download Detail Purchase - PDF
│        │    └── Get - Download Detail Purchase - Excel
│        ├── Sales/
│        │    ├── Get - Combobox Product (Sale Price)
│        │    ├── Get - Combobox Members
│        │    ├── Get - All Sales
│        │    ├── Post - Create Sale
│        │    ├── Put - Update Sale by SaleID
│        │    ├── Delete - Delete Sale by SaleID
│        │    ├── Get - All Sale Items by SaleID
│        │    ├── Post - Create Sale Item
│        │    ├── Put - Update Sale Item by SaleID
│        │    ├── Delete - Delete Sale Item by SaleID
│        │    ├── Get - Cetak / Print Struk Penjualan
│        │    ├── Get - Download Sale - PDF
│        │    ├── Get - Download Sale - Excel
│        │    ├── Get - Download Detail Sale - PDF
│        │    └── Get - Download Detail Sale - Excel
│        ├── Kopi Resep/
│        │    ├── Get - Combobox Product (Sale Price)
│        │    ├── Get - Combobox Members
│        │    ├── Get - All Duplicate Receipts
│        │    ├── Post - Create Duplicate Receipt
│        │    ├── Put - Update Duplicate Receipt by DuplicateReceiptID
│        │    ├── Delete - Delete Duplicate Receipt by DuplicateReceiptID
│        │    ├── Get - All Duplicate Receipt Items by DuplicateReceiptID
│        │    ├── Post - Create Duplicate Receipt Item
│        │    ├── Put - Update Duplicate Receipt Item by DuplicateReceiptID
│        │    ├── Delete - Delete Duplicate Receipt Item by DuplicateReceiptID
│        │    ├── Get - Cetak / Print Struk Penjualan (Duplicate Receipt)
│        │    ├── Get - Download Duplicate Receipt - PDF
│        │    ├── Get - Download Duplicate Receipt - Excel
│        │    ├── Get - Download Detail Duplicate Receipt - PDF
│        │    └── Get - Download Detail Duplicate Receipt - Excel
│        ├── Buy or Purchase Returns/
│        │    ├── Get - Combobox Purchases
│        │    ├── Get - Combobox Items (from selected Purchase)
│        │    ├── Post - Create Buy Return
│        │    ├── Get - All Buy Returns
│        │    ├── Get - Cetak / Print Struk Retur Pembelian (Buy Return)
│        │    ├── Get - Download Buy Return - PDF
│        │    ├── Get - Download Buy Return - Excel
│        │    ├── Get - Download Detail Buy Return - PDF
│        │    └── Get - Download Detail Buy Return - Excel
│        ├── Sale Returns/
│        │    ├── Get - Combobox Sales
│        │    ├── Get - Combobox Items (from selected Sale)
│        │    ├── Post - Create Sale Return
│        │    ├── Get - All Sale Returns
│        │    ├── Get - Cetak / Print Struk Retur Penjualan (Sale Return)
│        │    ├── Get - Download Sale Return - PDF
│        │    ├── Get - Download Sale Return - Excel
│        │    ├── Get - Download Detail Sale Return - PDF
│        │    └── Get - Download Detail Sale Return - Excel
│        ├── Expenses/
│        │    ├── Get - All Expenses
│        │    ├── Post - Create Expense by ExpensID
│        │    ├── Put - Update Expense by ExpensID
│        │    ├── Delete - Delete Expense by ExpensID
│        │    ├── Get - Download Expenses - PDF
│        │    └── Get - Download Expenses - Excel
│        └── Another Incomes/
│             ├── Get - All Another Incomes
│             ├── Post - Create Another Income by AnotherIncomeID
│             ├── Put - Update Another Income by AnotherIncomeID
│             ├── Delete - Delete Another Income by AnotherIncomeID
│             ├── Get - Download Another Incomes - PDF
│             └── Get - Download Another Incomes - Excel
├── Audit & Finances/
│        ├── First Stocks/
│        │    ├── Get - All First Stocks
│        │    ├── Post - Create First Stock
│        │    ├── Delete - Delete First Stock by FirstStockID
│        │    ├── Get - Cetak / Print First Stock by FirstStockID
│        │    ├── Get - Download First Stock PDF
│        │    ├── Get - Download First Stock Excel
│        │    └── Get - Download Detail First Stock PDF by FirstStockID
│        └── Opnames/
│             ├── Mobile/
│             ├── Desktop/
│             ├── Get - Download Opname PDF
│             ├── Get - Download Opname Excel
│             └── Get - Download Detail Opname by OpnameID
├── Reports
│        ├── Get Neraca Saldo
│        ├── Provit by Mounth
│        ├── Get All Asset
│        ├── Download Daily Assets Excel
│        └── Download Neraca Saldo Excel
└── Complements/
         ├── Dashboard/
         │    ├── Get - Monthly Profit Report - Main Chart
         │    ├── Get - Provit Today - Second Line
         │    ├── Get - Weekly Profit Report
         │    ├── Get - Today Report - Total Provit, Trano, ABV
         │    ├── Get - Fast Moving Products
         │    ├── Get - Slow Moving Products
         │    ├── Get - Near Expired
         │    ├── Get - Download Fast Moving Products Excel
         │    ├── Get - Download Slow Moving Products Excel
         │    └── Get - Download Near Expired
         ├── Menus
         └── Test
```

---

## Struktur API saat ini di repo

Berikut endpoint yang **sudah ada di rewrite saat ini**.

### Health
- `GET /health`

### Auth
- `POST /api/login`
- `GET /api/list_branches`
- `GET /api/menus`
- `POST /api/set_branch`
- `GET /api/profile`
- `POST /api/logout`

### Branches
- `GET /api/branches`
- `GET /api/branches/:id`

### User Branches
- `GET /api/user-branches`
- `GET /api/user-branches/:user_id/:branch_id`

### User Management
- `GET /api/users`
- `GET /api/detail-users/:id`

### Products
- `POST /api/products`
- `GET /api/sales-products-combo`
- `GET /api/purchase-products-combo`
- `GET /api/cmb-product-opname`

### Purchases
- `POST /api/purchases`

### Sales
- `POST /api/sales`

### Opnames
- `POST /api/opnames`
- `GET /api/opnames/:id`
- `POST /api/opname-items`
- `POST /api/opname-items-all`

---

## Catatan detail endpoint yang sudah disempurnakan

### Auth flow
Sudah mendukung alur 2-step login:
1. `POST /api/login`
2. `GET /api/list_branches`
3. `POST /api/set_branch`
4. endpoint lanjutan memakai token branch-context

### Menus
- saat ini masih memakai `menus.json`
- belum dipindahkan ke tabel database
- filtering menu berdasarkan `user_role` dari token

### Branches list
`GET /api/branches` saat ini sudah mendukung:
- `search`
- `page`
- `limit`
- `meta` pagination

### Users list
`GET /api/users` saat ini sudah mendukung:
- `search`
- `page`
- `limit`
- `meta` pagination

### Sale safety improvement
Implementasi sale di rewrite **tidak mewarisi bug legacy** yang terlalu percaya pada `price` dari client. Tujuannya supaya pricing lebih aman dan dikontrol server.

---

## Contoh response contract

Project ini memakai response contract umum seperti:

### Success
```json
{
  "status": "success",
  "message": "...",
  "data": {}
}
```

### Success with meta
```json
{
  "status": "success",
  "message": "...",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_data": 100,
    "last_page": 10
  }
}
```

### Error
```json
{
  "status": "error",
  "message": "...",
  "error": "..."
}
```

---

## Testing & validasi runtime

### Build validation

```bash
go build ./...
```

### Menjalankan smoke test
Script yang tersedia:
- `scripts/smoke_test.sh`
- `scripts/auth_smoke_test.sh`

Contoh:

```bash
cd fiber-apotek-clean && \
BASE_URL="http://127.0.0.1:1113" \
USERNAME="your-username" \
PASSWORD="your-password" \
BRANCH_ID="your-branch-id" \
./scripts/auth_smoke_test.sh
```

### Catatan runtime penting
Dalam praktik sebelumnya, validasi runtime untuk project ini lebih konsisten jika app dijalankan dari **Terminal GUI**, bukan dari runtime exec sandbox, karena environment koneksi DB dev bisa berbeda.

Lihat juga:
- `docs/runtime-smoke-test.md`

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
