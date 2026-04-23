#!/usr/bin/env bash
set -euo pipefail

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "ERROR: command '$1' tidak ditemukan" >&2
    exit 1
  }
}

require_cmd curl
require_cmd jq

BASE_URL="${BASE_URL:-http://127.0.0.1:1113}"
USERNAME="${USERNAME:-}"
PASSWORD="${PASSWORD:-}"
BRANCH_ID="${BRANCH_ID:-BRC250118132203}"
PRODUCT_ID="${PRODUCT_ID:-PRD25050451578}"
TEST_DATE="${TEST_DATE:-2026-04-23}"
DESCRIPTION_CREATE="${DESCRIPTION_CREATE:-Copy resep dokter smoke test}"
DESCRIPTION_UPDATE="${DESCRIPTION_UPDATE:-Copy resep dokter smoke test update}"
PAYMENT="${PAYMENT:-cash}"
QTY="${QTY:-1}"
MEMBER_ID="${MEMBER_ID:-}"

if [[ -z "$USERNAME" || -z "$PASSWORD" ]]; then
  echo "ERROR: USERNAME dan PASSWORD wajib diisi" >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

step() {
  echo
  echo "==> $1"
}

api_json() {
  local method="$1"
  local url="$2"
  local auth="${3:-}"
  local data="${4:-}"
  local body_file="$TMP_DIR/body.json"
  local code_file="$TMP_DIR/code.txt"

  if [[ -n "$data" ]]; then
    if [[ -n "$auth" ]]; then
      curl -sS -X "$method" "$url" \
        -H 'Content-Type: application/json' \
        -H "Authorization: Bearer $auth" \
        -d "$data" \
        -o "$body_file" -w '%{http_code}' > "$code_file"
    else
      curl -sS -X "$method" "$url" \
        -H 'Content-Type: application/json' \
        -d "$data" \
        -o "$body_file" -w '%{http_code}' > "$code_file"
    fi
  else
    if [[ -n "$auth" ]]; then
      curl -sS -X "$method" "$url" \
        -H "Authorization: Bearer $auth" \
        -o "$body_file" -w '%{http_code}' > "$code_file"
    else
      curl -sS -X "$method" "$url" \
        -o "$body_file" -w '%{http_code}' > "$code_file"
    fi
  fi

  HTTP_CODE="$(cat "$code_file")"
  BODY="$(cat "$body_file")"
}

assert_http_200() {
  if [[ "$HTTP_CODE" != "200" ]]; then
    echo "ERROR: expected HTTP 200, got $HTTP_CODE" >&2
    echo "$BODY" | jq . 2>/dev/null || echo "$BODY"
    exit 1
  fi
}

print_body() {
  echo "$BODY" | jq . 2>/dev/null || echo "$BODY"
}

step "Health"
api_json GET "$BASE_URL/health"
assert_http_200
print_body

step "Login"
LOGIN_PAYLOAD=$(jq -nc --arg username "$USERNAME" --arg password "$PASSWORD" '{username:$username,password:$password}')
api_json POST "$BASE_URL/api/login" "" "$LOGIN_PAYLOAD"
assert_http_200
print_body
LOGIN_TOKEN="$(echo "$BODY" | jq -r '.data')"

step "Set branch"
BRANCH_PAYLOAD=$(jq -nc --arg branch_id "$BRANCH_ID" '{branch_id:$branch_id}')
api_json POST "$BASE_URL/api/set_branch" "$LOGIN_TOKEN" "$BRANCH_PAYLOAD"
assert_http_200
print_body
TOKEN="$(echo "$BODY" | jq -r '.data')"

step "Get product before create"
api_json GET "$BASE_URL/api/products/$PRODUCT_ID" "$TOKEN"
assert_http_200
print_body
STOCK_BEFORE_CREATE="$(echo "$BODY" | jq -r '.data.stock // empty')"

step "Create duplicate receipt"
CREATE_PAYLOAD=$(jq -nc \
  --arg member_id "$MEMBER_ID" \
  --arg description "$DESCRIPTION_CREATE" \
  --arg duplicate_receipt_date "$TEST_DATE" \
  --arg payment "$PAYMENT" \
  --arg product_id "$PRODUCT_ID" \
  --argjson qty "$QTY" \
  '{duplicate_receipt:{member_id:$member_id,description:$description,duplicate_receipt_date:$duplicate_receipt_date,payment:$payment},items:[{product_id:$product_id,qty:$qty}]}' )
