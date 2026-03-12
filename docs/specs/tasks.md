# Implementation Plan: DocMind RAG System

## Overview

This implementation plan organizes the DocMind RAG System development into 8 iterations. Each iteration delivers a complete, testable increment that leaves the project in a deployable, git-safe state. The system consists of two microservices (Go and Python) communicating via gRPC, with PostgreSQL, MinIO, Qdrant, and Ollama as infrastructure dependencies.

## Iteration Strategy

- Iteration 1: Infrastructure & Security Foundation
- Iteration 2: Go Service - Authentication & User Management
- Iteration 3: Go Service - Document Upload & Storage
- Iteration 4: Python Service - Document Processing & Chunking
- Iteration 5: Python Service - Embeddings & Vector Storage
- Iteration 6: gRPC Integration Between Services
- Iteration 7: Python Service - RAG Query & LLM Integration
- Iteration 8: Testing, Documentation & Production Readiness

---

## ITERATION 1: Infrastructure & Security Foundation

**Goal:** Set up Docker infrastructure, databases, security configuration, and project structure. After this iteration, all services can start and connect to dependencies.

### Tasks

- [ ] 1.1 Create Docker Compose infrastructure
  - Create `docker-compose.yml` with services: PostgreSQL, MinIO, Qdrant, Ollama
  - Configure PostgreSQL with environment variables (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)
  - Configure MinIO with access key and secret key
  - Configure Qdrant with persistent volume
  - Configure Ollama with llama2 model
  - Expose only necessary ports (PostgreSQL: 5432, MinIO: 9000/9001, Qdrant: 6333)
  - Add health checks for all services
  - _Requirements: 14.1, 14.2, 14.3, 14.4_

- [ ] 1.2 Create comprehensive .gitignore file
  - Add .env files to .gitignore
  - Add vendor/, node_modules/, __pycache__/ directories
  - Add binary files (go service executable)
  - Add IDE-specific files (.vscode/, .idea/)
  - Add log files and temporary files
  - Add MinIO data directory
  - _Requirements: 13.1, 13.2_

- [ ] 1.3 Create .env.example template
  - Document all required environment variables with descriptions
  - Include PostgreSQL connection settings
  - Include MinIO credentials and endpoint
  - Include Qdrant endpoint
  - Include Ollama endpoint
  - Include JWT secret placeholder
  - Include service ports (Go: 8080, Python: 8081)
  - Add comments explaining each variable
  - _Requirements: 13.3, 13.4, 13.5_

- [ ] 1.4 Create Go service project structure

  - Create directory structure: cmd/server/, internal/handlers/, internal/models/, internal/middleware/, internal/grpc_client/
  - Initialize Go module with `go mod init docmind-go`
  - Create main.go entry point with basic HTTP server
  - Add Dockerfile for Go service
  - _Requirements: 14.5_

- [ ] 1.5 Create Python service project structure
  - Create directory structure: app/, app/services/, app/models/, app/grpc_server/
  - Create requirements.txt with FastAPI, gRPC, sentence-transformers, qdrant-client, PyPDF2, python-docx, langchain
  - Create main.py entry point with FastAPI app
  - Add Dockerfile for Python service
  - _Requirements: 14.5_

- [ ] 1.6 Create environment configuration loaders
  - Implement Go config loader that reads .env file
  - Validate required environment variables on startup
  - Implement Python config loader using python-dotenv
  - Validate required environment variables on startup
  - Fail with clear error messages if required variables are missing
  - _Requirements: 13.6, 21.2_

- [ ]* 1.7 Write property test for configuration validation
  - **Property 37: Required Configuration Validation**
  - **Validates: Requirements 21.2**

- [ ] 1.8 Initialize PostgreSQL schema
  - Create SQL migration file for users table
  - Create SQL migration file for documents table
  - Add indexes on user email and document user_id
  - Run migrations on container startup
  - _Requirements: 1.1, 3.4_

- [ ] 1.9 Create MinIO bucket initialization
  - Create initialization script to create "docmind-documents" bucket
  - Set bucket policy to private
  - Run on container startup
  - _Requirements: 3.3_

- [ ] 1.10 Update README.md with setup instructions
  - Document prerequisites (Docker, Docker Compose)
  - Document how to copy .env.example to .env
  - Document how to start services with docker-compose up
  - Document how to verify all services are running
  - Add architecture diagram reference
  - _Requirements: 14.1_

- [ ] 1.11 Checkpoint - Verify infrastructure
  - Ensure all Docker containers start successfully
  - Verify PostgreSQL accepts connections
  - Verify MinIO console is accessible
  - Verify Qdrant API responds
  - Verify Ollama is ready
  - Ensure .gitignore prevents committing secrets
  - Ask user if questions arise

---

## ITERATION 2: Go Service - Authentication & User Management

**Goal:** Implement user registration, login, JWT authentication, and user CRUD operations. After this iteration, users can register and authenticate.

### Tasks

- [ ] 2.1 Implement User model and database layer
  - Create User struct in Go with fields: id, email, password_hash, created_at, updated_at
  - Implement database functions: CreateUser, GetUserByEmail, GetUserByID, UpdateUser, DeleteUser
  - Use prepared statements to prevent SQL injection
  - Hash passwords using bcrypt before storing
  - _Requirements: 1.1, 1.3, 16.1_


