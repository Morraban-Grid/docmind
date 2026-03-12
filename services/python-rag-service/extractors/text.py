import logging
from .base import TextExtractor

logger = logging.getLogger(__name__)


class TextExtractorPlain(TextExtractor):
    """Extract text from plain text files"""
    
    def extract(self, file_path: str) -> str:
        """Extract text from TXT file"""
        if not self.validate_file(file_path):
            raise ValueError(f"Invalid text file: {file_path}")
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                text = f.read()
            
            logger.info(f"Extracted {len(text)} characters from TXT")
            return text
        except UnicodeDecodeError:
            logger.error("Failed to decode text file as UTF-8")
            raise ValueError("File encoding not supported (UTF-8 required)")
        except Exception as e:
            logger.error(f"Text extraction failed: {e}")
            raise ValueError(f"Failed to extract text: {str(e)}")
