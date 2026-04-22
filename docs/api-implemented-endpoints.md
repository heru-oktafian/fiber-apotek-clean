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

### POST `/api/products`
**Header:**
```http
Content-Type: application/json
Authorization: Bearer <TOKEN_2>
```

**Body:**
- mengikuti struktur create product yang dipakai usecase saat ini

**Catatan:**
- untuk testing nyata, sesuaikan field dengan schema product di DB dev

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

## 7. Purchases

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

## 8. Sales

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

## 9. Opnames

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