- [ ]* 2.2 Write property tests for User model
  - **Property 1: Email Uniqueness Enforcement**
  - **Validates: Requirements 1.2**

- [ ]* 2.3 Write property test for password hashing
  - **Property 2: Password Hashing**
  - **Validates: Requirements 1.3**

- [ ] 2.4 Implement JWT token generation and validation
  - Create JWT utility functions: GenerateToken, ValidateToken, ExtractUserID
  - Use HS256 algorithm with secret from environment
  - Set token expiration to 24 hours
  - Include user_id in token claims
  - _Requirements: 2.1, 2.2, 2.4_

- [ ]* 2.5 Write property tests for JWT operations
  - **Property 4: Valid Credentials Generate JWT**
  - **Validates: Requirements 2.1, 2.2**

- [ ]* 2.6 Write property test for expired tokens
  - **Property 6: Expired Token Rejection**
  - **Validates: Requirements 2.5**

- [ ] 2.7 Implement authentication middleware
  - Create middleware that extracts JWT from Authorization header
  - Validate token and extract user_id
  - Attach user_id to request context
  - Return 401 for missing or invalid tokens
  - _Requirements: 2.5, 2.6_

- [ ] 2.8 Implement user registration endpoint
  - Create POST /api/auth/register endpoint
  - Validate email format and password length (min 8 chars)
  - Check email uniqueness before creating user
  - Return 400 for validation errors, 409 for duplicate email
  - Return user object without password_hash
  - _Requirements: 1.2, 1.4, 1.5, 16.1, 16.2_

- [ ]* 2.9 Write property test for password omission
  - **Property 3: Password Omission in Responses**
  - **Validates: Requirements 1.5**

- [ ]* 2.10 Write property tests for email validation
  - **Property 29: Email Format Validation**
  - **Validates: Requirements 16.1**

- [ ]* 2.11 Write property test for password validation
  - **Property 30: Password Length Validation**
  - **Validates: Requirements 16.2**

- [ ] 2.12 Implement login endpoint
  - Create POST /api/auth/login endpoint
  - Validate credentials against database
  - Return JWT token on success
  - Return 401 for invalid credentials
  - _Requirements: 2.1, 2.3_

- [ ]* 2.13 Write property test for invalid credentials
  - **Property 5: Invalid Credentials Rejection**
  - **Validates: Requirements 2.3**

- [ ] 2.14 Implement user CRUD endpoints
  - Create GET /api/users/me endpoint (authenticated)
  - Create PUT /api/users/me endpoint (authenticated)
  - Create DELETE /api/users/me endpoint (authenticated)
  - Ensure password_hash is never returned in responses
  - _Requirements: 1.4, 1.5_

- [ ] 2.15 Add structured logging
  - Implement logger using Go's log/slog package
  - Log all authentication attempts (success/failure)
  - Log user creation, updates, deletions
  - Include timestamp, level, user_id in logs
  - Write logs to stdout
  - _Requirements: 15.1, 15.2, 15.6_

- [ ] 2.16 Implement error handling
  - Create standardized error response format (JSON with error message)
  - Return appropriate HTTP status codes
  - Log all errors with stack traces
  - _Requirements: 11.7, 15.3_


- [ ] 2.17 Update README.md with API documentation
  - Document registration endpoint with example request/response
  - Document login endpoint with example request/response
  - Document user endpoints with authentication requirements
  - Add curl examples for testing
  - _Requirements: 11.1_

- [ ] 2.18 Checkpoint - Test authentication flow
  - Ensure user can register successfully
  - Ensure duplicate email registration fails
  - Ensure user can login and receive JWT
  - Ensure invalid credentials are rejected
  - Ensure protected endpoints require valid JWT
  - Ensure all tests pass
  - Ask user if questions arise

---

## ITERATION 3: Go Service - Document Upload & Storage

**Goal:** Implement document upload, MinIO storage, metadata persistence, and access control. After this iteration, authenticated users can upload and manage documents.

### Tasks

- [ ] 3.1 Implement Document model and database layer
  - Create Document struct with fields: id, user_id, filename, file_size, mime_type, storage_path, status, created_at, updated_at
  - Implement database functions: CreateDocument, GetDocumentByID, GetDocumentsByUserID, UpdateDocumentStatus, DeleteDocument
  - Add index on user_id for efficient queries
  - _Requirements: 3.4, 4.1_

- [ ]* 3.2 Write property test for document ID uniqueness
  - **Property 11: Document ID Uniqueness**
  - **Validates: Requirements 3.5**

- [ ] 3.3 Implement MinIO client wrapper
  - Create MinIO client initialization from environment config
  - Implement UploadFile function that stores file and returns storage_path
  - Implement DownloadFile function that retrieves file by storage_path
  - Implement DeleteFile function
  - Handle connection errors gracefully
  - _Requirements: 3.3_

- [ ]* 3.4 Write property test for upload storage round-trip
  - **Property 9: Upload Storage Round-Trip**
  - **Validates: Requirements 3.3**

