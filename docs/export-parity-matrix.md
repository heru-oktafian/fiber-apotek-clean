# Export Parity Matrix

Dokumen ini memetakan endpoint **download Excel/PDF** dari repo legacy `fiber-apotek` agar implementasi di `fiber-apotek-clean` bisa dikerjakan terstruktur.

## Prinsip
- Parity export tidak boleh dianggap bonus, karena ini bagian operasional utama.
- Implementasi sebaiknya dibuat **shared export layer** dulu, bukan handler per domain yang liar.
- Prioritas awal: domain yang CRUD-nya di clean repo sudah ada.

---

## 1. Endpoint export di repo legacy

### Masters
- `GET /api/products/excel`
- `GET /api/products/pdf`
- `GET /api/product-label/:id` (PDF label)
- `GET /api/units/excel`
- `GET /api/units/pdf`
- `GET /api/product-categories/excel`
- `GET /api/product-categories/pdf`
- `GET /api/unit-conversions/excel`
- `GET /api/unit-conversions/pdf`
- `GET /api/suppliers/excel`
- `GET /api/suppliers/pdf`
- `GET /api/supplier-categories/excel`
- `GET /api/supplier-categories/pdf`
- `GET /api/member-categories/excel`
- `GET /api/member-categories/pdf`
- `GET /api/members/excel`
- `GET /api/members/pdf`

### Audits
- `GET /api/first-stocks/excel`
- `GET /api/first-stocks/pdf`
- `GET /api/opnames/excel`
- `GET /api/opnames/pdf`
- `GET /api/first-stock-items/excel`
- `GET /api/first-stock-items/pdf`
- `GET /api/opname-items/excel`
- `GET /api/opname-items/pdf`

### Transactions
- `GET /api/purchases/excel`
- `GET /api/purchases/pdf`
- `GET /api/purchase-items/excel`
- `GET /api/purchase-items/pdf`
- `GET /api/sales/excel`
- `GET /api/sales/pdf`
- `GET /api/sale-items/excel`
- `GET /api/sale-items/pdf`
- `GET /api/duplicate-receipts/excel`
- `GET /api/duplicate-receipts/pdf`
- `GET /api/duplicate-receipt-items/excel`
- `GET /api/duplicate-receipt-items/pdf`
- `GET /api/buy-returns/excel`
- `GET /api/buy-returns/pdf`
- `GET /api/buy-return-items/excel`
- `GET /api/buy-return-items/pdf`
- `GET /api/sale-returns/excel`
- `GET /api/sale-returns/pdf`
- `GET /api/sale-return-items/excel`
- `GET /api/sale-return-items/pdf`
- `GET /api/expenses/excel`
- `GET /api/expenses/pdf`
- `GET /api/another-incomes/excel`
- `GET /api/another-incomes/pdf`

### Systems / Reports tambahan
- `GET /api/daily-assets/excel`
- `GET /api/defectas/excel`
- `GET /api/dashboard/neared-report/excel`
- `GET /api/dashboard/top-selling-report/excel`
- `GET /api/dashboard/least-selling-report/excel`
- `GET /api/reports/neraca-saldo/excel`

---

## 2. Status parity terhadap repo clean saat ini

### Sudah punya CRUD utama di clean, jadi siap dipasangi export dulu
#### Masters
- Products
- Units
- Product Categories
- Suppliers
- Supplier Categories
- Member Categories

#### Transactions
- Purchases
- Purchase Items
- Sales
- Sale Items

#### Audit
- Opnames (header + item flow baseline sudah ada)

### Belum siap parity penuh karena domain inti belum selesai di clean
- Members
- Unit Conversions
- First Stocks
- Duplicate Receipts
- Buy Returns
- Sale Returns
- Expenses
- Another Incomes
- Daily Assets
- Defecta
- Neared Report
- Top Selling Report
- Least Selling Report
- Neraca Saldo
- Product Label PDF

---

## 3. Rekomendasi implementasi bertahap

### Fase A, fondasi export shared
Buat layer bersama untuk:
- file naming
- content type
- content disposition
- helper download response
- query filter parsing (`month`, `search`, `id`, dll)
- excel builder wrapper
- pdf builder wrapper

### Fase B, export untuk domain yang CRUD-nya sudah hidup
Urutan yang direkomendasikan:
1. Products
2. Units
3. Product Categories
4. Suppliers
5. Supplier Categories
6. Member Categories
7. Purchases
8. Purchase Items
9. Sales
10. Sale Items
11. Opnames
12. Opname Items

### Fase C, lanjut domain yang belum selesai
Setelah domain inti selesai, baru tambahkan export untuk:
- members
- first stocks
- returns
- expenses
- another incomes
- dashboard/reports

---

## 4. Keputusan teknis yang direkomendasikan
- Jangan copy-paste mentah semua service export legacy.
- Ambil **kontrak endpoint + bentuk output** dari legacy, tapi implementasi ulang dengan struktur clean.
- Untuk Excel, gunakan builder shared berbasis `excelize`.
- Untuk PDF, buat renderer shared yang konsisten untuk tabel, title, summary, dan metadata filter.
- README dan `docs/api-implemented-endpoints.md` harus ikut menandai endpoint export yang sudah benar-benar hidup.

---

## 5. Next action yang paling masuk akal
Setelah audit ini, langkah terbaik adalah:
1. desain shared export foundation di repo clean
2. implement export untuk domain yang sudah complete CRUD dulu
3. baru sinkronkan dokumentasi implemented export endpoint per milestone