api_json POST "$BASE_URL/api/duplicate-receipts" "$TOKEN" "$CREATE_PAYLOAD"
assert_http_200
print_body
DUPLICATE_RECEIPT_ID="$(echo "$BODY" | jq -r '.data.id // empty')"

step "List duplicate receipts"
api_json GET "$BASE_URL/api/duplicate-receipts?page=1&limit=10&month=$(date +%m)" "$TOKEN"
assert_http_200
print_body

step "Get duplicate receipt detail"
api_json GET "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" "$TOKEN"
assert_http_200
print_body

step "Get product after create"
api_json GET "$BASE_URL/api/products/$PRODUCT_ID" "$TOKEN"
assert_http_200
print_body
STOCK_AFTER_CREATE="$(echo "$BODY" | jq -r '.data.stock // empty')"

step "List duplicate receipt items"
api_json GET "$BASE_URL/api/duplicate-receipts-items/all/$DUPLICATE_RECEIPT_ID" "$TOKEN"
assert_http_200
print_body
ITEM_ID="$(echo "$BODY" | jq -r '.data[0].id // .data.id // empty')"

step "Create duplicate receipt item"
CREATE_ITEM_PAYLOAD=$(jq -nc \
  --arg duplicate_receipt_id "$DUPLICATE_RECEIPT_ID" \
  --arg product_id "$PRODUCT_ID" \
  --argjson qty 1 \
  '{duplicate_receipt_id:$duplicate_receipt_id,product_id:$product_id,qty:$qty}')
api_json POST "$BASE_URL/api/duplicate-receipts-items" "$TOKEN" "$CREATE_ITEM_PAYLOAD"
assert_http_200
print_body
ITEM_ID="$(echo "$BODY" | jq -r '.data.id // empty')"

step "Update duplicate receipt item"
UPDATE_ITEM_PAYLOAD=$(jq -nc \
  --arg product_id "$PRODUCT_ID" \
  --argjson qty 2 \
  '{product_id:$product_id,qty:$qty}')
api_json PUT "$BASE_URL/api/duplicate-receipts-items/$ITEM_ID" "$TOKEN" "$UPDATE_ITEM_PAYLOAD"
assert_http_200
print_body

step "Update duplicate receipt header"
UPDATE_PAYLOAD=$(jq -nc \
  --arg member_id "$MEMBER_ID" \
  --arg description "$DESCRIPTION_UPDATE" \
  --arg payment "$PAYMENT" \
  '{member_id:$member_id,description:$description,payment:$payment}')
api_json PUT "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" "$TOKEN" "$UPDATE_PAYLOAD"
assert_http_200
print_body

step "Delete duplicate receipt item"
api_json DELETE "$BASE_URL/api/duplicate-receipts-items/$ITEM_ID" "$TOKEN"
assert_http_200
print_body

step "Delete duplicate receipt"
api_json DELETE "$BASE_URL/api/duplicate-receipts/$DUPLICATE_RECEIPT_ID" "$TOKEN"
assert_http_200
print_body

step "Get product after delete"
api_json GET "$BASE_URL/api/products/$PRODUCT_ID" "$TOKEN"
assert_http_200
print_body
STOCK_AFTER_DELETE="$(echo "$BODY" | jq -r '.data.stock // empty')"

step "Logout"
api_json POST "$BASE_URL/api/logout" "$TOKEN"
assert_http_200
print_body

echo
printf 'Duplicate receipt smoke test selesai.\nID: %s\nStock sebelum create: %s\nStock sesudah create: %s\nStock sesudah delete: %s\n' \
  "$DUPLICATE_RECEIPT_ID" "$STOCK_BEFORE_CREATE" "$STOCK_AFTER_CREATE" "$STOCK_AFTER_DELETE"

if [[ -n "$STOCK_BEFORE_CREATE" && -n "$STOCK_AFTER_DELETE" && "$STOCK_BEFORE_CREATE" != "$STOCK_AFTER_DELETE" ]]; then
  echo "WARNING: stock sesudah delete berbeda dari stock awal. Cek rollback stok manual ya." >&2
fi
