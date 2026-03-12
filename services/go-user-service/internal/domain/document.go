package domain

import (
	"time"
)

// Document represents a document uploaded by a user
type Document struct {
	DocumentID string    `json:"document_id" db:"document_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Filename   string    `json:"filename" db:"filename"`
	FileSize   int64     `json:"file_size" db:"file_size"`
	MimeType   string    `json:"mime_type" db:"mime_type"`
	StoragePath string   `json:"storage_path" db:"storage_path"`
	Status     string    `json:"status" db:"status"` // pending_indexing, indexed, failed
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// DocumentRepository defines the interface for document data operations
type DocumentRepository interface {
	CreateDocument(doc *Document) error
	GetDocumentByID(documentID string) (*Document, error)
	GetDocumentsByUserID(userID string, limit, offset int) ([]*Document, error)
	CountDocumentsByUserID(userID string) (int, error)
	UpdateDocumentStatus(documentID, status string) error
	DeleteDocument(documentID string) error
}
