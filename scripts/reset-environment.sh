#!/bin/bash
# Reset environment script
# WARNING: This will delete all data and containers

set -e

echo "========================================="
echo "⚠️  WARNING: Environment Reset"
echo "========================================="
echo ""
echo "This will:"
echo "  - Stop all DocMind containers"
echo "  - Remove all containers"
echo "  - Remove all volumes (ALL DATA WILL BE LOST)"
echo "  - Remove all images"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Reset cancelled."
    exit 0
fi

echo ""
echo "🛑 Stopping all services..."
docker compose -f deployments/docker/docker-compose.yml down

echo "🗑️  Removing volumes..."
docker compose -f deployments/docker/docker-compose.yml down -v

echo "🗑️  Removing images..."
docker compose -f deployments/docker/docker-compose.yml down -v --rmi all

echo ""
echo "✅ Environment reset complete!"
echo ""
echo "To start fresh, run: ./scripts/bootstrap.sh"
echo ""
