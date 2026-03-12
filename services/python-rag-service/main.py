import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from config import settings, logger
from routes import health, processing


def create_app() -> FastAPI:
    """Create and configure FastAPI application"""
    
    app = FastAPI(
        title="DocMind RAG Service",
        description="Document processing and chunking service",
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
    
    # Include routers
    app.include_router(health.router, tags=["health"])
    app.include_router(processing.router, prefix="/api", tags=["processing"])
    
    logger.info("DocMind RAG Service initialized", extra={
        "host": settings.SERVER_HOST,
        "port": settings.SERVER_PORT,
        "debug": settings.DEBUG
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
