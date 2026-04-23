# Implemented API Endpoints

Dokumen ini adalah daftar endpoint yang **sudah dibuat di repo** dan menjadi referensi operasional untuk testing, integrasi frontend, dan pengecekan progress rewrite.

> Prinsip kerja dokumen ini:
> setiap milestone endpoint yang sudah usable harus ikut diperbarui di sini.

## Base URL

Contoh base URL dev yang saat ini dipakai:

```text
http://200.200.200.20:1113
```

## Diagram struktur API yang sudah implemented

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
│   ├── Buy Returns/
│   │   ├── Get - /api/buy-returns
│   │   ├── Get - /api/buy-returns/excel
│   │   ├── Get - /api/buy-returns/pdf
│   │   ├── Post - /api/buy-returns
│   │   ├── Get - /api/buy-returns/:id
│   │   ├── Get - /api/buy-return-items/excel
│   │   ├── Get - /api/buy-return-items/pdf
│   │   ├── Get - /api/cmb-purchases
│   │   └── Get - /api/cmb-prod-buy-returns
│   ├── Sale Returns/
│   │   ├── Get - /api/sale-returns
│   │   ├── Get - /api/sale-returns/excel
│   │   ├── Get - /api/sale-returns/pdf
│   │   ├── Post - /api/sale-returns
│   │   ├── Get - /api/sale-returns/:id
│   │   ├── Get - /api/sale-return-items/excel
│   │   ├── Get - /api/sale-return-items/pdf
│   │   ├── Get - /api/cmb-sales
│   │   └── Get - /api/cmb-prod-sale-returns
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
│       ├── Get - /api/duplicate-receipts-details
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

## Alur token

### TOKEN_1
Didapat dari:
- `POST /api/login`

Dipakai untuk:
- `GET /api/list_branches`
- `POST /api/set_branch`

### TOKEN_2
Didapat dari:
- `POST /api/set_branch`

Dipakai untuk:
- `GET /api/profile`
- `GET /api/menus`
- `POST /api/logout`
- endpoint system/master/transaction lain yang sudah aktif

---

## 1. Health

### GET `/health`
**Header:**
- tidak perlu auth

**Body:**
- tidak ada

**Tujuan:**
- cek API hidup

---

## 2. Auth

### POST `/api/login`
**Header:**
```http
Content-Type: application/json
```

**Body:**
```json
{
  "username": "vita_fauzi",
  "password": "Sigala1102"
}
```

**Output penting:**
- menghasilkan `TOKEN_1`

---

### GET `/api/list_branches`
**Header:**
```http
Authorization: Bearer <TOKEN_1>
```

**Body:**
- tidak ada

---

### POST `/api/set_branch`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_1>
```

**Body:**
```json
{
  "branch_id": "BRC250118132203"
}
```

**Output penting:**
- menghasilkan `TOKEN_2`

---

### GET `/api/profile`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Body:**
- tidak ada

---

### GET `/api/menus`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Body:**
- tidak ada

**Catatan:**
- saat ini menus masih JSON-based (`menus.json`)

---

### POST `/api/logout`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Body:**
- tidak ada

---

## 3. Branches

### GET `/api/branches`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

**Contoh:**
```http
GET /api/branches?page=1&limit=10&search=ziida
```

**Catatan:**
- sudah mendukung pagination + search
- response memakai `meta`

---

### GET `/api/branches/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID branch

**Contoh:**
```http
GET /api/branches/BRC250118132203
```

---

### POST `/api/branches`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "branch_name": "Cabang Testing API",
  "address": "Jl. Testing No. 1",
  "phone": "08123456789",
  "email": "testing@apotek.local",
  "branch_status": "active"
}
```

**Catatan:**
- field wajib minimal: `branch_name`
- `branch_status` default ke `active` jika kosong
- ID branch digenerate server dengan prefix `BRC`

---

### DELETE `/api/branches/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID branch

**Catatan:**
- delete akan ditolak jika branch masih dipakai di `user_branches`
- guard ini sengaja ditambahkan agar branch yang masih punya relasi user tidak bisa dihapus

---

## 4. User Branches

### GET `/api/user-branches`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Body:**
- tidak ada

---

