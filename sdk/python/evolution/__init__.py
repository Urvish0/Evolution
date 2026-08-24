"""
Evolution Python SDK — AI-Native Version Control Platform.
"Version Intelligence, Not Code."
"""

from pathlib import Path

from evolution.capture import RecordContextManager, record, track
from evolution.exceptions import (
    ArtifactNotFoundError,
    CommandExecutionError,
    EvolutionError,
    ManifestNotFoundError,
    ManifestValidationError,
    RepositoryAlreadyExistsError,
    RepositoryNotFoundError,
)
from evolution.models import (
    BaseArtifact,
    EvaluationResult,
    EvaluationScore,
    Execution,
    Manifest,
    ManifestArtifacts,
    MemoryArtifact,
    ModelConfigArtifact,
    PolicyArtifact,
    PromptArtifact,
    RetrievalArtifact,
    TokenUsage,
    ToolArtifact,
    artifact_from_dict,
    compute_blob_hash,
)
from evolution.repository import CommitInfo, RepoStatus, Repository

__version__ = "0.8.0"


def init(path: Path | str = ".", name: str = "ai-intelligence") -> Repository:
    """Convenience function to initialize a new Evolution repository."""
    return Repository.init(path=path, name=name)


def open(path: Path | str = ".") -> Repository:
    """Convenience function to open an existing Evolution repository."""
    return Repository.open(path=path)


__all__ = [
    "ArtifactNotFoundError",
    "BaseArtifact",
    "CommandExecutionError",
    "CommitInfo",
    "EvaluationResult",
    "EvaluationScore",
    "EvolutionError",
    "Execution",
    "Manifest",
    "ManifestArtifacts",
    "ManifestNotFoundError",
    "ManifestValidationError",
    "MemoryArtifact",
    "ModelConfigArtifact",
    "PolicyArtifact",
    "PromptArtifact",
    "RecordContextManager",
    "RepoStatus",
    "Repository",
    "RepositoryAlreadyExistsError",
    "RepositoryNotFoundError",
    "RetrievalArtifact",
    "TokenUsage",
    "ToolArtifact",
    "__version__",
    "artifact_from_dict",
    "compute_blob_hash",
    "init",
    "open",
    "record",
    "track",
]
