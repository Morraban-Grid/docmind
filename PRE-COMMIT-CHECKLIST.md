# Pre-Commit Security Checklist

Before running `git add .` and `git push`, verify:

## ✅ Security Checks

- [ ] `.env` file exists locally but is NOT in git
- [ ] `.env.example` has placeholder values (not real secrets)
- [ ] `.kiro/` directory is ignored by git
- [ ] No hardcoded passwords in code
- [ ] No API keys in code
- [ ] No tokens in code
- [ ] `JWT_SECRET` is strong (32+ characters)
- [ ] Database passwords are strong
- [ ] MinIO passwords are strong

## ✅ File Checks

- [ ] `docs/specs/` contains requirements.md, design.md, tasks.md
- [ ] `README.md` is updated
- [ ] `SECURITY.md` exists
- [ ] `CONTRIBUTING.md` exists
- [ ] `CHANGELOG.md` is updated
- [ ] All scripts in `scripts/` are executable

## ✅ Docker Checks

- [ ] `docker-compose.yml` has no `version:` line (obsolete)
- [ ] All services have health checks
- [ ] No hardcoded credentials in docker-compose.yml

## ✅ Git Checks

Run these commands to verify:

```bash
# 1. Check .env is ignored
git status | grep ".env"
# Should only show .env.example

# 2. Check .kiro is ignored
git status | grep ".kiro"
# Should show nothing

# 3. Check for secrets
git diff | grep -i "password\|secret\|token\|api"
# Should only show references to env variables

# 4. View what will be committed
git status
```

## ✅ Final Verification

```bash
# Run verification script
./scripts/verify-setup.sh

# Check git status
git status

# If everything looks good:
git add .
git commit -m "feat: iteration 1 - infrastructure and security foundation"
git push origin main
```

## 🚨 If You Find Issues

**If .env is in git:**
```bash
git rm --cached .env
git commit -m "fix: remove .env from git"
```

**If .kiro is in git:**
```bash
git rm -r --cached .kiro/
git commit -m "fix: remove .kiro from git"
```

**If you committed secrets:**
1. Change ALL passwords immediately
2. Generate new JWT_SECRET
3. Consider git history rewrite (advanced)

---

**Remember**: Once you push, it's public. Double-check everything!
