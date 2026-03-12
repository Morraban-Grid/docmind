# DocMind System Architecture

## Overview

DocMind is a microservices-based RAG (Retrieval-Augmented Generation) system for private document reasoning. The system allows users to upload documents, which are then processed, embedded, and stored in a vector database for semantic search and question answering.

## System Components

### 1. Go User Service
**Responsibilities:**
- User authentication and authorization (JWT)
- Document upload and management
- File storage in MinIO
- Access control and permissions
- gRPC client for Python service

**Technology:**
- Language: Go 1.21
- Framework: Gin
- Database: PostgreSQL (via database/sql)
- Object Storage: MinIO
- Authentication: JWT

### 2. Python RAG Service
**Responsibilities:**
- Document text extraction (PDF, DOCX, TXT, MD)
- Text chunking with LangChain
- Embedding generation (sentence-transformers)
- Vector storage in Qdrant
- Semantic search
- RAG query processing with Ollama LLM
- gRPC server for Go service

**Technology:**
- Language: Python 3.12
- Framework: FastAPI
- Embeddings: sentence-transformers (all-MiniLM-L6-v2)
- Vector DB: Qdrant
- LLM: Ollama (llama2)
- Document Processing: LangChain, PyPDF2, python-docx

### 3. Infrastructure Services

**PostgreSQL 16.4:**
- User metadata storage
- Document metadata storage
- Relational data management

**Qdrant v1.11.3:**
- Vector embeddings storage
- Semantic similarity search
- User-isolated collections

**MinIO:**
- S3-compatible object storage
- Original document files storage
- Scalable file management

**Ollama:**
- Local LLM inference
- llama2 model
- Response generation

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Client                              │
│                    (Web/Mobile App)                         │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Go User Service                           │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  • Authentication (JWT)                              │  │
│  │  • User Management                                   │  │
│  │  • Document Upload                                   │  │
│  │  • Access Control                                    │  │
│  └──────────────────────────────────────────────────────┘  │
└─────┬──────────────┬──────────────┬──────────────┬─────────┘
      │              │              │              │
      │ SQL          │ S3 API       │ gRPC         │
      │              │              │              │
┌─────▼─────┐  ┌─────▼─────┐  ┌────▼──────┐  ┌───▼──────────┐
│PostgreSQL │  │   MinIO   │  │  Python   │  │   Qdrant     │
│           │  │           │  │    RAG    │  │              │
│  Users    │  │Documents  │  │  Service  │  │  Vectors     │
│Documents  │  │  Storage  │  │           │  │              │
└───────────┘  └───────────┘  └────┬──────┘  └──────────────┘
                                    │
                                    │ HTTP API
                                    │
                              ┌─────▼─────┐
                              │  Ollama   │
                              │   LLM     │
                              │  (llama2) │
                              └───────────┘
```

## Data Flow

### Document Upload Flow
1. User authenticates with Go service (JWT)
2. User uploads document via REST API
3. Go service validates file (type, size)
4. Go service stores file in MinIO
5. Go service saves metadata in PostgreSQL
6. Go service calls Python service via gRPC (IndexDocument)
7. Python service extracts text from document
8. Python service chunks text (1000 chars, 200 overlap)
9. Python service generates embeddings
10. Python service stores vectors in Qdrant
11. Python service returns success to Go service
12. Go service updates document status to "indexed"

### Query Flow
1. User authenticates with Go service (JWT)
2. User submits query via REST API
3. Go service calls Python service via gRPC (QueryDocument)
4. Python service generates query embedding
5. Python service searches Qdrant for similar vectors
6. Python service retrieves top 5 relevant chunks
7. Python service constructs prompt with context
8. Python service calls Ollama LLM for response
9. Python service returns answer with sources
10. Go service returns response to user

## Security Architecture

### Authentication
- JWT tokens with HS256 signing
- 24-hour token expiration
- Secure password hashing with bcrypt (cost factor 12)

### Authorization
- User-level document isolation
- Document ownership verification on all operations
- User-specific Qdrant collections

### Data Protection
- Environment variables for secrets (.env)
- .gitignore prevents credential commits
- No hardcoded credentials in code
- SQL injection prevention (prepared statements)
- File upload validation (type, size)

## Scalability Considerations

### Horizontal Scaling
- Go service: Stateless, can scale horizontally
- Python service: Can add multiple replicas
- PostgreSQL: Can use read replicas
- Qdrant: Can partition by user_id ranges

### Performance Optimization
- Connection pooling (PostgreSQL: 5-25 connections)
- Batch embedding generation (32 chunks)
- LLM concurrency limiting (3 concurrent requests)
- Async I/O in Python service

## Deployment

### Docker Compose (Development)
- All services in single docker-compose.yml
- Shared Docker network
- Volume persistence for data
- Health checks for all services

### Kubernetes (Production - Future)
- Separate deployments for each service
- Horizontal Pod Autoscaling
- Persistent Volume Claims
- Service mesh for inter-service communication

## Monitoring and Observability

### Health Checks
- Go service: /health endpoint
- Python service: /health endpoint
- Database connectivity checks
- Dependency status reporting

### Logging
- Structured logging (JSON format)
- Log levels: DEBUG, INFO, WARNING, ERROR
- Request ID tracing across services
- Stdout logging for Docker capture

### Metrics (Future)
- Prometheus metrics endpoint
- Request count, latency, error rate
- Resource utilization
- Custom business metrics

## Technology Choices Rationale

### Why Go for User Service?
- Fast compilation and execution
- Excellent concurrency support
- Strong typing and error handling
- Great for API services

### Why Python for RAG Service?
- Rich ML/AI ecosystem
- LangChain for RAG pipelines
- sentence-transformers for embeddings
- Easy integration with ML models

### Why Qdrant?
- Purpose-built for vector search
- High performance
- Easy to use
- Free and open source

### Why Ollama?
- Local LLM inference (no API costs)
- Easy model management
- Good performance
- Privacy-preserving

### Why MinIO?
- S3-compatible API
- Self-hosted
- Scalable
- Free and open source

## Future Enhancements

1. **Multi-tenancy**: Organization-level isolation
2. **Advanced Search**: Hybrid search (vector + keyword)
3. **Document Versioning**: Track document changes
4. **Collaborative Features**: Document sharing
5. **Advanced Analytics**: Usage statistics, popular queries
6. **Model Fine-tuning**: Custom embedding models
7. **Multi-language Support**: i18n for UI and documents
8. **Real-time Updates**: WebSocket for live notifications