### GET `/api/user-branches/:user_id/:branch_id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path params:**
- `user_id`
- `branch_id`

**Contoh:**
```http
GET /api/user-branches/USR250118132201/BRC250118132203
```

---

### POST `/api/user-branches`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body:**
```json
{
  "user_id": "USR250118132201",
  "branch_id": "BRC250118132203"
}
```

**Catatan:**
- dipakai untuk menambahkan branch ke user
- akan ditolak jika user atau branch tidak ditemukan
- akan ditolak jika relasi user-branch sudah ada

---

## 5. User Management

### GET `/api/users`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

**Contoh:**
```http
GET /api/users?page=1&limit=10&search=vita
```

**Catatan:**
- password tidak ikut keluar di response
- response memakai `meta`

---

### GET `/api/detail-users/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID user

**Contoh:**
```http
GET /api/detail-users/USR250118132201
```

**Catatan:**
- return `user` + `detail_branches`

---

### POST `/api/users`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "username": "tester_api",
  "name": "Tester API",
  "password": "123456",
  "user_role": "operator",
  "user_status": "active"
}
```

**Catatan:**
- field wajib: `username`, `name`, `password`, `user_role`
- `user_status` default ke `inactive` jika kosong
- password di-hash oleh server

---

### PUT `/api/users/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID user

**Body contoh:**
```json
{
  "name": "Tester API Update",
  "user_role": "cashier",
  "user_status": "inactive"
}
```

**Contoh ganti password:**
```json
{
  "password": "654321"
}
```

---

## 6. Products

### GET `/api/products`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### POST `/api/products`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "sku": "PRD-CUSTOM-001",
  "name": "Paracetamol 500mg",
  "alias": "Paracetamol",
  "description": "Obat penurun panas",
  "ingredient": "Paracetamol",
  "dosage": "3x1 sehari",
  "side_affection": "Mual ringan",
  "unit_id": "UNT250118132755",
  "purchase_price": 5000,
  "sales_price": 7000,
  "alternate_price": 6500,
  "expired_date": "2027-12-31T00:00:00Z",
  "product_category_id": 1
}
```

**Catatan:**
- jika `sku` kosong, server akan fallback ke ID product
- `stock` default awal tetap `0`
- branch mengikuti branch context dari token

---

### GET `/api/products/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID product

---

### PUT `/api/products/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "sku": "PRD-CUSTOM-001",
  "name": "Paracetamol 500mg Update",
  "alias": "Paracetamol Update",
  "description": "Obat penurun panas update",
  "ingredient": "Paracetamol",
  "dosage": "2x1 sehari",
  "side_affection": "Mengantuk",
  "unit_id": "UNT250118132755",
  "purchase_price": 5500,
  "sales_price": 7500,
  "alternate_price": 7000,
  "expired_date": "2028-12-31T00:00:00Z",
  "product_category_id": 1
}
```

---

### DELETE `/api/products/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID product

---

### GET `/api/products/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data products dalam format Excel

---

### GET `/api/products/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data products dalam format PDF

---

### GET `/api/sales-products-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

**Contoh:**
```http
GET /api/sales-products-combo?search=par
```

---

### GET `/api/purchase-products-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

**Contoh:**
```http
GET /api/purchase-products-combo?search=par
```

---

### GET `/api/cmb-product-opname`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

**Contoh:**
```http
GET /api/cmb-product-opname?search=par
```

---

## 7. Suppliers

### GET `/api/suppliers`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### GET `/api/suppliers/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data suppliers dalam format Excel

---

### GET `/api/suppliers/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data suppliers dalam format PDF

---

### GET `/api/suppliers/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID supplier

---

### POST `/api/suppliers`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Supplier Testing",
  "phone": "08123456789",
  "address": "Jl. Supplier No. 1",
  "pic": "Budi",
  "supplier_category_id": 1
}
```

**Catatan:**
- field wajib minimal: `name`, `supplier_category_id`
- branch mengikuti branch context dari token

---

### PUT `/api/suppliers/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Supplier Testing Update",
  "phone": "08129999999",
  "address": "Jl. Supplier Update",
  "pic": "Andi",
  "supplier_category_id": 1
}
```

---

### DELETE `/api/suppliers/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID supplier

---

### GET `/api/suppliers-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

---

