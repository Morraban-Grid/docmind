#!/bin/bash
# Bootstrap script for DocMind project
# This script sets up the development environment

set -e

echo "========================================="
echo "DocMind Project Bootstrap"
echo "========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi
echo "✅ Docker is installed"

# Check if Docker Compose is installed
if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi
echo "✅ Docker Compose is installed"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  IMPORTANT: Please edit .env file and change the following:"
    echo "   - POSTGRES_PASSWORD"
    echo "   - MINIO_ROOT_PASSWORD"
    echo "   - JWT_SECRET (use a strong random string, min 32 characters)"
    echo ""
    echo "Press Enter to continue after editing .env file..."
    read
else
    echo "✅ .env file already exists"
fi

# Make scripts executable
echo "🔧 Making scripts executable..."
chmod +x scripts/*.sh

# Pull Docker images
echo "📦 Pulling Docker images..."
docker compose -f deployments/docker/docker-compose.yml pull

# Build services
echo "🏗️  Building services..."
docker compose -f deployments/docker/docker-compose.yml build

echo ""
echo "========================================="
echo "✅ Bootstrap complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Review and edit .env file with your configuration"
echo "2. Start services: make up"
echo "3. View logs: make logs"
echo "4. Check status: make ps"
echo ""
echo "For more commands, run: make help"
echo ""
