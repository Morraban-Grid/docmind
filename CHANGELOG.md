# Changelog

All notable changes to the DocMind RAG System will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Iteration 2 - Go Service Authentication & User Management (2026-03-11)

#### Added
- Complete user authentication system with JWT tokens
- User registration endpoint with email validation
- User login endpoint with bcrypt password hashing (cost 12)
- Protected user endpoints (GET/PUT/DELETE /api/users/me)
- JWT middleware for route protection
- Recovery middleware for panic handling
- Structured logging with slog
- PostgreSQL user repository with prepared statements
- Domain models for User and custom errors
- Health check endpoint
- Integration tests for authentication flow

#### Security
- JWT authentication with HS256 algorithm
- Bcrypt password hashing with cost factor 12
- No hardcoded credentials (all from environment variables)
- Input validation for email and password
- SQL injection protection with prepared statements
- Secure error handling without exposing sensitive information

#### Technical Details
- Go service structure following clean architecture
- Gin web framework for HTTP routing
- PostgreSQL for user data storage
- JWT token expiration configurable via environment
- Comprehensive error handling and logging

---

### Iteration 1 - Infrastructure & Security Foundation (2026-03-10)

#### Added
- Docker Compose infrastructure setup
  - PostgreSQL database
  - Qdrant vector database
  - MinIO object storage
  - Ollama LLM service
- Go service skeleton with configuration management
- Python RAG service skeleton with FastAPI
- Database migrations for users and documents tables
- Comprehensive `.gitignore` for security
- Environment variable management (`.env`, `.env.example`)
- Security documentation (`SECURITY.md`)
- Pre-commit checklist (`PRE-COMMIT-CHECKLIST.md`)
- Project documentation (`README.md`)
- GitHub Actions workflows for security scanning

#### Security
- Secret scanning with Gitleaks
- Dependency scanning with CodeQL
- Protected `.env` and `.kiro/` directories
- No credentials in version control
- Placeholder values in `.env.example`

#### Infrastructure
- Multi-service Docker Compose setup
- Health checks for all services
- Volume persistence for databases
- Network isolation between services
- Port mapping for local development

---

## Version History

- **Iteration 2** (2026-03-11): Go Service Authentication & User Management
- **Iteration 1** (2026-03-10): Infrastructure & Security Foundation
