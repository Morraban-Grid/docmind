import logging
import asyncio
from concurrent import futures

import grpc

from grpc_server.rag_service import RAGServicer
from config import settings

logger = logging.getLogger(__name__)


class RAGServiceGrpcServer:
    """gRPC server for RAG service"""
    
    def __init__(self, host: str = "0.0.0.0", port: int = 50051):
        """Initialize gRPC server"""
        self.host = host
        self.port = port
        self.server = None
        self.servicer = RAGServicer()
        logger.info(f"RAGServiceGrpcServer initialized on {host}:{port}")
    
    async def start(self):
        """Start gRPC server"""
        try:
            # Create server
            self.server = grpc.aio.server(
                futures.ThreadPoolExecutor(max_workers=10)
            )
            
            # Add servicer
            # Note: This will be properly implemented after proto generation
            # For now, we're setting up the structure
            
            # Add port
            self.server.add_insecure_port(f"{self.host}:{self.port}")
            
            # Start server
            await self.server.start()
            logger.info(f"gRPC server started on {self.host}:{self.port}")
            
            # Keep server running
            await self.server.wait_for_termination()
        except Exception as e:
            logger.error(f"Failed to start gRPC server: {e}")
            raise
    
    async def stop(self):
        """Stop gRPC server"""
        if self.server:
            logger.info("Stopping gRPC server")
            await self.server.stop(grace=5)
            logger.info("gRPC server stopped")


async def run_grpc_server():
    """Run gRPC server"""
    server = RAGServiceGrpcServer(
        host=settings.SERVER_HOST,
        port=50051
    )
    await server.start()


if __name__ == "__main__":
    asyncio.run(run_grpc_server())
