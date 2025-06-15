#!/bin/bash

echo "🚀 Setting up HR Claims Management System"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start PostgreSQL
echo "📦 Starting PostgreSQL database..."
docker compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Setup backend
echo "🔧 Setting up Go backend..."
cd backend

# Download Go dependencies
echo "📥 Downloading Go dependencies..."
go mod download

# Run database migrations
echo "🗄️  Running database migrations..."
go run main.go &
BACKEND_PID=$!
sleep 3
kill $BACKEND_PID

cd ..

# Setup frontend
echo "🎨 Setting up Vue.js frontend..."
cd frontend

# Install Node.js dependencies
echo "📥 Installing Node.js dependencies..."
npm install

cd ..

echo "✅ Setup complete!"
echo ""
echo "🚀 To start the application:"
echo "   Backend:  cd backend && go run main.go"
echo "   Frontend: cd frontend && npm run dev"
echo ""
echo "📱 Application URLs:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8000"
echo ""
echo "🗄️  Database:"
echo "   Host: localhost:5432"
echo "   Database: hrcs"
echo "   Username: postgres"
echo "   Password: postgres"