- [ ] 3.5 Implement file validation
  - Validate file extension (PDF, TXT, DOCX, MD only)
  - Validate file size (max 10MB)
  - Validate file is not empty
  - Return 400 with descriptive error for validation failures
  - _Requirements: 3.1, 3.2, 16.3_

- [ ]* 3.6 Write property tests for file validation
  - **Property 7: Supported File Extensions**
  - **Validates: Requirements 3.1**

- [ ]* 3.7 Write property test for unsupported extensions
  - **Property 8: Unsupported File Extension Rejection**
  - **Validates: Requirements 3.2**

- [ ]* 3.8 Write property test for file size validation
  - **Property 31: File Size Validation**
  - **Validates: Requirements 16.3**

- [ ] 3.9 Implement document upload endpoint
  - Create POST /api/documents endpoint (authenticated)
  - Accept multipart/form-data file upload
  - Validate file extension and size
  - Generate unique document_id (UUID)
  - Upload file to MinIO with path: {user_id}/{document_id}/{filename}
  - Store document metadata in PostgreSQL with status "pending_indexing"
  - Return document metadata (id, filename, file_size, upload_date, storage_path)
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ]* 3.10 Write property test for document metadata persistence
  - **Property 10: Document Metadata Persistence**
  - **Validates: Requirements 3.4**


- [ ] 3.11 Implement document access control
  - Create middleware to verify document ownership
  - Check that document.user_id matches authenticated user_id
  - Return 403 for unauthorized access attempts
  - Return 404 for non-existent documents
  - _Requirements: 4.1, 4.2_

- [ ]* 3.12 Write property tests for document access control
  - **Property 12: Document Access Authorization**
  - **Validates: Requirements 4.1**

- [ ]* 3.13 Write property test for cross-user access denial
  - **Property 13: Cross-User Access Denial**
  - **Validates: Requirements 4.2**

- [ ] 3.14 Implement document retrieval endpoints
  - Create GET /api/documents endpoint (authenticated) - list user's documents with pagination
  - Create GET /api/documents/{id} endpoint (authenticated) - get document metadata
  - Create GET /api/documents/{id}/download endpoint (authenticated) - download file from MinIO
  - Enforce access control on all endpoints
  - _Requirements: 4.1, 4.3_

- [ ] 3.15 Implement pagination for document list
  - Accept query parameters: page (default 1), page_size (default 20)
  - Return pagination metadata: total_items, total_pages, current_page, page_size
  - Return empty list for non-existent pages
  - _Requirements: 23.1, 23.2, 23.4_

- [ ]* 3.16 Write property test for pagination metadata
  - **Property 35: Pagination Metadata Presence**
  - **Validates: Requirements 23.2**

- [ ] 3.17 Implement document deletion endpoint
  - Create DELETE /api/documents/{id} endpoint (authenticated)
  - Verify document ownership
  - Delete file from MinIO
  - Delete document metadata from PostgreSQL
  - Return 404 if document doesn't exist after deletion
  - _Requirements: 4.4_

- [ ]* 3.18 Write property test for document deletion completeness
  - **Property 14: Document Deletion Completeness**
  - **Validates: Requirements 4.4**

- [ ] 3.19 Implement UUID validation
  - Validate document_id format in all endpoints
  - Return 400 for invalid UUID format
  - _Requirements: 16.5_

- [ ]* 3.20 Write property test for UUID validation
  - **Property 33: UUID Validation**
  - **Validates: Requirements 16.5**

- [ ] 3.21 Add logging for document operations
  - Log all document uploads with user_id, filename, file_size
  - Log all document access attempts
  - Log all document deletions
  - Include document_id in all log entries
  - _Requirements: 15.1, 15.2_

- [ ] 3.22 Update README.md with document API documentation
  - Document upload endpoint with multipart/form-data example
  - Document list, get, download, delete endpoints
  - Document pagination parameters
  - Add curl examples
  - _Requirements: 11.1_

- [ ] 3.23 Checkpoint - Test document management
  - Ensure authenticated user can upload documents
  - Ensure unsupported file types are rejected
  - Ensure files over 10MB are rejected
  - Ensure user can list their documents
  - Ensure user can download their documents
  - Ensure user cannot access other users' documents
  - Ensure user can delete their documents
  - Ensure all tests pass
  - Ask user if questions arise

---


## ITERATION 4: Python Service - Document Processing & Chunking

**Goal:** Implement text extraction from multiple file formats and document chunking. After this iteration, the Python service can process documents into chunks.

### Tasks

- [ ] 4.1 Set up Python FastAPI application
  - Create FastAPI app instance in main.py
  - Add health check endpoint GET /health
  - Configure CORS if needed
  - Add structured logging using Python logging module
  - Load configuration from environment variables
  - _Requirements: 22.1_

- [ ] 4.2 Implement text extractor for PDF files
  - Create TextExtractor class with extract_pdf method
  - Use PyPDF2 to extract text from PDF files
  - Preserve paragraph structure
  - Handle multi-page PDFs
  - Return error if PDF contains only images without text
  - _Requirements: 5.1, 5.5, 17.1, 17.2, 17.3, 17.5_

