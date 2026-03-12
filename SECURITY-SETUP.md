# Security Setup Guide - DocMind

## ⚠️ CRITICAL: Before First Commit

This guide ensures you never commit sensitive information to version control.

## 1. Environment Variables Setup

### Step 1: Copy the example file
```bash
cp .env.example .env
```

### Step 2: Generate secure secrets

**NEVER use the example values in production or even development!**

Generate secure random secrets:

```bash
# Generate JWT Secret (Linux/Mac)
openssl rand -base64 32

# Generate JWT Secret (Windows PowerShell)
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))

# Generate PostgreSQL Password
openssl rand -base64 24

# Generate MinIO Password
openssl rand -base64 24
```

### Step 3: Update .env file

Edit your `.env` file and replace ALL placeholder values:

```env
# PostgreSQL - CHANGE THESE!
POSTGRES_USER=your_custom_username
POSTGRES_PASSWORD=<paste-generated-password-here>
POSTGRES_DB=your_custom_db_name

# MinIO - CHANGE THESE!
MINIO_ROOT_USER=your_custom_minio_user
MINIO_ROOT_PASSWORD=<paste-generated-password-here>

# JWT - CHANGE THIS!
JWT_SECRET=<paste-generated-secret-here>
```

## 2. Verify .gitignore

Ensure `.gitignore` contains:

```gitignore
# CRITICAL - Never commit these
.env
.env.*
!.env.example
```

## 3. Pre-Commit Checklist

Before EVERY commit, verify:

```bash
# 1. Check no .env files are staged
git status | grep -E "\.env$"
# Should return nothing

# 2. Check for hardcoded secrets
git diff --cached | grep -iE "(password|secret|token|api_key)\s*=\s*['\"][^$]"
# Should return nothing

# 3. Run Gitleaks locally (optional but recommended)
docker run -v $(pwd):/path zricethezav/gitleaks:latest detect --source="/path" -v
```

## 4. What to NEVER Commit

❌ **NEVER commit:**
- `.env` files (except `.env.example`)
- Actual passwords, tokens, or API keys
- Database credentials
- JWT secrets
- MinIO/AWS access keys
- Private keys or certificates
- Any file containing real credentials

✅ **Safe to commit:**
- `.env.example` with placeholder values like `changeme_*`
- Configuration templates
- Documentation
- Code that reads from environment variables

## 5. If You Accidentally Commit Secrets

If you accidentally commit secrets:

1. **Immediately rotate all exposed credentials**
2. **Remove from Git history:**
   ```bash
   # Use BFG Repo-Cleaner or git-filter-repo
   git filter-repo --path .env --invert-paths
   ```
3. **Force push (if safe):**
   ```bash
   git push --force
   ```
4. **Notify your team**

## 6. GitHub Actions Security

The repository includes automated security scanning:

- **Gitleaks**: Scans for secrets in commits
- **Trivy**: Scans for vulnerabilities in dependencies

These run on every push and pull request.

## 7. Docker Compose Security

The `docker-compose.yml` file requires environment variables:

```bash
# Start services (will fail if .env is missing or incomplete)
docker-compose -f deployments/docker/docker-compose.yml up

# The compose file does NOT have default passwords
# You MUST provide them via .env file
```

## 8. Production Deployment

For production:

1. Use a secrets management service (AWS Secrets Manager, HashiCorp Vault, etc.)
2. Never store secrets in environment variables on the host
3. Use strong, unique passwords (minimum 32 characters)
4. Rotate secrets regularly
5. Enable audit logging
6. Use TLS/SSL for all connections

## 9. Security Contacts

If you discover a security vulnerability:

1. **DO NOT** open a public issue
2. Email: [security contact - add your email]
3. Include: description, impact, reproduction steps

## 10. Compliance Checklist

Before pushing to main:

- [ ] No `.env` files in git status
- [ ] All passwords are strong and unique
- [ ] JWT secret is at least 32 characters
- [ ] No hardcoded credentials in code
- [ ] `.gitignore` is properly configured
- [ ] Gitleaks scan passes
- [ ] All secrets are in `.env` file (not committed)

---

**Remember: Security is everyone's responsibility. When in doubt, ask!**
