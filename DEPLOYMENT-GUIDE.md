# DocMind RAG System - Deployment Guide

## System Requirements

### Minimum Requirements
- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM
- 10GB disk space
- Linux, macOS, or Windows with WSL2

### Recommended Requirements
- Docker 24.0+
- Docker Compose 2.20+
- 8GB RAM
- 50GB disk space
- Linux (Ubuntu 22.04 LTS or later)

## Pre-Deployment Checklist

- [ ] Docker and Docker Compose installed
- [ ] `.env` file created from `.env.example`
- [ ] All required environment variables set
- [ ] Sufficient disk space available
- [ ] Network connectivity verified
- [ ] Ports 8080, 8001, 5432, 9000, 6333, 11434 available

## Environment Configuration

### 1. Create .env File

```bash
cp .env.example .env
```

### 2. Configure Environment Variables

Edit `.env` with production values:

```bash
# Go Service
SERVER_PORT=8080

# Database
DATABASE_URL=postgres://user:password@postgres:5432/docmind

# JWT
JWT_SECRET=your-secret-key-here-minimum-32-characters

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_BUCKET=documents
MINIO_USE_SSL=false

# Qdrant
QDRANT_HOST=qdrant
QDRANT_PORT=6333

# Ollama
OLLAMA_HOST=ollama
OLLAMA_PORT=11434

# Python Service
PYTHON_GRPC_HOST=python-rag-service
PYTHON_GRPC_PORT=50051
```

### 3. Security Considerations

**IMPORTANT**: Never commit `.env` file to version control.

- Use strong JWT_SECRET (minimum 32 characters)
- Use strong database password
- Use strong MinIO credentials
- Rotate credentials regularly
- Use environment-specific values

## Deployment Steps

### 1. Start Services

```bash
docker-compose up -d
```

### 2. Verify Services

```bash
# Check all containers running
docker-compose ps

# Expected output:
# NAME                    STATUS
# docmind-postgres        Up
# docmind-minio           Up
# docmind-qdrant          Up
# docmind-ollama          Up
# docmind-go-service      Up
# docmind-python-service  Up
```

### 3. Verify Health Checks

```bash
# Go service health
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "healthy",
#   "service": "docmind-go-service",
#   "dependencies": [...]
# }
```

### 4. Run Integration Tests

```bash
# Run tests to verify complete flow
python -m pytest tests/integration/ -v
```

## Service Management

### Start Services

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f go-user-service
docker-compose logs -f python-rag-service
```

### Restart Service

```bash
docker-compose restart go-user-service
```

### Rebuild Images

```bash
docker-compose build --no-cache
docker-compose up -d
```

## Database Management

### Run Migrations

Migrations run automatically on container startup.

### Backup Database

```bash
# Backup PostgreSQL
docker-compose exec postgres pg_dump -U user docmind > backup.sql

# Backup Qdrant
docker-compose exec qdrant tar czf /tmp/qdrant-backup.tar.gz /qdrant/storage
docker cp docmind-qdrant:/tmp/qdrant-backup.tar.gz ./qdrant-backup.tar.gz
```

### Restore Database

```bash
# Restore PostgreSQL
docker-compose exec -T postgres psql -U user docmind < backup.sql

# Restore Qdrant
docker cp qdrant-backup.tar.gz docmind-qdrant:/tmp/
docker-compose exec qdrant tar xzf /tmp/qdrant-backup.tar.gz -C /
```

## Monitoring

### Health Checks

```bash
# Go service
curl http://localhost:8080/health

# Python service
curl http://localhost:8001/health
```

### View Metrics

```bash
# Check container resource usage
docker stats

# Check logs for errors
docker-compose logs | grep ERROR
```

### Performance Monitoring

Monitor these metrics:
- Request latency
- Error rate
- Database connection pool
- Memory usage
- Disk usage

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker-compose logs service-name

# Common issues:
# - Port already in use
# - Insufficient disk space
# - Database connection failed
# - Environment variables missing
```

### Database Connection Failed

```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Verify DATABASE_URL in .env
```

### gRPC Connection Failed

```bash
# Verify Python service is running
docker-compose ps python-rag-service

# Check Python service logs
docker-compose logs python-rag-service

# Verify PYTHON_GRPC_HOST and PYTHON_GRPC_PORT
```

### Ollama Not Available

```bash
# Verify Ollama is running
docker-compose ps ollama

# Check Ollama logs
docker-compose logs ollama

# Verify OLLAMA_HOST and OLLAMA_PORT
```

### High Memory Usage

```bash
# Check container memory
docker stats

# Reduce batch size in config
# Reduce concurrent connections
# Increase available RAM
```

### Slow Query Performance

```bash
# Check Qdrant status
curl http://localhost:6333/health

# Check database indexes
docker-compose exec postgres psql -U user -d docmind -c "\d+ documents"

# Monitor query logs
docker-compose logs python-rag-service | grep "Query"
```

## Production Deployment

### Pre-Production Checklist

- [ ] All tests passing
- [ ] Security audit completed
- [ ] Performance tested
- [ ] Backup strategy verified
- [ ] Monitoring configured
- [ ] Logging configured
- [ ] Error handling verified
- [ ] Documentation reviewed

### Production Configuration

```bash
# Use strong credentials
JWT_SECRET=<generate-strong-secret>
MINIO_ROOT_PASSWORD=<strong-password>
POSTGRES_PASSWORD=<strong-password>

# Use production endpoints
MINIO_ENDPOINT=minio.production.com:9000
QDRANT_HOST=qdrant.production.com
OLLAMA_HOST=ollama.production.com

# Enable SSL/TLS
MINIO_USE_SSL=true
```

### Scaling

To scale services:

```bash
# Scale Go service to 3 instances
docker-compose up -d --scale go-user-service=3

# Use load balancer (nginx, HAProxy) in front
```

### Backup Strategy

- Daily PostgreSQL backups
- Daily Qdrant backups
- Store backups off-site
- Test restore procedures monthly
- Document recovery procedures

## Rollback Procedure

If deployment fails:

```bash
# Stop current deployment
docker-compose down

# Restore from backup
docker-compose exec -T postgres psql -U user docmind < backup.sql

# Start services
docker-compose up -d

# Verify health
curl http://localhost:8080/health
```

## Support

For issues or questions:
1. Check logs: `docker-compose logs`
2. Review troubleshooting section
3. Check API documentation
4. Review deployment guide
5. Contact support team

## Next Steps

1. Configure monitoring and alerting
2. Set up automated backups
3. Configure log aggregation
4. Set up CI/CD pipeline
5. Plan capacity scaling
