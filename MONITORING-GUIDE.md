# DocMind RAG System - Monitoring Guide

## Overview

Comprehensive monitoring strategy for DocMind RAG System including health checks, metrics, logging, and alerting.

## Health Checks

### Go Service Health

**Endpoint**: `GET http://localhost:8080/health`

**Response**:
```json
{
  "status": "healthy",
  "service": "docmind-go-service",
  "dependencies": [
    {
      "name": "database",
      "status": "healthy",
      "available": true
    },
    {
      "name": "grpc",
      "status": "healthy",
      "available": true
    }
  ]
}
```

### Python Service Health

**Endpoint**: `GET http://localhost:8001/health`

**Response**:
```json
{
  "status": "healthy",
  "service": "python-rag-service",
  "dependencies": [
    {
      "name": "qdrant",
      "status": "healthy",
      "available": true
    },
    {
      "name": "ollama",
      "status": "healthy",
      "available": true
    }
  ]
}
```

## Metrics

### Key Metrics to Monitor

**Request Metrics**:
- Request count (total, by endpoint, by status)
- Request latency (p50, p95, p99)
- Error rate (by endpoint, by error type)
- Request size (average, max)

**Database Metrics**:
- Connection pool usage
- Query latency
- Query count
- Slow queries

**Vector Store Metrics**:
- Search latency
- Upsert latency
- Collection size
- Memory usage

**LLM Metrics**:
- Generation latency
- Token count
- Error rate
- Timeout count

**System Metrics**:
- CPU usage
- Memory usage
- Disk usage
- Network I/O

### Prometheus Metrics

Optional Prometheus endpoint for metrics collection:

```
GET http://localhost:8080/metrics
```

Metrics format:
```
# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/health",status="200"} 1000

# HELP http_request_duration_seconds HTTP request duration
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{endpoint="/documents",le="0.1"} 500
```

## Logging

### Log Levels

- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages for potential issues
- **ERROR**: Error messages for failures
- **FATAL**: Fatal errors requiring immediate attention

### Log Format

```
2026-03-12T10:30:00Z [INFO] request_id=req-12345 user_id=user-uuid endpoint=POST /api/documents message="Document uploaded successfully"
```

### Log Locations

**Go Service**:
```bash
docker-compose logs go-user-service
```

**Python Service**:
```bash
docker-compose logs python-rag-service
```

### Important Log Messages

**Authentication**:
```
[INFO] User login successful: user_id=user-uuid
[WARN] Failed login attempt: email=user@example.com
[ERROR] Invalid JWT token: request_id=req-12345
```

**Document Operations**:
```
[INFO] Document uploaded: document_id=doc-uuid file_size=1024000
[INFO] Document indexed: document_id=doc-uuid chunk_count=50
[ERROR] Document indexing failed: document_id=doc-uuid error=...
```

**Query Operations**:
```
[INFO] Query received: user_id=user-uuid query_length=50
[INFO] Search completed: chunk_count=5 latency=250ms
[ERROR] Query failed: user_id=user-uuid error=...
```

**System Events**:
```
[INFO] Service started: version=1.0.0
[WARN] Database connection lost: reconnecting...
[ERROR] Ollama service unavailable: status=503
```

## Alerting

### Alert Conditions

**Critical**:
- Service down (health check failing)
- Database connection lost
- Disk space < 1GB
- Memory usage > 90%
- Error rate > 5%

**Warning**:
- Response latency > 5 seconds
- Error rate > 1%
- Database connection pool > 80%
- Disk space < 5GB
- Memory usage > 70%

**Info**:
- Service restart
- Configuration change
- Backup completed
- Maintenance window

### Alert Channels

Configure alerts to:
- Email
- Slack
- PagerDuty
- SMS (for critical)

## Dashboards

### Key Dashboard Panels

**System Health**:
- Service status (up/down)
- Dependency status
- Error rate
- Response time

**Request Metrics**:
- Requests per second
- Request latency (p50, p95, p99)
- Error rate by endpoint
- Top slow endpoints

**Database**:
- Connection pool usage
- Query latency
- Slow queries
- Database size

**Vector Store**:
- Search latency
- Upsert latency
- Collection size
- Memory usage

**LLM**:
- Generation latency
- Token count
- Error rate
- Timeout count

## Monitoring Tools

### Docker Stats

```bash
# Real-time container metrics
docker stats

# Specific container
docker stats go-user-service
```

### Log Aggregation

For production, use:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Splunk
- Datadog
- CloudWatch

