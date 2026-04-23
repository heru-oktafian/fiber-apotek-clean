# Runtime Smoke Test

Panduan ini untuk validasi cepat baseline `fiber-apotek-clean` terhadap environment dev.

## Prasyarat

- Jalankan app dari **Terminal GUI**, bukan runtime exec, jika koneksi DB dev sensitif konteks.
- Pastikan `.env` mengarah ke DB dev, bukan production.
- Jika repo lama masih jalan di port `1112`, set `SERVER_PORT=1113` untuk repo baru.

## Menjalankan app

```bash
cd ~/Projects/Go/fiber-apotek-clean
go run ./cmd/api
```

## Variabel bantu shell

```bash
BASE_URL="http://127.0.0.1:1113"
USERNAME="<isi-username-dev>"
PASSWORD="<isi-password-dev>"
BRANCH_ID="BRC250118132203"
SUPPLIER_ID="SPL250207144606"
PURCHASE_PRODUCT_ID="PRD25050451578"
OPNAME_PRODUCT_ID="PRD054724Q21ODS"
```

Jika app jalan di port lama, ganti `1113` menjadi `1112`.

---

## 1. Health

```bash
curl -s "$BASE_URL/health"
```

Expected:
- HTTP 200
- body mengandung `message: ok`

---

## 2. Login

```bash
LOGIN_JSON=$(curl -s -X POST "$BASE_URL/api/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

echo "$LOGIN_JSON"
```

Ambil token login:

```bash
LOGIN_TOKEN=$(echo "$LOGIN_JSON" | jq -r '.data')
echo "$LOGIN_TOKEN"
```

Expected:
- HTTP 200
- `.data` berisi token string

---

## 3. List branches

```bash
BRANCH_LIST_JSON=$(curl -s "$BASE_URL/api/list_branches" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "$BRANCH_LIST_JSON" | jq
```

Expected:
- HTTP 200
- `.status = "success"`
- `.message = "User Branch found"`
- `.data` berisi array branch user

---

## 4. Set branch

```bash
BRANCH_JSON=$(curl -s -X POST "$BASE_URL/api/set_branch" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $LOGIN_TOKEN" \
  -d "{\"branch_id\":\"$BRANCH_ID\"}")

echo "$BRANCH_JSON"
```

Ambil token branch:

```bash
TOKEN=$(echo "$BRANCH_JSON" | jq -r '.data')
echo "$TOKEN"
```

Expected:
- HTTP 200
- `.status = "success"`
- `.data` berisi token branch baru
- token hasil set branch mengandung `branch_id`

---

## 5. Profile

```bash
PROFILE_JSON=$(curl -s "$BASE_URL/api/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "$PROFILE_JSON" | jq
```

Expected:
- HTTP 200
- `.status = "success"`
- `.message` berbentuk `Otoritas : <role>`
- `.data.branch_id` sesuai branch yang dipilih

---

## 6. Combo endpoints

### Sales combo
```bash
curl -s "$BASE_URL/api/sales-products-combo?search=par" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Purchase combo
```bash
curl -s "$BASE_URL/api/purchase-products-combo?search=par" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Opname combo
```bash
curl -s "$BASE_URL/api/cmb-product-opname?search=par" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200
- data tidak kosong
- item relevan dengan branch dev

---

## 7. Purchase create

```bash
PURCHASE_JSON=$(curl -s -X POST "$BASE_URL/api/purchases" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"purchase\": {
      \"supplier_id\": \"$SUPPLIER_ID\",
      \"purchase_date\": \"2026-04-15\",
      \"payment\": \"cash\"
    },
    \"purchase_items\": [
      {
        \"product_id\": \"$PURCHASE_PRODUCT_ID\",
        \"unit_id\": \"UNT250118132755\",
        \"price\": 6433,
        \"qty\": 1,
        \"expired_date\": \"2027-12-31\"
      }
    ]
  }")

