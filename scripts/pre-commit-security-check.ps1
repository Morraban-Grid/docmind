# Pre-commit security check script for DocMind (PowerShell version)
# This script helps prevent committing sensitive information

$ErrorActionPreference = "Stop"

Write-Host "🔒 Running pre-commit security checks..." -ForegroundColor Cyan

$ERRORS = 0

# Check 1: Ensure no .env files are staged
Write-Host "`n📋 Check 1: Verifying no .env files are staged..." -ForegroundColor Yellow
$stagedEnvFiles = git diff --cached --name-only | Select-String -Pattern "\.env$" | Where-Object { $_ -notmatch ".env.example" }
if ($stagedEnvFiles) {
    Write-Host "❌ ERROR: .env file(s) are staged for commit!" -ForegroundColor Red
    Write-Host "   Run: git reset HEAD .env" -ForegroundColor Red
    $ERRORS++
} else {
    Write-Host "✓ No .env files staged" -ForegroundColor Green
}

# Check 2: Look for potential hardcoded secrets in staged files
Write-Host "`n📋 Check 2: Scanning for hardcoded secrets..." -ForegroundColor Yellow
$secrets = git diff --cached | Select-String -Pattern "(password|secret|token|api_key|apikey)\s*=\s*['\`"][^$\{]" -CaseSensitive:$false | Where-Object { $_ -notmatch "changeme" -and $_ -notmatch ".env.example" }
if ($secrets) {
    Write-Host "❌ ERROR: Potential hardcoded secrets found!" -ForegroundColor Red
    Write-Host "   Review the output above and remove any hardcoded credentials" -ForegroundColor Red
    $ERRORS++
} else {
    Write-Host "✓ No hardcoded secrets detected" -ForegroundColor Green
}

# Check 3: Verify .gitignore exists and contains .env
Write-Host "`n📋 Check 3: Verifying .gitignore configuration..." -ForegroundColor Yellow
if (-not (Test-Path .gitignore)) {
    Write-Host "❌ ERROR: .gitignore file not found!" -ForegroundColor Red
    $ERRORS++
} elseif (-not (Get-Content .gitignore | Select-String -Pattern "^\.env$")) {
    Write-Host "❌ ERROR: .gitignore does not exclude .env files!" -ForegroundColor Red
    $ERRORS++
} else {
    Write-Host "✓ .gitignore properly configured" -ForegroundColor Green
}

# Check 4: Verify .env.example doesn't contain real secrets
Write-Host "`n📋 Check 4: Verifying .env.example contains only placeholders..." -ForegroundColor Yellow
if (Test-Path .env.example) {
    $realValues = git diff --cached .env.example | Select-String -Pattern "^\+" | Select-String -Pattern "(password|secret|token)" -CaseSensitive:$false | Where-Object { $_ -notmatch "(changeme|your-|example-|placeholder|<|>|\[|\])" }
    if ($realValues) {
        Write-Host "⚠️  WARNING: .env.example may contain real values instead of placeholders" -ForegroundColor Yellow
        Write-Host "   Please review and use placeholder values like 'changeme_*'" -ForegroundColor Yellow
    } else {
        Write-Host "✓ .env.example contains only placeholders" -ForegroundColor Green
    }
}

# Check 5: Look for common secret patterns
Write-Host "`n📋 Check 5: Scanning for common secret patterns..." -ForegroundColor Yellow
$secretPatterns = @(
    "AKIA[0-9A-Z]{16}",  # AWS Access Key
    "ghp_[0-9a-zA-Z]{36}",  # GitHub Personal Access Token
    "sk_live_[0-9a-zA-Z]{24}",  # Stripe Live Key
    "-----BEGIN (RSA |DSA )?PRIVATE KEY-----"  # Private Keys
)

$foundPatterns = $false
foreach ($pattern in $secretPatterns) {
    $matches = git diff --cached | Select-String -Pattern $pattern
    if ($matches) {
        Write-Host "❌ ERROR: Potential secret pattern detected: $pattern" -ForegroundColor Red
        $ERRORS++
        $foundPatterns = $true
    }
}

if (-not $foundPatterns) {
    Write-Host "✓ No common secret patterns detected" -ForegroundColor Green
}

# Check 6: Verify docker-compose.yml doesn't have hardcoded credentials
Write-Host "`n📋 Check 6: Checking docker-compose files..." -ForegroundColor Yellow
$dockerSecrets = git diff --cached | Select-String -Pattern "docker-compose.*\.yml" | Select-String -Pattern "(password|secret).*:.*[^$\{]" | Where-Object { $_ -notmatch "changeme" }
if ($dockerSecrets) {
    Write-Host "⚠️  WARNING: docker-compose file may contain hardcoded credentials" -ForegroundColor Yellow
    Write-Host "   Use environment variables instead: `${VARIABLE_NAME}" -ForegroundColor Yellow
}

# Summary
Write-Host "`n================================" -ForegroundColor Cyan
if ($ERRORS -eq 0) {
    Write-Host "✅ All security checks passed!" -ForegroundColor Green
    Write-Host "================================" -ForegroundColor Cyan
    exit 0
} else {
    Write-Host "❌ $ERRORS security check(s) failed!" -ForegroundColor Red
    Write-Host "================================" -ForegroundColor Cyan
    Write-Host "`nPlease fix the issues above before committing." -ForegroundColor Yellow
    Write-Host "If you're certain these are false positives, you can:" -ForegroundColor Yellow
    Write-Host "  1. Review the changes carefully" -ForegroundColor Yellow
    Write-Host "  2. Update .gitleaksignore if needed" -ForegroundColor Yellow
    Write-Host "  3. Commit with --no-verify (NOT RECOMMENDED)" -ForegroundColor Yellow
    exit 1
}