## 8. Units

### GET `/api/units`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### GET `/api/units/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID unit

---

### POST `/api/units`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Strip"
}
```

**Catatan:**
- field wajib: `name`
- branch mengikuti branch context dari token

---

### PUT `/api/units/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Box"
}
```

---

### DELETE `/api/units/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID unit

---

### GET `/api/units/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data units dalam format Excel

---

### GET `/api/units/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data units dalam format PDF

---

### GET `/api/cmb-units`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

---

## 9. Product Categories

### GET `/api/product-categories`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### POST `/api/product-categories`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Tablet"
}
```

**Catatan:**
- field wajib: `name`
- ID category mengikuti schema legacy, yaitu auto increment numeric (`uint`)
- branch mengikuti branch context dari token

---

### GET `/api/product-categories/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data product categories dalam format Excel

---

### GET `/api/product-categories/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data product categories dalam format PDF

---

### GET `/api/product-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID product category

---

### PUT `/api/product-categories/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Kaplet"
}
```

---

### DELETE `/api/product-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID product category

---

### GET `/api/product-categories-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

---

## 10. Supplier Categories

### GET `/api/supplier-categories`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### POST `/api/supplier-categories`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Obat Pabrik"
}
```

**Catatan:**
- field wajib: `name`
- ID supplier category mengikuti schema legacy, yaitu auto increment numeric (`uint`)
- branch mengikuti branch context dari token

---

### GET `/api/supplier-categories/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data supplier categories dalam format Excel

---

### GET `/api/supplier-categories/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data supplier categories dalam format PDF

---

### GET `/api/supplier-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID supplier category

---

### PUT `/api/supplier-categories/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Distributor Utama"
}
```

---

### DELETE `/api/supplier-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID supplier category

---

### GET `/api/supplier-categories-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

## 11. Member Categories

### GET `/api/another-incomes`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)
- `page` (opsional, default `1`)
- `limit` (opsional, default `10`)

**Catatan:**
- list another incomes dengan dukungan filter pencarian dan bulan

---

### GET `/api/another-incomes/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/another-incomes/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### POST `/api/another-incomes`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "income_date": "2026-04-23",
  "description": "Pendapatan tambahan jasa racik",
  "total_income": 50000,
  "payment": "paid_by_cash"
}
```

---

### PUT `/api/another-incomes/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "income_date": "2026-04-23",
  "description": "Pendapatan tambahan update",
  "total_income": 60000,
  "payment": "paid_by_cash"
}
```

---

### DELETE `/api/another-incomes/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/expenses`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)
- `page` (opsional, default `1`)
- `limit` (opsional, default `10`)

**Catatan:**
- list expenses dengan dukungan filter pencarian dan bulan

---

### GET `/api/expenses/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/expenses/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### POST `/api/expenses`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "expense_date": "2026-04-23",
  "description": "Biaya operasional harian",
  "total_expense": 25000,
  "payment": "paid_by_cash"
}
```

---

### PUT `/api/expenses/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "expense_date": "2026-04-23",
  "description": "Biaya operasional update",
  "total_expense": 30000,
  "payment": "paid_by_cash"
}
```

---

### DELETE `/api/expenses/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/buy-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)
- `page` (opsional, default `1`)
- `limit` (opsional, default `10`)

---

### POST `/api/buy-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "buy_return": {
    "purchase_id": "PURxxxx",
    "return_date": "2026-04-23",
    "payment": "paid_by_cash"
  },
  "buy_return_items": [
    {
      "product_id": "PRDxxxx",
      "qty": 1,
      "expired_date": "2026-12-31"
    }
  ]
}
```

---

### GET `/api/buy-returns/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/sale-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)
- `page` (opsional, default `1`)
- `limit` (opsional, default `10`)

---

### POST `/api/sale-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "sale_return": {
    "sale_id": "SALxxxx",
    "return_date": "2026-04-23",
    "payment": "paid_by_cash"
  },
  "sale_return_items": [
    {
      "product_id": "PRDxxxx",
      "qty": 1,
      "expired_date": "2026-12-31"
    }
  ]
}
```

---

### GET `/api/sale-returns/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/cmb-purchases`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/cmb-prod-buy-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `purchase_id` (wajib)

---

### GET `/api/cmb-sales`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/cmb-prod-sale-returns`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `sale_id` (wajib)

