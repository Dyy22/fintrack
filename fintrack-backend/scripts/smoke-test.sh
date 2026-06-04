#!/usr/bin/env sh
set -eu

API_BASE="${API_BASE:-http://localhost:8080/api/v1}"
EMAIL="${SMOKE_EMAIL:-smoke+$(date +%s)@example.com}"
PASSWORD="${SMOKE_PASSWORD:-securepassword}"

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

request() {
  method="$1"
  url="$2"
  data="${3:-}"
  token="${4:-}"
  expected_status="$5"
  output_file="$6"

  tmp_status="$(mktemp)"
  if [ -n "$data" ] && [ -n "$token" ]; then
    curl -sS -o "$output_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d "$data" > "$tmp_status"
  elif [ -n "$data" ]; then
    curl -sS -o "$output_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$data" > "$tmp_status"
  elif [ -n "$token" ]; then
    curl -sS -o "$output_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -H "Authorization: Bearer $token" > "$tmp_status"
  else
    curl -sS -o "$output_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" > "$tmp_status"
  fi

  status="$(cat "$tmp_status")"
  rm -f "$tmp_status"
  if [ "$status" != "$expected_status" ]; then
    echo "Expected $method $url to return $expected_status, got $status" >&2
    echo "Response body:" >&2
    cat "$output_file" >&2
    echo >&2
    exit 1
  fi
}

json_get() {
  path="$1"
  file="$2"
  python3 - "$path" "$file" <<'PY'
import json
import sys

path = sys.argv[1].split('.')
with open(sys.argv[2], 'r', encoding='utf-8') as f:
    data = json.load(f)
for part in path:
    if part.isdigit():
        data = data[int(part)]
    else:
        data = data[part]
print(data)
PY
}

json_find_category_id() {
  name="$1"
  file="$2"
  python3 - "$name" "$file" <<'PY'
import json
import sys

name = sys.argv[1]
with open(sys.argv[2], 'r', encoding='utf-8') as f:
    data = json.load(f)
for category in data.get('categories', []):
    if category.get('name') == name:
        print(category['id'])
        sys.exit(0)
raise SystemExit(f'category not found: {name}')
PY
}

json_assert_number() {
  path="$1"
  expected="$2"
  file="$3"
  python3 - "$path" "$expected" "$file" <<'PY'
import json
import sys

path = sys.argv[1].split('.')
expected = float(sys.argv[2])
with open(sys.argv[3], 'r', encoding='utf-8') as f:
    data = json.load(f)
for part in path:
    if part.isdigit():
        data = data[int(part)]
    else:
        data = data[part]
actual = float(data)
if actual != expected:
    raise SystemExit(f'expected {".".join(path)}={expected}, got {actual}')
PY
}

require_command curl
require_command python3

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

echo "Smoke test API base: $API_BASE"
echo "Smoke test user: $EMAIL"

request GET "$API_BASE/health" "" "" 200 "$workdir/health.json"

request POST "$API_BASE/auth/register" "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" "" 201 "$workdir/register.json"
request POST "$API_BASE/auth/login" "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" "" 200 "$workdir/login.json"
TOKEN="$(json_get token "$workdir/login.json")"

request GET "$API_BASE/account-types" "" "$TOKEN" 200 "$workdir/account_types.json"

request POST "$API_BASE/accounts" "{\"name\":\"BCA Savings\",\"account_type_id\":1,\"balance\":5000000}" "$TOKEN" 201 "$workdir/bank_account.json"
BANK_ACCOUNT_ID="$(json_get id "$workdir/bank_account.json")"

request POST "$API_BASE/accounts" "{\"name\":\"Cash\",\"account_type_id\":3,\"balance\":500000}" "$TOKEN" 201 "$workdir/cash_account.json"
CASH_ACCOUNT_ID="$(json_get id "$workdir/cash_account.json")"

request GET "$API_BASE/categories?type=expense" "" "$TOKEN" 200 "$workdir/categories.json"
FOOD_CATEGORY_ID="$(json_find_category_id Food "$workdir/categories.json")"

request POST "$API_BASE/transactions" "{\"account_id\":\"$BANK_ACCOUNT_ID\",\"category_id\":\"$FOOD_CATEGORY_ID\",\"type\":\"expense\",\"amount\":50000,\"description\":\"Lunch\",\"date\":\"2026-06-03T12:00:00Z\"}" "$TOKEN" 201 "$workdir/expense_transaction.json"

request POST "$API_BASE/transactions" "{\"account_id\":\"$BANK_ACCOUNT_ID\",\"transfer_account_id\":\"$CASH_ACCOUNT_ID\",\"type\":\"transfer\",\"amount\":100000,\"description\":\"ATM withdrawal\",\"date\":\"2026-06-03T13:00:00Z\"}" "$TOKEN" 201 "$workdir/transfer_transaction.json"

request GET "$API_BASE/accounts" "" "$TOKEN" 200 "$workdir/accounts.json"
request GET "$API_BASE/transactions?start_date=2026-06-01&end_date=2026-06-30" "" "$TOKEN" 200 "$workdir/transactions.json"
request GET "$API_BASE/reports/net-worth" "" "$TOKEN" 200 "$workdir/net_worth.json"
json_assert_number net_worth 5450000 "$workdir/net_worth.json"

request GET "$API_BASE/reports/spending-by-category?start_date=2026-06-01&end_date=2026-06-30" "" "$TOKEN" 200 "$workdir/spending.json"
json_assert_number total_spending 50000 "$workdir/spending.json"
json_assert_number categories.0.amount 50000 "$workdir/spending.json"

printf '\nSmoke test passed.\n'