- [ ]* 4.3 Write unit tests for PDF extraction
  - Test single-page PDF extraction
  - Test multi-page PDF extraction
  - Test PDF with only images (should error)
  - Test corrupt PDF (should error)
  - _Requirements: 17.4, 17.5_

- [ ] 4.4 Implement text extractor for DOCX files
  - Add extract_docx method to TextExtractor
  - Use python-docx to extract text
  - Preserve paragraph structure
  - Handle different encodings
  - Return error for corrupt files
  - _Requirements: 5.1, 5.5, 18.1, 18.2, 18.4_

- [ ]* 4.5 Write unit tests for DOCX extraction
  - Test DOCX with multiple paragraphs
  - Test DOCX with tables
  - Test corrupt DOCX (should error)
  - _Requirements: 18.3, 18.4_

- [ ] 4.6 Implement text extractor for TXT files
  - Add extract_txt method to TextExtractor
  - Handle UTF-8 and UTF-16 encodings automatically
  - Preserve original content exactly
  - _Requirements: 5.1, 5.2_

- [ ]* 4.7 Write property test for text extraction preservation
  - **Property 15: Text Extraction Preservation**
  - **Validates: Requirements 5.2**

- [ ] 4.7 Implement text extractor for Markdown files
  - Add extract_markdown method to TextExtractor
  - Preserve Markdown formatting (headers, lists, code blocks)
  - Handle UTF-8 and UTF-16 encodings
  - _Requirements: 5.1, 19.1, 19.2, 19.4_

- [ ]* 4.8 Write unit tests for Markdown extraction
  - Test Markdown with headers and lists
  - Test Markdown with code blocks
  - Test different encodings
  - _Requirements: 19.3, 19.4_

- [ ] 4.9 Implement document chunker
  - Create DocumentChunker class using LangChain's RecursiveCharacterTextSplitter
  - Set chunk_size to 1000 characters
  - Set chunk_overlap to 200 characters
  - Generate unique chunk_id for each chunk (UUID)
  - Add metadata: chunk_index (sequential 0, 1, 2...), original_position
  - Ensure all chunks are ≤ 1200 characters
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ]* 4.10 Write property test for chunk size limit
  - **Property 16: Chunk Size Limit**
  - **Validates: Requirements 6.3**

- [ ]* 4.11 Write property test for chunk metadata
  - **Property 17: Chunk Metadata Presence**
  - **Validates: Requirements 6.4**

- [ ]* 4.12 Write property test for chunk ID uniqueness and order
  - **Property 18: Chunk ID Uniqueness and Order**
  - **Validates: Requirements 6.5**


- [ ] 4.13 Create document processing pipeline
  - Create ProcessingPipeline class that orchestrates extraction and chunking
  - Accept file_path and file_type as input
  - Route to appropriate extractor based on file_type
  - Pass extracted text to chunker
  - Return list of chunks with metadata
  - Handle errors at each stage with descriptive messages
  - _Requirements: 5.1, 6.1_

- [ ] 4.14 Add input validation
  - Validate file_type is one of: pdf, txt, docx, md
  - Validate file exists and is readable
  - Return 400 with descriptive error for validation failures
  - _Requirements: 16.1_

- [ ]* 4.15 Write property test for validation error format
  - **Property 34: Validation Error Response Format**
  - **Validates: Requirements 16.6**

- [ ] 4.16 Add logging for processing operations
  - Log start of document processing with document_id
  - Log extraction completion with character count
  - Log chunking completion with chunk count
  - Log any errors with full context
  - Write logs to stdout
  - _Requirements: 15.1, 15.2, 15.6_

- [ ] 4.17 Implement resource cleanup
  - Ensure file handles are closed after processing
  - Clean up temporary files if created
  - Handle exceptions to prevent resource leaks
  - _Requirements: 20.1, 20.2_

- [ ] 4.18 Update README.md with Python service documentation
  - Document text extraction capabilities
  - Document supported file formats
  - Document chunking strategy and parameters
  - Add examples of processing output
  - _Requirements: 11.1_

- [ ] 4.19 Checkpoint - Test document processing
  - Ensure PDF files are extracted correctly
  - Ensure DOCX files are extracted correctly
  - Ensure TXT files are extracted correctly
  - Ensure Markdown files are extracted correctly
  - Ensure documents are chunked with correct size and overlap
  - Ensure chunks have proper metadata
  - Ensure all tests pass
  - Ask user if questions arise

---

## ITERATION 5: Python Service - Embeddings & Vector Storage

**Goal:** Implement embedding generation and Qdrant vector storage. After this iteration, document chunks can be converted to embeddings and stored in Qdrant.

### Tasks

- [ ] 5.1 Implement embedding generator
  - Create EmbeddingGenerator class using sentence-transformers
  - Load model: all-MiniLM-L6-v2
  - Cache model in memory after first load
  - Implement generate_embedding method that accepts text and returns 384-dimensional vector
  - Normalize vectors to unit length (L2 norm = 1.0)
  - Handle empty chunks by returning zero vector
  - _Requirements: 7.1, 7.2, 7.3, 7.5, 24.4_

