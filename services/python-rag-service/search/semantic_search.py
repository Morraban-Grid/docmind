import logging
from typing import List, Dict, Any

from embeddings.generator import EmbeddingGenerator
from vector_store.qdrant_client import QdrantClient

logger = logging.getLogger(__name__)


class SearchResult:
    """Represents a search result"""
    
    def __init__(self, document_id: str, chunk_text: str, similarity: float, chunk_index: int):
        self.document_id = document_id
        self.chunk_text = chunk_text
        self.similarity = similarity
        self.chunk_index = chunk_index
    
    def to_dict(self) -> Dict[str, Any]:
        return {
            "document_id": self.document_id,
            "chunk_text": self.chunk_text,
            "similarity": self.similarity,
            "chunk_index": self.chunk_index
        }


class SemanticSearch:
    """Semantic search using embeddings and Qdrant"""
    
    def __init__(self, similarity_threshold: float = 0.7, top_k: int = 5):
        """Initialize semantic search"""
        self.similarity_threshold = similarity_threshold
        self.top_k = top_k
        logger.info(f"SemanticSearch initialized with threshold={similarity_threshold}, top_k={top_k}")
    
    def search(self, query: str, user_id: str) -> List[SearchResult]:
        """
        Search for relevant chunks using semantic similarity
        
        Args:
            query: Query text
            user_id: User ID for filtering results
            
        Returns:
            List of SearchResult objects
        """
        if not query or not query.strip():
            logger.warning("Empty query provided for search")
            return []
        
        try:
            logger.info(f"Searching for query: {query[:100]}... for user {user_id}")
            
            # Generate embedding for query
            query_embedding = EmbeddingGenerator.generate_embedding(query)
            logger.debug(f"Generated query embedding with dimension {query_embedding.dimension}")
            
            # Search in Qdrant
            search_results = QdrantClient.search(
                query_vector=query_embedding.vector,
                user_id=user_id,
                limit=self.top_k,
                score_threshold=self.similarity_threshold
            )
            
            # Convert to SearchResult objects
            results = []
            for result in search_results:
                search_result = SearchResult(
                    document_id=result["payload"]["document_id"],
                    chunk_text=result["payload"]["chunk_text"],
                    similarity=result["score"],
                    chunk_index=result["payload"]["chunk_index"]
                )
                results.append(search_result)
            
            logger.info(f"Search returned {len(results)} results for user {user_id}")
            return results
        except Exception as e:
            logger.error(f"Search failed: {e}")
            raise ValueError(f"Search failed: {str(e)}")
