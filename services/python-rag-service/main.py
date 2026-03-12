import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from config import settings, logger
from routes import health, processing, indexing
from vector_store.qdrant_client import QdrantClient


def create_app() -> FastAPI:
    """Create and configure FastAPI application"""
    
    app = FastAPI(
        title="DocMind RAG Service",
        description="Document processing, chunking, and embeddings service",
        version="1.0.0"
    )
    
    # Add CORS middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    
    # Initialize Qdrant collection on startup
    @app.on_event("startup")
    async def startup_event():
        logger.info("Starting DocMind RAG Service")
        try:
            QdrantClient.create_collection()
            logger.info("Qdrant collection initialized")
        except Exception as e:
            logger.warning(f"Failed to initialize Qdrant collection: {e}")
    
    @app.on_event("shutdown")
    async def shutdown_event():
        logger.info("Shutting down DocMind RAG Service")
    
    # Include routers
    app.include_router(health.router, tags=["health"])
    app.include_router(processing.router, prefix="/api", tags=["processing"])
    app.include_router(indexing.router, prefix="/api", tags=["indexing"])
    
    logger.info("DocMind RAG Service initialized", extra={
        "host": settings.SERVER_HOST,
        "port": settings.SERVER_PORT,
        "debug": settings.DEBUG,
        "qdrant_host": settings.QDRANT_HOST,
        "qdrant_port": settings.QDRANT_PORT
    })
    
    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        app,
        host=settings.SERVER_HOST,
        port=settings.SERVER_PORT,
        log_level=settings.LOG_LEVEL.lower()
    )
