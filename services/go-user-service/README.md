# DocMind Go User Service

Go microservice responsible for user management, authentication, document upload, and access control.

## Features

- User registration and authentication (JWT)
- Document upload and management
- MinIO integration for file storage
- PostgreSQL for metadata storage
- gRPC client for Python RAG service communication

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 16+
- MinIO

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
- `MINIO_ACCESS_KEY`: MinIO access key
- `MINIO_SECRET_KEY`: MinIO secret key
- `GO_SERVICE_PORT`: Server port (default: 8080)

## API Endpoints

### Authentication

#### Register User
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}

# Response (201 Created)
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}

# Example with curl
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123","name":"John Doe"}'
```

**Validation Rules:**
- Email must be valid format
- Password must be at least 8 characters
- Email must be unique (returns 409 if already exists)

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}

# Response (200 OK)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}

# Example with curl
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}'
```

**Authentication:**
- Returns JWT token valid for 24 hours
- Token must be included in Authorization header for protected endpoints
- Format: `Authorization: Bearer <token>`

### User Management

All user endpoints require authentication via JWT token.

#### Get Current User
```bash
GET /api/users/me
Authorization: Bearer <token>

# Response (200 OK)
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}

# Example with curl
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Current User
```bash
PUT /api/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe"
}

# Response (200 OK)
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "Jane Doe",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:30:00Z"
}

# Example with curl
curl -X PUT http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe"}'
```

#### Delete Current User
```bash
DELETE /api/users/me
Authorization: Bearer <token>

# Response (200 OK)
{
  "message": "User deleted successfully"
}

# Example with curl
curl -X DELETE http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Error Responses

All errors follow a standardized JSON format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` (400): Invalid input data
- `INVALID_CREDENTIALS` (401): Authentication failed
- `TOKEN_MISSING` (401): No authorization header provided
- `TOKEN_INVALID` (401): Invalid or expired JWT token
- `ACCESS_DENIED` (403): User lacks permission for resource
- `NOT_FOUND` (404): Resource not found
- `EMAIL_ALREADY_EXISTS` (409): Email already registered
- `INTERNAL_ERROR` (500): Server error

## Security

- Passwords are hashed using bcrypt with cost factor 12
- JWT tokens use HS256 algorithm
- All passwords are validated (minimum 8 characters)
- Email addresses are validated and normalized
- SQL injection prevention via prepared statements
- Password hashes are never returned in API responses

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/service -v
```

## Logging

The service uses structured logging (log/slog) with JSON output to stdout.

Log entries include:
- Timestamp
- Log level (INFO, WARN, ERROR)
- User ID (when applicable)
- Request context
- Error details

Example log entry:
```json
{
  "time": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "msg": "user created",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com"
}
```
