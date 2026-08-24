"""
Execution recording context manager for tracing AI invocations.
"""

from __future__ import annotations

import time
from pathlib import Path
from typing import Any

from evolution.capture.introspect import extract_llm_response
from evolution.models.execution import Execution, TokenUsage
from evolution.repository import Repository


class RecordContextManager:
    """Context manager for tracing and recording an AI execution run."""

    def __init__(
        self,
        repo: Repository | Path | str | None = None,
        inputs: str = "",
        commit_id: str | None = None,
        metadata: dict[str, Any] | None = None,
    ):
        if isinstance(repo, Repository):
            self.repo = repo
        elif repo is not None:
            self.repo = Repository.open(repo)
        else:
            try:
                self.repo = Repository.open(".")
            except Exception:
                self.repo = None

        self.inputs = inputs
        self.outputs = ""
        self.commit_id = commit_id
        self.metadata = metadata or {}
        self.tokens = TokenUsage()
        self.duration_ms: int = 0
        self.execution: Execution | None = None
        self._start_time: float = 0.0

    def set_input(self, inputs: str) -> RecordContextManager:
        """Sets or updates the input query string."""
        self.inputs = inputs
        return self

    def set_output(self, outputs: Any) -> RecordContextManager:
        """Sets or updates the output string, auto-extracting from LLM response objects if needed."""
        text, tokens, _ = extract_llm_response(outputs)
        self.outputs = text
        if tokens.total_tokens > 0:
            self.tokens = tokens
        return self

    def set_tokens(
        self,
        prompt_tokens: int = 0,
        completion_tokens: int = 0,
        total_tokens: int = 0,
    ) -> RecordContextManager:
        """Sets token consumption metrics."""
        self.tokens = TokenUsage(
            prompt_tokens=prompt_tokens,
            completion_tokens=completion_tokens,
            total_tokens=total_tokens or (prompt_tokens + completion_tokens),
        )
        return self

    def set_metadata(self, key: str, value: Any) -> RecordContextManager:
        """Sets a metadata key-value pair."""
        self.metadata[key] = value
        return self

    def __enter__(self) -> RecordContextManager:
        self._start_time = time.perf_counter()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        elapsed = time.perf_counter() - self._start_time
        self.duration_ms = max(1, int(elapsed * 1000))

        if exc_type is not None and not self.outputs:
            self.outputs = f"Error ({exc_type.__name__}): {exc_val}"
            self.metadata["error"] = str(exc_val)
            self.metadata["error_type"] = exc_type.__name__

        if self.repo is not None:
            self.execution = self.repo.record_execution(
                inputs=self.inputs,
                outputs=self.outputs,
                duration_ms=self.duration_ms,
                prompt_tokens=self.tokens.prompt_tokens,
                completion_tokens=self.tokens.completion_tokens,
                commit_id=self.commit_id,
                metadata=self.metadata,
            )


def record(
    repo: Repository | Path | str | None = None,
    inputs: str = "",
    commit_id: str | None = None,
    metadata: dict[str, Any] | None = None,
) -> RecordContextManager:
    """Convenience factory returning an execution recording context manager."""
    return RecordContextManager(repo=repo, inputs=inputs, commit_id=commit_id, metadata=metadata)
