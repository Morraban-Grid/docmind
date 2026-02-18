# API Specification

## Go Service

### POST /auth/register

Registers a new user.

Request:
{
  "email": "user@example.com",
  "password": "securepassword"
}

Response:
{
  "message": "User registered successfully"
}

---

### POST /auth/login

Authenticates a user and returns JWT.

---

### POST /documents/upload

Uploads a document.

Headers:
Authorization: Bearer <token>

Response:
{
  "document_id": "uuid",
  "status": "uploaded"
}

---

## Python RAG Service

### POST /ingest

Triggers document ingestion.

Request:
{
  "document_id": "uuid"
}

---

### POST /query

Queries the RAG system.

Request:
{
  "query": "What does the contract say about termination?"
}

Response:
{
  "answer": "..."
}
