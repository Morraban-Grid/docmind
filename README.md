# DocMind - RAG Document Reasoning System

DocMind is a professional private document reasoning system based on RAG (Retrieval-Augmented Generation) architecture. The system allows authenticated users to upload documents, transform them into vector representations through embeddings, and perform natural language queries that combine semantic retrieval with contextualized response generation.

## 🏗️ Architecture

The system consists of two main microservices:

- **Go Service**: User management, authentication, file upload, and access control
- **Python RAG Service**: Complete RAG pipeline (chunking, embeddings, semantic retrieval, generation)

### Technology Stack

**Infrastructure:**
- PostgreSQL 16.4 - Relational database for metadata
- Qdrant v1.11.3 - Vector database for embeddings
- MinIO - S3-compatible object storage
- Ollama - Local LLM (llama2)
- Docker & Docker Compose - Containerization

**Go Service:**
- Gin - Web framework
- GORM - ORM
- JWT - Authentication
- gRPC - Inter-service communication

**Python Service:**
- FastAPI - Web framework
- LangChain - RAG pipeline
- sentence-transformers - Embeddings (all-MiniLM-L6-v2)
- Qdrant Client - Vector storage
- PyPDF2, python-docx - Document parsing

## 🚀 Quick Start

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM minimum
- 10GB disk space

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/docmind.git
cd docmind
```

2. **Create environment configuration**
```bash
cp .env.example .env
```

3. **Edit .env file with your configuration**
```bash
# IMPORTANT: Change these values for security
nano .env

# At minimum, change:
# - POSTGRES_PASSWORD
# - MINIO_ROOT_PASSWORD
# - JWT_SECRET (use a strong random string, min 32 characters)
```

4. **Start all services**
```bash
docker compose -f deployments/docker/docker-compose.yml up -d
```

5. **Verify all services are running**
```bash
docker compose -f deployments/docker/docker-compose.yml ps
```

You should see all services in "healthy" state:
- docmind-postgres
- docmind-qdrant
- docmind-minio
- docmind-ollama

### Accessing Services

- **Go Service API**: http://localhost:8080
- **Python Service API**: http://localhost:8000
- **MinIO Console**: http://localhost:9001
- **Qdrant Dashboard**: http://localhost:6333/dashboard

## 📖 Development Status

### ✅ Iteration 1: Infrastructure & Security (COMPLETED)
- Docker Compose infrastructure
- PostgreSQL database with migrations
- MinIO object storage
- Qdrant vector database
- Ollama LLM
- Security configuration (.gitignore, .env.example)
- Project structure for Go and Python services

### 🚧 Upcoming Iterations
- Iteration 2: Authentication & User Management
- Iteration 3: Document Upload & Storage
- Iteration 4: Document Processing & Chunking
- Iteration 5: Embeddings & Vector Storage
- Iteration 6: gRPC Integration
- Iteration 7: RAG Query & LLM Integration
- Iteration 8: Testing & Production Readiness

## 🔒 Security

**IMPORTANT**: Never commit sensitive information to version control.

- All secrets must be stored in `.env` file (which is gitignored)
- Use strong passwords and secrets in production
- JWT_SECRET must be at least 32 characters
- Change default credentials before deployment

## 📝 API Documentation

API documentation will be available at:
- Go Service: http://localhost:8080/swagger (Coming in Iteration 2)
- Python Service: http://localhost:8000/docs (Coming in Iteration 4)

## 🧪 Testing

```bash
# Run Go tests
cd services/go-user-service
go test ./...

# Run Python tests
cd services/python-rag-service
pytest
```

## 🛠️ Development

### Project Structure

```
docmind/
├── .kiro/specs/              # Specification documents
├── deployments/
│   ├── docker/               # Docker Compose files
│   └── k8s/                  # Kubernetes manifests (future)
├── docs/                     # Documentation
├── proto/                    # gRPC protocol definitions
├── scripts/                  # Utility scripts
├── services/
│   ├── go-user-service/      # Go microservice
│   │   ├── cmd/server/       # Application entry point
│   │   ├── internal/         # Internal packages
│   │   └── migrations/       # Database migrations
│   └── python-rag-service/   # Python microservice
│       ├── app/              # Application code
│       └── tests/            # Test files
├── .env.example              # Environment template
├── .gitignore                # Git ignore rules
└── README.md                 # This file
```

### Running Services Locally

**Go Service:**
```bash
cd services/go-user-service
go run cmd/server/main.go
```

**Python Service:**
```bash
cd services/python-rag-service
pip install -r requirements.txt
uvicorn app.main:app --reload
```

## 📊 Monitoring

View logs for all services:
```bash
docker compose -f deployments/docker/docker-compose.yml logs -f
```

View logs for specific service:
```bash
docker compose -f deployments/docker/docker-compose.yml logs -f go-service
docker compose -f deployments/docker/docker-compose.yml logs -f python-rag-service
```

## 🐛 Troubleshooting

### Services won't start
1. Check Docker is running: `docker info`
2. Check ports are not in use: `netstat -an | grep -E '(8080|8000|5432|6333|9000|11434)'`
3. Check logs: `docker compose logs`

### Database connection errors
1. Verify PostgreSQL is healthy: `docker compose ps postgres`
2. Check DATABASE_URL in .env file
3. Verify migrations ran: `docker compose logs postgres`

### MinIO connection errors
1. Verify MinIO is healthy: `docker compose ps minio`
2. Check MINIO_ENDPOINT in .env file
3. Access MinIO console: http://localhost:9001

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## 📧 Contact

For questions or support, please open an issue on GitHub.

---

**Note**: This project uses only free and open-source tools. No paid services or API keys are required.
