import logging
import os
import requests
from typing import Optional

logger = logging.getLogger(__name__)


class OllamaClient:
    """Client for Ollama LLM service"""
    
    def __init__(self, host: str = "localhost", port: int = 11434, model: str = "llama2"):
        """Initialize Ollama client"""
        self.host = host
        self.port = port
        self.model = model
        self.base_url = f"http://{host}:{port}"
        self.timeout = 120  # 2 minutes timeout for LLM generation
        
        logger.info(f"OllamaClient initialized: {self.base_url}, model={model}")
    
    def health_check(self) -> bool:
        """Check if Ollama service is available"""
        try:
            response = requests.get(
                f"{self.base_url}/api/tags",
                timeout=5
            )
            available = response.status_code == 200
            logger.debug(f"Ollama health check: {'available' if available else 'unavailable'}")
            return available
        except Exception as e:
            logger.error(f"Ollama health check failed: {e}")
            return False
    
    def generate_response(self, prompt: str, max_tokens: int = 500, temperature: float = 0.7) -> str:
        """
        Generate response using Ollama LLM
        
        Args:
            prompt: Prompt text
            max_tokens: Maximum tokens to generate
            temperature: Temperature for generation (0.0-1.0)
            
        Returns:
            Generated response text
        """
        if not prompt or not prompt.strip():
            logger.warning("Empty prompt provided for generation")
            return ""
        
        try:
            logger.info(f"Generating response with model {self.model}")
            
            # Prepare request
            payload = {
                "model": self.model,
                "prompt": prompt,
                "stream": False,
                "options": {
                    "temperature": temperature,
                    "num_predict": max_tokens
                }
            }
            
            # Call Ollama API
            response = requests.post(
                f"{self.base_url}/api/generate",
                json=payload,
                timeout=self.timeout
            )
            
            if response.status_code != 200:
                logger.error(f"Ollama API error: {response.status_code}")
                raise ValueError(f"Ollama API error: {response.status_code}")
            
            # Extract response
            result = response.json()
            generated_text = result.get("response", "")
            
            logger.info(f"Generated response with {len(generated_text)} characters")
            return generated_text
        except requests.Timeout:
            logger.error("Ollama request timeout")
            raise ValueError("LLM generation timeout")
        except Exception as e:
            logger.error(f"LLM generation failed: {e}")
            raise ValueError(f"LLM generation failed: {str(e)}")
