# Implemented vs Next Target

Dokumen ini merangkum kondisi **endpoint yang sudah hidup** di repo `fiber-apotek-clean` dibandingkan dengan **target domain berikutnya** yang perlu dikerjakan.

Tujuannya sederhana:
- memisahkan mana yang sudah implemented dan usable
- memperjelas milestone berikutnya
- membantu tracking parity tanpa harus membaca semua detail endpoint satu per satu

## Snapshot singkat

### Implemented sekarang
- Health
- Auth flow utama
- Branches read + write parsial
- User branches read
- User management read + write parsial
- Product create + combo endpoints
- Purchase create baseline
- Sale create baseline
- Opname baseline

### Next target yang paling dekat
1. Sys > User Management > adding branch by user ID
2. Masters > Suppliers
3. Masters > Units / Categories / Products full CRUD
4. Transactions full CRUD bertahap

---

## 1. Sys

### Auth
**Implemented:**
- `POST /api/login`
- `GET /api/list_branches`
- `POST /api/set_branch`
- `GET /api/profile`
- `GET /api/menus`
- `POST /api/logout`

**Next target:**
- tidak urgent untuk fase ini
- fokus auth saat ini lebih ke stabilitas kontrak

### Branches
**Implemented:**
- `GET /api/branches`
- `GET /api/branches/:id`
- `POST /api/branches`
- `DELETE /api/branches/:id`

**Catatan:**
- delete branch ditolak jika branch masih dipakai di `user_branches`

**Next target:**
- `PUT /api/branches/:id`
- verifikasi kontrak response terhadap frontend legacy jika diperlukan

### User Branches
**Implemented:**
- `GET /api/user-branches`
- `GET /api/user-branches/:user_id/:branch_id`
- `POST /api/user-branches`

**Next target:**
- kemungkinan delete relasi branch dari user
- jika dibutuhkan, tambahkan update flow relasi user-branch

### User Management
**Implemented:**
- `GET /api/users`
- `GET /api/detail-users/:id`
- `POST /api/users`
- `PUT /api/users/:id`

**Catatan:**
- password di-hash server-side
- response tidak membocorkan password
- add branch by user ID sekarang ditangani lewat `POST /api/user-branches`

**Next target:**
- delete user jika memang diperlukan parity
- rapikan flow user + user_branches agar cohesive

---

## 2. Masters

### Products
**Implemented:**
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

**Next target:**
- product label PDF
- combobox category/unit jika frontend memang butuh endpoint terpisah
- validasi runtime end-to-end untuk full CRUD product + export

### Suppliers
**Implemented:**
- `GET /api/suppliers`
- `GET /api/suppliers/excel`
- `GET /api/suppliers/pdf`
- `GET /api/suppliers/:id`
- `POST /api/suppliers`
- `PUT /api/suppliers/:id`
- `DELETE /api/suppliers/:id`
- `GET /api/suppliers-combo`

**Next target:**
- refine visual export bila ingin makin mirip legacy
- validasi runtime export end-to-end

### Units
**Implemented:**
- `GET /api/units`
- `GET /api/units/excel`
- `GET /api/units/pdf`
- `GET /api/units/:id`
- `POST /api/units`
- `PUT /api/units/:id`
- `DELETE /api/units/:id`
- `GET /api/cmb-units`

**Next target:**
- validasi kontrak frontend jika nanti ada kebutuhan khusus
- refine styling export agar makin dekat ke legacy bila diperlukan

### Product Categories
**Implemented:**
- `GET /api/product-categories`
- `GET /api/product-categories/excel`
- `GET /api/product-categories/pdf`
- `POST /api/product-categories`
- `GET /api/product-categories/:id`
- `PUT /api/product-categories/:id`
- `DELETE /api/product-categories/:id`
- `GET /api/product-categories-combo`

**Catatan:**
- ID product category mengikuti schema legacy, yaitu numeric auto increment

**Next target:**
- refine visual export bila ingin makin mirip legacy
- validasi runtime export end-to-end

### Supplier Categories
**Implemented:**
- `GET /api/supplier-categories`
- `GET /api/supplier-categories/excel`
- `GET /api/supplier-categories/pdf`
- `POST /api/supplier-categories`
- `GET /api/supplier-categories/:id`
- `PUT /api/supplier-categories/:id`
- `DELETE /api/supplier-categories/:id`
- `GET /api/supplier-categories-combo`

**Catatan:**
- ID supplier category mengikuti schema legacy, yaitu numeric auto increment

**Next target:**
- refine visual export bila ingin makin mirip legacy
- validasi runtime export end-to-end

