import logging
from .pdf import PDFExtractor
from .docx import DOCXExtractor
from .text import TextExtractorPlain
from .markdown import MarkdownExtractor
from .base import TextExtractor

logger = logging.getLogger(__name__)


class ExtractorFactory:
    """Factory for creating appropriate text extractors"""
    
    _extractors = {
        'pdf': PDFExtractor,
        'docx': DOCXExtractor,
        'txt': TextExtractorPlain,
        'md': MarkdownExtractor,
    }
    
    @classmethod
    def get_extractor(cls, file_extension: str) -> TextExtractor:
        """Get extractor for file type"""
        ext = file_extension.lower().lstrip('.')
        
        if ext not in cls._extractors:
            raise ValueError(f"Unsupported file type: {ext}")
        
        logger.info(f"Creating extractor for {ext}")
        return cls._extractors[ext]()
    
    @classmethod
    def get_supported_extensions(cls) -> list:
        """Get list of supported file extensions"""
        return list(cls._extractors.keys())
