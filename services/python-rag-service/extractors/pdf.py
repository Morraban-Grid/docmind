import logging
from pypdf import PdfReader
from .base import TextExtractor

logger = logging.getLogger(__name__)


class PDFExtractor(TextExtractor):
    """Extract text from PDF files"""
    
    def extract(self, file_path: str) -> str:
        """Extract text from PDF"""
        if not self.validate_file(file_path):
            raise ValueError(f"Invalid PDF file: {file_path}")
        
        try:
            text = []
            with open(file_path, 'rb') as f:
                reader = PdfReader(f)
                for page_num, page in enumerate(reader.pages):
                    page_text = page.extract_text()
                    if page_text:
                        text.append(page_text)
            
            result = "\n".join(text)
            logger.info(f"Extracted {len(result)} characters from PDF")
            return result
        except Exception as e:
            logger.error(f"PDF extraction failed: {e}")
            raise ValueError(f"Failed to extract PDF: {str(e)}")
