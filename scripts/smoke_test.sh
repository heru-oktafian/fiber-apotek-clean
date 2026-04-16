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
SUPPLIER_ID="${SUPPLIER_ID:-SPL250207144606}"
PURCHASE_PRODUCT_ID="${PURCHASE_PRODUCT_ID:-PRD25050451578}"
PURCHASE_UNIT_ID="${PURCHASE_UNIT_ID:-UNT250118132755}"
OPNAME_PRODUCT_ID="${OPNAME_PRODUCT_ID:-PRD054724Q21ODS}"
OPNAME_PRICE="${OPNAME_PRICE:-11400}"
PURCHASE_PRICE="${PURCHASE_PRICE:-6433}"
SEARCH_TERM="${SEARCH_TERM:-par}"
TEST_DATE="${TEST_DATE:-2026-04-15}"
EXPIRED_DATE="${EXPIRED_DATE:-2027-12-31}"

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

debug_env() {
  echo "BASE_URL=$BASE_URL"
  echo "USERNAME=$USERNAME"
  echo "BRANCH_ID=$BRANCH_ID"
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

step "Debug context"
debug_env

step "Health"
api_json GET "$BASE_URL/health"
assert_http_200
print_body

step "Login"
LOGIN_PAYLOAD=$(jq -nc --arg username "$USERNAME" --arg password "$PASSWORD" '{username:$username,password:$password}')
echo "LOGIN_USERNAME=$(echo "$LOGIN_PAYLOAD" | jq -r '.username')"
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

step "Sales combo"
api_json GET "$BASE_URL/api/sales-products-combo?search=$SEARCH_TERM" "$TOKEN"
assert_http_200
print_body

step "Purchase combo"
api_json GET "$BASE_URL/api/purchase-products-combo?search=$SEARCH_TERM" "$TOKEN"
assert_http_200
print_body

step "Opname combo"
api_json GET "$BASE_URL/api/cmb-product-opname?search=$SEARCH_TERM" "$TOKEN"
assert_http_200
print_body

step "Create purchase"
PURCHASE_PAYLOAD=$(jq -nc \
  --arg supplier_id "$SUPPLIER_ID" \
  --arg purchase_date "$TEST_DATE" \
  --arg payment "cash" \
  --arg product_id "$PURCHASE_PRODUCT_ID" \
  --arg unit_id "$PURCHASE_UNIT_ID" \
  --arg expired_date "$EXPIRED_DATE" \
  --argjson price "$PURCHASE_PRICE" \
  --argjson qty 1 \
  '{purchase:{supplier_id:$supplier_id,purchase_date:$purchase_date,payment:$payment},purchase_items:[{product_id:$product_id,unit_id:$unit_id,price:$price,qty:$qty,expired_date:$expired_date}]}')
api_json POST "$BASE_URL/api/purchases" "$TOKEN" "$PURCHASE_PAYLOAD"
assert_http_200
print_body
PURCHASE_ID="$(echo "$BODY" | jq -r '.data.purchase.id // empty')"

step "Create sale"
SALE_PAYLOAD=$(jq -nc \
  --arg product_id "$PURCHASE_PRODUCT_ID" \
  --argjson price 999999 \
  --argjson qty 1 \
  '{sale:{payment:"cash",discount:0},sale_items:[{product_id:$product_id,price:$price,qty:$qty}]}')
api_json POST "$BASE_URL/api/sales" "$TOKEN" "$SALE_PAYLOAD"
assert_http_200
print_body
SALE_ID="$(echo "$BODY" | jq -r '.data.sale.id // empty')"

step "Create opname header"
OPNAME_PAYLOAD=$(jq -nc --arg description "uji rewrite opname" --arg opname_date "$TEST_DATE" '{description:$description,opname_date:$opname_date}')
api_json POST "$BASE_URL/api/opnames" "$TOKEN" "$OPNAME_PAYLOAD"
assert_http_200
print_body
OPNAME_ID="$(echo "$BODY" | jq -r '.data.id')"

step "Create opname item"
OPNAME_ITEM_PAYLOAD=$(jq -nc \
  --arg opname_id "$OPNAME_ID" \
  --arg product_id "$OPNAME_PRODUCT_ID" \
  --arg expired_date "$EXPIRED_DATE" \
  --argjson qty 3 \
  --argjson price "$OPNAME_PRICE" \
  '{opname_id:$opname_id,product_id:$product_id,qty:$qty,price:$price,expired_date:$expired_date}')
api_json POST "$BASE_URL/api/opname-items" "$TOKEN" "$OPNAME_ITEM_PAYLOAD"
assert_http_200
print_body

step "Get opname items all"
OPNAME_ITEMS_PAYLOAD=$(jq -nc --arg opname_id "$OPNAME_ID" '{opname_id:$opname_id}')
api_json POST "$BASE_URL/api/opname-items-all" "$TOKEN" "$OPNAME_ITEMS_PAYLOAD"
assert_http_200
print_body

step "Get opname detail"
api_json GET "$BASE_URL/api/opnames/$OPNAME_ID" "$TOKEN"
assert_http_200
print_body

step "Logout"
api_json POST "$BASE_URL/api/logout" "$TOKEN"
assert_http_200
print_body

echo
printf 'Smoke test selesai.\nPurchase ID: %s\nSale ID: %s\nOpname ID: %s\n' "$PURCHASE_ID" "$SALE_ID" "$OPNAME_ID"
