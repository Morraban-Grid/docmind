import logging
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from indexing.indexer import DocumentIndexer, IndexingRequest, IndexingResponse
from config import settings

logger = logging.getLogger(__name__)
router = APIRouter()

# Initialize indexer
indexer = DocumentIndexer(
    chunk_size=settings.CHUNK_SIZE,
    chunk_overlap=settings.CHUNK_OVERLAP
)


class IndexResponse(BaseModel):
    """Response model for indexing endpoint"""
    document_id: str
    chunk_count: int
    embedding_count: int
    status: str


@router.post("/index", response_model=IndexResponse)
async def index_document(
    document_id: str,
    user_id: str,
    file_path: str,
    file_type: str
):
    """Index document for embeddings and vector storage"""
    try:
        # Validate inputs
        if not document_id or not document_id.strip():
            raise HTTPException(status_code=400, detail="document_id is required")
        if not user_id or not user_id.strip():
            raise HTTPException(status_code=400, detail="user_id is required")
        if not file_path or not file_path.strip():
            raise HTTPException(status_code=400, detail="file_path is required")
        if not file_type or not file_type.strip():
            raise HTTPException(status_code=400, detail="file_type is required")
        
        # Create indexing request
        request = IndexingRequest(
            document_id=document_id,
            user_id=user_id,
            file_path=file_path,
            file_type=file_type
        )
        
        # Index document
        result = await indexer.index(request)
        
        logger.info(f"Document {document_id} indexed successfully")
        
        return IndexResponse(
            document_id=result.document_id,
            chunk_count=result.chunk_count,
            embedding_count=result.embedding_count,
            status=result.status
        )
    except ValueError as e:
        logger.error(f"Validation error: {e}")
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Indexing error: {e}")
        raise HTTPException(status_code=500, detail="Document indexing failed")
