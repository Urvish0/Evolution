"""
Typed AI Artifact models conforming to Intelligence Manifest Specification v1.0.
"""

from __future__ import annotations

import hashlib
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Literal


ArtifactType = Literal["prompt", "memory", "retrieval", "tool", "model_config", "policy"]


def compute_blob_hash(content: bytes) -> str:
    """Computes SHA-256 hash using Evolution's Git-compatible blob header format:
    SHA-256('blob <len>\0<content>').
    """
    header = f"blob {len(content)}\0".encode("utf-8")
    hasher = hashlib.sha256()
    hasher.update(header + content)
    return hasher.hexdigest()


@dataclass
class BaseArtifact:
    """Base class for all typed AI artifacts."""
    type: ArtifactType
    name: str
    path: str = ""
    hash: str = ""
    description: str = ""

    def compute_hash(self, workspace_root: Path | str | None = None) -> str:
        """Computes the SHA-256 content hash from the underlying file if path exists.
        Updates self.hash and returns it.
        """
        if not self.path:
            return self.hash

        file_path = Path(self.path)
        if workspace_root and not file_path.is_absolute():
            file_path = Path(workspace_root) / file_path

        if file_path.is_file():
            content = file_path.read_bytes()
            self.hash = compute_blob_hash(content)
        return self.hash

    def to_dict(self) -> dict[str, Any]:
        """Serializes the artifact to a dictionary matching the v1.0 manifest schema."""
        d: dict[str, Any] = {
            "type": self.type,
            "name": self.name,
            "path": self.path,
        }
        if self.hash:
            d["hash"] = self.hash
        if self.description:
            d["description"] = self.description
        return d


@dataclass
class PromptArtifact(BaseArtifact):
    """Prompt template, system instructions, or few-shot examples."""
    type: ArtifactType = field(default="prompt", init=False)
    role: Literal["system", "user", "assistant", "few_shot"] = "system"
    format: Literal["text", "template", "jinja2", "mustache"] = "text"

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        d["role"] = self.role
        d["format"] = self.format
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> PromptArtifact:
        return cls(
            name=data["name"],
            path=data.get("path", ""),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            role=data.get("role", "system"),
            format=data.get("format", "text"),
        )


@dataclass
class MemoryArtifact(BaseArtifact):
    """Conversation history and context management strategy."""
    type: ArtifactType = field(default="memory", init=False)
    strategy: Literal["buffer_window", "summary", "vector", "graph"] = "buffer_window"
    max_tokens: int | None = None

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        d["strategy"] = self.strategy
        if self.max_tokens is not None:
            d["max_tokens"] = self.max_tokens
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> MemoryArtifact:
        return cls(
            name=data["name"],
            path=data.get("path", ""),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            strategy=data.get("strategy", "buffer_window"),
            max_tokens=data.get("max_tokens"),
        )


@dataclass
class RetrievalArtifact(BaseArtifact):
    """Vector database and semantic search retrieval configuration."""
    type: ArtifactType = field(default="retrieval", init=False)
    source: Literal["pinecone", "chroma", "weaviate", "local", "elasticsearch"] = "chroma"
    chunk_size: int | None = None
    top_k: int | None = None

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        d["source"] = self.source
        if self.chunk_size is not None:
            d["chunk_size"] = self.chunk_size
        if self.top_k is not None:
            d["top_k"] = self.top_k
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> RetrievalArtifact:
        return cls(
            name=data["name"],
            path=data.get("path", ""),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            source=data.get("source", "chroma"),
            chunk_size=data.get("chunk_size"),
            top_k=data.get("top_k"),
        )


@dataclass
class ToolArtifact(BaseArtifact):
    """External tool, API, or function calling definition."""
    type: ArtifactType = field(default="tool", init=False)
    provider: str = ""
    auth_required: bool = False

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        if self.provider:
            d["provider"] = self.provider
        d["auth_required"] = self.auth_required
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> ToolArtifact:
        return cls(
            name=data["name"],
            path=data.get("path", ""),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            provider=data.get("provider", ""),
            auth_required=data.get("auth_required", False),
        )


@dataclass
class ModelConfigArtifact(BaseArtifact):
    """LLM provider, model name, and inference parameters."""
    type: ArtifactType = field(default="model_config", init=False)
    model: str = "gpt-4o"
    provider: Literal["openai", "anthropic", "google", "local", "mistral", "cohere", "aws_bedrock"] = "openai"
    temperature: float | None = 0.7
    max_tokens: int | None = None
    top_p: float | None = None

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        d["model"] = self.model
        d["provider"] = self.provider
        if self.temperature is not None:
            d["temperature"] = self.temperature
        if self.max_tokens is not None:
            d["max_tokens"] = self.max_tokens
        if self.top_p is not None:
            d["top_p"] = self.top_p
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> ModelConfigArtifact:
        return cls(
            name=data.get("name", "primary-model"),
            path=data.get("path", "config/model.json"),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            model=data.get("model", "gpt-4o"),
            provider=data.get("provider", "openai"),
            temperature=data.get("temperature", 0.7),
            max_tokens=data.get("max_tokens"),
            top_p=data.get("top_p"),
        )


@dataclass
class PolicyArtifact(BaseArtifact):
    """Safety guardrail, compliance rule, or output policy."""
    type: ArtifactType = field(default="policy", init=False)
    enforcement: Literal["strict", "warn", "log"] = "strict"

    def to_dict(self) -> dict[str, Any]:
        d = super().to_dict()
        d["enforcement"] = self.enforcement
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> PolicyArtifact:
        return cls(
            name=data["name"],
            path=data.get("path", ""),
            hash=data.get("hash", ""),
            description=data.get("description", ""),
            enforcement=data.get("enforcement", "strict"),
        )


# Mapping from type string to artifact class
ARTIFACT_CLASS_MAP = {
    "prompt": PromptArtifact,
    "memory": MemoryArtifact,
    "retrieval": RetrievalArtifact,
    "tool": ToolArtifact,
    "model_config": ModelConfigArtifact,
    "policy": PolicyArtifact,
}


def artifact_from_dict(data: dict[str, Any]) -> BaseArtifact:
    """Instantiates the appropriate typed artifact subclass from a dictionary."""
    art_type = data.get("type", "")
    cls = ARTIFACT_CLASS_MAP.get(art_type)
    if not cls:
        raise ValueError(f"Unknown artifact type: {art_type}")
    return cls.from_dict(data)
