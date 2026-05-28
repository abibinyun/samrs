#!/bin/bash

# SAMRS Backend API Test Runner
# Usage: ./run-all-tests.sh

set -e

BASE_URL="http://localhost:8090"
RESULTS_FILE="docs/backend-testing/08-test-results/test-run-$(date +%Y-%m-%d-%H%M%S).md"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASS=0
FAIL=0
TOTAL=0

# Helper functions
log_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASS++))
    ((TOTAL++))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAIL++))
    ((TOTAL++))
}

# Get token
get_token() {
    curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"password123"}' | jq -r '.data.token'
}

echo "========================================="
echo "SAMRS Backend API Test Runner"
echo "========================================="
echo ""

# Get auth token
log_test "Getting authentication token..."
TOKEN=$(get_token)
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    log_fail "Failed to get token"
    exit 1
fi
log_pass "Token obtained"
echo ""

# =========================================
# PHASE 1: AUTHENTICATION
# =========================================
echo "========================================="
echo "PHASE 1: AUTHENTICATION"
echo "========================================="
echo ""

# Test 1.1: Public Ping
log_test "1.1 Public Ping"
RESPONSE=$(curl -s "$BASE_URL/ping")
if echo "$RESPONSE" | jq -e '.status == "online"' > /dev/null; then
    log_pass "Public ping working"
else
    log_fail "Public ping failed"
fi

# Test 1.2: Auth Me
log_test "1.2 Auth Me"
RESPONSE=$(curl -s "$BASE_URL/api/v1/auth/me" -H "Authorization: Bearer $TOKEN")
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
    log_pass "Auth me working"
else
    log_fail "Auth me failed"
fi

# Test 1.3: Secure Ping (with token)
log_test "1.3 Secure Ping (with token)"
RESPONSE=$(curl -s "$BASE_URL/api/v1/secure-ping" -H "Authorization: Bearer $TOKEN")
if echo "$RESPONSE" | jq -e '.tenant_id' > /dev/null; then
    log_pass "Secure ping with token working"
else
    log_fail "Secure ping with token failed"
fi

# Test 1.4: Secure Ping (without token)
log_test "1.4 Secure Ping (without token)"
RESPONSE=$(curl -s "$BASE_URL/api/v1/secure-ping")
if echo "$RESPONSE" | jq -e '.success == false' > /dev/null; then
    log_pass "Secure ping without token correctly rejected"
else
    log_fail "Secure ping without token should be rejected"
fi

echo ""

# =========================================
# PHASE 2: ROOM MANAGEMENT
# =========================================
echo "========================================="
echo "PHASE 2: ROOM MANAGEMENT"
echo "========================================="
echo ""

# Test 2.1: Create Room
log_test "2.1 Create Room"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/rooms" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Room Auto","code":"TRA-001","location":"Auto Test"}')
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
    ROOM_ID=$(echo "$RESPONSE" | jq -r '.data.id')
    log_pass "Room created: $ROOM_ID"
else
    log_fail "Room creation failed"
    ROOM_ID=""
fi

# Test 2.2: Get All Rooms
log_test "2.2 Get All Rooms"
RESPONSE=$(curl -s "$BASE_URL/api/v1/rooms" -H "Authorization: Bearer $TOKEN")
if echo "$RESPONSE" | jq -e '.data | length > 0' > /dev/null; then
    COUNT=$(echo "$RESPONSE" | jq '.data | length')
    log_pass "Rooms listed: $COUNT rooms"
else
    log_fail "Failed to list rooms"
fi

# Test 2.3: Get Room By ID
if [ -n "$ROOM_ID" ]; then
    log_test "2.3 Get Room By ID"
    RESPONSE=$(curl -s "$BASE_URL/api/v1/rooms/$ROOM_ID" -H "Authorization: Bearer $TOKEN")
    if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
        log_pass "Room retrieved by ID"
    else
        log_fail "Failed to get room by ID"
    fi
fi

# Test 2.4: Update Room
if [ -n "$ROOM_ID" ]; then
    log_test "2.4 Update Room"
    RESPONSE=$(curl -s -X PATCH "$BASE_URL/api/v1/rooms/$ROOM_ID" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"name":"Test Room Updated","code":"TRA-001-UPD","location":"Updated Location"}')
    if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
        log_pass "Room updated"
    else
        log_fail "Failed to update room"
    fi
fi

# Test 2.5: Delete Room
if [ -n "$ROOM_ID" ]; then
    log_test "2.5 Delete Room"
    RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/rooms/$ROOM_ID" -H "Authorization: Bearer $TOKEN")
    if echo "$RESPONSE" | jq -e '.success == true' > /dev/null; then
        log_pass "Room deleted"
    else
        log_fail "Failed to delete room"
    fi
fi

echo ""

# =========================================
# SUMMARY
# =========================================
echo "========================================="
echo "TEST SUMMARY"
echo "========================================="
echo "Total Tests: $TOTAL"
echo -e "${GREEN}Passed: $PASS${NC}"
echo -e "${RED}Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    exit 1
fi
