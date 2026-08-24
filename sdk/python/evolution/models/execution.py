"""
Execution recording models conforming to Intelligence Manifest Specification v1.0.
"""

from __future__ import annotations

import uuid
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any


@dataclass
class TokenUsage:
    """Token consumption metrics for an execution run."""
    prompt_tokens: int = 0
    completion_tokens: int = 0
    total_tokens: int = 0

    def __post_init__(self):
        if self.total_tokens == 0:
            self.total_tokens = self.prompt_tokens + self.completion_tokens

    def to_dict(self) -> dict[str, Any]:
        return {
            "prompt_tokens": self.prompt_tokens,
            "completion_tokens": self.completion_tokens,
            "total_tokens": self.total_tokens,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> TokenUsage:
        return cls(
            prompt_tokens=data.get("prompt_tokens", 0),
            completion_tokens=data.get("completion_tokens", 0),
            total_tokens=data.get("total_tokens", 0),
        )


@dataclass
class Execution:
    """Records a single AI system invocation at a specific commit snapshot."""
    commit_id: str
    inputs: str
    outputs: str
    duration_ms: int = 0
    tokens: TokenUsage = field(default_factory=TokenUsage)
    id: str = field(default_factory=lambda: str(uuid.uuid4()))
    timestamp: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    metadata: dict[str, Any] = field(default_factory=dict)

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "commit_id": self.commit_id,
            "inputs": self.inputs,
            "outputs": self.outputs,
            "duration_ms": self.duration_ms,
            "tokens": self.tokens.to_dict(),
            "timestamp": self.timestamp,
            "metadata": self.metadata,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Execution:
        tokens_data = data.get("tokens", {})
        tokens = TokenUsage.from_dict(tokens_data) if isinstance(tokens_data, dict) else TokenUsage()
        return cls(
            id=data.get("id", str(uuid.uuid4())),
            commit_id=data.get("commit_id", ""),
            inputs=data.get("inputs", ""),
            outputs=data.get("outputs", ""),
            duration_ms=data.get("duration_ms", 0),
            tokens=tokens,
            timestamp=data.get("timestamp", datetime.now(timezone.utc).isoformat()),
            metadata=data.get("metadata", {}),
        )
