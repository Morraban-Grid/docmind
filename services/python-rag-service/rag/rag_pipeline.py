import logging
from typing import List, Dict, Any
from pydantic import BaseModel

from search.semantic_search import SemanticSearch, SearchResult
from llm.ollama_client import OllamaClient

logger = logging.getLogger(__name__)


class RAGResponse(BaseModel):
    """RAG response model"""
    answer: str
    sources: List[str]
    chunk_count: int


class RAGPipeline:
    """RAG pipeline: search + LLM generation"""
    
    def __init__(self, ollama_host: str = "localhost", ollama_port: int = 11434):
        """Initialize RAG pipeline"""
        self.search = SemanticSearch(similarity_threshold=0.7, top_k=5)
        self.llm = OllamaClient(host=ollama_host, port=ollama_port, model="llama2")
        logger.info("RAGPipeline initialized")
    
    def query(self, query: str, user_id: str) -> RAGResponse:
        """
        Execute RAG query
        
        Args:
            query: User query
            user_id: User ID for filtering
            
        Returns:
            RAGResponse with answer and sources
        """
        if not query or not query.strip():
            logger.warning("Empty query provided")
            raise ValueError("Query cannot be empty")
        
        try:
            logger.info(f"RAG query from user {user_id}: {query[:100]}...")
            
            # Step 1: Semantic search
            logger.info("Step 1: Performing semantic search")
            search_results = self.search.search(query, user_id)
            
            if not search_results:
                logger.warning(f"No relevant chunks found for query")
                return RAGResponse(
                    answer="No relevant information found in your documents.",
                    sources=[],
                    chunk_count=0
                )
            
            logger.info(f"Found {len(search_results)} relevant chunks")
            
            # Step 2: Construct prompt
            logger.info("Step 2: Constructing prompt")
            prompt = self._construct_prompt(query, search_results)
            
            # Step 3: Generate response
            logger.info("Step 3: Generating response with LLM")
            answer = self.llm.generate_response(prompt)
            
            # Step 4: Extract sources
            sources = list(set([result.document_id for result in search_results]))
            
            logger.info(f"RAG query completed: {len(sources)} sources, {len(answer)} chars")
            
            return RAGResponse(
                answer=answer,
                sources=sources,
                chunk_count=len(search_results)
            )
        except Exception as e:
            logger.error(f"RAG query failed: {e}")
            raise
    
    def _construct_prompt(self, query: str, search_results: List[SearchResult]) -> str:
        """
        Construct prompt for LLM
        
        Args:
            query: User query
            search_results: Search results with context
            
        Returns:
            Formatted prompt
        """
        # Limit to top 5 chunks
        chunks = search_results[:5]
        
        # Format context
        context_parts = []
        for i, result in enumerate(chunks, 1):
            context_parts.append(f"[Source {i}] {result.chunk_text}")
        
        context = "\n\n".join(context_parts)
        
        # Construct prompt
        prompt = f"""Based on the following context from documents, answer the question. If the answer is not in the context, say "I don't have enough information to answer this question."

Context:
{context}

Question: {query}

Answer:"""
        
        logger.debug(f"Constructed prompt with {len(context)} chars of context")
        return prompt
    
    def health_check(self) -> Dict[str, bool]:
        """Check health of RAG pipeline dependencies"""
        return {
            "ollama": self.llm.health_check()
        }
