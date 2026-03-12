"""
Configuration management for Python RAG Service
Loads and validates environment variables
"""
import os
from typing import Optional
from pydantic_settings import BaseSettings
from pydantic import Field, validator


class Settings(BaseSettings):
    """Application settings loaded from environment variables"""
    
    # Server configuration
    python_service_port: int = Field(default=8000, env="PYTHON_SERVICE_PORT")
    python_service_grpc_port: int = Field(default=50051, env="PYTHON_SERVICE_GRPC_PORT")
    
    # Qdrant configuration
    qdrant_host: str = Field(default="qdrant", env="QDRANT_HOST")
    qdrant_port: int = Field(default=6333, env="QDRANT_PORT")
    qdrant_url: Optional[str] = Field(default=None, env="QDRANT_URL")
    
    # Ollama configuration
    ollama_host: str = Field(default="ollama", env="OLLAMA_HOST")
    ollama_port: int = Field(default=11434, env="OLLAMA_PORT")
    ollama_url: Optional[str] = Field(default=None, env="OLLAMA_URL")
    ollama_model: str = Field(default="llama2", env="OLLAMA_MODEL")
    
    # Embedding configuration
    embedding_model: str = Field(
        default="sentence-transformers/all-MiniLM-L6-v2",
        env="EMBEDDING_MODEL"
    )
    embedding_dimension: int = Field(default=384, env="EMBEDDING_DIMENSION")
    
    # Application configuration
    app_env: str = Field(default="development", env="APP_ENV")
    log_level: str = Field(default="INFO", env="LOG_LEVEL")
    
    # Processing configuration
    chunk_size: int = Field(default=1000, env="CHUNK_SIZE")
    chunk_overlap: int = Field(default=200, env="CHUNK_OVERLAP")
    batch_size: int = Field(default=32, env="BATCH_SIZE")
    
    # Search configuration
    top_k: int = Field(default=5, env="TOP_K")
    similarity_threshold: float = Field(default=0.7, env="SIMILARITY_THRESHOLD")
    
    @validator("qdrant_url", pre=True, always=True)
    def build_qdrant_url(cls, v, values):
        """Build Qdrant URL if not provided"""
        if v:
            return v
        host = values.get("qdrant_host", "qdrant")
        port = values.get("qdrant_port", 6333)
        return f"http://{host}:{port}"
    
    @validator("ollama_url", pre=True, always=True)
    def build_ollama_url(cls, v, values):
        """Build Ollama URL if not provided"""
        if v:
            return v
        host = values.get("ollama_host", "ollama")
        port = values.get("ollama_port", 11434)
        return f"http://{host}:{port}"
    
    @validator("log_level")
    def validate_log_level(cls, v):
        """Validate log level"""
        valid_levels = ["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"]
        if v.upper() not in valid_levels:
            raise ValueError(f"LOG_LEVEL must be one of {valid_levels}")
        return v.upper()
    
    class Config:
        env_file = ".env"
        case_sensitive = False


# Global settings instance
settings = Settings()


def validate_required_settings():
    """Validate that all required settings are present"""
    errors = []
    
    # Check Qdrant configuration
    if not settings.qdrant_url:
        errors.append("QDRANT_URL or QDRANT_HOST/QDRANT_PORT is required")
    
    # Check Ollama configuration
    if not settings.ollama_url:
        errors.append("OLLAMA_URL or OLLAMA_HOST/OLLAMA_PORT is required")
    
    # Check embedding model
    if not settings.embedding_model:
        errors.append("EMBEDDING_MODEL is required")
    
    if errors:
        error_message = "Configuration validation failed:\n" + "\n".join(f"  - {err}" for err in errors)
        raise ValueError(error_message)


# Validate settings on import
try:
    validate_required_settings()
except ValueError as e:
    import logging
    logging.error(str(e))
    raise