echo "$PURCHASE_JSON" | jq
```

Expected:
- HTTP 200
- `purchase` dan `purchase_items` terbentuk
- stock produk naik

Catatan:
- `unit_id` di atas perlu disesuaikan bila unit produk dev berbeda.

---

## 8. Sale create

```bash
SALE_JSON=$(curl -s -X POST "$BASE_URL/api/sales" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"sale\": {
      \"payment\": \"cash\",
      \"discount\": 0
    },
    \"sale_items\": [
      {
        \"product_id\": \"$PURCHASE_PRODUCT_ID\",
        \"price\": 999999,
        \"qty\": 1
      }
    ]
  }")

echo "$SALE_JSON" | jq
```

Expected:
- HTTP 200 jika stock cukup
- harga final item mengikuti **harga server**, bukan `999999`
- stock turun

Ini checkpoint penting karena rewrite baru harus lebih aman daripada legacy.

---

## 9. Opname header

```bash
OPNAME_JSON=$(curl -s -X POST "$BASE_URL/api/opnames" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "description": "uji rewrite opname",
    "opname_date": "2026-04-15"
  }')

echo "$OPNAME_JSON" | jq
```

Ambil ID opname:

```bash
OPNAME_ID=$(echo "$OPNAME_JSON" | jq -r '.data.id')
echo "$OPNAME_ID"
```

---

## 10. Duplicate Receipts Batch 1

### Variabel bantu tambahan

```bash
PRODUCT_ID="PRD25050451578"
TEST_DATE="2026-04-23"
```

### Opsi cepat: pakai script siap jalan

```bash
cd ~/Projects/Go/fiber-apotek-clean
chmod +x scripts/duplicate_receipt_smoke_test.sh
USERNAME="$USERNAME" PASSWORD="$PASSWORD" BRANCH_ID="$BRANCH_ID" PRODUCT_ID="$PRODUCT_ID" ./scripts/duplicate_receipt_smoke_test.sh
```

Script di atas akan menjalankan urutan:
- health
- login
- set branch
- cek stock produk awal
- create duplicate receipt
- list duplicate receipts
- get detail duplicate receipt
- list duplicate receipt items
- create item duplicate receipt
- update item duplicate receipt
- update header duplicate receipt
- delete item duplicate receipt
- delete duplicate receipt
- cek stock produk setelah delete
- logout

### Opsi manual: create duplicate receipt

```bash
DUPLICATE_RECEIPT_JSON=$(curl -s -X POST "$BASE_URL/api/duplicate-receipts" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"duplicate_receipt\": {
      \"member_id\": \"\",
      \"description\": \"Copy resep dokter smoke test\",
      \"duplicate_receipt_date\": \"$TEST_DATE\",
      \"payment\": \"cash\"
    },
    \"items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"qty\": 1
      }
    ]
  }")

echo "$DUPLICATE_RECEIPT_JSON" | jq
```

Ambil ID duplicate receipt:

```bash
DUPLICATE_RECEIPT_ID=$(echo "$DUPLICATE_RECEIPT_JSON" | jq -r '.data.id')
echo "$DUPLICATE_RECEIPT_ID"
```

Expected:
- HTTP 200
- header duplicate receipt terbentuk
- stock produk turun
- total memakai harga server-side, bukan harga client

### List duplicate receipts

```bash
curl -s "$BASE_URL/api/duplicate-receipts?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Get duplicate receipt detail

```bash
curl -s "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### List duplicate receipt items

```bash
curl -s "$BASE_URL/api/duplicate-receipts-items/all/$DUPLICATE_RECEIPT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Create duplicate receipt item

```bash
curl -s -X POST "$BASE_URL/api/duplicate-receipts-items" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"duplicate_receipt_id\": \"$DUPLICATE_RECEIPT_ID\",
    \"product_id\": \"$PRODUCT_ID\",
    \"qty\": 1
  }" | jq
```

Expected:
- HTTP 200
- item duplicate receipt bertambah atau merge ke item produk yang sama
- stock turun
- total/profit header ikut naik

### Update duplicate receipt item