- [ ]* 5.2 Write property test for embedding dimensionality
  - **Property 19: Embedding Dimensionality**
  - **Validates: Requirements 7.2**

- [ ]* 5.3 Write property test for embedding normalization
  - **Property 20: Embedding Normalization**
  - **Validates: Requirements 7.5**

- [ ] 5.4 Implement batch embedding generation
  - Add generate_embeddings_batch method for processing multiple chunks
  - Process in batches of 32 chunks to optimize performance
  - Return list of embeddings in same order as input
  - _Requirements: 24.1, 24.2_


- [ ] 5.5 Implement Qdrant client wrapper
  - Create QdrantClient wrapper class
  - Initialize connection from environment config (QDRANT_HOST, QDRANT_PORT)
  - Implement create_collection method with vector size 384 and cosine distance
  - Implement collection_exists method
  - Handle connection errors and automatic reconnections
  - _Requirements: 8.1, 8.2, 8.5_

- [ ]* 5.6 Write property test for collection existence
  - **Property 21: Collection Existence After Indexing**
  - **Validates: Requirements 8.2**

- [ ] 5.7 Implement vector storage operations
  - Implement upsert_embeddings method that stores vectors with payload
  - Payload must include: document_id, user_id, chunk_text, chunk_index
  - Use chunk_id as point ID in Qdrant
  - Implement batch upsert for efficiency
  - _Requirements: 8.3, 8.4_

- [ ]* 5.8 Write property test for embedding payload completeness
  - **Property 22: Embedding Payload Completeness**
  - **Validates: Requirements 8.3**

- [ ] 5.9 Implement vector deletion operations
  - Implement delete_by_document_id method
  - Use Qdrant filter to delete all points with matching document_id
  - Return count of deleted vectors
  - _Requirements: 4.5_

- [ ] 5.10 Implement document indexing endpoint
  - Create REST endpoint POST /api/process (for testing, will be replaced by gRPC)
  - Accept document_id, user_id, file_path, file_type
  - Run processing pipeline (extract → chunk → embed → store)
  - Store all embeddings in Qdrant with metadata
  - Return success with chunk count and embedding count
  - Handle errors at each stage
  - _Requirements: 3.6, 7.4, 8.4_

- [ ] 5.11 Add error handling for embedding failures
  - Catch model loading errors
  - Catch embedding generation errors
  - Catch Qdrant connection errors
  - Return descriptive error messages
  - Log all errors with context
  - _Requirements: 15.3, 15.4_

- [ ] 5.12 Implement health check with dependency status
  - Update /health endpoint to check Qdrant connection
  - Return 200 if all dependencies are available
  - Return 503 if Qdrant is unavailable
  - Include dependency status in response
  - _Requirements: 22.2, 22.3, 22.5_

- [ ] 5.13 Add logging for embedding operations
  - Log embedding generation start with chunk count
  - Log embedding generation completion with timing
  - Log vector storage operations
  - Log any errors with full context
  - _Requirements: 15.1, 15.2_

- [ ] 5.14 Implement resource cleanup for embeddings
  - Ensure model resources are released on shutdown
  - Close Qdrant connections properly
  - Handle reconnections on service restart
  - _Requirements: 20.3, 20.4, 20.5_

- [ ] 5.15 Update README.md with embedding documentation
  - Document embedding model and dimensions
  - Document vector storage strategy
  - Document Qdrant collection structure
  - Add examples of embedding operations
  - _Requirements: 11.1_

- [ ] 5.16 Checkpoint - Test embedding and storage
  - Ensure embeddings are generated with correct dimensions
  - Ensure embeddings are normalized
  - Ensure embeddings are stored in Qdrant with correct payload
  - Ensure batch processing works efficiently
  - Ensure Qdrant collection is created automatically
  - Ensure all tests pass
  - Ask user if questions arise

---


## ITERATION 6: gRPC Integration Between Services

**Goal:** Implement gRPC communication between Go and Python services. After this iteration, Go service can trigger document indexing via gRPC.

### Tasks

- [ ] 6.1 Define gRPC protocol buffer schema
  - Create proto/rag_service.proto file
  - Define IndexDocumentRequest message (document_id, user_id, file_path, file_type)
  - Define IndexDocumentResponse message (success, chunk_count, error_message)
  - Define DeleteDocumentRequest message (document_id)
  - Define DeleteDocumentResponse message (success, deleted_count)
  - Define RAGService service with IndexDocument and DeleteDocument methods
  - _Requirements: 12.1, 12.2_

- [ ] 6.2 Generate gRPC code for Go
  - Install protoc compiler and Go gRPC plugin
  - Generate Go code from proto file
  - Add generated code to internal/grpc_client/
  - Update go.mod with gRPC dependencies
  - _Requirements: 12.1_

- [ ] 6.3 Generate gRPC code for Python
  - Install grpcio-tools
  - Generate Python code from proto file
  - Add generated code to app/grpc_server/
  - Update requirements.txt with gRPC dependencies
  - _Requirements: 12.1_

