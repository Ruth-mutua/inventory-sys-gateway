#!/bin/bash

# Test script for Go Inventory System API

BASE_URL="http://localhost:8000"

echo "🧪 Testing Go Inventory System API"
echo "=================================="

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s -X GET "$BASE_URL/health"
echo -e "\n"

# Test user registration
echo "2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser"
  }')
echo "$REGISTER_RESPONSE"
echo -e "\n"

# Extract token from registration response
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get token from registration"
    exit 1
fi

echo "✅ Registration successful, token: ${TOKEN:0:20}..."

# Test login
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo "$LOGIN_RESPONSE"
echo -e "\n"

# Test getting current user profile
echo "4. Testing get current user profile..."
curl -s -X GET "$BASE_URL/users/me" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n"

# Test creating an order
echo "5. Testing create order..."
curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "product_name": "Laptop",
    "quantity": 1,
    "total_price": 999.99
  }'
echo -e "\n"

# Test getting all orders
echo "6. Testing get all orders..."
curl -s -X GET "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n"

# Test getting all users
echo "7. Testing get all users..."
curl -s -X GET "$BASE_URL/users" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n"

echo "✅ All tests completed!" 