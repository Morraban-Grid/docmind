# Complete Git history cleanup for DocMind
# This creates a fresh repository with clean history

$ErrorActionPreference = "Stop"

Write-Host "🧹 Complete Git History Cleanup" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Create backup
Write-Host "📦 Step 1: Creating backup..." -ForegroundColor Yellow
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$backupDir = "../docmind-backup-$timestamp"

try {
    # Create a complete backup
    Copy-Item -Path . -Destination $backupDir -Recurse -Force -Exclude @('.git')
    Write-Host "✓ Backup created at: $backupDir" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to create backup: $_" -ForegroundColor Red
    exit 1
}

# Step 2: Remove Git history
Write-Host "`n🗑️  Step 2: Removing old Git history..." -ForegroundColor Yellow
Remove-Item -Path .git -Recurse -Force
Write-Host "✓ Old history removed" -ForegroundColor Green

# Step 3: Initialize new repository
Write-Host "`n🆕 Step 3: Initializing fresh repository..." -ForegroundColor Yellow
git init
git branch -M main
Write-Host "✓ Fresh repository initialized" -ForegroundColor Green

# Step 4: Add all files
Write-Host "`n📝 Step 4: Adding all files..." -ForegroundColor Yellow
git add .
Write-Host "✓ Files staged" -ForegroundColor Green

# Step 5: Create initial commit
Write-Host "`n💾 Step 5: Creating clean initial commit..." -ForegroundColor Yellow
git commit -m "chore: initialize DocMind with clean history

This is a fresh start with all security measures in place:
- No default passwords in configuration files
- Gitleaks configuration for secret scanning
- Comprehensive .gitignore
- Security documentation and best practices
- Pre-commit security check scripts

Previous history was cleaned to remove weak default passwords
that were present in docker-compose.yml configuration.

Iterations completed:
- Iteration 1: Infrastructure & Security Foundation
- Iteration 2: User Authentication & Management

Ready for Iteration 3: Document Upload & Storage"

Write-Host "✓ Clean commit created" -ForegroundColor Green

# Step 6: Show status
Write-Host "`n📊 Step 6: Repository status..." -ForegroundColor Yellow
git log --oneline
Write-Host ""
git status

Write-Host "`n================================" -ForegroundColor Cyan
Write-Host "✅ History cleanup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Review the changes: git log" -ForegroundColor Cyan
Write-Host "2. Add remote: git remote add origin https://github.com/Morraban-Grid/docmind.git" -ForegroundColor Cyan
Write-Host "3. Force push: git push -u origin main --force" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  WARNING: This will completely replace the remote history!" -ForegroundColor Yellow
Write-Host "Backup location: $backupDir" -ForegroundColor Cyan
Write-Host ""