---

### GET `/api/buy-returns/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/buy-returns/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/buy-return-items/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `buy_return_id` (wajib)

---

### GET `/api/buy-return-items/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `buy_return_id` (wajib)

---

### GET `/api/sale-returns/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/sale-returns/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/sale-return-items/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `sale_return_id` (wajib)

---

### GET `/api/sale-return-items/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `sale_return_id` (wajib)

---

### GET `/api/first-stocks`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `search` (opsional)
- `month` (opsional, format `YYYY-MM`)
- `page` (opsional, default `1`)
- `limit` (opsional, default `10`)

---

### GET `/api/first-stocks/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### GET `/api/first-stocks/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

---

### POST `/api/first-stocks`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "description": "Stok awal pembukaan sistem",
  "first_stock_date": "2026-04-23"
}
```

---

### PUT `/api/first-stocks/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

---

### DELETE `/api/first-stocks/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/first-stock-with-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/first-stock-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/first-stock-items/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `first_stock_id` (wajib)

---

### GET `/api/first-stock-items/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `first_stock_id` (wajib)

---

### POST `/api/first-stock-items`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

**Body:**
```json
{
  "first_stock_id": "FSTxxxx",
  "product_id": "PRDxxxx",
  "unit_id": "UNTxxxx",
  "qty": 10,
  "expired_date": "2026-12-31"
}
```

---

### PUT `/api/first-stock-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
Content-Type: application/json
```

---

### DELETE `/api/first-stock-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

---

### GET `/api/member-categories`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### GET `/api/member-categories/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data member categories dalam format Excel

---

### GET `/api/member-categories/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Catatan:**
- download data member categories dalam format PDF

---

### GET `/api/member-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID member category

---

### POST `/api/member-categories`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Gold",
  "points_conversion_rate": 500
}
```

**Catatan:**
- field wajib minimal: `name`
- field bisnis penting: `points_conversion_rate`
- ID member category mengikuti schema legacy, yaitu auto increment numeric (`uint`)
- branch mengikuti branch context dari token

---

### PUT `/api/member-categories/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "name": "Platinum",
  "points_conversion_rate": 1000
}
```

---

### DELETE `/api/member-categories/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = numeric ID member category

---

### GET `/api/member-categories-combo`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `search`

---

## 12. Purchases

### GET `/api/purchases`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### GET `/api/purchases/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

**Catatan:**
- download data purchase header dalam format Excel

---

### GET `/api/purchases/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

**Catatan:**
- download data purchase header dalam format PDF

---

### GET `/api/purchases/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID purchase

---

### POST `/api/purchases`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "purchase": {
    "supplier_id": "SPL250207144606",
    "purchase_date": "2026-04-15",
    "payment": "cash"
  },
  "purchase_items": [
    {
      "product_id": "PRD25050451578",
      "unit_id": "UNT250118132755",
      "price": 6433,
      "qty": 1,
      "expired_date": "2027-12-31"
    }
  ]
}
```

---

### PUT `/api/purchases/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "supplier_id": "SPL250207144606",
  "purchase_date": "2026-04-20",
  "payment": "cash"
}
```

---

### DELETE `/api/purchases/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID purchase

---

### GET `/api/purchase-items/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `purchase_id` (wajib)

**Catatan:**
- download detail item purchase dalam format Excel

---

### GET `/api/purchase-items/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `purchase_id` (wajib)

**Catatan:**
- download detail item purchase dalam format PDF

---

### GET `/api/purchase-items/all/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID purchase

---

### POST `/api/purchase-items`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "purchase_id": "PUR250423000001",
  "product_id": "PRD25050451578",
  "unit_id": "UNT250118132755",
  "price": 7000,
  "qty": 2,
  "expired_date": "2027-12-31"
}
```

---

### PUT `/api/purchase-items/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "product_id": "PRD25050451578",
  "unit_id": "UNT250118132755",
  "price": 7500,
  "qty": 3,
  "expired_date": "2027-12-31"
}
```

---

### DELETE `/api/purchase-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID purchase item

---

## 13. Sales

### GET `/api/sales`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query params opsional:**
- `page`
- `limit`
- `search`

