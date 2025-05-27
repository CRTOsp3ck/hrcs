#!/bin/bash

# HR Claims Management System - Database Seeding Script

echo "🌱 HR Claims Management System - Database Seeder"
echo "================================================"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Check if PostgreSQL is running
if ! nc -z localhost 5432 2>/dev/null; then
    echo "❌ Error: PostgreSQL is not running on localhost:5432"
    echo "💡 Tip: Run 'docker-compose up -d postgres' to start PostgreSQL"
    exit 1
fi

# Parse command line arguments
CLEAR_DATA=false
HELP=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --clear|-c)
            CLEAR_DATA=true
            shift
            ;;
        --help|-h)
            HELP=true
            shift
            ;;
        *)
            echo "❌ Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

if [ "$HELP" = true ]; then
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --clear, -c    Clear all existing data before seeding"
    echo "  --help, -h     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0             # Seed database (skip if data exists)"
    echo "  $0 --clear     # Clear and reseed database"
    echo ""
    echo "This script will create:"
    echo "  - Admin and normal users with default password 'password123'"
    echo "  - Various claim types (Travel, Medical, Office Supplies, etc.)"
    echo "  - User groups (Engineering, Sales, Marketing, etc.)"
    echo "  - Approval levels for each user group"
    echo "  - Sample claims in different statuses"
    exit 0
fi

echo "🔧 Setting up environment..."

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
fi

# Navigate to backend directory
cd backend

echo "📦 Installing Go dependencies..."
go mod download

echo "🗄️  Running database migrations..."
go run main.go &
SERVER_PID=$!
sleep 3
kill $SERVER_PID 2>/dev/null || true

echo "🌱 Seeding database..."
if [ "$CLEAR_DATA" = true ]; then
    echo "🧹 Clearing existing data first..."
    go run cmd/seed/main.go -clear
else
    go run cmd/seed/main.go
fi

# Check if seeding was successful
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Database seeding completed successfully!"
    echo ""
    echo "🚀 You can now start the application:"
    echo "   Backend:  cd backend && go run main.go"
    echo "   Frontend: cd frontend && npm run dev"
    echo ""
    echo "🌐 Application URLs:"
    echo "   Frontend: http://localhost:3000"
    echo "   Backend:  http://localhost:8000"
else
    echo "❌ Database seeding failed!"
    exit 1
fi