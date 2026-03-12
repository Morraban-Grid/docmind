# Design Document - DocMind RAG System

## Overview

DocMind is a professional document reasoning system based on RAG (Retrieval-Augmented Generation) architecture. The system enables authenticated users to upload private documents, transform them into vector representations through embeddings, and perform natural language queries that combine semantic retrieval with contextualized response generation.

The system is designed as a microservices architecture with two main services:

- **Go Service**: Handles user management, authentication, file uploads, and access control
- **Python Service**: Implements the complete RAG pipeline (chunking, embeddings, semantic retrieval, generation)

Key architectural decisions:
- **Communication**: gRPC for inter-service communication to minimize latency
- **Deployment**: Fully local infrastructure using only free tools
- **Security**: JWT-based authentication, password hashing, user-level document isolation
- **Scalability**: Stateless services, horizontal scaling capability
- **Storage**: PostgreSQL for metadata, MinIO for objects, Qdrant for vectors

## Architecture

### System Architecture Diagram

```mermaid
graph TB
    Client[Client Application]
    
    subgraph "Go Service"
        API[REST API<br/>Gin Framework]
        Auth[Authentication<br/>JWT]
        FileHandler[File Handler]
        UserMgmt[User Management]
        DB_Client[Database Client<br/>GORM]
        Storage_Client[MinIO Client]
        GRPC_Client[gRPC Client]
    end
    
    subgraph "Python Service"
        GRPC_Server[gRPC Server<br/>FastAPI]
        TextExtractor[Text Extractor]
        Chunker[Document Chunker<br/>LangChain]
        Embedder[Embedding Generator<br/>sentence-transformers]
        VectorDB_Client[Qdrant Client]
        LLM_Client[Ollama Client<br/>llama2]
    end
    
    subgraph "Data Layer"
        PostgreSQL[(PostgreSQL<br/>Metadata)]
        MinIO[(MinIO<br/>Object Storage)]
        Qdrant[(Qdrant<br/>Vector Store)]
        Ollama[Ollama<br/>LLM Service]
    end
    
    Client -->|HTTP/REST| API
    API --> Auth
    API --> FileHandler
    API --> UserMgmt
    UserMgmt --> DB_Client
    FileHandler --> Storage_Client
    FileHandler --> GRPC_Client
    
    GRPC_Client -->|gRPC| GRPC_Server
    GRPC_Server --> TextExtractor
    TextExtractor --> Chunker
    Chunker --> Embedder
    Embedder --> VectorDB_Client
    GRPC_Server --> VectorDB_Client
    GRPC_Server --> LLM_Client
    
    DB_Client --> PostgreSQL
    Storage_Client --> MinIO
    VectorDB_Client --> Qdrant
    LLM_Client --> Ollama
```

### Deployment Architecture

```mermaid
graph TB
    subgraph "Docker Compose Environment"
        subgraph "Application Services"
            GoService[go-service:8080]
            PythonService[python-service:50051]
        end
        
        subgraph "Storage Services"
            PostgreSQL[postgresql:5432]
            MinIO[minio:9000]
            Qdrant[qdrant:6333]
            Ollama[ollama:11434]
        end
        
        GoService -->|TCP| PostgreSQL
        GoService -->|HTTP| MinIO
        GoService -->|gRPC| PythonService
        PythonService -->|HTTP| Qdrant
        PythonService -->|HTTP| Ollama
    end
    
    Host[Host Machine] -->|HTTP:8080| GoService
```

### Sequence Diagrams

#### Document Upload Flow

```mermaid
sequenceDiagram
    participant Client
    participant GoService
    participant PostgreSQL
    participant MinIO
    participant PythonService
    participant Qdrant
    
    Client->>GoService: POST /api/documents (file, JWT)
    GoService->>GoService: Validate JWT
    GoService->>GoService: Validate file extension
    GoService->>GoService: Validate file size (<10MB)
    GoService->>MinIO: Upload file
    MinIO-->>GoService: storage_path
    GoService->>PostgreSQL: Insert document metadata
    PostgreSQL-->>GoService: document_id
    GoService->>PythonService: IndexDocument(document_id, user_id, storage_path)
    
    PythonService->>MinIO: Download file
    MinIO-->>PythonService: file content
    PythonService->>PythonService: Extract text
    PythonService->>PythonService: Chunk document
    PythonService->>PythonService: Generate embeddings
    PythonService->>Qdrant: Store vectors with metadata
    Qdrant-->>PythonService: success
    PythonService-->>GoService: IndexDocumentResponse(success=true)
    
    GoService-->>Client: 201 Created {document_id}
```