---

### GET `/api/sales/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

**Catatan:**
- download data sales header dalam format Excel

---

### GET `/api/sales/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `month` (opsional, format `YYYY-MM`)

**Catatan:**
- download data sales header dalam format PDF

---

### GET `/api/sales/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID sale

---

### POST `/api/sales`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "sale": {
    "payment": "cash",
    "discount": 0
  },
  "sale_items": [
    {
      "product_id": "PRD25050451578",
      "price": 999999,
      "qty": 1
    }
  ]
}
```

**Catatan:**
- implementasi rewrite menjaga arah harga server-side lebih aman daripada legacy

---

### PUT `/api/sales/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "member_id": "MBR000001",
  "discount": 1000,
  "payment": "cash"
}
```

---

### DELETE `/api/sales/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID sale

---

### GET `/api/sale-items/excel`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `sale_id` (wajib)

**Catatan:**
- download detail item sale dalam format Excel

---

### GET `/api/sale-items/pdf`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query:**
- `sale_id` (wajib)

**Catatan:**
- download detail item sale dalam format PDF

---

### GET `/api/sale-items/all/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID sale

---

### POST `/api/sale-items`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "sale_id": "SAL250423000001",
  "product_id": "PRD25050451578",
  "qty": 2
}
```

**Catatan:**
- harga item tetap dihitung server-side dari `products.sales_price`

---

### PUT `/api/sale-items/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "product_id": "PRD25050451578",
  "qty": 3
}
```

**Catatan:**
- harga item tetap dihitung server-side dari `products.sales_price`

---

### DELETE `/api/sale-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID sale item

---

## 14. Duplicate Receipts

### GET `/api/duplicate-receipts`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Query param opsional:**
- `search`
- `page`
- `limit`
- `month`

**Catatan:**
- duplicate receipt diperlakukan sebagai sale-like transaction berbasis resep dokter
- list menampilkan header duplicate receipt aktif di branch berjalan

---

### POST `/api/duplicate-receipts`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "duplicate_receipt": {
    "member_id": "MBR250423000001",
    "description": "Copy resep dokter umum",
    "duplicate_receipt_date": "2026-04-23",
    "payment": "cash"
  },
  "items": [
    {
      "product_id": "PRD25050451578",
      "qty": 1
    }
  ]
}
```

**Catatan:**
- harga item dihitung server-side dari `products.sales_price`
- create header, item, stock movement, transaction report, dan daily profit sudah dibungkus transaction
- bila `member_id` kosong, service bisa fallback ke default member

---

### GET `/api/duplicate-receipts/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt

---

### PUT `/api/duplicate-receipts/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt

**Body contoh:**
```json
{
  "member_id": "MBR250423000001",
  "description": "Copy resep update",
  "payment": "cash"
}
```

**Catatan:**
- update saat ini fokus ke field header
- total dan profit dihitung ulang dari item yang sudah tersimpan

---

### DELETE `/api/duplicate-receipts/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt

**Catatan:**
- rollback stok item dan delete header/item/report sudah dibungkus transaction

---

### GET `/api/duplicate-receipts-items/all/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt

**Catatan:**
- mengembalikan item duplicate receipt lengkap dengan nama produk dan unit

---

### POST `/api/duplicate-receipts-items`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "duplicate_receipt_id": "DUR250423000001",
  "product_id": "PRD25050451578",
  "qty": 1
}
```

**Catatan:**
- harga dan subtotal item dihitung server-side
- create item ikut mengubah stok, total header, profit header, transaction report, dan daily profit
- create item sekarang sudah dibungkus transaction safety

---

### PUT `/api/duplicate-receipts-items/:id`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt item

**Body contoh:**
```json
{
  "product_id": "PRD25050451578",
  "qty": 2
}
```

**Catatan:**
- update item menghitung ulang delta stok dan recalculation total/profit header
- update item sekarang sudah dibungkus transaction safety

---

### DELETE `/api/duplicate-receipts-items/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID duplicate receipt item

**Catatan:**
- delete item melakukan rollback stok lalu sinkronkan ulang total/profit/report
- delete item sekarang sudah dibungkus transaction safety

---

## 15. Opnames

