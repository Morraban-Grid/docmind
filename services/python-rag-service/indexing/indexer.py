import logging
import uuid
from typing import List
from pydantic import BaseModel

from extractors.factory import ExtractorFactory
from chunking.chunker import DocumentChunker
from embeddings.generator import EmbeddingGenerator
from vector_store.qdrant_client import QdrantClient

logger = logging.getLogger(__name__)


class IndexingRequest(BaseModel):
    """Request to index a document"""
    document_id: str
    user_id: str
    file_path: str
    file_type: str


class IndexingResponse(BaseModel):
    """Response from document indexing"""
    document_id: str
    chunk_count: int
    embedding_count: int
    status: str


class DocumentIndexer:
    """Index documents: extract → chunk → embed → store"""
    
    def __init__(self, chunk_size: int = 1000, chunk_overlap: int = 200):
        """Initialize indexer"""
        self.chunker = DocumentChunker(chunk_size, chunk_overlap)
        self.extractor_factory = ExtractorFactory()
        logger.info("DocumentIndexer initialized")
    
    async def index(self, request: IndexingRequest) -> IndexingResponse:
        """Index document through complete pipeline"""
        try:
            logger.info(f"Starting indexing for document {request.document_id}")
            
            # Validate file type
            if not request.file_type.lower().lstrip('.') in self.extractor_factory.get_supported_extensions():
                raise ValueError(f"Unsupported file type: {request.file_type}")
            
            # Extract text
            logger.info(f"Extracting text from {request.file_type}")
            extractor = self.extractor_factory.get_extractor(request.file_type)
            text = extractor.extract(request.file_path)
            logger.info(f"Extracted {len(text)} characters")
            
            # Chunk text
            logger.info("Chunking text")
            chunks = self.chunker.chunk(text)
            logger.info(f"Created {len(chunks)} chunks")
            
            # Generate embeddings
            logger.info("Generating embeddings")
            chunk_texts = [chunk.content for chunk in chunks]
            embeddings = EmbeddingGenerator.generate_embeddings_batch(chunk_texts)
            logger.info(f"Generated {len(embeddings)} embeddings")
            
            # Store in Qdrant
            logger.info("Storing embeddings in Qdrant")
            points = []
            for i, embedding in enumerate(embeddings):
                points.append({
                    "id": str(uuid.uuid4()),
                    "vector": embedding.vector,
                    "payload": {
                        "document_id": request.document_id,
                        "user_id": request.user_id,
                        "chunk_text": embedding.text,
                        "chunk_index": i
                    }
                })
            
            embedding_count = QdrantClient.upsert_embeddings(points)
            logger.info(f"Stored {embedding_count} embeddings in Qdrant")
            
            logger.info(f"Successfully indexed document {request.document_id}")
            
            return IndexingResponse(
                document_id=request.document_id,
                chunk_count=len(chunks),
                embedding_count=embedding_count,
                status="success"
            )
        except Exception as e:
            logger.error(f"Indexing failed: {e}")
            raise
