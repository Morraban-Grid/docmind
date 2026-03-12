#!/bin/bash
# Pre-commit security check script for DocMind
# This script helps prevent committing sensitive information

set -e

echo "🔒 Running pre-commit security checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0

# Check 1: Ensure no .env files are staged
echo "📋 Check 1: Verifying no .env files are staged..."
if git diff --cached --name-only | grep -E "\.env$" | grep -v ".env.example"; then
    echo -e "${RED}❌ ERROR: .env file(s) are staged for commit!${NC}"
    echo "   Run: git reset HEAD .env"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✓ No .env files staged${NC}"
fi

# Check 2: Look for potential hardcoded secrets in staged files
echo "📋 Check 2: Scanning for hardcoded secrets..."
if git diff --cached | grep -iE "(password|secret|token|api_key|apikey)\s*=\s*['\"][^$\{]" | grep -v "changeme" | grep -v ".env.example"; then
    echo -e "${RED}❌ ERROR: Potential hardcoded secrets found!${NC}"
    echo "   Review the output above and remove any hardcoded credentials"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✓ No hardcoded secrets detected${NC}"
fi

# Check 3: Verify .gitignore exists and contains .env
echo "📋 Check 3: Verifying .gitignore configuration..."
if [ ! -f .gitignore ]; then
    echo -e "${RED}❌ ERROR: .gitignore file not found!${NC}"
    ERRORS=$((ERRORS + 1))
elif ! grep -q "^\.env$" .gitignore; then
    echo -e "${RED}❌ ERROR: .gitignore does not exclude .env files!${NC}"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✓ .gitignore properly configured${NC}"
fi

# Check 4: Verify .env.example doesn't contain real secrets
echo "📋 Check 4: Verifying .env.example contains only placeholders..."
if [ -f .env.example ]; then
    if git diff --cached .env.example | grep -E "^\+" | grep -iE "(password|secret|token)" | grep -vE "(changeme|your-|example-|placeholder|<|>|\[|\])"; then
        echo -e "${YELLOW}⚠️  WARNING: .env.example may contain real values instead of placeholders${NC}"
        echo "   Please review and use placeholder values like 'changeme_*'"
    else
        echo -e "${GREEN}✓ .env.example contains only placeholders${NC}"
    fi
fi

# Check 5: Look for common secret patterns
echo "📋 Check 5: Scanning for common secret patterns..."
SECRET_PATTERNS=(
    "AKIA[0-9A-Z]{16}"  # AWS Access Key
    "ghp_[0-9a-zA-Z]{36}"  # GitHub Personal Access Token
    "sk_live_[0-9a-zA-Z]{24}"  # Stripe Live Key
    "-----BEGIN (RSA |DSA )?PRIVATE KEY-----"  # Private Keys
)

for pattern in "${SECRET_PATTERNS[@]}"; do
    if git diff --cached | grep -E "$pattern"; then
        echo -e "${RED}❌ ERROR: Potential secret pattern detected: $pattern${NC}"
        ERRORS=$((ERRORS + 1))
    fi
done

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ No common secret patterns detected${NC}"
fi

# Check 6: Verify docker-compose.yml doesn't have hardcoded credentials
echo "📋 Check 6: Checking docker-compose files..."
if git diff --cached | grep -E "docker-compose.*\.yml" | grep -E "(password|secret).*:.*[^$\{]" | grep -v "changeme"; then
    echo -e "${YELLOW}⚠️  WARNING: docker-compose file may contain hardcoded credentials${NC}"
    echo "   Use environment variables instead: \${VARIABLE_NAME}"
fi

# Summary
echo ""
echo "================================"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ All security checks passed!${NC}"
    echo "================================"
    exit 0
else
    echo -e "${RED}❌ $ERRORS security check(s) failed!${NC}"
    echo "================================"
    echo ""
    echo "Please fix the issues above before committing."
    echo "If you're certain these are false positives, you can:"
    echo "  1. Review the changes carefully"
    echo "  2. Update .gitleaksignore if needed"
    echo "  3. Commit with --no-verify (NOT RECOMMENDED)"
    exit 1
fi
