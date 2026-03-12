#!/bin/bash
# Verification script to check if the environment is properly set up

set -e

echo "========================================="
echo "DocMind Setup Verification"
echo "========================================="
echo ""

ERRORS=0

# Check Docker
echo "🔍 Checking Docker..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    echo "✅ Docker installed: $DOCKER_VERSION"
else
    echo "❌ Docker is not installed"
    ERRORS=$((ERRORS + 1))
fi

# Check Docker Compose
echo "🔍 Checking Docker Compose..."
if docker compose version &> /dev/null; then
    COMPOSE_VERSION=$(docker compose version)
    echo "✅ Docker Compose installed: $COMPOSE_VERSION"
else
    echo "❌ Docker Compose is not installed"
    ERRORS=$((ERRORS + 1))
fi

# Check .env file
echo "🔍 Checking .env file..."
if [ -f .env ]; then
    echo "✅ .env file exists"
    
    # Check for default passwords
    if grep -q "changeme" .env; then
        echo "⚠️  WARNING: .env file contains default passwords"
        echo "   Please change POSTGRES_PASSWORD, MINIO_ROOT_PASSWORD, and JWT_SECRET"
    fi
    
    # Check JWT_SECRET length
    JWT_SECRET=$(grep "^JWT_SECRET=" .env | cut -d'=' -f2)
    if [ ${#JWT_SECRET} -lt 32 ]; then
        echo "⚠️  WARNING: JWT_SECRET is less than 32 characters"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "❌ .env file not found"
    echo "   Run: cp .env.example .env"
    ERRORS=$((ERRORS + 1))
fi

# Check .gitignore
echo "🔍 Checking .gitignore..."
if [ -f .gitignore ]; then
    if grep -q "^\.env$" .gitignore; then
        echo "✅ .gitignore protects .env file"
    else
        echo "❌ .gitignore does not protect .env file"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "❌ .gitignore file not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Go installation (optional)
echo "🔍 Checking Go (optional for local development)..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "✅ Go installed: $GO_VERSION"
else
    echo "ℹ️  Go not installed (optional - only needed for local Go development)"
fi

# Check Python installation (optional)
echo "🔍 Checking Python (optional for local development)..."
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version)
    echo "✅ Python installed: $PYTHON_VERSION"
else
    echo "ℹ️  Python not installed (optional - only needed for local Python development)"
fi

# Check disk space
echo "🔍 Checking disk space..."
AVAILABLE_SPACE=$(df -h . | awk 'NR==2 {print $4}')
echo "ℹ️  Available disk space: $AVAILABLE_SPACE"

# Check if services are running
echo "🔍 Checking if services are running..."
if docker compose -f deployments/docker/docker-compose.yml ps | grep -q "Up"; then
    echo "✅ Some services are running"
    docker compose -f deployments/docker/docker-compose.yml ps
else
    echo "ℹ️  No services are currently running"
    echo "   To start services: make up"
fi

echo ""
echo "========================================="
if [ $ERRORS -eq 0 ]; then
    echo "✅ Setup verification passed!"
    echo "========================================="
    echo ""
    echo "You're ready to start DocMind!"
    echo "Run: make up"
else
    echo "❌ Setup verification failed with $ERRORS error(s)"
    echo "========================================="
    echo ""
    echo "Please fix the errors above before starting."
fi
echo ""