#### Query Document Flow

```mermaid
sequenceDiagram
    participant Client
    participant GoService
    participant PythonService
    participant Qdrant
    participant Ollama
    
    Client->>GoService: POST /api/query (query, JWT)
    GoService->>GoService: Validate JWT
    GoService->>GoService: Extract user_id from JWT
    GoService->>PythonService: QueryDocument(user_id, query)
    
    PythonService->>PythonService: Generate query embedding
    PythonService->>Qdrant: Search similar vectors (filter by user_id)
    Qdrant-->>PythonService: Top 5 chunks (score > 0.7)
    
    alt Chunks found
        PythonService->>PythonService: Build context from chunks
        PythonService->>Ollama: Generate response with context
        Ollama-->>PythonService: Generated answer
        PythonService-->>GoService: QueryDocumentResponse(answer, sources)
    else No chunks found
        PythonService-->>GoService: QueryDocumentResponse("No information found")
    end
    
    GoService-->>Client: 200 OK {answer, sources}
```

## Components and Interfaces

### Go Service Components

#### REST API Endpoints

**Authentication**
```
POST /api/auth/register
Request: {email: string, password: string, name: string}
Response: {user_id: string, message: string}
Status: 201 Created | 400 Bad Request | 409 Conflict

POST /api/auth/login
Request: {email: string, password: string}
Response: {token: string, user_id: string}
Status: 200 OK | 401 Unauthorized

GET /api/auth/me
Headers: Authorization: Bearer <token>
Response: {user_id: string, email: string, name: string, created_at: timestamp}
Status: 200 OK | 401 Unauthorized
```

**User Management**
```
GET /api/users/:id
Headers: Authorization: Bearer <token>
Response: {user_id: string, email: string, name: string, created_at: timestamp}
Status: 200 OK | 401 Unauthorized | 404 Not Found

PUT /api/users/:id
Headers: Authorization: Bearer <token>
Request: {name?: string, email?: string}
Response: {user_id: string, email: string, name: string}
Status: 200 OK | 400 Bad Request | 401 Unauthorized | 404 Not Found

DELETE /api/users/:id
Headers: Authorization: Bearer <token>
Response: {message: string}
Status: 200 OK | 401 Unauthorized | 404 Not Found
```

**Document Management**
```
POST /api/documents
Headers: Authorization: Bearer <token>
Content-Type: multipart/form-data
Request: file (PDF|TXT|DOCX|MD, max 10MB)
Response: {document_id: string, filename: string, status: string}
Status: 201 Created | 400 Bad Request | 401 Unauthorized

GET /api/documents
Headers: Authorization: Bearer <token>
Query: ?page=1&page_size=20
Response: {
  documents: [{document_id, filename, upload_date, file_size, status}],
  pagination: {total_items, total_pages, current_page, page_size}
}
Status: 200 OK | 401 Unauthorized

GET /api/documents/:id
Headers: Authorization: Bearer <token>
Response: {document_id, filename, upload_date, file_size, status, storage_path}
Status: 200 OK | 401 Unauthorized | 403 Forbidden | 404 Not Found

DELETE /api/documents/:id
Headers: Authorization: Bearer <token>
Response: {message: string}
Status: 200 OK | 401 Unauthorized | 403 Forbidden | 404 Not Found
```

**Query**
```
POST /api/query
Headers: Authorization: Bearer <token>
Request: {query: string}
Response: {
  answer: string,
  sources: [{document_id: string, content: string, score: float}]
}
Status: 200 OK | 400 Bad Request | 401 Unauthorized
```

**Health Check**
```
GET /health
Response: {
  status: "healthy" | "degraded" | "unhealthy",
  services: {
    database: "up" | "down",
    storage: "up" | "down",
    python_service: "up" | "down"
  },
  timestamp: string
}
Status: 200 OK | 503 Service Unavailable
```

#### Authentication Middleware

The authentication middleware validates JWT tokens on all protected routes:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        // Validate token signature and expiration
        // Extract user_id from claims
        // Set user_id in context
        // Call next handler or return 401
    }
}
```

JWT Token Structure:
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "exp": 1234567890,
  "iat": 1234567890
}
```

