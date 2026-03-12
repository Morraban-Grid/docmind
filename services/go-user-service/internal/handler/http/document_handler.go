package http

import (
	"net/http"
	"strconv"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/domain"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/service"
	"github.com/gin-gonic/gin"
)

// DocumentHandler handles document-related HTTP requests
type DocumentHandler struct {
	documentService *service.DocumentService
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// UploadDocument handles document upload
// POST /api/documents
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		infrastructure.LogError("failed to parse file upload", err, nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Upload document
	req := service.UploadDocumentRequest{
		File:   file,
		UserID: userID.(string),
	}

	doc, err := h.documentService.UploadDocument(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// GetDocument retrieves a document by ID
// GET /api/documents/:id
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, _ := c.Get("user_id")

	doc, err := h.documentService.GetDocument(documentID, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments retrieves all documents for the authenticated user
// GET /api/documents
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	docs, err := h.documentService.ListDocuments(userID.(string), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, docs)
}

// DownloadDocument downloads a document file
// GET /api/documents/:id/download
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, _ := c.Get("user_id")

	reader, filename, err := h.documentService.DownloadDocument(c.Request.Context(), documentID, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}
	defer reader.Close()

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	// Stream file to response
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", reader, nil)
}

// DeleteDocument deletes a document
// DELETE /api/documents/:id
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	userID, _ := c.Get("user_id")

	err := h.documentService.DeleteDocument(c.Request.Context(), documentID, userID.(string))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// handleError handles different types of errors and returns appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *domain.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
	case *domain.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": e.Message})
	case *domain.AuthorizationError:
		c.JSON(http.StatusForbidden, gin.H{"error": e.Message})
	case *domain.InternalError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
