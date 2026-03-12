# DocMind Python RAG Service

Document processing, chunking, embeddings, and vector storage service for the DocMind RAG system.

## Overview

The Python RAG Service handles:
- Text extraction from multiple file formats (PDF, DOCX, TXT, Markdown)
- Document chunking using LangChain
- Embedding generation using sentence-transformers
- Vector storage in Qdrant
- Semantic search and RAG query capabilities
- Processing pipeline orchestration
- Comprehensive logging and error handling

## Architecture

```
main.py                 # FastAPI application entry point
├── routes/
│   ├── health.py      # Health check with dependency status
│   ├── processing.py  # Document processing endpoints
│   └── indexing.py    # Document indexing endpoints
├── extractors/
│   ├── base.py        # Base extractor class
│   ├── pdf.py         # PDF text extraction
│   ├── docx.py        # DOCX text extraction
│   ├── text.py        # Plain text extraction
│   ├── markdown.py    # Markdown text extraction
│   └── factory.py     # Extractor factory pattern
├── chunking/
│   └── chunker.py     # LangChain-based document chunking
├── embeddings/
│   └── generator.py   # Sentence-transformers embedding generation
├── vector_store/
│   └── qdrant_client.py # Qdrant vector database client
├── indexing/
│   └── indexer.py     # Document indexing pipeline
├── processing/
│   └── pipeline.py    # Document processing pipeline
└── config.py          # Configuration management
```

## Setup

### Prerequisites
- Python 3.9+
- pip or poetry
- Qdrant running (Docker recommended)

### Installation

1. Create virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

3. Configure environment:
```bash
cp .env.example .env
# Edit .env with your configuration
```

## Running the Service

### Development
```bash
python main.py
```

### Production
```bash
uvicorn main:app --host 0.0.0.0 --port 8001 --workers 4
```

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "service": "python-rag-service",
  "dependencies": [
    {
      "name": "qdrant",
      "status": "healthy",
      "available": true
    }
  ]
}
```

### Process Document
```
POST /api/process
```

Parameters:
- `document_id` (query): Unique document identifier
- `user_id` (query): User identifier
- `file` (form-data): Document file (PDF, DOCX, TXT, or Markdown)

Response:
```json
{
  "document_id": "uuid",
  "total_chunks": 42,
  "status": "success"
}
```

### Index Document
```
POST /api/index
```

Parameters:
- `document_id` (query): Unique document identifier
- `user_id` (query): User identifier
- `file_path` (query): Path to document file
- `file_type` (query): File extension (pdf, docx, txt, md)

Response:
```json
{
  "document_id": "uuid",
  "chunk_count": 42,
  "embedding_count": 42,
  "status": "success"
}
```

## Supported File Formats

| Format | Extension | Extractor |
|--------|-----------|-----------|
| PDF | .pdf | PDFExtractor |
| DOCX | .docx | DOCXExtractor |
| Plain Text | .txt | TextExtractorPlain |
| Markdown | .md | MarkdownExtractor |

## Configuration

Environment variables (see `.env.example`):

### Server
- `SERVER_HOST`: Server bind address (default: 0.0.0.0)
- `SERVER_PORT`: Server port (default: 8001)
- `DEBUG`: Debug mode (default: false)

### Services
- `GO_SERVICE_URL`: Go service URL for integration

### Document Processing
- `MAX_FILE_SIZE`: Maximum file size in bytes (default: 10MB)
- `ALLOWED_EXTENSIONS`: Comma-separated list of allowed extensions

### Chunking
- `CHUNK_SIZE`: Document chunk size (default: 1000 characters)
- `CHUNK_OVERLAP`: Chunk overlap for context (default: 200 characters)

### Qdrant Vector Store
- `QDRANT_HOST`: Qdrant server host (default: localhost)
- `QDRANT_PORT`: Qdrant server port (default: 6333)

### Embeddings
- `EMBEDDING_MODEL`: Sentence-transformers model (default: all-MiniLM-L6-v2)
- `EMBEDDING_DIMENSION`: Embedding vector dimension (default: 384)
- `EMBEDDING_BATCH_SIZE`: Batch size for embedding generation (default: 32)

### Logging
- `LOG_LEVEL`: Logging level (default: INFO)

## Chunking Strategy

The service uses LangChain's `RecursiveCharacterTextSplitter` with:
- Configurable chunk size (default: 1000 characters)
- Configurable overlap (default: 200 characters)
- Hierarchical splitting: paragraphs → sentences → words → characters

## Embedding Generation

The service uses `sentence-transformers` with:
- Model: `all-MiniLM-L6-v2` (384-dimensional embeddings)
- Normalization: L2 norm = 1.0 (unit vectors)
- Batch processing: 32 texts per batch for efficiency
- Lazy loading: Model loaded on first use

## Vector Storage

The service uses Qdrant with:
- Collection: `docmind_chunks`
- Distance metric: Cosine similarity
- Vector dimension: 384
- Payload: document_id, user_id, chunk_text, chunk_index

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200`: Success
- `400`: Bad request (invalid file, unsupported format)
- `500`: Server error

Error responses include descriptive messages for debugging.

## Logging

Comprehensive logging at multiple levels:
- INFO: Operation summaries
- DEBUG: Detailed processing steps
- ERROR: Failures and exceptions

Logs include context information for tracing.

## Development

### Adding New File Format Support

1. Create extractor in `extractors/`:
```python
from extractors.base import TextExtractor

class NewFormatExtractor(TextExtractor):
    def extract(self, file_path: str) -> str:
        # Implementation
        pass
```

2. Register in `extractors/factory.py`:
```python
_extractors = {
    'newformat': NewFormatExtractor,
    # ...
}
```

### Changing Embedding Model

Update `config.py`:
```python
EMBEDDING_MODEL: str = "your-model-name"
EMBEDDING_DIMENSION: int = 768  # Update dimension
```

### Modifying Chunking Strategy

Update `config.py`:
```python
CHUNK_SIZE: int = 2000
CHUNK_OVERLAP: int = 400
```

## Testing

Run tests:
```bash
pytest tests/
```

## Integration with Go Service

The Python service integrates with the Go User Service:
- Receives document metadata from Go service
- Processes documents asynchronously
- Returns chunk and embedding information
- Stores embeddings in Qdrant for semantic search

## Performance Considerations

- Temporary files are cleaned up after processing
- Streaming for large file uploads
- Configurable chunk sizes for memory optimization
- Batch embedding generation for efficiency
- Lazy model loading to reduce startup time
- Logging overhead is minimal in production

## Security

- File size validation (10MB default limit)
- File type validation
- Temporary file cleanup
- No sensitive data in logs
- CORS enabled for cross-service communication
- User isolation in vector storage (user_id filtering)

## License

Part of DocMind RAG System
