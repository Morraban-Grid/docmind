# DocMind Python RAG Service

Python microservice responsible for document processing, embedding generation, semantic search, and RAG query processing.

## Features

- Document text extraction (PDF, DOCX, TXT, MD)
- Document chunking with LangChain
- Embedding generation with sentence-transformers
- Vector storage with Qdrant
- Semantic search
- RAG query with Ollama LLM
- gRPC server for Go service communication

## Development

### Prerequisites

- Python 3.12+
- Qdrant
- Ollama

### Setup

1. Create virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

3. Run the service:
```bash
uvicorn app.main:app --reload
```

## Configuration

All configuration is done via environment variables. See `.env.example` in the project root.

## API Endpoints

- GET /health - Health check
- POST /api/process - Process document (Coming in Iteration 4)
- POST /api/query - Query documents (Coming in Iteration 7)

## gRPC Services

- IndexDocument - Index a document for search
- QueryDocument - Query indexed documents
- DeleteDocument - Delete document embeddings

## Testing

```bash
pytest
```

## Models

- Embedding Model: sentence-transformers/all-MiniLM-L6-v2 (384 dimensions)
- LLM: Ollama llama2
