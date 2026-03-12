# gRPC Client for RAG Service

This package contains the gRPC client for communicating with the Python RAG service.

## Proto Generation

To generate Go code from proto files:

```bash
# Install protoc compiler and Go plugins
# On macOS: brew install protobuf
# On Linux: apt-get install protobuf-compiler

# Install Go gRPC plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code
protoc --go_out=. --go-grpc_out=. proto/rag_service.proto
```

## Generated Files

After running protoc, the following files will be generated:
- `pb/rag_service.pb.go` - Protocol buffer definitions
- `pb/rag_service_grpc.pb.go` - gRPC service definitions

## Usage

```go
import "github.com/Morraban-Grid/docmind/services/go-user-service/internal/grpc_client"

// Create client
client, err := grpc_client.NewRAGClient(logger)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Call IndexDocument
success, chunkCount, embeddingCount, err := client.IndexDocument(
    ctx,
    documentID,
    userID,
    filePath,
    fileType,
)
```

## Configuration

Environment variables:
- `PYTHON_GRPC_HOST` - Host of Python gRPC service (default: localhost)
- `PYTHON_GRPC_PORT` - Port of Python gRPC service (default: 50051)

## Timeouts

- IndexDocument: 30 seconds
- DeleteDocument: 10 seconds
- QueryDocument: 45 seconds
- HealthCheck: 5 seconds

## Error Handling

All methods return errors that should be handled appropriately:
- Connection errors are logged and returned
- Timeout errors are handled with context cancellation
- Service errors are returned in the response message
