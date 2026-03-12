import logging
from fastapi import APIRouter
from pydantic import BaseModel

from vector_store.qdrant_client import QdrantClient
from llm.ollama_client import OllamaClient

logger = logging.getLogger(__name__)
router = APIRouter()

# Initialize Ollama client for health checks
ollama_client = OllamaClient()


class DependencyStatus(BaseModel):
    """Status of a dependency"""
    name: str
    status: str
    available: bool


class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    service: str
    dependencies: list[DependencyStatus]


@router.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint with dependency status"""
    
    # Check Qdrant
    qdrant_available = QdrantClient.health_check()
    qdrant_status = "healthy" if qdrant_available else "unhealthy"
    
    # Check Ollama
    ollama_available = ollama_client.health_check()
    ollama_status = "healthy" if ollama_available else "unhealthy"
    
    # Overall status
    all_healthy = qdrant_available and ollama_available
    overall_status = "healthy" if all_healthy else "degraded"
    
    logger.debug(f"Health check: {overall_status}")
    
    return {
        "status": overall_status,
        "service": "python-rag-service",
        "dependencies": [
            {
                "name": "qdrant",
                "status": qdrant_status,
                "available": qdrant_available
            },
            {
                "name": "ollama",
                "status": ollama_status,
                "available": ollama_available
            }
        ]
    }
