import logging
import asyncio
from typing import Optional

from indexing.indexer import DocumentIndexer, IndexingRequest
from vector_store.qdrant_client import QdrantClient
from config import settings

logger = logging.getLogger(__name__)


class RAGServicer:
    """gRPC service implementation for RAG operations"""
    
    def __init__(self):
        """Initialize RAG service"""
        self.indexer = DocumentIndexer(
            chunk_size=settings.CHUNK_SIZE,
            chunk_overlap=settings.CHUNK_OVERLAP
        )
        logger.info("RAGServicer initialized")
    
    async def index_document(self, request) -> dict:
        """
        Index a document for RAG
        
        Args:
            request: IndexDocumentRequest with document_id, user_id, file_path, file_type
            
        Returns:
            dict with success, chunk_count, embedding_count, error_message
        """
        try:
            logger.info(f"Indexing document {request.document_id} for user {request.user_id}")
            
            # Create indexing request
            indexing_request = IndexingRequest(
                document_id=request.document_id,
                user_id=request.user_id,
                file_path=request.file_path,
                file_type=request.file_type
            )
            
            # Index document
            result = await self.indexer.index(indexing_request)
            
            logger.info(f"Document {request.document_id} indexed successfully")
            
            return {
                "success": True,
                "chunk_count": result.chunk_count,
                "embedding_count": result.embedding_count,
                "error_message": ""
            }
        except Exception as e:
            logger.error(f"Indexing failed: {e}")
            return {
                "success": False,
                "chunk_count": 0,
                "embedding_count": 0,
                "error_message": str(e)
            }
    
    async def delete_document(self, request) -> dict:
        """
        Delete document embeddings from Qdrant
        
        Args:
            request: DeleteDocumentRequest with document_id
            
        Returns:
            dict with success, deleted_count, error_message
        """
        try:
            logger.info(f"Deleting embeddings for document {request.document_id}")
            
            # Delete from Qdrant
            deleted_count = QdrantClient.delete_by_document_id(request.document_id)
            
            logger.info(f"Deleted embeddings for document {request.document_id}")
            
            return {
                "success": True,
                "deleted_count": deleted_count,
                "error_message": ""
            }
        except Exception as e:
            logger.error(f"Deletion failed: {e}")
            return {
                "success": False,
                "deleted_count": 0,
                "error_message": str(e)
            }
    
    async def query_document(self, request) -> dict:
        """
        Query documents using RAG (placeholder for Iteration 7)
        
        Args:
            request: QueryDocumentRequest with query, user_id
            
        Returns:
            dict with answer, sources, chunk_count, error_message
        """
        try:
            logger.info(f"Query received from user {request.user_id}: {request.query}")
            
            # Placeholder - will be implemented in Iteration 7
            return {
                "answer": "RAG query not yet implemented",
                "sources": [],
                "chunk_count": 0,
                "error_message": "RAG query functionality coming in Iteration 7"
            }
        except Exception as e:
            logger.error(f"Query failed: {e}")
            return {
                "answer": "",
                "sources": [],
                "chunk_count": 0,
                "error_message": str(e)
            }
