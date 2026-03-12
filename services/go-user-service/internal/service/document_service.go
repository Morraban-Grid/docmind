package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/domain"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
	"github.com/google/uuid"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var (
	AllowedExtensions = map[string]bool{
		".pdf":  true,
		".txt":  true,
		".docx": true,
		".md":   true,
	}

	MimeTypes = map[string]string{
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".md":   "text/markdown",
	}
)

// DocumentService handles document business logic
type DocumentService struct {
	docRepo     domain.DocumentRepository
	minioClient *infrastructure.MinIOClient
}

// NewDocumentService creates a new document service
func NewDocumentService(docRepo domain.DocumentRepository, minioClient *infrastructure.MinIOClient) *DocumentService {
	return &DocumentService{
		docRepo:     docRepo,
		minioClient: minioClient,
	}
}

// UploadDocumentRequest represents a document upload request
type UploadDocumentRequest struct {
	File   *multipart.FileHeader
	UserID string
}

// DocumentResponse represents a document response
type DocumentResponse struct {
	DocumentID  string    `json:"document_id"`
	Filename    string    `json:"filename"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	Status      string    `json:"status"`
	StoragePath string    `json:"storage_path"`
	UploadDate  time.Time `json:"upload_date"`
}

// DocumentListResponse represents a paginated list of documents
type DocumentListResponse struct {
	Documents   []*DocumentResponse `json:"documents"`
	TotalItems  int                 `json:"total_items"`
	TotalPages  int                 `json:"total_pages"`
	CurrentPage int                 `json:"current_page"`
	PageSize    int                 `json:"page_size"`
}

// ValidateFile validates file extension and size
func ValidateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > MaxFileSize {
		return domain.NewValidationError(
			fmt.Sprintf("File size exceeds maximum allowed size of %d MB", MaxFileSize/(1024*1024)),
			nil,
		)
	}

	// Check if file is empty
	if file.Size == 0 {
		return domain.NewValidationError("File is empty", nil)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedExtensions[ext] {
		return domain.NewValidationError(
			fmt.Sprintf("Unsupported file type: %s. Allowed types: PDF, TXT, DOCX, MD", ext),
			nil,
		)
	}

	return nil
}

// ValidateUUID validates UUID format
func ValidateUUID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return domain.NewValidationError("Invalid document ID format", err)
	}
	return nil
}

// UploadDocument handles document upload
func (s *DocumentService) UploadDocument(ctx context.Context, req UploadDocumentRequest) (*DocumentResponse, error) {
	// Validate file
	if err := ValidateFile(req.File); err != nil {
		return nil, err
	}

	// Generate unique document ID
	documentID := uuid.New().String()

	// Get file extension and mime type
	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	mimeType := MimeTypes[ext]

	// Open uploaded file
	file, err := req.File.Open()
	if err != nil {
		infrastructure.LogError("failed to open uploaded file", err, map[string]interface{}{
			"filename": req.File.Filename,
		})
		return nil, domain.NewInternalError("Failed to process uploaded file")
	}
	defer file.Close()

	// Create storage path: {user_id}/{document_id}/{filename}
	objectName := fmt.Sprintf("%s/%s/%s", req.UserID, documentID, req.File.Filename)

	// Upload to MinIO
	storagePath, err := s.minioClient.UploadFile(ctx, objectName, file, req.File.Size, mimeType)
	if err != nil {
		infrastructure.LogError("failed to upload to MinIO", err, map[string]interface{}{
			"document_id": documentID,
			"user_id":     req.UserID,
		})
		return nil, domain.NewInternalError("Failed to store document")
	}

	// Create document record
	doc := &domain.Document{
		DocumentID:  documentID,
		UserID:      req.UserID,
		Filename:    req.File.Filename,
		FileSize:    req.File.Size,
		MimeType:    mimeType,
		StoragePath: storagePath,
		Status:      "pending_indexing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.docRepo.CreateDocument(doc); err != nil {
		// Attempt to clean up MinIO file
		_ = s.minioClient.DeleteFile(ctx, objectName)
		return nil, err
	}

	infrastructure.LogInfo("document uploaded successfully", map[string]interface{}{
		"document_id": documentID,
		"user_id":     req.UserID,
		"filename":    req.File.Filename,
		"file_size":   req.File.Size,
	})

	return &DocumentResponse{
		DocumentID:  doc.DocumentID,
		Filename:    doc.Filename,
		FileSize:    doc.FileSize,
		MimeType:    doc.MimeType,
		Status:      doc.Status,
		StoragePath: doc.StoragePath,
		UploadDate:  doc.CreatedAt,
	}, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(documentID, userID string) (*DocumentResponse, error) {
	// Validate UUID
	if err := ValidateUUID(documentID); err != nil {
		return nil, err
	}

	doc, err := s.docRepo.GetDocumentByID(documentID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if doc.UserID != userID {
		infrastructure.LogInfo("unauthorized document access attempt", map[string]interface{}{
			"document_id":  documentID,
			"requester_id": userID,
			"owner_id":     doc.UserID,
		})
		return nil, domain.NewAuthorizationError("You do not have permission to access this document")
	}

	return &DocumentResponse{
		DocumentID:  doc.DocumentID,
		Filename:    doc.Filename,
		FileSize:    doc.FileSize,
		MimeType:    doc.MimeType,
		Status:      doc.Status,
		StoragePath: doc.StoragePath,
		UploadDate:  doc.CreatedAt,
	}, nil
}

// ListDocuments retrieves all documents for a user with pagination
func (s *DocumentService) ListDocuments(userID string, page, pageSize int) (*DocumentListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Get documents
	docs, err := s.docRepo.GetDocumentsByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := s.docRepo.CountDocumentsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	docResponses := make([]*DocumentResponse, len(docs))
	for i, doc := range docs {
		docResponses[i] = &DocumentResponse{
			DocumentID:  doc.DocumentID,
			Filename:    doc.Filename,
			FileSize:    doc.FileSize,
			MimeType:    doc.MimeType,
			Status:      doc.Status,
			StoragePath: doc.StoragePath,
			UploadDate:  doc.CreatedAt,
		}
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	return &DocumentListResponse{
		Documents:   docResponses,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}

// DownloadDocument retrieves a document file from storage
func (s *DocumentService) DownloadDocument(ctx context.Context, documentID, userID string) (io.ReadCloser, string, error) {
	// Validate UUID
	if err := ValidateUUID(documentID); err != nil {
		return nil, "", err
	}

	doc, err := s.docRepo.GetDocumentByID(documentID)
	if err != nil {
		return nil, "", err
	}

	// Check ownership
	if doc.UserID != userID {
		infrastructure.LogInfo("unauthorized document download attempt", map[string]interface{}{
			"document_id":  documentID,
			"requester_id": userID,
			"owner_id":     doc.UserID,
		})
		return nil, "", domain.NewAuthorizationError("You do not have permission to download this document")
	}

	// Extract object name from storage path
	objectName := fmt.Sprintf("%s/%s/%s", doc.UserID, doc.DocumentID, doc.Filename)

	// Download from MinIO
	object, err := s.minioClient.DownloadFile(ctx, objectName)
	if err != nil {
		infrastructure.LogError("failed to download document", err, map[string]interface{}{
			"document_id": documentID,
		})
		return nil, "", domain.NewInternalError("Failed to download document")
	}

	infrastructure.LogInfo("document downloaded", map[string]interface{}{
		"document_id": documentID,
		"user_id":     userID,
	})

	return object, doc.Filename, nil
}

// DeleteDocument removes a document
func (s *DocumentService) DeleteDocument(ctx context.Context, documentID, userID string) error {
	// Validate UUID
	if err := ValidateUUID(documentID); err != nil {
		return err
	}

	doc, err := s.docRepo.GetDocumentByID(documentID)
	if err != nil {
		return err
	}

	// Check ownership
	if doc.UserID != userID {
		infrastructure.LogInfo("unauthorized document deletion attempt", map[string]interface{}{
			"document_id":  documentID,
			"requester_id": userID,
			"owner_id":     doc.UserID,
		})
		return domain.NewAuthorizationError("You do not have permission to delete this document")
	}

	// Extract object name from storage path
	objectName := fmt.Sprintf("%s/%s/%s", doc.UserID, doc.DocumentID, doc.Filename)

	// Delete from MinIO
	if err := s.minioClient.DeleteFile(ctx, objectName); err != nil {
		infrastructure.LogError("failed to delete file from MinIO", err, map[string]interface{}{
			"document_id": documentID,
		})
		// Continue with database deletion even if MinIO deletion fails
	}

	// Delete from database
	if err := s.docRepo.DeleteDocument(documentID); err != nil {
		return err
	}

	infrastructure.LogInfo("document deleted successfully", map[string]interface{}{
		"document_id": documentID,
		"user_id":     userID,
	})

	return nil
}
