# Architecture Overview

## System Design

DocMind follows a microservices architecture composed of two independent services:

User Client
│
▼
Go API Service
│
├── PostgreSQL (metadata)
├── MinIO (file storage)
│
▼
Python RAG Service
│
▼
Qdrant (Vector Database)


---

## Service Responsibilities

### Go Service

- User registration and authentication
- JWT issuance
- Document upload handling
- Metadata storage
- Secure access control

### Python RAG Service

- Document parsing
- Chunking
- Embeddings generation
- Vector storage
- Semantic retrieval
- Response generation

---

## RAG Flow

1. User uploads document
2. Document stored in MinIO
3. Python service ingests document
4. Text is chunked
5. Embeddings generated
6. Vectors stored in Qdrant

Query flow:

1. User submits query
2. Query converted to embedding
3. Top-k relevant chunks retrieved
4. Context assembled
5. LLM generates response
