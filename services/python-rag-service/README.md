# DocMind Python RAG Service

Document processing and chunking service for the DocMind RAG system.

## Overview

The Python RAG Service handles:
- Text extraction from multiple file formats (PDF, DOCX, TXT, Markdown)
- Document chunking using LangChain
- Processing pipeline orchestration
- Comprehensive logging and error handling

## Architecture

```
main.py                 # FastAPI application entry point
├── routes/
│   ├── health.py      # Health check endpoint
│   └── processing.py  # Document processing endpoints
├── extractors/
│   ├── base.py        # Base extractor class
│   ├── pdf.py         # PDF text extraction
│   ├── docx.py        # DOCX text extraction
│   ├── text.py        # Plain text extraction
│   ├── markdown.py    # Markdown text extraction
│   └── factory.py     # Extractor factory pattern
├── chunking/
│   └── chunker.py     # LangChain-based document chunking
├── processing/
│   └── pipeline.py    # Document processing pipeline
└── config.py          # Configuration management
```

## Setup

### Prerequisites
- Python 3.9+
- pip or poetry

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
  "service": "python-rag-service"
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

## Supported File Formats

| Format | Extension | Extractor |
|--------|-----------|-----------|
| PDF | .pdf | PDFExtractor |
| DOCX | .docx | DOCXExtractor |
| Plain Text | .txt | TextExtractorPlain |
| Markdown | .md | MarkdownExtractor |

## Configuration

Environment variables (see `.env.example`):

- `SERVER_HOST`: Server bind address (default: 0.0.0.0)
- `SERVER_PORT`: Server port (default: 8001)
- `DEBUG`: Debug mode (default: false)
- `GO_SERVICE_URL`: Go service URL for integration
- `MAX_FILE_SIZE`: Maximum file size in bytes (default: 10MB)
- `CHUNK_SIZE`: Document chunk size (default: 1000)
- `CHUNK_OVERLAP`: Chunk overlap for context (default: 200)
- `LOG_LEVEL`: Logging level (default: INFO)

## Chunking Strategy

The service uses LangChain's `RecursiveCharacterTextSplitter` with:
- Configurable chunk size (default: 1000 characters)
- Configurable overlap (default: 200 characters)
- Hierarchical splitting: paragraphs → sentences → words → characters

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

### Testing

Run tests:
```bash
pytest tests/
```

## Integration with Go Service

The Python service integrates with the Go User Service:
- Receives document metadata from Go service
- Processes documents asynchronously
- Returns chunk information for indexing

## Performance Considerations

- Temporary files are cleaned up after processing
- Streaming for large file uploads
- Configurable chunk sizes for memory optimization
- Logging overhead is minimal in production

## Security

- File size validation (10MB default limit)
- File type validation
- Temporary file cleanup
- No sensitive data in logs
- CORS enabled for cross-service communication

## License

Part of DocMind RAG System