```bash
curl -s -X PUT "$BASE_URL/api/duplicate-receipts-items/$ITEM_ID" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"product_id\": \"$PRODUCT_ID\",
    \"qty\": 2
  }" | jq
```

Expected:
- HTTP 200
- qty item berubah
- delta stock dan recalculation total/profit ikut sinkron

### Update header duplicate receipt

```bash
curl -s -X PUT "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "member_id": "",
    "description": "Copy resep dokter smoke test update",
    "payment": "cash"
  }' | jq
```

Expected:
- HTTP 200
- field header berubah
- total/profit tetap konsisten dari item tersimpan

### Delete duplicate receipt item

```bash
curl -s -X DELETE "$BASE_URL/api/duplicate-receipts-items/$ITEM_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200
- item terhapus
- stock rollback sesuai qty item
- total/profit header ikut turun

### Delete duplicate receipt

```bash
curl -s -X DELETE "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200
- header duplicate receipt terhapus
- stock produk rollback

## 11. Opname item

```bash
OPNAME_ITEM_JSON=$(curl -s -X POST "$BASE_URL/api/opname-items" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"opname_id\": \"$OPNAME_ID\",
    \"product_id\": \"$OPNAME_PRODUCT_ID\",
    \"qty\": 3,
    \"price\": 11400,
    \"expired_date\": \"2027-12-31\"
  }")

echo "$OPNAME_ITEM_JSON" | jq
```

Expected:
- HTTP 200
- item tersimpan
- `qty_exist`, `sub_total_exist`, `sub_total` masuk akal
- product stock berubah ke qty opname

---

## 11. Opname items all

```bash
curl -s -X POST "$BASE_URL/api/opname-items-all" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"opname_id\":\"$OPNAME_ID\"}" | jq
```

Expected:
- HTTP 200
- semua item untuk header tampil

---

## 12. Opname detail

```bash
curl -s "$BASE_URL/api/opnames/$OPNAME_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200
- header + items tampil
- `total_opname` konsisten dengan data item

---

## 13. Logout

```bash
curl -s -X POST "$BASE_URL/api/logout" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200

---

## Checklist observasi

Catat hasil berikut saat smoke test:

- apakah login dan set branch konsisten
- apakah combo endpoint mengembalikan data branch yang benar
- apakah `purchase` menambah stock
- apakah `sale` memakai harga server, bukan harga client
- apakah `opname-items` update stock dan total header dengan benar
- apakah response contract cukup dekat dengan legacy atau masih perlu adaptasi
- apakah ada panic, 500, atau mismatch field penting

## Verifikasi Auth fase 1 pada 2026-04-16

Sudah tervalidasi manual via Postman terhadap app yang berjalan dari Terminal GUI:

- `POST /api/login` ✅
- `GET /api/list_branches` ✅
- `POST /api/set_branch` ✅
- `GET /api/profile` ✅
- `POST /api/logout` ✅
- `GET /api/menus` ✅

Catatan hasil verifikasi:
- response contract sudah konsisten memiliki `status`, `message`, dan `data`
- `set_branch` terbukti menghasilkan token kedua yang berisi konteks branch
- mismatch awal pada `default_member` ternyata berasal dari data DB dev yang kosong, bukan bug endpoint
- beberapa field seperti `branch_name`, `sipa_name`, dan `member_name` tetap bergantung pada isi data DB dev saat pengujian
- endpoint `menus` berhasil memfilter menu berdasarkan `user_role` dari token branch-context
- response `menus` sangat panjang sehingga validasi lebih cocok dilakukan via Postman daripada dikirim penuh ke Telegram
- runtime validation tetap mengandalkan Terminal GUI sebagai sumber kebenaran

## Known gaps saat ini

- purchase masih menerima `price` dari request
- beberapa nilai response masih bisa berbeda dari contoh legacy bila data DB dev berbeda
- branch ownership guard belum lengkap di seluruh endpoint opname
- belum ada delete/update item endpoint di rewrite
- belum ada report sync seperti legacy