Token expiration: 24 hours

#### File Handler

Responsibilities:
- Validate file extensions (PDF, TXT, DOCX, MD)
- Validate file size (max 10MB)
- Generate unique storage paths
- Upload to MinIO
- Create metadata records in PostgreSQL
- Trigger indexing via gRPC

#### gRPC Client

The Go service acts as a gRPC client to communicate with the Python service:

```protobuf
service RAGService {
    rpc IndexDocument (IndexDocumentRequest) returns (IndexDocumentResponse);
    rpc QueryDocument (QueryDocumentRequest) returns (QueryDocumentResponse);
}

message IndexDocumentRequest {
    string document_id = 1;
    string user_id = 2;
    string storage_path = 3;
}

message IndexDocumentResponse {
    bool success = 1;
    string message = 2;
}

message QueryDocumentRequest {
    string user_id = 1;
    string query = 2;
}

message QueryDocumentResponse {
    string answer = 1;
    repeated RetrievedChunk sources = 2;
}

message RetrievedChunk {
    string document_id = 1;
    string content = 2;
    float score = 3;
}
```

Timeout configuration:
- IndexDocument: 30 seconds
- QueryDocument: 45 seconds

### Python Service Components

#### gRPC Server

FastAPI-based gRPC server implementing the RAGService interface.

#### Text Extractor

Extracts text from different file formats:

- **PDF**: PyPDF2 or pdfplumber
- **TXT**: Direct read with UTF-8 encoding
- **DOCX**: python-docx library
- **MD**: Direct read preserving structure

#### Document Chunker

Uses LangChain's RecursiveCharacterTextSplitter:

Configuration:
- chunk_size: 1000 characters
- chunk_overlap: 200 characters
- max_chunk_size: 1200 characters

Output:
- List of text chunks
- Metadata: chunk_index, original_position

#### Embedding Generator

Uses sentence-transformers model: `all-MiniLM-L6-v2`

Specifications:
- Output dimensions: 384
- Batch size: 32 chunks
- Normalization: L2 norm = 1.0

Special cases:
- Empty chunks: zero vector

#### Vector Database Client

Qdrant client for vector storage and retrieval:

Collection configuration:
- Distance metric: Cosine similarity
- Vector size: 384
- Payload schema: {document_id, user_id, chunk_text, chunk_index}

Operations:
- Create collection (if not exists)
- Upsert vectors with payload
- Search with filters
- Delete by document_id

#### LLM Client

Ollama client for response generation:

Model: llama2

Configuration:
- Temperature: 0.7
- Max tokens: 500

Prompt template:
```
Context: {retrieved_chunks}

Question: {user_query}

Answer based only on the provided context. If the context doesn't contain relevant information, say so.
```

## Data Models

### PostgreSQL Schema

#### Users Table

