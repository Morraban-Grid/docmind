from abc import ABC, abstractmethod
from typing import Optional
import logging

logger = logging.getLogger(__name__)


class TextExtractor(ABC):
    """Base class for text extractors"""
    
    @abstractmethod
    def extract(self, file_path: str) -> str:
        """Extract text from file"""
        pass
    
    def validate_file(self, file_path: str) -> bool:
        """Validate file exists and is readable"""
        try:
            with open(file_path, 'rb') as f:
                f.read(1)
            return True
        except Exception as e:
            logger.error(f"File validation failed: {e}")
            return False
