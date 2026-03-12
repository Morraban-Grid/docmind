# DocMind RAG System - Developer Guide

## Project Structure

```
docmind/
├── services/
│   ├── go-user-service/          # Go REST API service
│   │   ├── cmd/server/           # Entry point
│   │   ├── internal/
│   │   │   ├── handler/http/     # HTTP handlers
│   │   │   ├── service/          # Business logic
│   │   │   ├── repository/       # Data access
│   │   │   ├── domain/           # Domain models
│   │   │   ├── middleware/       # HTTP middleware
│   │   │   ├── grpc_client/      # gRPC client
│   │   │   └── infrastructure/   # Infrastructure
│   │   ├── migrations/           # Database migrations
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── go.sum
│   └── python-rag-service/       # Python RAG service
│       ├── main.py               # Entry point
│       ├── config.py             # Configuration
│       ├── routes/               # API routes
│       ├── search/               # Semantic search
│       ├── llm/                  # LLM client
│       ├── rag/                  # RAG pipeline
│       ├── embeddings/           # Embeddings
│       ├── extractors/           # Text extractors
│       ├── chunking/             # Document chunking
│       ├── vector_store/         # Vector storage
│       ├── indexing/             # Indexing pipeline
│       ├── grpc_server/          # gRPC server
│       ├── Dockerfile
│       ├── requirements.txt
│       └── .env.example
├── proto/                        # Protocol buffers
│   └── rag_service.proto
├── tests/                        # Test suite
│   ├── integration/              # Integration tests
│   ├── performance/              # Performance tests
│   ├── resilience/               # Resilience tests
│   ├── conftest.py
│   └── fixtures.py
├── docker-compose.yml            # Docker Compose
├── .env.example                  # Environment template
├── .gitignore                    # Git ignore rules
├── README.md                     # Project README
├── API-DOCUMENTATION.md          # API docs
├── DEPLOYMENT-GUIDE.md           # Deployment guide
├── DEVELOPER-GUIDE.md            # This file
└── CHANGELOG.md                  # Version history
```

## Development Setup

### Prerequisites

- Go 1.21+
- Python 3.10+
- Docker and Docker Compose
- Git
- protoc compiler

### Local Development

```bash
# Clone repository
git clone https://github.com/Morraban-Grid/docmind.git
cd docmind

# Copy environment file
cp .env.example .env

# Start infrastructure
docker-compose up -d postgres minio qdrant ollama

# Go service
cd services/go-user-service
go mod download
go run cmd/server/main.go

# Python service (in another terminal)
cd services/python-rag-service
pip install -r requirements.txt
python main.py
```

## Code Organization

### Go Service

**Handler Layer** (`internal/handler/http/`)
- HTTP request/response handling
- Input validation
- Error responses
- Logging

**Service Layer** (`internal/service/`)
- Business logic
- Orchestration
- Error handling

**Repository Layer** (`internal/repository/`)
- Database access
- Query building
- Transaction management

**Domain Layer** (`internal/domain/`)
- Domain models
- Business rules
- Validation

**Middleware** (`internal/middleware/`)
- Authentication
- Logging
- Error recovery
- CORS

### Python Service

**Routes** (`routes/`)
- FastAPI endpoints
- Request/response handling
- Input validation

**Search** (`search/`)
- Semantic search implementation
- Vector similarity

**LLM** (`llm/`)
- Ollama client
- Response generation

**RAG** (`rag/`)
- RAG pipeline orchestration
- Prompt construction

**Embeddings** (`embeddings/`)
- Embedding generation
- Model management

**Extractors** (`extractors/`)
- Text extraction
- Format-specific handling

**Chunking** (`chunking/`)
- Document chunking
- Metadata management

**Vector Store** (`vector_store/`)
- Qdrant client
- Vector operations

## Adding New Features

### Add New File Format Support

1. Create extractor in `services/python-rag-service/extractors/`
2. Implement `extract_<format>()` method
3. Add to `TextExtractor` factory
4. Update `ALLOWED_EXTENSIONS` in config
5. Add tests
6. Update documentation

### Add New Endpoint

**Go Service**:
1. Create handler in `internal/handler/http/`
2. Add route in `cmd/server/main.go`
3. Add service logic in `internal/service/`
4. Add tests
5. Update API documentation