### Unit Conversions
**Implemented:**
- belum ada

**Next target:**
- bisa menyusul sebelum products full CRUD kalau frontend benar-benar butuh

---

## 3. Transactions

### Purchases
**Implemented:**
- `GET /api/purchases`
- `GET /api/purchases/excel`
- `GET /api/purchases/pdf`
- `GET /api/purchases/:id`
- `POST /api/purchases`
- `PUT /api/purchases/:id`
- `DELETE /api/purchases/:id`
- `GET /api/purchase-items/excel`
- `GET /api/purchase-items/pdf`
- `GET /api/purchase-items/all/:id`
- `POST /api/purchase-items`
- `PUT /api/purchase-items/:id`
- `DELETE /api/purchase-items/:id`

**Next target:**
- refine visual export bila ingin makin mirip legacy
- validasi runtime end-to-end purchase flow

### Sales
**Implemented:**
- `GET /api/sales`
- `GET /api/sales/excel`
- `GET /api/sales/pdf`
- `GET /api/sales/:id`
- `POST /api/sales`
- `PUT /api/sales/:id`
- `DELETE /api/sales/:id`
- `GET /api/sale-items/excel`
- `GET /api/sale-items/pdf`
- `GET /api/sale-items/all/:id`
- `POST /api/sale-items`
- `PUT /api/sale-items/:id`
- `DELETE /api/sale-items/:id`

**Catatan:**
- rewrite tidak mengikuti bug legacy yang terlalu percaya harga client
- harga sale item tetap dihitung server-side

**Next target:**
- sales detail summary endpoint parity bila masih diperlukan
- refine visual export bila ingin makin mirip legacy
- validasi runtime end-to-end sales flow

### Another Incomes
**Implemented:**
- `GET /api/another-incomes`
- `POST /api/another-incomes`
- `PUT /api/another-incomes/:id`
- `DELETE /api/another-incomes/:id`

**Catatan:**
- list mendukung `search`, `page`, `limit`, dan `month`
- create/update/delete tersinkron ke `transaction_reports` dengan transaction type `income`

**Next target:**
- export Excel/PDF
- validasi runtime end-to-end

### Expenses
**Implemented:**
- `GET /api/expenses`
- `POST /api/expenses`
- `PUT /api/expenses/:id`
- `DELETE /api/expenses/:id`

**Catatan:**
- list mendukung `search`, `page`, `limit`, dan `month`
- create/update/delete tersinkron ke `transaction_reports` dengan transaction type `expense`

**Next target:**
- export Excel/PDF
- validasi runtime end-to-end

### Duplicate Receipts / Buy Returns / Sale Returns
**Implemented:**
- belum ada

**Next target:**
- mulai bertahap setelah `another_incomes` dan `expenses`

---

## 4. Audit & Finances

### Opnames
**Implemented:**
- `POST /api/opnames`
- `GET /api/opnames/:id`
- `POST /api/opname-items`
- `POST /api/opname-items-all`

**Next target:**
- list opname
- delete/update bila memang dibutuhkan parity
- dukungan mobile opname yang lebih eksplisit
- export PDF/Excel

### First Stocks
**Implemented:**
- belum ada

**Catatan bisnis penting:**
- first stock memengaruhi stok produk
- tetapi secara bisnis dianggap nol, bukan pembelian/pengeluaran biasa

**Next target:**
- implement sebagai domain sendiri, jangan disamakan dengan purchase biasa

---

## 5. Reports / Dashboard / Export

**Implemented:**
- belum ada laporan besar yang selesai di rewrite

**Next target:**
- reports domain
- dashboard
- export Excel/PDF di tiap domain penting

---

### Member Categories
**Implemented:**
- `GET /api/member-categories`
- `GET /api/member-categories/excel`
- `GET /api/member-categories/pdf`
- `GET /api/member-categories/:id`
- `POST /api/member-categories`
- `PUT /api/member-categories/:id`
- `DELETE /api/member-categories/:id`
- `GET /api/member-categories-combo`

**Catatan:**
- field bisnis penting: `points_conversion_rate`
- ID member category mengikuti schema legacy, yaitu numeric auto increment

**Next target:**
- refine visual export bila ingin makin mirip legacy
- validasi runtime export end-to-end

---

## Rekomendasi urutan kerja berikutnya

Urutan yang paling sehat saat ini:

1. **Returns / Another Incomes / Expenses**
2. **First Stocks / Reports / Export lanjutan**

---

## Dokumen terkait

- `docs/api-implemented-endpoints.md`
- `docs/feature-matrix.md`
- `docs/project-map.md`
- `README.md`
