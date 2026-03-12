package httphandler

import (
	"net/http"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/grpc_client"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	ragClient *grpc_client.RAGClient
}

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type QueryResponse struct {
	Answer     string   `json:"answer"`
	Sources    []string `json:"sources"`
	ChunkCount int      `json:"chunk_count"`
}

func NewQueryHandler(ragClient *grpc_client.RAGClient) *QueryHandler {
	return &QueryHandler{
		ragClient: ragClient,
	}
}

// QueryDocuments handles POST /api/query
func (h *QueryHandler) QueryDocuments(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		infrastructure.LogError("user_id not found in context", nil)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		infrastructure.LogError("Invalid query request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Query == "" {
		infrastructure.LogError("Query cannot be empty", nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
		return
	}

	infrastructure.LogInfo("Query request from user", map[string]interface{}{
		"user_id": userID,
		"query":   req.Query[:min(len(req.Query), 100)],
	})

	// Call gRPC service
	result, err := h.ragClient.QueryDocument(req.Query, userID)
	if err != nil {
		infrastructure.LogError("Query failed", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query processing failed"})
		return
	}

	infrastructure.LogInfo("Query completed successfully", map[string]interface{}{
		"user_id":     userID,
		"chunk_count": result.ChunkCount,
	})

	c.JSON(http.StatusOK, QueryResponse{
		Answer:     result.Answer,
		Sources:    result.Sources,
		ChunkCount: result.ChunkCount,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
