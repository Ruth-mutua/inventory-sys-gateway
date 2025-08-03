#!/bin/bash

# Startup script for Go Inventory System

echo "🚀 Starting Go Inventory System..."
echo "=================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Create data directories
mkdir -p data

# Start services in background
echo "🔧 Starting services..."

# Start Auth Service
echo "Starting Auth Service on port 8083..."
go run services/auth/main.go &
AUTH_PID=$!

# Start Users Service
echo "Starting Users Service on port 8081..."
go run services/users/main.go &
USERS_PID=$!

# Start Orders Service
echo "Starting Orders Service on port 8082..."
go run services/orders/main.go &
ORDERS_PID=$!

# Wait a moment for services to start
sleep 3

# Start Gateway
echo "Starting API Gateway on port 8000..."
go run gateway/main.go &
GATEWAY_PID=$!

echo "✅ All services started!"
echo ""
echo "🌐 Services running:"
echo "   - API Gateway: http://localhost:8000"
echo "   - Auth Service: http://localhost:8083"
echo "   - Users Service: http://localhost:8081"
echo "   - Orders Service: http://localhost:8082"
echo ""
echo "📋 Test the API:"
echo "   ./test_api.sh"
echo ""
echo "🛑 To stop all services, press Ctrl+C"

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping services..."
    kill $AUTH_PID $USERS_PID $ORDERS_PID $GATEWAY_PID 2>/dev/null
    echo "✅ All services stopped"
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Wait for all background processes
wait 