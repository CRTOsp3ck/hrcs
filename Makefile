# HR Claims Management System - Makefile

.PHONY: help setup seed seed-clear dev-backend dev-frontend dev build clean docker-up docker-down

# Default target
help:
	@echo "HR Claims Management System"
	@echo "=========================="
	@echo ""
	@echo "Available commands:"
	@echo "  make setup        - Initial project setup"
	@echo "  make seed         - Seed database with sample data"
	@echo "  make seed-clear   - Clear and reseed database"
	@echo "  make dev-backend  - Start backend development server"
	@echo "  make dev-frontend - Start frontend development server" 
	@echo "  make dev          - Start both backend and frontend"
	@echo "  make build        - Build the application"
	@echo "  make docker-up    - Start PostgreSQL with Docker"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make clean        - Clean build artifacts"
	@echo ""

# Initial setup
setup:
	@echo "🚀 Setting up HR Claims Management System..."
	@./setup.sh

# Database seeding
seed:
	@echo "🌱 Seeding database..."
	@cd backend && go run cmd/seed/main.go

seed-clear:
	@echo "🧹 Clearing and seeding database..."
	@cd backend && go run cmd/seed/main.go -clear

# Development servers
dev-backend:
	@echo "🔧 Starting backend server..."
	@cd backend && go run main.go

dev-frontend:
	@echo "🎨 Starting frontend server..."
	@cd frontend && npm run dev

dev:
	@echo "🚀 Starting both backend and frontend..."
	@echo "Backend will be available at http://localhost:8000"
	@echo "Frontend will be available at http://localhost:3000"
	@make -j2 dev-backend dev-frontend

# Build
build:
	@echo "🔨 Building backend..."
	@cd backend && go build -o ../bin/hrcs-backend main.go
	@echo "🔨 Building frontend..."
	@cd frontend && npm run build
	@echo "✅ Build complete!"

# Docker commands
docker-up:
	@echo "🐳 Starting PostgreSQL..."
	@docker-compose up -d postgres

docker-down:
	@echo "🐳 Stopping Docker services..."
	@docker-compose down

# Clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf frontend/dist/
	@cd backend && go clean
	@echo "✅ Clean complete!"

# Install dependencies
deps:
	@echo "📦 Installing backend dependencies..."
	@cd backend && go mod download
	@echo "📦 Installing frontend dependencies..."
	@cd frontend && npm install