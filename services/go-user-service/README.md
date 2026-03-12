# DocMind Go User Service

Go microservice responsible for user management, authentication, document upload, and access control.

## Features

- User registration and authentication (JWT)
- Document upload and management
- MinIO integration for file storage
- PostgreSQL for metadata storage
- gRPC client for Python RAG service communication
- Comprehensive error handling and logging
- Document access control and pagination

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 16+
- MinIO
- Docker & Docker Compose (for running services)

### Setup

1. Install dependencies:
```bash
go mod download
```

2. Run migrations:
```bash
# Migrations will be applied automatically on startup
```

3. Run the service:
```bash
go run cmd/server/main.go
```

## Configuration

All configuration is done via environment variables. See `.env.example` in the project root.

Required environment variables:
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT token signing (minimum 32 characters)
- `MINIO_ENDPOINT`: MinIO server endpoint
- `MINIO_ROOT_USER`: MinIO root user
- `MINIO_ROOT_PASSWORD`: MinIO root password
- `MINIO_BUCKET`: MinIO bucket name
- `GO_SERVICE_PORT`: Server port (default: 8080)

## API Endpoints

### Health Check

#### GET /health
Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "service": "docmind-go-service"
}
```

### Authentication

#### POST /api/auth/register
Register a new user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2026-03-11T12:00:00Z",
  "updated_at": "2026-03-11T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid email or password format
- `409 Conflict`: Email already registered

**Example with curl:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}'
```

#### POST /api/auth/login
Authenticate a user and receive JWT token.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2026-03-11T12:00:00Z",
    "updated_at": "2026-03-11T12:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request format
- `401 Unauthorized`: Invalid credentials

**Example with curl:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}'
```

### User Management

All user endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

#### GET /api/users/me
Get the current authenticated user's profile.

**Response (200 OK):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2026-03-11T12:00:00Z",
  "updated_at": "2026-03-11T12:00:00Z"
}
```

**Example with curl:**
```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <token>"
```

#### PUT /api/users/me
Update the current user's profile.

**Request:**
```json
{
  "email": "newemail@example.com"
}
```

**Response (200 OK):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "created_at": "2026-03-11T12:00:00Z",
  "updated_at": "2026-03-11T12:00:00Z"
}
```

**Example with curl:**
```bash
curl -X PUT http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"newemail@example.com"}'
```

#### DELETE /api/users/me
Delete the current user account.

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

**Example with curl:**
```bash
curl -X DELETE http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <token>"
```

### Document Management

All document endpoints require authentication.

#### POST /api/documents
Upload a new document.

**Request:**
- Content-Type: multipart/form-data
- Field name: `file`
- Supported formats: PDF, TXT, DOCX, MD
- Maximum file size: 10MB

**Response (201 Created):**
```json
{
  "document_id": "660e8400-e29b-41d4-a716-446655440001",
  "filename": "document.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "status": "pending_indexing",
  "storage_path": "docmind-documents/550e8400-e29b-41d4-a716-446655440000/660e8400-e29b-41d4-a716-446655440001/document.pdf",
  "upload_date": "2026-03-11T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid file type or size exceeds limit
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Storage error

**Example with curl:**
```bash
curl -X POST http://localhost:8080/api/documents \
  -H "Authorization: Bearer <token>" \
  -F "file=@document.pdf"
```

#### GET /api/documents
List all documents for the authenticated user with pagination.

**Query Parameters:**
- `page` (optional, default: 1): Page number
- `page_size` (optional, default: 20): Items per page (max: 100)

**Response (200 OK):**
```json
{
  "documents": [
    {
      "document_id": "660e8400-e29b-41d4-a716-446655440001",
      "filename": "document.pdf",
      "file_size": 1024000,
      "mime_type": "application/pdf",
      "status": "pending_indexing",
      "storage_path": "docmind-documents/550e8400-e29b-41d4-a716-446655440000/660e8400-e29b-41d4-a716-446655440001/document.pdf",
      "upload_date": "2026-03-11T12:00:00Z"
    }
  ],
  "total_items": 1,
  "total_pages": 1,
  "current_page": 1,
  "page_size": 20
}
```

**Example with curl:**
```bash
curl -X GET "http://localhost:8080/api/documents?page=1&page_size=20" \
  -H "Authorization: Bearer <token>"
```

#### GET /api/documents/{id}
Get metadata for a specific document.

**Response (200 OK):**
```json
{
  "document_id": "660e8400-e29b-41d4-a716-446655440001",
  "filename": "document.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "status": "pending_indexing",
  "storage_path": "docmind-documents/550e8400-e29b-41d4-a716-446655440000/660e8400-e29b-41d4-a716-446655440001/document.pdf",
  "upload_date": "2026-03-11T12:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid document ID format
- `403 Forbidden`: Document belongs to another user
- `404 Not Found`: Document not found

**Example with curl:**
```bash
curl -X GET http://localhost:8080/api/documents/660e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <token>"
```

#### GET /api/documents/{id}/download
Download a document file.

**Response (200 OK):**
- Returns the file with appropriate Content-Type header
- Content-Disposition header set for download

**Error Responses:**
- `400 Bad Request`: Invalid document ID format
- `403 Forbidden`: Document belongs to another user
- `404 Not Found`: Document not found

**Example with curl:**
```bash
curl -X GET http://localhost:8080/api/documents/660e8400-e29b-41d4-a716-446655440001/download \
  -H "Authorization: Bearer <token>" \
  -o document.pdf
```

#### DELETE /api/documents/{id}
Delete a document and its associated file.

**Response (200 OK):**
```json
{
  "message": "Document deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid document ID format
- `403 Forbidden`: Document belongs to another user
- `404 Not Found`: Document not found

**Example with curl:**
```bash
curl -X DELETE http://localhost:8080/api/documents/660e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <token>"
```

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request format or validation error
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Access denied (e.g., accessing another user's document)
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate email)
- `500 Internal Server Error`: Server error

## Logging

The service uses structured logging with the following information:
- Timestamp
- Log level (INFO, ERROR, WARNING)
- Message
- Context (user_id, document_id, etc.)

All logs are written to stdout.

## Testing

Run tests with:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
services/go-user-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── domain/
│   │   ├── user.go              # User domain model
│   │   ├── document.go          # Document domain model
│   │   └── errors.go            # Custom error types
│   ├── handler/
│   │   └── http/
│   │       ├── auth_handler.go  # Authentication endpoints
│   │       ├── user_handler.go  # User management endpoints
│   │       ├── document_handler.go # Document endpoints
│   │       └── health_handler.go   # Health check endpoint
│   ├── infrastructure/
│   │   ├── db.go                # Database initialization
│   │   ├── jwt.go               # JWT token management
│   │   ├── logger.go            # Logging setup
│   │   └── minio.go             # MinIO client wrapper
│   ├── middleware/
│   │   ├── auth_middleware.go   # JWT authentication
│   │   ├── logging.go           # Request logging
│   │   └── recovery.go          # Panic recovery
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_postgres.go # User repository
│   │   │   └── document_postgres.go # Document repository
│   │   ├── user_repository.go   # User repository interface
│   │   └── document_repository.go # Document repository interface
│   └── service/
│       ├── auth_service.go      # Authentication logic
│       ├── user_service.go      # User management logic
│       └── document_service.go  # Document management logic
├── migrations/
│   ├── 001_create_users_table.sql
│   ├── 002_create_documents_table.sql
│   └── 003_create_documents_table.sql
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Dockerfile                   # Docker image definition
├── Makefile                     # Build and development commands
└── README.md                    # This file
```

## License

MIT License - See LICENSE file in project root