### POST `/api/opnames`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body:**
```json
{
  "description": "uji rewrite opname",
  "opname_date": "2026-04-15"
}
```

---

### GET `/api/opnames/:id`
**Header:**
```http
Authorization: Bearer <TOKEN_2>
```

**Path param:**
- `id` = ID opname

**Contoh:**
```http
GET /api/opnames/OPN739440AJFT4G
```

---

### POST `/api/opname-items`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body contoh:**
```json
{
  "opname_id": "OPN739440AJFT4G",
  "product_id": "PRD054724Q21ODS",
  "qty": 3,
  "price": 11400,
  "expired_date": "2027-12-31"
}
```

---

### POST `/api/opname-items-all`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body:**
```json
{
  "opname_id": "OPN739440AJFT4G"
}
```

---

## Ringkasan endpoint implemented saat ini

- `GET /health`
- `POST /api/login`
- `GET /api/list_branches`
- `GET /api/menus`
- `POST /api/set_branch`
- `GET /api/profile`
- `POST /api/logout`
- `GET /api/branches`
- `GET /api/branches/:id`
- `POST /api/branches`
- `DELETE /api/branches/:id`
- `GET /api/user-branches`
- `GET /api/user-branches/:user_id/:branch_id`
- `POST /api/user-branches`
- `GET /api/users`
- `GET /api/detail-users/:id`
- `POST /api/users`
- `PUT /api/users/:id`
- `GET /api/products`
- `GET /api/products/excel`
- `GET /api/products/pdf`
- `POST /api/products`
- `GET /api/products/:id`
- `PUT /api/products/:id`
- `DELETE /api/products/:id`
- `GET /api/sales-products-combo`
- `GET /api/purchase-products-combo`
- `GET /api/cmb-product-opname`
- `GET /api/suppliers`
- `GET /api/suppliers/excel`
- `GET /api/suppliers/pdf`
- `GET /api/suppliers/:id`
- `POST /api/suppliers`
- `PUT /api/suppliers/:id`
- `DELETE /api/suppliers/:id`
- `GET /api/suppliers-combo`
- `GET /api/units`
- `GET /api/units/excel`
- `GET /api/units/pdf`
- `GET /api/units/:id`
- `POST /api/units`
- `PUT /api/units/:id`
- `DELETE /api/units/:id`
- `GET /api/cmb-units`
- `GET /api/product-categories`
- `GET /api/product-categories/excel`
- `GET /api/product-categories/pdf`
- `POST /api/product-categories`
- `GET /api/product-categories/:id`
- `PUT /api/product-categories/:id`
- `DELETE /api/product-categories/:id`
- `GET /api/product-categories-combo`
- `GET /api/supplier-categories`
- `GET /api/supplier-categories/excel`
- `GET /api/supplier-categories/pdf`
- `POST /api/supplier-categories`
- `GET /api/supplier-categories/:id`
- `PUT /api/supplier-categories/:id`
- `DELETE /api/supplier-categories/:id`
- `GET /api/supplier-categories-combo`
- `GET /api/member-categories`
- `GET /api/member-categories/excel`
- `GET /api/member-categories/pdf`
- `GET /api/member-categories/:id`
- `POST /api/member-categories`
- `PUT /api/member-categories/:id`
- `DELETE /api/member-categories/:id`
- `GET /api/member-categories-combo`
- `GET /api/purchases`
- `GET /api/purchases/:id`
- `POST /api/purchases`
- `PUT /api/purchases/:id`
- `DELETE /api/purchases/:id`
- `GET /api/purchase-items/all/:id`
- `POST /api/purchase-items`
- `PUT /api/purchase-items/:id`
- `DELETE /api/purchase-items/:id`
- `GET /api/sales`
- `GET /api/sales/:id`
- `POST /api/sales`
- `PUT /api/sales/:id`
- `DELETE /api/sales/:id`
- `GET /api/sale-items/all/:id`
- `POST /api/sale-items`
- `PUT /api/sale-items/:id`
- `DELETE /api/sale-items/:id`
- `GET /api/duplicate-receipts`
- `GET /api/duplicate-receipts-details`
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
- `POST /api/opnames`
- `GET /api/opnames/:id`
- `POST /api/opname-items`
- `POST /api/opname-items-all`
