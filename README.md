# DocMind

AI document reasoning system using RAG architecture.

DocMind is a production-oriented microservices-based platform that enables users to upload private documents and query them using natural language. The system leverages Retrieval-Augmented Generation (RAG) to transform unstructured documents into structured semantic knowledge.

---

## 🚀 Overview

DocMind is composed of two main microservices:

- **Go Service** → Authentication, document management, metadata handling
- **Python RAG Service** → Document ingestion, embeddings generation, semantic retrieval, LLM-based reasoning

The system is fully dockerized and runs locally using open-source tools.

---

## 🧠 Architecture

- Microservices architecture
- Vector database for semantic search
- JWT-based authentication
- Docker-based infrastructure
- Modular RAG pipeline

See `docs/architecture.md` for detailed diagrams.

---

## 🛠 Tech Stack

### Backend (API & Management)
- Go
- Gin
- PostgreSQL
- MinIO

### AI & Retrieval
- Python
- FastAPI
- LangChain
- Qdrant
- Sentence Transformers

### Infrastructure
- Docker
- Docker Compose

---

## 📦 Setup

Instructions will be provided after infrastructure configuration.

---

## 📌 Roadmap

- [ ] Infrastructure setup
- [ ] Authentication service
- [ ] Document ingestion pipeline
- [ ] Semantic retrieval
- [ ] Integration testing
- [ ] Production hardening

---

## 📄 License

MIT License
