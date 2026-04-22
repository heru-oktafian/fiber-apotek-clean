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
- `POST /api/products`
- `GET /api/sales-products-combo`
- `GET /api/purchase-products-combo`
- `GET /api/cmb-product-opname`

**Next target:**
- `GET /api/products`
- `GET /api/products/:id`
- `PUT /api/products/:id`
- `DELETE /api/products/:id`
- combobox category/unit
- export label/pdf/excel

### Suppliers
**Implemented:**
- `GET /api/suppliers`
- `GET /api/suppliers/:id`
- `POST /api/suppliers`
- `PUT /api/suppliers/:id`
- `DELETE /api/suppliers/:id`
- `GET /api/suppliers-combo`

**Next target:**
- supplier categories
- combobox supplier category
- export PDF/Excel

### Units
**Implemented:**
- `GET /api/units`
- `GET /api/units/:id`
- `POST /api/units`
- `PUT /api/units/:id`
- `DELETE /api/units/:id`
- `GET /api/cmb-units`

**Next target:**
- validasi kontrak frontend jika nanti ada kebutuhan khusus
- export PDF/Excel bila domain ini butuh parity penuh

### Product Categories
**Implemented:**
- `GET /api/product-categories`
- `POST /api/product-categories`
- `GET /api/product-categories/:id`
- `PUT /api/product-categories/:id`
- `DELETE /api/product-categories/:id`
- `GET /api/product-categories-combo`

**Catatan:**
- ID product category mengikuti schema legacy, yaitu numeric auto increment

**Next target:**
- validasi kontrak frontend bila ada format khusus
- export PDF/Excel kalau nanti dibutuhkan parity penuh

### Supplier Categories
**Implemented:**
- `GET /api/supplier-categories`
- `POST /api/supplier-categories`
- `GET /api/supplier-categories/:id`
- `PUT /api/supplier-categories/:id`
- `DELETE /api/supplier-categories/:id`
- `GET /api/supplier-categories-combo`

**Catatan:**
- ID supplier category mengikuti schema legacy, yaitu numeric auto increment

**Next target:**
- validasi kontrak frontend bila perlu
- export PDF/Excel jika nanti dibutuhkan parity penuh

### Unit Conversions
**Implemented:**
- belum ada

**Next target:**
- bisa menyusul sebelum products full CRUD kalau frontend benar-benar butuh

---

## 3. Transactions

### Purchases
**Implemented:**
- `POST /api/purchases`

**Next target:**
- list purchases
- detail purchase
- update/delete purchase
- item-level CRUD
- export/print

### Sales
**Implemented:**
- `POST /api/sales`

**Catatan:**
- rewrite tidak mengikuti bug legacy yang terlalu percaya harga client

**Next target:**
- list sales
- detail sale
- update/delete sale
- item-level CRUD
- export/print

### Duplicate Receipts / Buy Returns / Sale Returns / Expenses / Another Incomes
**Implemented:**
- belum ada

**Next target:**
- mulai bertahap setelah purchases/sales lebih matang
- tetap ingat `another incomes` dan `expenses` adalah domain penting user

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

## Rekomendasi urutan kerja berikutnya

Urutan yang paling sehat saat ini:

1. **Member Categories**
2. **Products full CRUD**
3. **Purchases & Sales full CRUD**
4. **Returns / Another Incomes / Expenses**
5. **First Stocks / Reports / Export**

---

## Dokumen terkait

- `docs/api-implemented-endpoints.md`
- `docs/feature-matrix.md`
- `docs/project-map.md`
- `README.md`
