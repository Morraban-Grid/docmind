import logging
from typing import List
from langchain_text_splitters import RecursiveCharacterTextSplitter
from pydantic import BaseModel

logger = logging.getLogger(__name__)


class Chunk(BaseModel):
    """Represents a text chunk"""
    content: str
    chunk_index: int
    total_chunks: int


class DocumentChunker:
    """Chunk documents using LangChain"""
    
    def __init__(self, chunk_size: int = 1000, chunk_overlap: int = 200):
        """Initialize chunker with parameters"""
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap
        self.splitter = RecursiveCharacterTextSplitter(
            chunk_size=chunk_size,
            chunk_overlap=chunk_overlap,
            separators=["\n\n", "\n", " ", ""]
        )
        logger.info(f"DocumentChunker initialized with size={chunk_size}, overlap={chunk_overlap}")
    
    def chunk(self, text: str) -> List[Chunk]:
        """Split text into chunks"""
        if not text or not text.strip():
            logger.warning("Empty text provided for chunking")
            return []
        
        try:
            chunks_text = self.splitter.split_text(text)
            chunks = [
                Chunk(
                    content=chunk,
                    chunk_index=i,
                    total_chunks=len(chunks_text)
                )
                for i, chunk in enumerate(chunks_text)
            ]
            logger.info(f"Created {len(chunks)} chunks from text")
            return chunks
        except Exception as e:
            logger.error(f"Chunking failed: {e}")
            raise ValueError(f"Failed to chunk document: {str(e)}")
