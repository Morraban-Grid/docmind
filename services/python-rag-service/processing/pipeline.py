import logging
import os
import tempfile
from typing import List
from pydantic import BaseModel

from extractors.factory import ExtractorFactory
from chunking.chunker import DocumentChunker, Chunk

logger = logging.getLogger(__name__)


class ProcessingRequest(BaseModel):
    """Request to process a document"""
    document_id: str
    file_content: bytes
    file_extension: str
    user_id: str


class ProcessingResponse(BaseModel):
    """Response from document processing"""
    document_id: str
    chunks: List[Chunk]
    total_chunks: int
    status: str


class DocumentProcessingPipeline:
    """Orchestrate document processing pipeline"""
    
    def __init__(self, chunk_size: int = 1000, chunk_overlap: int = 200):
        """Initialize pipeline"""
        self.chunker = DocumentChunker(chunk_size, chunk_overlap)
        self.extractor_factory = ExtractorFactory()
        logger.info("DocumentProcessingPipeline initialized")
    
    async def process(self, request: ProcessingRequest) -> ProcessingResponse:
        """Process document through pipeline"""
        temp_file = None
        try:
            # Validate file extension
            if not request.file_extension.lower().lstrip('.') in self.extractor_factory.get_supported_extensions():
                raise ValueError(f"Unsupported file type: {request.file_extension}")
            
            # Write file to temporary location
            temp_file = self._write_temp_file(request.file_content)
            logger.info(f"Processing document {request.document_id} for user {request.user_id}")
            
            # Extract text
            extractor = self.extractor_factory.get_extractor(request.file_extension)
            text = extractor.extract(temp_file)
            
            # Chunk text
            chunks = self.chunker.chunk(text)
            
            logger.info(f"Successfully processed document {request.document_id} into {len(chunks)} chunks")
            
            return ProcessingResponse(
                document_id=request.document_id,
                chunks=chunks,
                total_chunks=len(chunks),
                status="success"
            )
        except Exception as e:
            logger.error(f"Pipeline processing failed: {e}")
            raise
        finally:
            # Cleanup temporary file
            if temp_file and os.path.exists(temp_file):
                try:
                    os.remove(temp_file)
                    logger.debug(f"Cleaned up temporary file: {temp_file}")
                except Exception as e:
                    logger.warning(f"Failed to cleanup temp file: {e}")
    
    def _write_temp_file(self, content: bytes) -> str:
        """Write content to temporary file"""
        with tempfile.NamedTemporaryFile(delete=False) as tmp:
            tmp.write(content)
            return tmp.name