- [ ] 6.4 Implement gRPC server in Python service
  - Create RAGServiceServicer class implementing RAGService
  - Implement IndexDocument method that calls processing pipeline
  - Implement DeleteDocument method that deletes embeddings from Qdrant
  - Start gRPC server on port 50051
  - Add graceful shutdown handling
  - _Requirements: 12.2, 12.3_

- [ ] 6.5 Implement gRPC client in Go service
  - Create gRPC client wrapper in internal/grpc_client/
  - Initialize connection to Python service from environment config
  - Implement IndexDocument method with 30-second timeout
  - Implement DeleteDocument method with 10-second timeout
  - Handle connection errors and retries (max 3 attempts)
  - _Requirements: 12.4, 12.5, 12.6, 12.7_

- [ ] 6.6 Integrate gRPC into document upload flow
  - After successful file upload to MinIO, invoke gRPC IndexDocument
  - Update document status to "indexed" on success
  - Update document status to "pending_indexing" on gRPC failure
  - Log gRPC invocation results
  - _Requirements: 3.6, 3.7_

- [ ] 6.7 Integrate gRPC into document deletion flow
  - Before deleting document from database, invoke gRPC DeleteDocument
  - Delete embeddings from Qdrant via gRPC
  - Continue with file and metadata deletion
  - Log gRPC invocation results
  - _Requirements: 4.5_

- [ ] 6.8 Add gRPC error handling
  - Handle gRPC connection timeouts
  - Handle gRPC service unavailable errors
  - Return appropriate HTTP status codes to client
  - Log all gRPC errors with context
  - _Requirements: 12.5, 15.3_

- [ ] 6.9 Add gRPC logging
  - Log all gRPC requests with method name and parameters
  - Log all gRPC responses with status and timing
  - Log gRPC connection events
  - _Requirements: 15.1, 15.2_

- [ ] 6.10 Update docker-compose.yml for gRPC
  - Add Python service to docker-compose.yml
  - Expose gRPC port 50051 internally (not to host)
  - Add Go service to docker-compose.yml
  - Configure service discovery via Docker network
  - Set environment variables for gRPC endpoints
  - _Requirements: 14.6_

- [ ] 6.11 Update .env.example with gRPC configuration
  - Add PYTHON_GRPC_HOST and PYTHON_GRPC_PORT
  - Add GO_SERVICE_PORT
  - Document gRPC timeout settings
  - _Requirements: 13.3_


- [ ] 6.12 Update README.md with gRPC documentation
  - Document gRPC service interface
  - Document message formats
  - Document timeout settings
  - Add architecture diagram showing gRPC communication
  - _Requirements: 11.1_

- [ ] 6.13 Checkpoint - Test gRPC integration
  - Ensure Go service can connect to Python gRPC server
  - Ensure document upload triggers gRPC IndexDocument call
  - Ensure document indexing completes successfully
  - Ensure document status is updated to "indexed"
  - Ensure document deletion triggers gRPC DeleteDocument call
  - Ensure embeddings are deleted from Qdrant
  - Ensure gRPC timeouts work correctly
  - Ensure all tests pass
  - Ask user if questions arise

---

## ITERATION 7: Python Service - RAG Query & LLM Integration

**Goal:** Implement semantic search and RAG query with Ollama LLM. After this iteration, users can query their documents and receive contextualized answers.

### Tasks

- [ ] 7.1 Implement semantic search
  - Create SemanticSearch class
  - Implement search method that accepts query text and user_id
  - Generate embedding for query using same model (all-MiniLM-L6-v2)
  - Search Qdrant with cosine similarity
  - Filter results by user_id
  - Return top 5 chunks with similarity ≥ 0.7
  - Return empty list if no chunks meet threshold
  - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

- [ ]* 7.2 Write property test for query embedding dimensionality
  - **Property 23: Query Embedding Dimensionality**
  - **Validates: Requirements 9.1**

- [ ]* 7.3 Write property test for search result limit
  - **Property 24: Search Result Limit**
  - **Validates: Requirements 9.2**

- [ ]* 7.4 Write property test for search result privacy
  - **Property 25: Search Result Privacy**
  - **Validates: Requirements 9.3**

- [ ]* 7.5 Write property test for similarity threshold
  - **Property 26: Search Result Similarity Threshold**
  - **Validates: Requirements 9.4**

- [ ] 7.6 Implement Ollama LLM client
  - Create OllamaClient class
  - Initialize connection from environment config (OLLAMA_HOST, OLLAMA_PORT)
  - Implement generate_response method that accepts prompt and context
  - Use llama2 model
  - Set max_tokens to 500
  - Set temperature to 0.7
  - Handle connection errors and timeouts
  - _Requirements: 10.1, 10.2_

- [ ] 7.7 Implement RAG prompt construction
  - Create prompt template that includes context chunks and user query
  - Format: "Based on the following context, answer the question. Context: {chunks}. Question: {query}"
  - Limit context to top 5 chunks
  - Include chunk sources in prompt
  - _Requirements: 10.3_