**Python Service**:
1. Create route in `routes/`
2. Add to main.py router
3. Add business logic
4. Add tests
5. Update API documentation

### Modify Chunking Strategy

1. Edit `services/python-rag-service/chunking/chunker.py`
2. Update `CHUNK_SIZE` and `CHUNK_OVERLAP` in config
3. Update tests
4. Verify embedding performance
5. Update documentation

### Change Embedding Model

1. Edit `services/python-rag-service/config.py`
2. Update `EMBEDDING_MODEL` and `EMBEDDING_DIMENSION`
3. Update Qdrant collection creation
4. Re-index all documents
5. Update tests
6. Update documentation

## Testing

### Run All Tests

```bash
pytest tests/ -v
```

### Run Integration Tests

```bash
pytest tests/integration/ -v
```

### Run Performance Tests

```bash
pytest tests/performance/ -v
```

### Run Resilience Tests

```bash
pytest tests/resilience/ -v
```

### Test Coverage

```bash
pytest --cov=services tests/
```

### Write New Tests

1. Create test file in appropriate directory
2. Use fixtures from `conftest.py`
3. Follow naming convention: `test_<feature>.py`
4. Use descriptive test names
5. Include docstrings
6. Run tests locally before committing

## Code Style

### Go

- Follow Go conventions
- Use `gofmt` for formatting
- Use `golint` for linting
- Use `go vet` for analysis
- Write tests for all functions

### Python

- Follow PEP 8
- Use `black` for formatting
- Use `pylint` for linting
- Use `mypy` for type checking
- Write docstrings for all functions

## Git Workflow

### Branch Naming

- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `docs/description` - Documentation
- `test/description` - Tests

### Commit Messages

```
<type>: <description>

<body>

<footer>
```

Types: feat, fix, docs, test, refactor, perf, chore

### Pull Request Process

1. Create feature branch
2. Make changes
3. Write tests
4. Update documentation
5. Create pull request
6. Code review
7. Merge to main

## Debugging

### Go Service

```bash
# Enable debug logging
LOG_LEVEL=DEBUG go run cmd/server/main.go

# Use delve debugger
dlv debug cmd/server/main.go
```

### Python Service

```bash
# Enable debug logging
LOG_LEVEL=DEBUG python main.py

# Use pdb debugger
import pdb; pdb.set_trace()
```

### Docker Debugging

```bash
# View logs
docker-compose logs -f service-name

# Execute command in container
docker-compose exec service-name bash

# Inspect container
docker inspect container-name
```

## Performance Optimization

### Database

- Use indexes on frequently queried columns
- Use prepared statements
- Monitor query performance
- Use connection pooling

### Embeddings

- Use batch processing
- Cache embeddings
- Use GPU acceleration if available
- Monitor memory usage

### API

- Use caching headers
- Implement pagination
- Use compression
- Monitor response times

## Security Best Practices

- Never commit secrets
- Use environment variables
- Validate all inputs
- Use prepared statements
- Hash passwords with bcrypt
- Implement rate limiting
- Use HTTPS in production
- Keep dependencies updated

## Documentation

### Code Comments

- Explain why, not what
- Use clear language
- Keep comments updated
- Document edge cases

### API Documentation

- Document all endpoints
- Include examples
- Document error codes
- Document authentication

### README

- Project overview
- Quick start guide
- Architecture diagram
- Contributing guidelines

## Deployment

### Development

```bash
docker-compose up -d
```

### Production

See DEPLOYMENT-GUIDE.md

## Troubleshooting

### Common Issues

**Import errors in Go**:
```bash
go mod tidy
go mod download
```

**Import errors in Python**:
```bash
pip install -r requirements.txt
```

**Port already in use**:
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

**Database connection failed**:
```bash
# Check PostgreSQL
docker-compose logs postgres

# Verify DATABASE_URL
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Python Documentation](https://docs.python.org/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [gRPC Documentation](https://grpc.io/docs/)
- [Docker Documentation](https://docs.docker.com/)

## Support

For questions or issues:
1. Check documentation
2. Review code comments
3. Check tests for examples
4. Ask team members
5. Create GitHub issue
