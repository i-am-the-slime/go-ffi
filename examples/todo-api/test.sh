#!/bin/bash

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing Todo API Server${NC}\n"

BASE_URL="http://localhost:3000"

# Helper function
test_endpoint() {
    echo -e "${YELLOW}➜ $1${NC}"
    echo "  $2"
    echo ""
    curl -s "$3" | python3 -m json.tool 2>/dev/null || curl -s "$3"
    echo -e "\n"
}

test_endpoint_with_data() {
    echo -e "${YELLOW}➜ $1${NC}"
    echo "  $2"
    echo ""
    curl -s -X "$3" "$4" -H "Content-Type: application/json" -d "$5" | python3 -m json.tool 2>/dev/null || curl -s -X "$3" "$4" -H "Content-Type: application/json" -d "$5"
    echo -e "\n"
}

# Check if server is running
if ! curl -s "$BASE_URL" > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Server not running. Start it with:${NC}"
    echo "  ./todo-server"
    echo ""
    echo "Then run this test script again."
    exit 1
fi

echo -e "${GREEN}✓ Server is running${NC}\n"

# Test 1: Welcome
test_endpoint \
    "Test 1: Welcome Message" \
    "GET /" \
    "$BASE_URL/"

# Test 2: List todos
test_endpoint \
    "Test 2: List All Todos" \
    "GET /todos" \
    "$BASE_URL/todos"

# Test 3: Get specific todo
test_endpoint \
    "Test 3: Get Todo #1" \
    "GET /todos/1" \
    "$BASE_URL/todos/1"

# Test 4: Create new todo
test_endpoint_with_data \
    "Test 4: Create New Todo" \
    "POST /todos" \
    "POST" \
    "$BASE_URL/todos" \
    '{"title":"Test Todo","description":"Created by test script","completed":false}'

# Test 5: Update todo
test_endpoint_with_data \
    "Test 5: Update Todo #2" \
    "PUT /todos/2" \
    "PUT" \
    "$BASE_URL/todos/2" \
    '{"completed":true,"title":"Build Go FFI - DONE!"}'

# Test 6: Get stats
test_endpoint \
    "Test 6: Get Statistics" \
    "GET /stats" \
    "$BASE_URL/stats"

# Test 7: Invalid JSON (error handling)
echo -e "${YELLOW}➜ Test 7: Error Handling (Invalid JSON)${NC}"
echo "  POST /todos (with invalid JSON)"
echo ""
curl -s -X POST "$BASE_URL/todos" -H "Content-Type: application/json" -d 'not valid json'
echo -e "\n"

# Test 8: Non-existent todo
echo -e "${YELLOW}➜ Test 8: Not Found (Todo #999)${NC}"
echo "  GET /todos/999"
echo ""
curl -s "$BASE_URL/todos/999"
echo -e "\n"

# Test 9: Delete todo
echo -e "${YELLOW}➜ Test 9: Delete Todo #3${NC}"
echo "  DELETE /todos/3"
echo ""
curl -s -X DELETE "$BASE_URL/todos/3"
echo -e "\n"

# Test 10: Verify deletion
test_endpoint \
    "Test 10: Verify Deletion (List Todos)" \
    "GET /todos" \
    "$BASE_URL/todos"

# Test 11: External API demo
test_endpoint \
    "Test 11: External API Demo" \
    "GET /external" \
    "$BASE_URL/external"

echo -e "${GREEN}✅ All tests completed!${NC}"

