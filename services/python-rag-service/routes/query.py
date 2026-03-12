import logging
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from rag.rag_pipeline import RAGPipeline, RAGResponse
from config import settings

logger = logging.getLogger(__name__)
router = APIRouter()

# Initialize RAG pipeline
rag_pipeline = RAGPipeline(
    ollama_host=settings.OLLAMA_HOST,
    ollama_port=settings.OLLAMA_PORT
)


class QueryRequest(BaseModel):
    """Query request model"""
    query: str


class QueryResponseModel(BaseModel):
    """Query response model"""
    answer: str
    sources: list
    chunk_count: int


@router.post("/query", response_model=QueryResponseModel)
async def query_documents(query_request: QueryRequest, user_id: str):
    """Query documents using RAG"""
    try:
        # Validate query
        if not query_request.query or not query_request.query.strip():
            raise HTTPException(status_code=400, detail="Query cannot be empty")
        
        if not user_id or not user_id.strip():
            raise HTTPException(status_code=400, detail="user_id is required")
        
        logger.info(f"Query from user {user_id}: {query_request.query[:100]}...")
        
        # Execute RAG query
        result = rag_pipeline.query(query_request.query, user_id)
        
        logger.info(f"Query completed successfully")
        
        return QueryResponseModel(
            answer=result.answer,
            sources=result.sources,
            chunk_count=result.chunk_count
        )
    except ValueError as e:
        logger.error(f"Validation error: {e}")
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Query error: {e}")
        raise HTTPException(status_code=500, detail="Query processing failed")
