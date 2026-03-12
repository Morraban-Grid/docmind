"""
DocMind Python RAG Service
Main application entry point
"""
import os
import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Create FastAPI application
app = FastAPI(
    title="DocMind RAG Service",
    description="Document processing and RAG query service",
    version="1.0.0"
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "python-rag-service",
        "version": "1.0.0"
    }

@app.on_event("startup")
async def startup_event():
    """Initialize services on startup"""
    logger.info("Starting DocMind Python RAG Service")
    logger.info(f"Service port: {os.getenv('PYTHON_SERVICE_PORT', '8000')}")
    logger.info(f"gRPC port: {os.getenv('PYTHON_SERVICE_GRPC_PORT', '50051')}")
    
    # TODO: Initialize Qdrant client
    # TODO: Initialize embedding model
    # TODO: Initialize Ollama client
    # TODO: Start gRPC server
    
    logger.info("Service initialization complete")

@app.on_event("shutdown")
async def shutdown_event():
    """Cleanup on shutdown"""
    logger.info("Shutting down DocMind Python RAG Service")
    
    # TODO: Close Qdrant connections
    # TODO: Release model resources
    # TODO: Stop gRPC server
    
    logger.info("Shutdown complete")

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PYTHON_SERVICE_PORT", "8000"))
    uvicorn.run(app, host="0.0.0.0", port=port)
