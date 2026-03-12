import logging
from fastapi import APIRouter, UploadFile, File, HTTPException
from pydantic import BaseModel

from processing.pipeline import DocumentProcessingPipeline, ProcessingRequest, ProcessingResponse
from config import settings

logger = logging.getLogger(__name__)
router = APIRouter()

# Initialize pipeline
pipeline = DocumentProcessingPipeline(
    chunk_size=settings.CHUNK_SIZE,
    chunk_overlap=settings.CHUNK_OVERLAP
)


class ProcessResponse(BaseModel):
    """Response model for processing endpoint"""
    document_id: str
    total_chunks: int
    status: str


@router.post("/process", response_model=ProcessResponse)
async def process_document(
    document_id: str,
    user_id: str,
    file: UploadFile = File(...)
):
    """Process document for chunking"""
    try:
        # Validate file size
        content = await file.read()
        if len(content) > settings.MAX_FILE_SIZE:
            raise HTTPException(
                status_code=400,
                detail=f"File too large. Maximum size: {settings.MAX_FILE_SIZE} bytes"
            )
        
        # Get file extension
        file_ext = file.filename.split('.')[-1] if file.filename else ""
        
        # Create processing request
        request = ProcessingRequest(
            document_id=document_id,
            file_content=content,
            file_extension=file_ext,
            user_id=user_id
        )
        
        # Process document
        result = await pipeline.process(request)
        
        logger.info(f"Document {document_id} processed successfully")
        
        return ProcessResponse(
            document_id=result.document_id,
            total_chunks=result.total_chunks,
            status=result.status
        )
    except ValueError as e:
        logger.error(f"Validation error: {e}")
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Processing error: {e}")
        raise HTTPException(status_code=500, detail="Document processing failed")
