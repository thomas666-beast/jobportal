#!/bin/bash
BASE=http://localhost:8080/api
FAIL=0

check() {
  local label="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then
    echo "  ✓ $label → $actual"
  else
    echo "  ✗ $label → expected $expected, got $actual"
    FAIL=$((FAIL+1))
  fi
}

echo "=== 1. Public endpoints ==="
check "GET /jobs"      200 "$(curl -s -o /dev/null -w '%{http_code}' $BASE/jobs)"
check "GET /companies" 200 "$(curl -s -o /dev/null -w '%{http_code}' $BASE/companies)"

echo ""
echo "=== 2. Auth ==="
LOGIN=$(curl -s -X POST $BASE/login \
  -H "Content-Type: application/json" \
  -d '{"email":"seeker1@seed.local","password":"secret123"}')
TOKEN=$(echo "$LOGIN" | jq -r .token)
if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
  echo "  ✓ Login → token received"
else
  echo "  ✗ Login failed"
  echo "$LOGIN" | jq
  FAIL=$((FAIL+1))
fi

echo ""
echo "=== 3. Authenticated endpoints ==="
check "GET /me"                 200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/me)"
check "GET /seeker/profile"     200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/seeker/profile)"
check "GET /seeker/applications" 200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/seeker/applications)"
check "GET /seeker/saved-jobs"  200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/seeker/saved-jobs)"
check "GET /notifications"      200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/notifications)"
check "GET /notifications/recent" 200 "$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer $TOKEN" $BASE/notifications/recent)"

echo ""
echo "=== 4. email_verified flag in /me ==="
ME=$(curl -s -H "Authorization: Bearer $TOKEN" $BASE/me)
echo "$ME" | jq -e '.email_verified == true' > /dev/null && echo "  ✓ email_verified = true" || { echo "  ✗ email_verified not true"; echo "$ME" | jq; FAIL=$((FAIL+1)); }

echo ""
echo "=== 5. Email endpoints ==="
check "POST /auth/forgot-password (known)"   200 "$(curl -s -o /dev/null -w '%{http_code}' -X POST $BASE/auth/forgot-password -H 'Content-Type: application/json' -d '{"email":"seeker1@seed.local"}')"
check "POST /auth/forgot-password (unknown)" 200 "$(curl -s -o /dev/null -w '%{http_code}' -X POST $BASE/auth/forgot-password -H 'Content-Type: application/json' -d '{"email":"nobody@example.com"}')"

echo ""
echo "=== 6. Rate limiting on register ==="
ok=0
blocked=0
for i in 1 2 3 4 5 6 7; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE/register \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"smoke_$i\",\"email\":\"smoke_$i@test.local\",\"password\":\"secret123\"}")
  echo "  attempt $i → $CODE"
  [ "$CODE" = "201" ] && ok=$((ok+1))
  [ "$CODE" = "429" ] && blocked=$((blocked+1))
done
if [ $ok -ge 5 ] && [ $blocked -ge 1 ]; then
  echo "  ✓ burst allowed then rate-limited"
else
  echo "  ✗ rate limiting not working as expected"
  FAIL=$((FAIL+1))
fi

echo ""
echo "=== 7. 401 for unauthenticated protected endpoint ==="
check "GET /me without token" 401 "$(curl -s -o /dev/null -w '%{http_code}' $BASE/me)"

echo ""
if [ $FAIL -eq 0 ]; then
  echo "✅ All checks passed"
else
  echo "❌ $FAIL check(s) failed"
  exit 1
fi
