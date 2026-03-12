package grpc_client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RAGClient wraps gRPC client for RAG service
type RAGClient struct {
	conn    *grpc.ClientConn
	client  interface{} // Will be properly typed after proto generation
	logger  *slog.Logger
	timeout time.Duration
}

// NewRAGClient creates a new RAG gRPC client
func NewRAGClient(logger *slog.Logger) (*RAGClient, error) {
	// Get connection details from environment
	host := os.Getenv("PYTHON_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PYTHON_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	logger.Info("Connecting to RAG gRPC service", slog.String("address", address))

	// Create connection with insecure credentials (for development)
	// In production, use TLS credentials
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 * 1024 * 1024), // 10MB
		),
	)
	if err != nil {
		logger.Error("Failed to connect to RAG service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to RAG service: %w", err)
	}

	logger.Info("Connected to RAG gRPC service")

	return &RAGClient{
		conn:    conn,
		logger:  logger,
		timeout: 30 * time.Second, // Default timeout for IndexDocument
	}, nil
}

// IndexDocument calls gRPC IndexDocument method
func (c *RAGClient) IndexDocument(ctx context.Context, documentID, userID, filePath, fileType string) (bool, int32, int32, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.logger.Info("Calling gRPC IndexDocument",
		slog.String("document_id", documentID),
		slog.String("user_id", userID),
	)

	// Call will be implemented after proto generation
	// For now, return placeholder
	c.logger.Info("IndexDocument call completed",
		slog.String("document_id", documentID),
	)

	return true, 0, 0, nil
}

// DeleteDocument calls gRPC DeleteDocument method
func (c *RAGClient) DeleteDocument(ctx context.Context, documentID string) (bool, int32, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c.logger.Info("Calling gRPC DeleteDocument",
		slog.String("document_id", documentID),
	)

	// Call will be implemented after proto generation
	// For now, return placeholder
	c.logger.Info("DeleteDocument call completed",
		slog.String("document_id", documentID),
	)

	return true, 0, nil
}

// QueryDocument calls gRPC QueryDocument method
func (c *RAGClient) QueryDocument(ctx context.Context, query, userID string) (string, []string, int32, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	c.logger.Info("Calling gRPC QueryDocument",
		slog.String("user_id", userID),
		slog.String("query", query),
	)

	// Call will be implemented after proto generation
	// For now, return placeholder
	c.logger.Info("QueryDocument call completed",
		slog.String("user_id", userID),
	)

	return "", []string{}, 0, nil
}

// Close closes the gRPC connection
func (c *RAGClient) Close() error {
	if c.conn != nil {
		c.logger.Info("Closing gRPC connection")
		return c.conn.Close()
	}
	return nil
}

// HealthCheck checks if RAG service is available
func (c *RAGClient) HealthCheck(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to get connection state
	state := c.conn.GetState()
	c.logger.Debug("RAG service connection state", slog.String("state", state.String()))

	return true
}
