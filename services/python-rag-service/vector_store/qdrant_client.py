import logging
import os
from typing import List, Optional, Dict, Any
from qdrant_client import QdrantClient as QdrantClientLib
from qdrant_client.models import Distance, VectorParams, PointStruct, Filter, FieldCondition, MatchValue
from pydantic import BaseModel

logger = logging.getLogger(__name__)


class VectorPayload(BaseModel):
    """Payload for vector storage"""
    document_id: str
    user_id: str
    chunk_text: str
    chunk_index: int


class QdrantClient:
    """Wrapper for Qdrant vector database client"""
    
    _client = None
    _collection_name = "docmind_chunks"
    _embedding_dim = 384
    
    @classmethod
    def _get_client(cls) -> QdrantClientLib:
        """Get or create Qdrant client"""
        if cls._client is None:
            host = os.getenv("QDRANT_HOST", "localhost")
            port = int(os.getenv("QDRANT_PORT", "6333"))
            
            logger.info(f"Connecting to Qdrant at {host}:{port}")
            cls._client = QdrantClientLib(host=host, port=port)
            logger.info("Connected to Qdrant successfully")
        
        return cls._client
    
    @classmethod
    def create_collection(cls) -> bool:
        """Create collection if it doesn't exist"""
        try:
            client = cls._get_client()
            
            # Check if collection exists
            if cls.collection_exists():
                logger.info(f"Collection '{cls._collection_name}' already exists")
                return True
            
            logger.info(f"Creating collection '{cls._collection_name}'")
            client.create_collection(
                collection_name=cls._collection_name,
                vectors_config=VectorParams(
                    size=cls._embedding_dim,
                    distance=Distance.COSINE
                )
            )
            logger.info(f"Collection '{cls._collection_name}' created successfully")
            return True
        except Exception as e:
            logger.error(f"Failed to create collection: {e}")
            raise ValueError(f"Failed to create collection: {str(e)}")
    
    @classmethod
    def collection_exists(cls) -> bool:
        """Check if collection exists"""
        try:
            client = cls._get_client()
            collections = client.get_collections()
            
            for collection in collections.collections:
                if collection.name == cls._collection_name:
                    logger.debug(f"Collection '{cls._collection_name}' exists")
                    return True
            
            logger.debug(f"Collection '{cls._collection_name}' does not exist")
            return False
        except Exception as e:
            logger.error(f"Failed to check collection existence: {e}")
            return False
    
    @classmethod
    def upsert_embeddings(cls, points: List[Dict[str, Any]]) -> int:
        """
        Upsert embeddings to Qdrant
        
        Points format:
        {
            "id": "chunk_id",
            "vector": [0.1, 0.2, ...],
            "payload": {
                "document_id": "doc_id",
                "user_id": "user_id",
                "chunk_text": "text",
                "chunk_index": 0
            }
        }
        """
        if not points:
            logger.warning("Empty points list provided for upsert")
            return 0
        
        try:
            client = cls._get_client()
            
            # Ensure collection exists
            if not cls.collection_exists():
                cls.create_collection()
            
            # Convert to PointStruct objects
            point_structs = []
            for point in points:
                point_structs.append(
                    PointStruct(
                        id=hash(point["id"]) & 0x7fffffff,  # Convert string ID to positive int
                        vector=point["vector"],
                        payload=point["payload"]
                    )
                )
            
            # Upsert points
            client.upsert(
                collection_name=cls._collection_name,
                points=point_structs
            )
            
            logger.info(f"Upserted {len(point_structs)} embeddings to Qdrant")
            return len(point_structs)
        except Exception as e:
            logger.error(f"Failed to upsert embeddings: {e}")
            raise ValueError(f"Failed to upsert embeddings: {str(e)}")
    
    @classmethod
    def search(cls, query_vector: List[float], user_id: str, limit: int = 5, score_threshold: float = 0.7) -> List[Dict[str, Any]]:
        """
        Search for similar embeddings
        
        Returns list of results with:
        {
            "id": "chunk_id",
            "score": 0.95,
            "payload": {...}
        }
        """
        try:
            client = cls._get_client()
            
            if not cls.collection_exists():
                logger.warning("Collection does not exist, returning empty results")
                return []
            
            # Search with user_id filter
            results = client.search(
                collection_name=cls._collection_name,
                query_vector=query_vector,
                query_filter=Filter(
                    must=[
                        FieldCondition(
                            key="payload.user_id",
                            match=MatchValue(value=user_id)
                        )
                    ]
                ),
                limit=limit,
                score_threshold=score_threshold
            )
            
            # Convert to dict format
            search_results = []
            for result in results:
                search_results.append({
                    "id": str(result.id),
                    "score": result.score,
                    "payload": result.payload
                })
            
            logger.info(f"Search returned {len(search_results)} results for user {user_id}")
            return search_results
        except Exception as e:
            logger.error(f"Search failed: {e}")
            raise ValueError(f"Search failed: {str(e)}")
    
    @classmethod
    def delete_by_document_id(cls, document_id: str) -> int:
        """Delete all embeddings for a document"""
        try:
            client = cls._get_client()
            
            if not cls.collection_exists():
                logger.warning("Collection does not exist, nothing to delete")
                return 0
            
            # Delete with filter
            result = client.delete(
                collection_name=cls._collection_name,
                points_selector=Filter(
                    must=[
                        FieldCondition(
                            key="payload.document_id",
                            match=MatchValue(value=document_id)
                        )
                    ]
                )
            )
            
            logger.info(f"Deleted embeddings for document {document_id}")
            return 1  # Qdrant doesn't return count, so we return 1 to indicate success
        except Exception as e:
            logger.error(f"Failed to delete embeddings: {e}")
            raise ValueError(f"Failed to delete embeddings: {str(e)}")
    
    @classmethod
    def health_check(cls) -> bool:
        """Check Qdrant connection health"""
        try:
            client = cls._get_client()
            client.get_collections()
            logger.debug("Qdrant health check passed")
            return True
        except Exception as e:
            logger.error(f"Qdrant health check failed: {e}")
            return False