```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

#### Documents Table

```sql
CREATE TABLE documents (
    document_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    storage_path VARCHAR(512) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending_indexing',
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    indexed_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_status ON documents(status);
```

Document status values:
- `pending_indexing`: Uploaded but not yet indexed
- `indexing`: Currently being processed
- `indexed`: Successfully indexed and queryable
- `failed`: Indexing failed

### MinIO Object Storage

Bucket structure:
```
docmind-documents/
  ├── {user_id}/
  │   ├── {document_id}.pdf
  │   ├── {document_id}.txt
  │   └── ...
```

### Qdrant Vector Store

Collection name: `docmind_embeddings`

Vector configuration:
```json
{
  "vectors": {
    "size": 384,
    "distance": "Cosine"
  }
}
```

Payload structure:
```json
{
  "document_id": "uuid",
  "user_id": "uuid",
  "chunk_text": "string",
  "chunk_index": "integer"
}
```

### Go Service Data Structures

```go
type User struct {
    UserID       string    `json:"user_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Email        string    `json:"email" gorm:"unique;not null"`
    PasswordHash string    `json:"-" gorm:"not null"`
    Name         string    `json:"name" gorm:"not null"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Document struct {
    DocumentID  string    `json:"document_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    UserID      string    `json:"user_id" gorm:"not null;type:uuid"`
    Filename    string    `json:"filename" gorm:"not null"`
    FileSize    int64     `json:"file_size" gorm:"not null"`
    StoragePath string    `json:"storage_path" gorm:"not null"`
    Status      string    `json:"status" gorm:"default:'pending_indexing'"`
    UploadDate  time.Time `json:"upload_date" gorm:"autoCreateTime"`
    IndexedAt   *time.Time `json:"indexed_at,omitempty"`
}

type JWTClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}
```

### Python Service Data Structures

```python
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class DocumentChunk:
    chunk_id: str
    document_id: str
    user_id: str
    chunk_text: str
    chunk_index: int
    embedding: Optional[List[float]] = None

@dataclass
class RetrievedChunk:
    document_id: str
    content: str
    score: float

@dataclass
class QueryResult:
    answer: str
    sources: List[RetrievedChunk]
```



## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

After analyzing all acceptance criteria, I've identified properties that can be validated through property-based testing. Some criteria relate to architecture or implementation details and are not suitable for property-based testing.

### Property 1: Email Uniqueness Enforcement

*For any* two user registration attempts with the same email address, the second registration SHALL fail with a conflict error.

**Validates: Requirements 1.2**

### Property 2: Password Hashing

*For any* user registration, the stored password hash SHALL never equal the plaintext password provided during registration.

**Validates: Requirements 1.3**

### Property 3: Password Omission in Responses

*For any* user query response, the password_hash field SHALL not be present in the JSON response body.

**Validates: Requirements 1.5**

### Property 4: Valid Credentials Generate JWT

*For any* valid user credentials (email and password), authentication SHALL produce a valid JWT token that can be decoded to reveal the user_id.

**Validates: Requirements 2.1, 2.2**

### Property 5: Invalid Credentials Rejection

*For any* invalid credentials (wrong email or wrong password), authentication SHALL fail with HTTP 401 status.

**Validates: Requirements 2.3**

### Property 6: Expired Token Rejection

*For any* expired or invalid JWT token, protected endpoints SHALL return HTTP 401 status.

**Validates: Requirements 2.5**

### Property 7: Supported File Extensions

*For any* file with extension PDF, TXT, DOCX, or MD, the upload SHALL be accepted (assuming valid authentication and size).

**Validates: Requirements 3.1**

### Property 8: Unsupported File Extension Rejection

*For any* file with an extension other than PDF, TXT, DOCX, or MD, the upload SHALL fail with HTTP 400 status.

**Validates: Requirements 3.2**

### Property 9: Upload Storage Round-Trip

*For any* successfully uploaded file, retrieving the file from object storage using the returned storage_path SHALL yield the original file content.

**Validates: Requirements 3.3**

### Property 10: Document Metadata Persistence

*For any* successful file upload, querying the database with the returned document_id SHALL return metadata containing the filename, user_id, upload_date, file_size, and storage_path.

**Validates: Requirements 3.4**

### Property 11: Document ID Uniqueness

*For any* two document uploads, their generated document_ids SHALL be different.

**Validates: Requirements 3.5**

### Property 12: Document Access Authorization

*For any* document access attempt, the request SHALL succeed only if the document's user_id matches the user_id in the JWT token.

**Validates: Requirements 4.1**

### Property 13: Cross-User Access Denial

*For any* attempt to access a document belonging to a different user, the request SHALL fail with HTTP 403 status.

**Validates: Requirements 4.2**

### Property 14: Document Deletion Completeness

*For any* document deletion, subsequent queries for that document_id SHALL return HTTP 404 status, and the file SHALL not exist in object storage.

**Validates: Requirements 4.4**

### Property 15: Text Extraction Preservation

*For any* UTF-8 text file, extracting text and comparing it to the original SHALL show that content is preserved.

**Validates: Requirements 5.2**

### Property 16: Chunk Size Limit

*For any* document chunking operation, all generated chunks SHALL have length ≤ 1200 characters.

**Validates: Requirements 6.3**

### Property 17: Chunk Metadata Presence

*For any* generated chunk, it SHALL have metadata including chunk_index and position information.

**Validates: Requirements 6.4**

### Property 18: Chunk ID Uniqueness and Order

*For any* document's chunks, they SHALL have unique chunk_ids and sequential chunk_index values (0, 1, 2, ...).

**Validates: Requirements 6.5**

### Property 19: Embedding Dimensionality

*For any* chunk converted to an embedding, the resulting vector SHALL have exactly 384 dimensions.

**Validates: Requirements 7.2**

### Property 20: Embedding Normalization

*For any* embedding vector stored in the system, its L2 norm SHALL equal 1.0 (±0.001 for floating point tolerance).

**Validates: Requirements 7.5**

### Property 21: Collection Existence After Indexing

*For any* document indexing operation, the Qdrant collection SHALL exist after the operation completes.

**Validates: Requirements 8.2**

### Property 22: Embedding Payload Completeness

*For any* embedding stored in Qdrant, retrieving it SHALL return payload containing document_id, user_id, chunk_text, and chunk_index fields.

**Validates: Requirements 8.3**

### Property 23: Query Embedding Dimensionality

*For any* user query, the generated query embedding SHALL have the same dimensionality (384) as document chunk embeddings.

**Validates: Requirements 9.1**

### Property 24: Search Result Limit

*For any* semantic search query, the number of returned chunks SHALL be ≤ 5.

**Validates: Requirements 9.2**

### Property 25: Search Result Privacy

*For any* query by a user, all returned chunks SHALL have user_id matching the requesting user's user_id.

**Validates: Requirements 9.3**

### Property 26: Search Result Similarity Threshold

*For any* semantic search results, all returned chunks SHALL have similarity score ≥ 0.7.

**Validates: Requirements 9.4**

### Property 27: Non-Empty Response Generation

*For any* successful LLM response generation with available context, the response SHALL be non-empty.

**Validates: Requirements 10.4**

### Property 28: Response Source References

*For any* query response that includes sources, each source SHALL include a valid document_id.

**Validates: Requirements 10.6**

### Property 29: Email Format Validation

*For any* email validation, strings not matching the email format pattern SHALL be rejected with HTTP 400 status.

**Validates: Requirements 16.1**

### Property 30: Password Length Validation

*For any* password validation, strings with length < 8 characters SHALL be rejected with HTTP 400 status.

**Validates: Requirements 16.2**

### Property 31: File Size Validation

*For any* file upload, files exceeding 10MB SHALL be rejected with HTTP 400 status.

**Validates: Requirements 16.3**

### Property 32: Empty Query Rejection

*For any* query validation, empty or whitespace-only queries SHALL be rejected with HTTP 400 status.

**Validates: Requirements 16.4**

### Property 33: UUID Validation

*For any* document_id validation, strings that are not valid UUIDs SHALL be rejected with HTTP 400 status.

**Validates: Requirements 16.5**

### Property 34: Validation Error Response Format

*For any* validation failure, the response SHALL have HTTP 400 status and include a descriptive error message in JSON format.

**Validates: Requirements 16.6**

### Property 35: Pagination Metadata Presence

*For any* paginated list response, the response SHALL include pagination metadata with total_items, total_pages, current_page, and page_size fields.

**Validates: Requirements 23.2**

### Property 36: Configuration Round-Trip

*For any* valid .env configuration, reading the file, modifying a value, writing it back, and reading again SHALL produce an equivalent configuration.

**Validates: Requirements 21.4**

### Property 37: Required Configuration Validation

*For any* configuration missing required environment variables, the system SHALL fail to start with a clear error message.

**Validates: Requirements 21.2**

### Property 38: Default Configuration Values

*For any* configuration missing optional environment variables, the system SHALL use documented default values.

**Validates: Requirements 21.3**

## Error Handling

### Error Response Format

All errors SHALL be returned in a consistent JSON format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

### Error Categories

#### Authentication Errors (401)
- `INVALID_CREDENTIALS`: Email or password incorrect
- `TOKEN_EXPIRED`: JWT token has expired
- `TOKEN_INVALID`: JWT token is malformed or has invalid signature
- `TOKEN_MISSING`: No authorization header provided

#### Authorization Errors (403)
- `ACCESS_DENIED`: User does not have permission to access resource
- `DOCUMENT_NOT_OWNED`: Attempting to access another user's document

#### Validation Errors (400)
- `INVALID_EMAIL`: Email format is invalid
- `PASSWORD_TOO_SHORT`: Password less than 8 characters
- `FILE_TOO_LARGE`: File exceeds 10MB limit
- `UNSUPPORTED_FILE_TYPE`: File extension not in [PDF, TXT, DOCX, MD]
- `INVALID_UUID`: Provided ID is not a valid UUID
- `EMPTY_QUERY`: Query string is empty or whitespace-only
- `MISSING_FIELD`: Required field missing from request

#### Resource Errors (404)
- `USER_NOT_FOUND`: User ID does not exist
- `DOCUMENT_NOT_FOUND`: Document ID does not exist

#### Conflict Errors (409)
- `EMAIL_ALREADY_EXISTS`: Email already registered

#### Server Errors (500)
- `DATABASE_ERROR`: Database operation failed
- `STORAGE_ERROR`: Object storage operation failed
- `GRPC_ERROR`: gRPC communication failed
- `INDEXING_ERROR`: Document indexing failed
- `EMBEDDING_ERROR`: Embedding generation failed
- `LLM_ERROR`: LLM response generation failed

### Error Handling Strategy

#### Go Service

```go
type AppError struct {
    Code       string                 `json:"code"`
    Message    string                 `json:"message"`
    Details    map[string]interface{} `json:"details,omitempty"`
    StatusCode int                    `json:"-"`
}

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            if appErr, ok := err.Err.(*AppError); ok {
                c.JSON(appErr.StatusCode, gin.H{"error": appErr})
            } else {
                c.JSON(500, gin.H{"error": AppError{
                    Code:    "INTERNAL_ERROR",
                    Message: "An unexpected error occurred",
                }})
            }
        }
    }
}
```

Retry strategy for gRPC calls:
- Max retries: 3
- Backoff: Exponential (1s, 2s, 4s)
- Retryable errors: UNAVAILABLE, DEADLINE_EXCEEDED

#### Python Service

```python
class ServiceError(Exception):
    def __init__(self, code: str, message: str, details: dict = None):
        self.code = code
        self.message = message
        self.details = details or {}

@app.exception_handler(ServiceError)
async def service_error_handler(request: Request, exc: ServiceError):
    return JSONResponse(
        status_code=500,
        content={
            "error": {
                "code": exc.code,
                "message": exc.message,
                "details": exc.details
            }
        }
    )
```

Retry strategy for external services:
- Qdrant: 3 retries with exponential backoff
- Ollama: 2 retries with 5s timeout each
- MinIO: 3 retries with exponential backoff

### Logging Strategy

All errors SHALL be logged with:
- Timestamp
- Error code
- Error message
- Stack trace (for 500 errors)
- Request context (user_id, request_id, endpoint)

Log levels:
- ERROR: All 500 errors, external service failures
- WARN: All 400 errors, authentication failures
- INFO: Successful operations, state changes
- DEBUG: Detailed operation traces (development only)

## Testing Strategy

### Dual Testing Approach

The system requires both unit testing and property-based testing for comprehensive coverage:

**Unit Tests**: Verify specific examples, edge cases, and error conditions
- Specific user registration scenarios
- Authentication with known credentials
- File upload with specific file types
- Error handling for specific failure modes
- Integration between components

**Property Tests**: Verify universal properties across all inputs
- All correctness properties defined in this document
- Comprehensive input coverage through randomization
- Validation of invariants and round-trip properties

Both approaches are complementary and necessary. Unit tests catch concrete bugs and validate specific scenarios, while property tests verify general correctness across the input space.

### Property-Based Testing Configuration

**Framework Selection**:
- **Go Service**: Use `gopter` library for property-based testing
- **Python Service**: Use `hypothesis` library for property-based testing

**Test Configuration**:
- Minimum iterations per property test: 100
- Each property test MUST reference its design document property
- Tag format: `Feature: docmind-rag-system, Property {number}: {property_text}`

**Example Property Test (Go)**:

```go
// Feature: docmind-rag-system, Property 11: Document ID Uniqueness
func TestProperty_DocumentIDUniqueness(t *testing.T) {
    parameters := gopter.DefaultTestParameters()
    parameters.MinSuccessfulTests = 100
    
    properties := gopter.NewProperties(parameters)
    
    properties.Property("Document IDs are always unique", prop.ForAll(
        func(count int) bool {
            ids := make(map[string]bool)
            for i := 0; i < count; i++ {
                id := generateDocumentID()
                if ids[id] {
                    return false // Duplicate found
                }
                ids[id] = true
            }
            return true
        },
        gen.IntRange(2, 100),
    ))
    
    properties.TestingRun(t)
}
```

**Example Property Test (Python)**:

```python
# Feature: docmind-rag-system, Property 19: Embedding Dimensionality
@given(text=st.text(min_size=1, max_size=1000))
@settings(max_examples=100)
def test_property_embedding_dimensionality(text):
    """For any chunk, embedding should have exactly 384 dimensions"""
    embedder = EmbeddingGenerator()
    embedding = embedder.generate_embedding(text)
    assert len(embedding) == 384
```

### Unit Testing Strategy

**Go Service Unit Tests**:
- User registration and authentication flows
- JWT token generation and validation
- File upload validation
- Document access control
- gRPC client error handling
- Database operations with test database
- MinIO operations with test bucket

**Python Service Unit Tests**:
- Text extraction from each file format
- Document chunking with known inputs
- Embedding generation with mock model
- Vector storage operations with test collection
- LLM response generation with mock LLM
- gRPC server endpoint handling

### Integration Testing

**End-to-End Flows**:
1. User registration → Login → Upload document → Query document
2. Document upload → Indexing → Semantic search → Response generation
3. Document deletion → Verify removal from all stores
4. Multi-user isolation → Verify privacy boundaries

**Test Environment**:
- Docker Compose test environment
- Separate test databases and collections
- Automated cleanup between tests

### Performance Testing

**Load Testing Scenarios**:
- Concurrent user registrations
- Simultaneous document uploads
- Parallel query processing
- Large document processing (near 10MB limit)

**Performance Targets**:
- Document upload: < 2s (excluding indexing)
- Query response: < 5s (including LLM generation)
- Authentication: < 100ms
- Semantic search: < 1s

### Security Testing

**Security Test Cases**:
- SQL injection attempts in all inputs
- JWT token tampering
- Cross-user document access attempts
- File upload with malicious content
- Oversized file uploads
- Path traversal in filenames
- XSS in document content

## Security Design

### Authentication

**Password Security**:
- Hashing algorithm: bcrypt with cost factor 12
- Minimum password length: 8 characters
- Password validation: Check against common password lists (optional)
- No password in logs or responses

**JWT Security**:
- Signing algorithm: HS256
- Secret key: 256-bit random key from environment variable
- Token expiration: 24 hours
- Claims: user_id, email, issued_at, expires_at
- Token refresh: Not implemented (re-login required)

### Authorization

**Access Control Model**:
- User-level isolation: Users can only access their own documents
- No role-based access control (all users have equal permissions)
- Document ownership: Enforced at database and application level

**Authorization Checks**:
```go
func CheckDocumentOwnership(userID, documentID string) error {
    doc, err := db.GetDocument(documentID)
    if err != nil {
        return ErrDocumentNotFound
    }
    if doc.UserID != userID {
        return ErrAccessDenied
    }
    return nil
}
```

### Data Privacy

**Vector Store Privacy**:
- All vector searches filtered by user_id
- No cross-user data leakage
- Payload includes user_id for filtering

**Object Storage Privacy**:
- Files organized by user_id in bucket structure
- Pre-signed URLs with short expiration (optional future enhancement)
- No public access to buckets

### Input Validation

**Validation Rules**:
- Email: RFC 5322 format validation
- Password: Minimum 8 characters, no maximum
- Filename: Alphanumeric, hyphens, underscores, dots only
- File size: Maximum 10MB
- Query: Non-empty, maximum 1000 characters
- UUIDs: Valid UUID v4 format

**Sanitization**:
- Strip leading/trailing whitespace from all string inputs
- Normalize email addresses to lowercase
- Remove null bytes from filenames
- Escape special characters in logs

### Network Security

**TLS/HTTPS**:
- Production deployment MUST use HTTPS
- TLS 1.2 or higher
- Certificate validation enabled

**gRPC Security**:
- Internal network only (not exposed to internet)
- Optional: mTLS for production deployments
- Timeout enforcement to prevent resource exhaustion

### Secrets Management

**Environment Variables**:
```bash
# Required secrets
JWT_SECRET=<256-bit-random-key>
POSTGRES_PASSWORD=<strong-password>
MINIO_ROOT_PASSWORD=<strong-password>

# Optional secrets
QDRANT_API_KEY=<api-key>  # If using cloud Qdrant
```

**Secret Storage**:
- Never commit secrets to version control
- Use .env file for local development (gitignored)
- Use secret management service for production (e.g., HashiCorp Vault, AWS Secrets Manager)

### Audit Logging

**Audit Events**:
- User registration and login
- Document upload and deletion
- Failed authentication attempts
- Authorization failures
- Configuration changes

**Log Format**:
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "event": "USER_LOGIN",
  "user_id": "uuid",
  "ip_address": "192.168.1.1",
  "success": true,
  "details": {}
}
```

## Deployment Architecture

### Docker Compose Configuration

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: docmind
      POSTGRES_USER: docmind
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U docmind"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 10s
      timeout: 5s
      retries: 5

  qdrant:
    image: qdrant/qdrant:latest
    volumes:
      - qdrant_data:/qdrant/storage
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6333/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  ollama:
    image: ollama/ollama:latest
    volumes:
      - ollama_data:/root/.ollama
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:11434/api/tags"]
      interval: 30s
      timeout: 10s
      retries: 3

  python-service:
    build:
      context: ./python-service
      dockerfile: Dockerfile
    environment:
      QDRANT_HOST: qdrant
      QDRANT_PORT: 6333
      OLLAMA_HOST: ollama
      OLLAMA_PORT: 11434
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: ${MINIO_ROOT_PASSWORD}
    depends_on:
      qdrant:
        condition: service_healthy
      ollama:
        condition: service_healthy
      minio:
        condition: service_healthy

  go-service:
    build:
      context: ./go-service
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://docmind:${POSTGRES_PASSWORD}@postgres:5432/docmind?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: ${MINIO_ROOT_PASSWORD}
      PYTHON_SERVICE_ADDR: python-service:50051
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_healthy
      python-service:
        condition: service_started

volumes:
  postgres_data:
  minio_data:
  qdrant_data:
  ollama_data:
```

### Environment Configuration

**.env.example**:
```bash
# Database
POSTGRES_PASSWORD=your_secure_password_here

# Object Storage
MINIO_ROOT_PASSWORD=your_secure_password_here

# Authentication
JWT_SECRET=your_256_bit_secret_key_here

# Optional: Logging
LOG_LEVEL=info

# Optional: Performance
EMBEDDING_BATCH_SIZE=32
MAX_WORKERS=4
```

### Startup Sequence

1. PostgreSQL starts and initializes database
2. MinIO starts and creates buckets
3. Qdrant starts and prepares vector storage
4. Ollama starts and downloads llama2 model
5. Python service starts and connects to dependencies
6. Go service starts and runs migrations
7. System ready for requests

### Resource Requirements

**Minimum Requirements**:
- CPU: 4 cores
- RAM: 8GB
- Disk: 20GB (plus document storage)

**Recommended Requirements**:
- CPU: 8 cores
- RAM: 16GB
- Disk: 100GB SSD

**Service Resource Allocation**:
- PostgreSQL: 512MB RAM
- MinIO: 512MB RAM
- Qdrant: 1GB RAM
- Ollama: 4GB RAM (for llama2 model)
- Python Service: 2GB RAM
- Go Service: 512MB RAM

### Monitoring and Observability

**Health Checks**:
- All services expose /health endpoints
- Docker healthchecks configured
- Startup dependencies enforced

**Metrics** (Future Enhancement):
- Prometheus metrics export
- Request latency histograms
- Error rate counters
- Resource utilization gauges

**Logging**:
- Structured JSON logs to stdout
- Log aggregation via Docker logging driver
- Centralized logging (optional: ELK stack, Loki)

### Backup Strategy

**Database Backup**:
```bash
# Daily PostgreSQL backup
docker exec postgres pg_dump -U docmind docmind > backup_$(date +%Y%m%d).sql
```

**Object Storage Backup**:
- MinIO bucket replication (optional)
- Periodic sync to external storage

**Vector Store Backup**:
- Qdrant snapshot API
- Periodic snapshots to persistent storage

### Disaster Recovery

**Recovery Time Objective (RTO)**: 1 hour
**Recovery Point Objective (RPO)**: 24 hours

**Recovery Procedure**:
1. Restore PostgreSQL from latest backup
2. Restore MinIO buckets from backup
3. Restore Qdrant snapshots
4. Restart all services
5. Verify system health
6. Re-index any documents uploaded since last backup

---

## Summary

This design document provides a comprehensive technical specification for the DocMind RAG System, including:

- Microservices architecture with Go and Python services
- Complete API specifications (REST and gRPC)
- Database schemas and data models
- 38 correctness properties for property-based testing
- Comprehensive error handling strategy
- Security design with authentication and authorization
- Docker-based deployment architecture
- Testing strategy combining unit and property-based tests

The system is designed to run entirely on local infrastructure using free tools while maintaining professional-grade security, scalability, and maintainability.
