# Script to clean sensitive data from Git history
# USE WITH CAUTION - This rewrites Git history

param(
    [switch]$DryRun = $false,
    [switch]$Force = $false
)

$ErrorActionPreference = "Stop"

Write-Host "🧹 Git History Cleanup Script for DocMind" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

if (-not $Force) {
    Write-Host "⚠️  WARNING: This script will REWRITE Git history!" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "This means:" -ForegroundColor Yellow
    Write-Host "  • All commit hashes will change" -ForegroundColor Yellow
    Write-Host "  • Anyone who cloned the repo will need to re-clone" -ForegroundColor Yellow
    Write-Host "  • You cannot undo this without the backup" -ForegroundColor Yellow
    Write-Host ""
    
    $response = Read-Host "Do you want to continue? (type 'YES' to proceed)"
    if ($response -ne "YES") {
        Write-Host "Aborted." -ForegroundColor Red
        exit 1
    }
}

# Step 1: Check if we're in a git repository
Write-Host "📋 Step 1: Verifying Git repository..." -ForegroundColor Cyan
if (-not (Test-Path .git)) {
    Write-Host "❌ ERROR: Not in a Git repository!" -ForegroundColor Red
    exit 1
}
Write-Host "✓ Git repository confirmed" -ForegroundColor Green

# Step 2: Check for uncommitted changes
Write-Host "`n📋 Step 2: Checking for uncommitted changes..." -ForegroundColor Cyan
$status = git status --porcelain
if ($status) {
    Write-Host "❌ ERROR: You have uncommitted changes!" -ForegroundColor Red
    Write-Host "Please commit or stash your changes first." -ForegroundColor Red
    git status
    exit 1
}
Write-Host "✓ Working directory clean" -ForegroundColor Green

# Step 3: Create backup
Write-Host "`n📋 Step 3: Creating backup..." -ForegroundColor Cyan
$backupDir = "../docmind-backup-$(Get-Date -Format 'yyyyMMdd-HHmmss')"
if ($DryRun) {
    Write-Host "[DRY RUN] Would create backup at: $backupDir" -ForegroundColor Yellow
} else {
    git clone --mirror . $backupDir
    Write-Host "✓ Backup created at: $backupDir" -ForegroundColor Green
}

# Step 4: Show what will be cleaned
Write-Host "`n📋 Step 4: Files to be cleaned from history..." -ForegroundColor Cyan
Write-Host "  • deployments/docker/docker-compose.yml (contains default passwords)" -ForegroundColor Yellow

# Step 5: Perform cleanup
Write-Host "`n📋 Step 5: Cleaning Git history..." -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "[DRY RUN] Would remove default passwords from docker-compose.yml history" -ForegroundColor Yellow
} else {
    # Method 1: Using git filter-branch (built-in, slower but available everywhere)
    Write-Host "Using git filter-branch to clean history..." -ForegroundColor Yellow
    
    # Remove the file from history
    git filter-branch --force --index-filter `
        "git rm --cached --ignore-unmatch deployments/docker/docker-compose.yml" `
        --prune-empty --tag-name-filter cat -- --all
    
    Write-Host "✓ History cleaned" -ForegroundColor Green
}

# Step 6: Clean up refs
Write-Host "`n📋 Step 6: Cleaning up references..." -ForegroundColor Cyan
if ($DryRun) {
    Write-Host "[DRY RUN] Would clean up refs" -ForegroundColor Yellow
} else {
    git for-each-ref --format="delete %(refname)" refs/original | git update-ref --stdin
    git reflog expire --expire=now --all
    git gc --prune=now --aggressive
    Write-Host "✓ References cleaned" -ForegroundColor Green
}

# Step 7: Re-add the corrected file
Write-Host "`n📋 Step 7: Re-adding corrected docker-compose.yml..." -ForegroundColor Cyan
if ($DryRun) {
    Write-Host "[DRY RUN] Would re-add corrected file" -ForegroundColor Yellow
} else {
    git add deployments/docker/docker-compose.yml
    git commit -m "fix: add docker-compose.yml without default passwords"
    Write-Host "✓ Corrected file added" -ForegroundColor Green
}

# Step 8: Instructions for force push
Write-Host "`n📋 Step 8: Next steps..." -ForegroundColor Cyan
Write-Host ""
if ($DryRun) {
    Write-Host "[DRY RUN] Completed dry run. No changes were made." -ForegroundColor Yellow
    Write-Host "Run without -DryRun to perform actual cleanup." -ForegroundColor Yellow
} else {
    Write-Host "✅ Git history has been cleaned locally!" -ForegroundColor Green
    Write-Host ""
    Write-Host "To complete the cleanup, you need to force push:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  git push origin --force --all" -ForegroundColor Cyan
    Write-Host "  git push origin --force --tags" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "⚠️  IMPORTANT: After force pushing:" -ForegroundColor Yellow
    Write-Host "  • Notify all collaborators" -ForegroundColor Yellow
    Write-Host "  • They must delete their local clones and re-clone" -ForegroundColor Yellow
    Write-Host "  • Or they can run: git fetch origin && git reset --hard origin/main" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Backup location: $backupDir" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Script completed!" -ForegroundColor Green
