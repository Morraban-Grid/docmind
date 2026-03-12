package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/domain"
	"github.com/Morraban-Grid/docmind/services/go-user-service/internal/infrastructure"
)

// DocumentRepository implements domain.DocumentRepository for PostgreSQL
type DocumentRepository struct {
	db *sql.DB
}

// NewDocumentRepository creates a new PostgreSQL document repository
func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// CreateDocument inserts a new document into the database
func (r *DocumentRepository) CreateDocument(doc *domain.Document) error {
	query := `
		INSERT INTO documents (document_id, user_id, filename, file_size, mime_type, storage_path, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(
		query,
		doc.DocumentID,
		doc.UserID,
		doc.Filename,
		doc.FileSize,
		doc.MimeType,
		doc.StoragePath,
		doc.Status,
		doc.CreatedAt,
		doc.UpdatedAt,
	)

	if err != nil {
		infrastructure.LogError("failed to create document", err, map[string]interface{}{
			"document_id": doc.DocumentID,
			"user_id":     doc.UserID,
		})
		return domain.NewInternalError("Failed to create document")
	}

	infrastructure.LogInfo("document created", map[string]interface{}{
		"document_id": doc.DocumentID,
		"user_id":     doc.UserID,
		"filename":    doc.Filename,
	})

	return nil
}

// GetDocumentByID retrieves a document by its ID
func (r *DocumentRepository) GetDocumentByID(documentID string) (*domain.Document, error) {
	query := `
		SELECT document_id, user_id, filename, file_size, mime_type, storage_path, status, created_at, updated_at
		FROM documents
		WHERE document_id = $1
	`

	doc := &domain.Document{}
	err := r.db.QueryRow(query, documentID).Scan(
		&doc.DocumentID,
		&doc.UserID,
		&doc.Filename,
		&doc.FileSize,
		&doc.MimeType,
		&doc.StoragePath,
		&doc.Status,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.NewNotFoundError("Document not found")
	}

	if err != nil {
		infrastructure.LogError("failed to get document", err, map[string]interface{}{
			"document_id": documentID,
		})
		return nil, domain.NewInternalError("Failed to retrieve document")
	}

	return doc, nil
}

// GetDocumentsByUserID retrieves all documents for a user with pagination
func (r *DocumentRepository) GetDocumentsByUserID(userID string, limit, offset int) ([]*domain.Document, error) {
	query := `
		SELECT document_id, user_id, filename, file_size, mime_type, storage_path, status, created_at, updated_at
		FROM documents
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		infrastructure.LogError("failed to get documents by user", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, domain.NewInternalError("Failed to retrieve documents")
	}
	defer rows.Close()

	documents := []*domain.Document{}
	for rows.Next() {
		doc := &domain.Document{}
		err := rows.Scan(
			&doc.DocumentID,
			&doc.UserID,
			&doc.Filename,
			&doc.FileSize,
			&doc.MimeType,
			&doc.StoragePath,
			&doc.Status,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			infrastructure.LogError("failed to scan document row", err, map[string]interface{}{
				"user_id": userID,
			})
			continue
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

// CountDocumentsByUserID counts total documents for a user
func (r *DocumentRepository) CountDocumentsByUserID(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM documents WHERE user_id = $1`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		infrastructure.LogError("failed to count documents", err, map[string]interface{}{
			"user_id": userID,
		})
		return 0, domain.NewInternalError("Failed to count documents")
	}

	return count, nil
}

// UpdateDocumentStatus updates the status of a document
func (r *DocumentRepository) UpdateDocumentStatus(documentID, status string) error {
	query := `
		UPDATE documents
		SET status = $1, updated_at = NOW()
		WHERE document_id = $2
	`

	result, err := r.db.Exec(query, status, documentID)
	if err != nil {
		infrastructure.LogError("failed to update document status", err, map[string]interface{}{
			"document_id": documentID,
			"status":      status,
		})
		return domain.NewInternalError("Failed to update document status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("Failed to verify document update")
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("Document not found")
	}

	infrastructure.LogInfo("document status updated", map[string]interface{}{
		"document_id": documentID,
		"status":      status,
	})

	return nil
}

// DeleteDocument removes a document from the database
func (r *DocumentRepository) DeleteDocument(documentID string) error {
	query := `DELETE FROM documents WHERE document_id = $1`

	result, err := r.db.Exec(query, documentID)
	if err != nil {
		infrastructure.LogError("failed to delete document", err, map[string]interface{}{
			"document_id": documentID,
		})
		return domain.NewInternalError("Failed to delete document")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.NewInternalError("Failed to verify document deletion")
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("Document not found")
	}

	infrastructure.LogInfo("document deleted from database", map[string]interface{}{
		"document_id": documentID,
	})

	return nil
}
