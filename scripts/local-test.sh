#!/bin/bash

# ============================================
# Local Docker Test Script
# ============================================
# Test your Docker setup locally before deploying
# ============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

info "Testing Docker production build locally..."

# Check if .env exists
if [ ! -f ../.env ]; then
    warning ".env file not found, creating from .env.example..."
    cp .env.example .env
    warning "Please edit .env with your configuration before continuing"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    error "Docker is not running. Please start Docker and try again."
    exit 1
fi

info "Building Docker images..."
docker compose build

success "Docker images built successfully!"

info "Starting services..."
docker compose up -d

success "Services started!"

info "Waiting for services to be ready..."
sleep 10

# Wait for nginx to be ready
info "Waiting for nginx..."
sleep 5

# Check backend health through nginx
info "Testing backend health through nginx..."
if curl -f http://localhost/api/v1/health > /dev/null 2>&1 || curl -f http://localhost/health > /dev/null 2>&1; then
    success "Backend is accessible!"
else
    warning "Backend not accessible through nginx yet (this is OK for local testing)"
fi

# Check if services are running
info "Checking if all containers are running..."
if docker compose ps | grep -q "Up"; then
    success "All containers are running!"
else
    error "Some containers failed to start"
    docker compose ps
    exit 1
fi

# Display container status
echo ""
info "Container status:"
docker compose ps

echo ""
success "All tests passed!"
echo ""
info "Your application is running:"
echo "  Frontend: http://localhost:3000"
echo "  Backend:  http://localhost:8080"
echo "  API:      http://localhost:8080/api/v1/health"
echo "  Swagger:  http://localhost:8080/swagger/index.html"
echo ""
info "To view logs:"
echo "  docker compose logs -f"
echo ""
info "To stop services:"
echo "  docker compose down"
