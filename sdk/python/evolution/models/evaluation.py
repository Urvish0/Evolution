"""
Evaluation report models conforming to Intelligence Manifest Specification v1.0.
"""

from __future__ import annotations

import uuid
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any


@dataclass
class EvaluationScore:
    """Score produced by a single evaluator."""
    name: str
    score: float  # Normalized 0.0 to 1.0
    unit: str = ""
    details: str = ""

    def to_dict(self) -> dict[str, Any]:
        d: dict[str, Any] = {
            "name": self.name,
            "score": self.score,
        }
        if self.unit:
            d["unit"] = self.unit
        if self.details:
            d["details"] = self.details
        return d

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> EvaluationScore:
        return cls(
            name=data.get("name", ""),
            score=float(data.get("score", 0.0)),
            unit=data.get("unit", ""),
            details=data.get("details", ""),
        )


@dataclass
class EvaluationResult:
    """Complete evaluation report for an AI execution."""
    commit_id: str
    execution_id: str
    overall_score: float
    scores: dict[str, EvaluationScore] = field(default_factory=dict)
    id: str = field(default_factory=lambda: str(uuid.uuid4()))
    timestamp: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())

    def to_dict(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "commit_id": self.commit_id,
            "execution_id": self.execution_id,
            "overall_score": self.overall_score,
            "scores": {k: v.to_dict() for k, v in self.scores.items()},
            "timestamp": self.timestamp,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> EvaluationResult:
        scores_raw = data.get("scores", {})
        scores = {}
        if isinstance(scores_raw, dict):
            for k, v in scores_raw.items():
                if isinstance(v, dict):
                    scores[k] = EvaluationScore.from_dict(v)
        return cls(
            id=data.get("id", str(uuid.uuid4())),
            commit_id=data.get("commit_id", ""),
            execution_id=data.get("execution_id", ""),
            overall_score=float(data.get("overall_score", 0.0)),
            scores=scores,
            timestamp=data.get("timestamp", datetime.now(timezone.utc).isoformat()),
        )