- [ ] 7.8 Implement RAG query pipeline
  - Create RAGPipeline class that orchestrates search and generation
  - Accept query text and user_id
  - Perform semantic search to retrieve relevant chunks
  - Construct prompt with retrieved context
  - Generate response using Ollama LLM
  - Return response with source document_ids
  - Handle empty search results gracefully
  - _Requirements: 10.1, 10.3, 10.4, 10.5, 10.6_

- [ ]* 7.9 Write property test for non-empty response
  - **Property 27: Non-Empty Response Generation**
  - **Validates: Requirements 10.4**

- [ ]* 7.10 Write property test for response source references
  - **Property 28: Response Source References**
  - **Validates: Requirements 10.6**


- [ ] 7.11 Add QueryDocument to gRPC interface
  - Update proto/rag_service.proto with QueryDocumentRequest (query, user_id)
  - Update proto/rag_service.proto with QueryDocumentResponse (answer, sources, chunk_count)
  - Regenerate gRPC code for Go and Python
  - Implement QueryDocument method in Python gRPC server
  - _Requirements: 12.2_

- [ ] 7.12 Implement query endpoint in Go service
  - Create POST /api/query endpoint (authenticated)
  - Accept query text in request body
  - Validate query is not empty or whitespace-only
  - Extract user_id from JWT token
  - Invoke Python service via gRPC QueryDocument with 45-second timeout
  - Return response with answer and sources
  - _Requirements: 10.5, 16.4_

- [ ]* 7.13 Write property test for empty query rejection
  - **Property 32: Empty Query Rejection**
  - **Validates: Requirements 16.4**

- [ ] 7.14 Add logging for query operations
  - Log all queries with user_id and query text
  - Log semantic search results with chunk count and similarity scores
  - Log LLM generation with timing
  - Log final response with source count
  - _Requirements: 15.1, 15.2_

- [ ] 7.15 Implement query error handling
  - Handle empty search results (no relevant chunks found)
  - Handle Ollama connection errors
  - Handle LLM generation timeouts
  - Return descriptive error messages
  - Log all errors with context
  - _Requirements: 15.3, 15.4_

- [ ] 7.16 Update health check with Ollama status
  - Add Ollama connection check to /health endpoint
  - Return 503 if Ollama is unavailable
  - Include Ollama status in response
  - _Requirements: 22.4, 22.5_

- [ ] 7.17 Implement concurrency handling
  - Use async/await for I/O operations in Python service
  - Configure FastAPI with multiple workers
  - Configure Go service with goroutines for concurrent requests
  - Test with multiple simultaneous queries
  - _Requirements: 25.1, 25.2, 25.3_

- [ ] 7.18 Update README.md with query documentation
  - Document query endpoint with examples
  - Document RAG pipeline flow
  - Document semantic search parameters
  - Document LLM configuration
  - Add example queries and responses
  - _Requirements: 11.1_

- [ ] 7.19 Checkpoint - Test RAG query flow
  - Ensure user can query documents after indexing
  - Ensure semantic search returns relevant chunks
  - Ensure LLM generates contextualized responses
  - Ensure responses include source references
  - Ensure empty queries are rejected
  - Ensure queries only return user's own documents
  - Ensure concurrent queries work correctly
  - Ensure all tests pass
  - Ask user if questions arise

---

## ITERATION 8: Testing, Documentation & Production Readiness

**Goal:** Comprehensive testing, documentation, security hardening, and production readiness. After this iteration, the system is ready for deployment.

### Tasks

