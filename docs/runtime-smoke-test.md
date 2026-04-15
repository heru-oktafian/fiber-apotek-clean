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

## 3. Set branch

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
- `.data` berisi token branch baru

---

## 4. Combo endpoints

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

## 5. Purchase create

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

## 6. Sale create

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

## 7. Opname header

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

## 8. Opname item

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

## 9. Opname items all

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

## 10. Opname detail

```bash
curl -s "$BASE_URL/api/opnames/$OPNAME_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Expected:
- HTTP 200
- header + items tampil
- `total_opname` konsisten dengan data item

---

## 11. Logout

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

## Known gaps saat ini

- purchase masih menerima `price` dari request
- response contract belum 1:1 legacy
- branch ownership guard belum lengkap di seluruh endpoint opname
- belum ada delete/update item endpoint di rewrite
- belum ada report sync seperti legacy
