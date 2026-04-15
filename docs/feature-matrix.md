# Feature Matrix - fiber-apotek-clean

Dokumen ini dipakai untuk melacak status fitur rewrite dibanding repo lama.

## Legend
- `done-baseline` = sudah ada baseline awal, belum full parity
- `partial` = ada sebagian, tapi masih banyak gap
- `not-started` = belum dimulai

| Domain | Legacy Route Group | Status | Catatan | File utama saat ini | Prioritas |
|---|---|---|---|---|---|
| Auth Core | `/api/login`, `/api/logout`, `/api/set_branch` | done-baseline | Login, set branch, logout sudah ada | `internal/usecase/auth/service.go` | selesai baseline |
| Auth Support | `/api/profile`, `/api/menus`, `/api/list_branches` | not-started | Penting untuk parity alur login penuh | belum ada | tinggi |
| Users | `/api/users` | not-started | Domain system core belum dibangun | belum ada | menengah |
| Branches | `/api/branches` | not-started | Domain system core belum dibangun | belum ada | menengah |
| User Branches | `/api/user-branches` | not-started | Penting untuk relasi user-branch | belum ada | menengah |
| Product Categories | `/api/product-categories` | not-started | Category master belum dibangun | belum ada | menengah |
| Supplier Categories | `/api/supplier-categories` | not-started | Category master belum dibangun | belum ada | menengah |
| Member Categories | `/api/member-categories` | not-started | Category master belum dibangun | belum ada | menengah |
| Products | `/api/products`, combo routes | partial | Create + combo ada, CRUD penuh belum | `internal/usecase/product/service.go` | menengah |
| Units | `/api/units` | not-started | Master unit belum dibangun | belum ada | menengah |
| Unit Conversions | `/api/unit-conversions` | not-started | Penting untuk purchase domain parity | belum ada | menengah |
| Suppliers | `/api/suppliers`, `/api/suppliers-combo` | not-started | Penting untuk domain purchase | belum ada | sangat tinggi |
| Members | `/api/members` | not-started | Dibutuhkan untuk parity sale lebih dalam | belum ada | menengah |
| Purchase | `/api/purchases`, `purchase-items` | partial | Baru create purchase, belum list/detail/update/delete | `internal/usecase/purchase/service.go` | tinggi |
| Buy Returns / Purchase Returns | `/api/buy-returns` | not-started | Dekat dengan purchase domain | belum ada | tinggi |
| Sale | `/api/sales`, `sale-items` | partial | Baru create sale, belum list/detail/update/delete | `internal/usecase/sale/service.go` | tinggi |
| Sale Returns | `/api/sale-returns` | not-started | Dekat dengan sale domain | belum ada | tinggi |
| Duplicate Receipts | `/api/duplicate-receipts` | not-started | Fitur turunan, bukan prioritas pertama | belum ada | menengah |
| Opname | `/api/opnames`, `opname-items` | partial | Header, item, detail, items-all ada; belum full parity | `internal/usecase/opname/service.go` | tinggi |
| First Stocks | `/api/first-stocks`, `first-stock-items` | not-started | Mempengaruhi stock awal, tapi secara bisnis dianggap nol | belum ada | sangat tinggi |
| Defectas | `/api/sys-defectas`, `sys-defecta-items` | not-started | Domain kontrol stok dan issue barang | belum ada | menengah |
| Another Incomes | `/api/another-incomes` | not-started | Pendapatan tambahan di luar jual beli | belum ada | tinggi |
| Expenses | `/api/expenses` | not-started | Pengeluaran operasional apotek | belum ada | tinggi |
| Daily Assets | `/api/daily_asset` | not-started | Operasional finansial / asset harian | belum ada | menengah |
| Dashboard | `/api/dashboard/*` | not-started | Belum disentuh | belum ada | rendah awal |
| Reports | `/api/report/*` | not-started | Belum disentuh | belum ada | rendah awal |
| Export Excel / PDF | export routes | not-started | Wajib dipikirkan untuk tiap fitur utama | belum ada | setelah core domain |
| Mobile Opname Support | `mobile-opnames*` | not-started | Endpoint support mobile belum disentuh | belum ada | rendah awal |

---

## Definisi bisnis penting

### First Stocks
- untuk memasukkan stok awal sebelum sistem digunakan
- memengaruhi `products.stock`
- secara bisnis dianggap nol, bukan pembelian dan bukan pengeluaran

### Another Incomes
- untuk mencatat pendapatan tambahan di luar jual beli utama

### Expenses
- untuk mencatat pengeluaran apotek, baik terencana maupun aksidental

### Export
- jangan dilupakan
- tiap domain/fungsi utama nantinya perlu dukungan export Excel dan PDF

---

## Checklist parity final

### 1. Auth
- login
- list_branches
- set_branch
- logout
- profile
- menus

### 2. Systems / Core
- users
- branches
- user_branches

### 3. Operational Finance
- another_incomes
- expenses
- daily_assets

### 4. Dashboard / Reports
- dashboard
- reports

### 5. Master Categories
- member_categories
- supplier_categories
- product_categories

### 6. Masters
- units
- unit_conversions
- products
- suppliers
- members

### 7. Transactions
- purchases
- buy_returns
- sales
- duplicate_receipts
- sale_returns

### 8. Audit & Stock Control
- opnames
- first_stocks
- defectas

### 9. Exports
- export Excel
- export PDF

### 10. Mobile / Support endpoints
- mobile opnames
- mobile opname item details
- mobile opname glimpses

---

## Urutan kerja yang direkomendasikan

1. Suppliers
2. First Stocks
3. Expenses
4. Another Incomes
5. Buy Returns / Purchase Returns
6. Sale Returns
7. Duplicate Receipts
8. Export Excel/PDF per domain
9. Reports / Dashboard / fitur pelengkap
