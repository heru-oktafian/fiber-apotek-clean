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
TOKEN_1="$(echo "$BODY" | jq -r '.data')"

step "List branches"
api_json GET "$BASE_URL/api/list_branches" "$TOKEN_1"
assert_http_200
print_body

step "Set branch"
SET_BRANCH_PAYLOAD=$(jq -nc --arg branch_id "$BRANCH_ID" '{branch_id:$branch_id}')
api_json POST "$BASE_URL/api/set_branch" "$TOKEN_1" "$SET_BRANCH_PAYLOAD"
assert_http_200
print_body
TOKEN_2="$(echo "$BODY" | jq -r '.data')"

step "Profile"
api_json GET "$BASE_URL/api/profile" "$TOKEN_2"
assert_http_200
print_body

step "Logout"
api_json POST "$BASE_URL/api/logout" "$TOKEN_2"
assert_http_200
print_body

echo
printf 'Auth smoke test selesai.\nTOKEN_1: %s\nTOKEN_2: %s\n' "$TOKEN_1" "$TOKEN_2"
