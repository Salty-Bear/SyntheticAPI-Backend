#!/bin/bash

# Syntra Demo Script
# This script demonstrates the complete Syntra tunnel system

echo "🚀 Syntra Demo - Reverse Tunnel System"
echo "======================================"
echo

# Check if required tools are available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

if ! command -v curl &> /dev/null; then
    echo "❌ curl is not installed. Please install curl first."
    exit 1
fi

# Build the server
echo "🔨 Building Syntra server..."
cd /Users/aryaman.raj/projects/Syntra
go build -o syntra-server main.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build server"
    exit 1
fi

# Build the CLI
echo "🔨 Building Syntra CLI..."
cd cli
go build -o syntra main.go
if [ $? -ne 0 ]; then
    echo "❌ Failed to build CLI"
    exit 1
fi

echo "✅ Build completed successfully!"
echo

# Start the server in background
echo "🌐 Starting Syntra server..."
cd ..
./syntra-server &
SERVER_PID=$!
echo "Server started with PID: $SERVER_PID"

# Wait for server to start
sleep 3

# Check if server is running
if ! ps -p $SERVER_PID > /dev/null; then
    echo "❌ Server failed to start"
    exit 1
fi

echo "✅ Server is running on http://localhost:8080"
echo

# Test the API endpoints
echo "🧪 Testing API endpoints..."

# Test tunnel API
echo "📡 Testing tunnel endpoints..."
curl -s http://localhost:8080/tunnel/active | jq . || echo "No active tunnels"

# Test testing API
echo "🔬 Testing test suite creation..."
SUITE_RESPONSE=$(curl -s -X POST http://localhost:8080/test/suites \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Test Suite",
    "base_url": "http://localhost:3000",
    "test_cases": [
      {
        "name": "Health Check",
        "method": "GET",
        "path": "/health",
        "expected": {
          "status_code": 200
        }
      }
    ]
  }')

echo "$SUITE_RESPONSE" | jq . || echo "Test suite creation response: $SUITE_RESPONSE"

echo
echo "🎉 Demo completed successfully!"
echo
echo "Next steps:"
echo "1. In one terminal, run: cd cli && ./syntra connect 3000"
echo "2. In another terminal, start your local API on port 3000"
echo "3. Use the tunnel URL to test your API from anywhere!"
echo
echo "API Endpoints:"
echo "- Tunnel Management: http://localhost:8080/tunnel/"
echo "- API Testing: http://localhost:8080/test/"
echo "- Active Tunnels: http://localhost:8080/tunnel/active"
echo
echo "Press Ctrl+C to stop the server (PID: $SERVER_PID)"

# Wait for user input
read -p "Press Enter to stop the demo server..."

# Clean up
echo "🧹 Cleaning up..."
kill $SERVER_PID 2>/dev/null
echo "✅ Demo completed!"