### Metrics Collection

For production, use:
- Prometheus + Grafana
- Datadog
- New Relic
- CloudWatch

## Performance Baselines

### Expected Performance

**Document Upload**:
- Latency: < 5 seconds
- Success rate: > 99%

**Document Indexing**:
- Latency: 10-30 seconds (depends on file size)
- Success rate: > 99%

**Document Query**:
- Latency: < 5 seconds
- Success rate: > 99%

**Search**:
- Latency: < 500ms
- Success rate: > 99%

**LLM Generation**:
- Latency: 5-30 seconds
- Success rate: > 95%

## Capacity Planning

### Resource Usage

**Go Service**:
- CPU: 100-500m per instance
- Memory: 200-500MB per instance
- Disk: 1GB (logs)

**Python Service**:
- CPU: 500m-2 per instance
- Memory: 1-2GB per instance
- Disk: 2GB (models, logs)

**PostgreSQL**:
- CPU: 200-500m
- Memory: 1-2GB
- Disk: 10-50GB (depends on documents)

**Qdrant**:
- CPU: 200-500m
- Memory: 2-4GB
- Disk: 5-20GB (depends on embeddings)

**Ollama**:
- CPU: 1-4 cores
- Memory: 4-8GB
- Disk: 5-10GB (models)

### Scaling Recommendations

**Horizontal Scaling**:
- Scale Go service to 3-5 instances
- Scale Python service to 2-3 instances
- Use load balancer (nginx, HAProxy)

**Vertical Scaling**:
- Increase RAM for better performance
- Use faster CPU for LLM inference
- Use SSD for database

## Troubleshooting

### High CPU Usage

```bash
# Check which process
docker stats

# Check logs for errors
docker-compose logs service-name | grep ERROR

# Possible causes:
# - Slow queries
# - Large batch processing
# - Memory pressure causing swapping
```

### High Memory Usage

```bash
# Check memory usage
docker stats

# Reduce batch size
# Reduce concurrent connections
# Increase available RAM
```

### High Disk Usage

```bash
# Check disk usage
df -h

# Check container logs
du -sh /var/lib/docker/containers/*/

# Rotate logs
docker-compose logs --tail 1000 > logs.txt
```

### Slow Queries

```bash
# Check database logs
docker-compose logs postgres | grep slow

# Check query performance
docker-compose exec postgres psql -U user -d docmind -c "EXPLAIN ANALYZE SELECT ..."

# Add indexes if needed
```

### Connection Failures

```bash
# Check service status
docker-compose ps

# Check logs
docker-compose logs service-name

# Verify network connectivity
docker-compose exec service-name ping other-service
```

## Maintenance

### Regular Tasks

**Daily**:
- Check health endpoints
- Review error logs
- Monitor resource usage

**Weekly**:
- Review performance metrics
- Check backup status
- Review security logs

**Monthly**:
- Capacity planning review
- Performance optimization
- Security audit

**Quarterly**:
- Disaster recovery drill
- Performance baseline update
- Architecture review

## Backup Monitoring

### Backup Status

```bash
# Check backup files
ls -lh backups/

# Verify backup integrity
tar -tzf backup.tar.gz > /dev/null && echo "OK" || echo "FAILED"
```

### Backup Alerts

Alert if:
- Backup fails
- Backup takes > 1 hour
- Backup size > expected
- Backup not completed in 24 hours

## Security Monitoring

### Security Events

Monitor for:
- Failed login attempts
- Unauthorized access attempts
- SQL injection attempts
- File upload violations
- API abuse

### Security Logs

```bash
# Failed logins
docker-compose logs go-user-service | grep "Failed login"

# Unauthorized access
docker-compose logs go-user-service | grep "Unauthorized"

# API errors
docker-compose logs go-user-service | grep "ERROR"
```

## Compliance

### Audit Logging

Log all:
- User authentication
- Document access
- Data modifications
- Administrative actions

### Data Retention

- Keep logs for 90 days
- Keep backups for 1 year
- Archive old logs to cold storage

## Reporting

### Daily Report

- Service uptime
- Error count
- Performance metrics
- Resource usage

### Weekly Report

- Performance trends
- Capacity usage
- Security events
- Backup status

### Monthly Report

- Performance analysis
- Capacity planning
- Cost analysis
- Recommendations

## Next Steps

1. Set up monitoring tools
2. Configure alerting
3. Create dashboards
4. Establish baselines
5. Document procedures
6. Train team
