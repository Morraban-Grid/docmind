# DocMind RAG System - API Documentation

## Overview

Complete API documentation for the DocMind RAG System. All endpoints are RESTful and return JSON responses.

## Base URL

```
http://localhost:8080/api
```

## Authentication

All protected endpoints require JWT authentication via the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

## Error Response Format

All errors return a consistent JSON format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "timestamp": "2026-03-12T10:30:00Z",
  "request_id": "req-12345-abcde"
}
```

## Authentication Endpoints

### Register User

**Endpoint**: `POST /auth/register`

**Request**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (201):
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "created_at": "2026-03-12T10:30:00Z"
}
```

**Errors**:
- 400: Invalid email or password
- 409: Email already exists

### Login

**Endpoint**: `POST /auth/login`

**Request**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (200):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 86400
}
```

**Errors**:
- 401: Invalid credentials

## User Endpoints

### Get Current User

**Endpoint**: `GET /users/me`

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "created_at": "2026-03-12T10:30:00Z",
  "updated_at": "2026-03-12T10:30:00Z"
}
```

### Update Current User

**Endpoint**: `PUT /users/me`

**Headers**: `Authorization: Bearer <token>`

**Request**:
```json
{
  "email": "newemail@example.com"
}
```

**Response** (200):
```json
{
  "id": "user-uuid",
  "email": "newemail@example.com",
  "updated_at": "2026-03-12T10:35:00Z"
}
```

### Delete Current User

**Endpoint**: `DELETE /users/me`

**Headers**: `Authorization: Bearer <token>`

**Response** (204): No content

## Document Endpoints

### Upload Document

**Endpoint**: `POST /documents`

**Headers**: `Authorization: Bearer <token>`

**Request**: multipart/form-data
- `file`: Document file (PDF, DOCX, TXT, MD)
- Max size: 10MB

**Response** (201):
```json
{
  "id": "doc-uuid",
  "filename": "document.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "status": "pending_indexing",
  "created_at": "2026-03-12T10:30:00Z"
}
```

**Errors**:
- 400: Invalid file type or size
- 401: Unauthorized
- 413: File too large

### List Documents

**Endpoint**: `GET /documents?page=1&page_size=20`

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20)

**Response** (200):
```json
{
  "documents": [
    {
      "id": "doc-uuid",
      "filename": "document.pdf",
      "file_size": 1024000,
      "status": "indexed",
      "created_at": "2026-03-12T10:30:00Z"
    }
  ],
  "pagination": {
    "total_items": 50,
    "total_pages": 3,
    "current_page": 1,
    "page_size": 20
  }
}
```

### Get Document

**Endpoint**: `GET /documents/{id}`

**Headers**: `Authorization: Bearer <token>`

**Response** (200):
```json
{
  "id": "doc-uuid",
  "filename": "document.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "status": "indexed",
  "created_at": "2026-03-12T10:30:00Z"
}
```

**Errors**:
- 404: Document not found
- 403: Unauthorized access

### Download Document

**Endpoint**: `GET /documents/{id}/download`

**Headers**: `Authorization: Bearer <token>`

**Response** (200): File content

**Errors**:
- 404: Document not found
- 403: Unauthorized access

### Delete Document

**Endpoint**: `DELETE /documents/{id}`

**Headers**: `Authorization: Bearer <token>`

**Response** (204): No content

**Errors**:
- 404: Document not found
- 403: Unauthorized access

## Query Endpoints

### Query Documents

**Endpoint**: `POST /query`

**Headers**: `Authorization: Bearer <token>`

**Request**:
```json
{
  "query": "What is the main topic of the document?"
}
```

**Response** (200):
```json
{
  "answer": "The document discusses...",
  "sources": ["doc-uuid-1", "doc-uuid-2"],
  "chunk_count": 3
}
```

**Errors**:
- 400: Empty query
- 401: Unauthorized
- 503: LLM service unavailable

## Health Endpoints

### Health Check

**Endpoint**: `GET /health`

**Response** (200):
```json
{
  "status": "healthy",
  "service": "docmind-go-service",
  "dependencies": [
    {
      "name": "database",
      "status": "healthy",
      "available": true
    },
    {
      "name": "grpc",
      "status": "healthy",
      "available": true
    }
  ]
}
```

## Status Codes

- `200`: Success
- `201`: Created
- `204`: No Content
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `409`: Conflict
- `413`: Payload Too Large
- `500`: Internal Server Error
- `503`: Service Unavailable

## Rate Limiting

No rate limiting implemented in MVP. Can be added in future iterations.

## Pagination

All list endpoints support pagination:
- `page`: Page number (1-indexed)
- `page_size`: Items per page (1-100, default 20)

## Sorting

Sorting not implemented in MVP. Can be added in future iterations.

## Filtering

Filtering not implemented in MVP. Can be added in future iterations.

## Versioning

API version: v1.0.0

No versioning in URL. Breaking changes will increment major version.

## CORS

CORS is enabled for all origins in development. Configure appropriately for production.

## Timeouts

- Document upload: 30 seconds
- Document query: 45 seconds
- gRPC calls: 30-45 seconds

## Request ID Tracing

All requests include a unique `request_id` for tracing:

```
X-Request-ID: req-12345-abcde
```

This ID is included in all log entries and error responses.