- [ ] 8.1 Create comprehensive integration tests
  - Test complete user registration → login → upload → query flow
  - Test multi-user isolation (users cannot access each other's documents)
  - Test document lifecycle (upload → index → query → delete)
  - Test error scenarios (invalid tokens, missing files, etc.)
  - _Requirements: 4.2, 10.5_

- [ ] 8.2 Create performance tests
  - Test document indexing performance with various file sizes
  - Test query response time with different document counts
  - Test concurrent user operations
  - Verify embedding generation meets performance requirements
  - _Requirements: 24.1, 24.2, 24.3, 25.3_


- [ ] 8.3 Security audit and hardening
  - Verify .gitignore prevents committing secrets
  - Verify no hardcoded credentials in code
  - Verify all passwords are hashed with bcrypt
  - Verify JWT tokens expire correctly
  - Verify SQL injection prevention (prepared statements)
  - Verify file upload validation prevents malicious files
  - Verify CORS configuration is appropriate
  - _Requirements: 13.1, 13.2, 1.3, 2.4, 16.1_

- [ ] 8.4 Create API documentation
  - Generate OpenAPI/Swagger documentation for Go REST API
  - Document all endpoints with request/response examples
  - Document authentication requirements
  - Document error codes and messages
  - Add Postman collection for testing
  - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6_

- [ ] 8.5 Create deployment documentation
  - Document system requirements (Docker, Docker Compose)
  - Document environment variable configuration
  - Document how to start/stop services
  - Document how to view logs
  - Document how to backup PostgreSQL and Qdrant data
  - Document troubleshooting common issues
  - _Requirements: 14.1, 15.5_

- [ ] 8.6 Create developer documentation
  - Document project structure for both services
  - Document how to add new file format support
  - Document how to modify chunking strategy
  - Document how to change embedding model
  - Document how to add new endpoints
  - Document testing strategy
  - _Requirements: 11.1_

- [ ] 8.7 Implement structured error responses
  - Ensure all errors return consistent JSON format
  - Include error code, message, and timestamp
  - Include request_id for tracing
  - Document all error codes
  - _Requirements: 11.7, 15.4_

- [ ] 8.8 Implement request ID tracing
  - Generate unique request_id for each API request
  - Include request_id in all log entries
  - Return request_id in error responses
  - Enable end-to-end request tracing
  - _Requirements: 15.4_

- [ ] 8.9 Add monitoring and observability
  - Add metrics endpoint for Prometheus (optional)
  - Log key metrics (request count, latency, error rate)
  - Add health check endpoints for all services
  - Document monitoring strategy
  - _Requirements: 22.1, 22.2, 22.3_

- [ ] 8.10 Optimize Docker images
  - Use multi-stage builds for smaller images
  - Use Alpine base images where possible
  - Remove unnecessary dependencies
  - Document image sizes
  - _Requirements: 14.5_

- [ ] 8.11 Create database migration strategy
  - Document how to run migrations
  - Create rollback scripts for migrations
  - Test migration on fresh database
  - _Requirements: 1.1_

- [ ] 8.12 Test resource cleanup and reconnection
  - Test service restart scenarios
  - Test database connection loss and reconnection
  - Test MinIO connection loss and reconnection
  - Test Qdrant connection loss and reconnection
  - Test Ollama connection loss and reconnection
  - Verify no resource leaks
  - _Requirements: 20.1, 20.2, 20.3, 20.4, 20.5_

- [ ] 8.13 Create example .env file with all variables
  - Ensure .env.example is complete and up-to-date
  - Add comments explaining each variable
  - Add example values (not real credentials)
  - Document required vs optional variables
  - _Requirements: 13.3, 13.4, 13.5_

- [ ] 8.14 Final security review
  - Review all code for security vulnerabilities
  - Verify input validation on all endpoints
  - Verify authentication on all protected endpoints
  - Verify authorization checks for document access
  - Verify no sensitive data in logs
  - _Requirements: 13.1, 16.1, 16.2, 16.3, 16.4, 16.5, 16.6_


- [ ] 8.15 Create comprehensive README.md
  - Add project overview and architecture diagram
  - Add quick start guide
  - Add API documentation links
  - Add troubleshooting section
  - Add contributing guidelines
  - Add license information
  - _Requirements: 11.1_

- [ ] 8.16 Create CHANGELOG.md
  - Document all iterations and features
  - Document breaking changes
  - Document bug fixes
  - Follow semantic versioning

- [ ] 8.17 Final checkpoint - Production readiness
  - Ensure all services start successfully with docker-compose up
  - Ensure complete user flow works end-to-end
  - Ensure all tests pass
  - Ensure documentation is complete
  - Ensure .gitignore prevents committing secrets
  - Ensure .env.example is up-to-date
  - Ensure no hardcoded credentials exist
  - Ensure all error scenarios are handled gracefully
  - Ensure logging is comprehensive
  - Ensure health checks work correctly
  - Ask user if questions arise

---

## Notes

- Tasks marked with `*` are optional property-based tests and can be skipped for faster MVP delivery
- Each task references specific requirements for traceability
- Checkpoints at the end of each iteration ensure incremental validation
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- Each iteration leaves the project in a git-safe, deployable state
- Security tasks are included in every iteration to maintain security throughout development
- The implementation uses Go for the API service and Python for the RAG pipeline as specified in the design

## Getting Started

To begin implementation:
1. Open this tasks.md file in your IDE
2. Click "Start task" next to any task item to begin
3. Complete tasks sequentially within each iteration
4. Commit and push after each iteration checkpoint
5. Verify .gitignore and .env.example are correct before committing

## Requirements Coverage

This task list covers all 25 requirements from the requirements document:
- Requirements 1-2: User Management & Authentication (Iteration 2)
- Requirements 3-4: Document Upload & Access Control (Iteration 3)
- Requirements 5-6: Text Extraction & Chunking (Iteration 4)
- Requirements 7-8: Embeddings & Vector Storage (Iteration 5)
- Requirement 9: Semantic Search (Iteration 7)
- Requirement 10: RAG Response Generation (Iteration 7)
- Requirement 11: REST API (All iterations)
- Requirement 12: gRPC Interface (Iteration 6)
- Requirement 13: Configuration & Secrets (Iteration 1, 8)
- Requirement 14: Dockerized Infrastructure (Iteration 1, 6)
- Requirement 15: Error Handling & Logging (All iterations)
- Requirement 16: Input Validation (All iterations)
- Requirements 17-19: File Format Processing (Iteration 4)
- Requirement 20: Resource Cleanup (Iterations 4, 5, 8)
- Requirement 21: Configuration Parsing (Iteration 1)
- Requirement 22: Health Monitoring (Iterations 4, 5, 7)
- Requirement 23: Pagination (Iteration 3)
- Requirement 24: Embedding Performance (Iteration 5)
- Requirement 25: Concurrency Handling (Iteration 7)
