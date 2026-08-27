"""
LlamaIndex framework adapter for extracting Intelligence Manifests.
"""

from __future__ import annotations

from typing import Any

from evolution.adapters.base import BaseAdapter
from evolution.models.artifacts import (
    ModelConfigArtifact,
    RetrievalArtifact,
)
from evolution.models.manifest import Manifest


class LlamaIndexAdapter(BaseAdapter):
    """Converts LlamaIndex indices, retrievers, query engines into an Intelligence Manifest."""

    @classmethod
    def from_llamaindex(cls, obj: Any, name: str = "llamaindex-agent", description: str = "") -> Manifest:
        return cls().to_manifest(obj, name=name, description=description)

    def to_manifest(self, obj: Any, name: str = "llamaindex-agent", description: str = "") -> Manifest:
        manifest = Manifest(name=name, description=description or "LlamaIndex RAG application state")

        # 1. Extract Retrieval settings
        source = "chroma"
        chunk_size: int | None = None
        top_k: int | None = None

        # Inspect storage_context / vector_store
        sc = getattr(obj, "storage_context", None)
        vs = getattr(sc, "vector_store", None) if sc else getattr(obj, "vector_store", None)
        if vs is not None:
            vs_name = vs.__class__.__name__.lower()
            if "pinecone" in vs_name:
                source = "pinecone"
            elif "weaviate" in vs_name:
                source = "weaviate"
            elif "elasticsearch" in vs_name or "elastic" in vs_name:
                source = "elasticsearch"
            elif "chroma" in vs_name:
                source = "chroma"
            else:
                source = "local"

        # Inspect retriever / query engine top_k
        retriever = getattr(obj, "retriever", None) or getattr(obj, "_retriever", None)
        if retriever is not None:
            top_k = getattr(retriever, "similarity_top_k", None) or getattr(retriever, "top_k", None)
        elif hasattr(obj, "similarity_top_k"):
            top_k = getattr(obj, "similarity_top_k")

        # Inspect chunk size
        service_context = getattr(obj, "service_context", None)
        if service_context is not None:
            node_parser = getattr(service_context, "node_parser", None)
            if node_parser is not None:
                chunk_size = getattr(node_parser, "chunk_size", None)
        elif hasattr(obj, "chunk_size"):
            chunk_size = getattr(obj, "chunk_size")

        manifest.add_artifact(RetrievalArtifact(
            name="llamaindex-retrieval",
            source=source,
            chunk_size=int(chunk_size) if chunk_size is not None else 512,
            top_k=int(top_k) if top_k is not None else 5,
            description="LlamaIndex vector retrieval configuration",
        ))

        # 2. Extract LLM configuration
        llm = getattr(obj, "llm", None)
        if service_context and not llm:
            llm = getattr(service_context, "llm", None)

        if llm is not None:
            model_name = getattr(llm, "model", getattr(llm, "model_name", "gpt-4o"))
            temperature = getattr(llm, "temperature", 0.1)

            provider = "openai"
            cls_name = llm.__class__.__name__.lower()
            if "anthropic" in cls_name or "claude" in str(model_name).lower():
                provider = "anthropic"
            elif "google" in cls_name or "gemini" in str(model_name).lower():
                provider = "google"

            manifest.add_artifact(ModelConfigArtifact(
                name="llamaindex-llm",
                model=str(model_name),
                provider=provider,
                temperature=float(temperature) if temperature is not None else 0.1,
            ))

        return manifest


def from_llamaindex(obj: Any, name: str = "llamaindex-rag", description: str = "") -> Manifest:
    """Convenience function to extract an Intelligence Manifest from a LlamaIndex object."""
    return LlamaIndexAdapter().to_manifest(obj, name=name, description=description)
