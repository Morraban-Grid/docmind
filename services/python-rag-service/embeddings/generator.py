import logging
import numpy as np
from typing import List
from sentence_transformers import SentenceTransformer
from pydantic import BaseModel

logger = logging.getLogger(__name__)


class Embedding(BaseModel):
    """Represents a text embedding"""
    text: str
    vector: List[float]
    dimension: int


class EmbeddingGenerator:
    """Generate embeddings using sentence-transformers"""
    
    _model = None
    _model_name = "all-MiniLM-L6-v2"
    _embedding_dim = 384
    
    @classmethod
    def _load_model(cls):
        """Load model lazily on first use"""
        if cls._model is None:
            logger.info(f"Loading embedding model: {cls._model_name}")
            cls._model = SentenceTransformer(cls._model_name)
            logger.info(f"Model loaded successfully. Dimension: {cls._embedding_dim}")
        return cls._model
    
    @classmethod
    def generate_embedding(cls, text: str) -> Embedding:
        """Generate embedding for a single text"""
        if not text or not text.strip():
            logger.warning("Empty text provided for embedding")
            # Return zero vector for empty text
            return Embedding(
                text=text,
                vector=[0.0] * cls._embedding_dim,
                dimension=cls._embedding_dim
            )
        
        try:
            model = cls._load_model()
            # Generate embedding
            embedding = model.encode(text, convert_to_numpy=True)
            
            # Normalize to unit length (L2 norm = 1.0)
            norm = np.linalg.norm(embedding)
            if norm > 0:
                embedding = embedding / norm
            
            # Convert to list
            vector = embedding.tolist()
            
            logger.debug(f"Generated embedding for text of length {len(text)}")
            
            return Embedding(
                text=text,
                vector=vector,
                dimension=len(vector)
            )
        except Exception as e:
            logger.error(f"Embedding generation failed: {e}")
            raise ValueError(f"Failed to generate embedding: {str(e)}")
    
    @classmethod
    def generate_embeddings_batch(cls, texts: List[str], batch_size: int = 32) -> List[Embedding]:
        """Generate embeddings for multiple texts in batches"""
        if not texts:
            logger.warning("Empty text list provided for batch embedding")
            return []
        
        try:
            model = cls._load_model()
            embeddings = []
            
            # Process in batches
            for i in range(0, len(texts), batch_size):
                batch = texts[i:i + batch_size]
                logger.debug(f"Processing batch {i // batch_size + 1} with {len(batch)} texts")
                
                # Generate embeddings for batch
                batch_embeddings = model.encode(batch, convert_to_numpy=True)
                
                # Normalize each embedding
                for j, embedding in enumerate(batch_embeddings):
                    norm = np.linalg.norm(embedding)
                    if norm > 0:
                        embedding = embedding / norm
                    
                    embeddings.append(Embedding(
                        text=batch[j],
                        vector=embedding.tolist(),
                        dimension=len(embedding)
                    ))
            
            logger.info(f"Generated {len(embeddings)} embeddings in {(len(texts) - 1) // batch_size + 1} batches")
            return embeddings
        except Exception as e:
            logger.error(f"Batch embedding generation failed: {e}")
            raise ValueError(f"Failed to generate batch embeddings: {str(e)}")
    
    @classmethod
    def get_embedding_dimension(cls) -> int:
        """Get embedding dimension"""
        return cls._embedding_dim
    
    @classmethod
    def get_model_name(cls) -> str:
        """Get model name"""
        return cls._model_name
