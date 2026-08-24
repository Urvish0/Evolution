"""
Intelligence Manifest Manager conforming to Intelligence Manifest Specification v1.0.
"""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

from evolution.exceptions import ManifestNotFoundError, ManifestValidationError
from evolution.models.artifacts import (
    BaseArtifact,
    MemoryArtifact,
    ModelConfigArtifact,
    PolicyArtifact,
    PromptArtifact,
    RetrievalArtifact,
    ToolArtifact,
    artifact_from_dict,
)

MANIFEST_FILE_NAME = "evolution.manifest.json"
DEFAULT_SPEC_VERSION = "1.0.0"


@dataclass
class ManifestArtifacts:
    """Container for grouped typed artifacts."""
    prompts: list[PromptArtifact] = field(default_factory=list)
    memory: list[MemoryArtifact] = field(default_factory=list)
    retrieval: list[RetrievalArtifact] = field(default_factory=list)
    tools: list[ToolArtifact] = field(default_factory=list)
    model_config: ModelConfigArtifact | None = None
    policies: list[PolicyArtifact] = field(default_factory=list)

    def all(self) -> list[BaseArtifact]:
        """Returns a flat list of all attached artifacts."""
        res: list[BaseArtifact] = []
        res.extend(self.prompts)
        res.extend(self.memory)
        res.extend(self.retrieval)
        res.extend(self.tools)
        if self.model_config:
            res.append(self.model_config)
        res.extend(self.policies)
        return res

    def to_dict(self) -> dict[str, Any]:
        """Serializes artifacts dictionary matching the schema."""
        d: dict[str, Any] = {}
        if self.prompts:
            d["prompts"] = [p.to_dict() for p in self.prompts]
        if self.memory:
            d["memory"] = [m.to_dict() for m in self.memory]
        if self.retrieval:
            d["retrieval"] = [r.to_dict() for r in self.retrieval]
        if self.tools:
            d["tools"] = [t.to_dict() for t in self.tools]
        if self.model_config:
            d["model_config"] = self.model_config.to_dict()
        if self.policies:
            d["policies"] = [p.to_dict() for p in self.policies]
        return d


@dataclass
class Manifest:
    """The Intelligence Manifest captures the complete operational state of an AI system."""
    name: str = "ai-intelligence"
    version: str = DEFAULT_SPEC_VERSION
    description: str = "AI system powered by Evolution version control"
    artifacts: ManifestArtifacts = field(default_factory=ManifestArtifacts)
    metadata: dict[str, Any] = field(default_factory=dict)

    def add_artifact(self, artifact: BaseArtifact) -> Manifest:
        """Adds or updates a typed artifact in the manifest."""
        if isinstance(artifact, PromptArtifact):
            self.artifacts.prompts = [p for p in self.artifacts.prompts if p.name != artifact.name] + [artifact]
        elif isinstance(artifact, MemoryArtifact):
            self.artifacts.memory = [m for m in self.artifacts.memory if m.name != artifact.name] + [artifact]
        elif isinstance(artifact, RetrievalArtifact):
            self.artifacts.retrieval = [r for r in self.artifacts.retrieval if r.name != artifact.name] + [artifact]
        elif isinstance(artifact, ToolArtifact):
            self.artifacts.tools = [t for t in self.artifacts.tools if t.name != artifact.name] + [artifact]
        elif isinstance(artifact, ModelConfigArtifact):
            self.artifacts.model_config = artifact
        elif isinstance(artifact, PolicyArtifact):
            self.artifacts.policies = [p for p in self.artifacts.policies if p.name != artifact.name] + [artifact]
        else:
            raise TypeError(f"Unsupported artifact class: {type(artifact)}")
        return self

    def get_artifact(self, name: str) -> BaseArtifact | None:
        """Finds an artifact by its unique name."""
        for art in self.artifacts.all():
            if art.name == name:
                return art
        return None

    def list_artifacts(self) -> list[BaseArtifact]:
        """Returns all artifacts as a flat list."""
        return self.artifacts.all()

    def remove_artifact(self, name: str) -> bool:
        """Removes an artifact by name. Returns True if removed, False if not found."""
        initial_len = len(self.artifacts.all())
        self.artifacts.prompts = [p for p in self.artifacts.prompts if p.name != name]
        self.artifacts.memory = [m for m in self.artifacts.memory if m.name != name]
        self.artifacts.retrieval = [r for r in self.artifacts.retrieval if r.name != name]
        self.artifacts.tools = [t for t in self.artifacts.tools if t.name != name]
        if self.artifacts.model_config and self.artifacts.model_config.name == name:
            self.artifacts.model_config = None
        self.artifacts.policies = [p for p in self.artifacts.policies if p.name != name]
        return len(self.artifacts.all()) < initial_len

    def compute_hashes(self, workspace_root: Path | str | None = None) -> None:
        """Auto-computes SHA-256 hashes for all attached artifacts with valid file paths."""
        for art in self.artifacts.all():
            art.compute_hash(workspace_root=workspace_root)

    def validate(self) -> None:
        """Validates manifest against v1.0 Specification.
        Raises ManifestValidationError on non-compliance.
        """
        from evolution.validator import validate_manifest
        validate_manifest(self)

    def to_dict(self) -> dict[str, Any]:
        """Serializes the manifest to a clean dict conforming to Spec v1.0."""
        d: dict[str, Any] = {
            "version": self.version,
            "name": self.name,
        }
        if self.description:
            d["description"] = self.description

        art_dict = self.artifacts.to_dict()
        if art_dict:
            d["artifacts"] = art_dict

        if self.metadata:
            d["metadata"] = self.metadata
        return d

    def to_json(self, indent: int = 2) -> str:
        """Returns formatted JSON string of the manifest."""
        return json.dumps(self.to_dict(), indent=indent, ensure_ascii=False)

    def save(self, destination: Path | str | None = None) -> Path:
        """Saves the manifest as evolution.manifest.json to the specified file or directory."""
        if destination is None:
            dest_path = Path(MANIFEST_FILE_NAME)
        else:
            dest_path = Path(destination)
            if dest_path.is_dir() or not dest_path.name.endswith(".json"):
                dest_path = dest_path / MANIFEST_FILE_NAME

        dest_path.parent.mkdir(parents=True, exist_ok=True)
        dest_path.write_text(self.to_json() + "\n", encoding="utf-8")
        return dest_path

    @classmethod
    def load(cls, source: Path | str | None = None) -> Manifest:
        """Loads a manifest from a file or workspace directory."""
        if source is None:
            src_path = Path(MANIFEST_FILE_NAME)
        else:
            src_path = Path(source)
            if src_path.is_dir():
                src_path = src_path / MANIFEST_FILE_NAME

        if not src_path.is_file():
            raise ManifestNotFoundError(f"Manifest not found at {src_path}")

        try:
            data = json.loads(src_path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as e:
            raise ManifestValidationError(f"Invalid JSON in manifest file: {e}")

        return cls.from_dict(data)

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> Manifest:
        """Parses a dictionary into a Manifest instance."""
        artifacts = ManifestArtifacts()
        raw_artifacts = data.get("artifacts", {})

        if isinstance(raw_artifacts, dict):
            for p in raw_artifacts.get("prompts", []):
                artifacts.prompts.append(PromptArtifact.from_dict(p))
            for m in raw_artifacts.get("memory", []):
                artifacts.memory.append(MemoryArtifact.from_dict(m))
            for r in raw_artifacts.get("retrieval", []):
                artifacts.retrieval.append(RetrievalArtifact.from_dict(r))
            for t in raw_artifacts.get("tools", []):
                artifacts.tools.append(ToolArtifact.from_dict(t))
            if "model_config" in raw_artifacts and isinstance(raw_artifacts["model_config"], dict):
                artifacts.model_config = ModelConfigArtifact.from_dict(raw_artifacts["model_config"])
            for pol in raw_artifacts.get("policies", []):
                artifacts.policies.append(PolicyArtifact.from_dict(pol))

        return cls(
            name=data.get("name", "ai-intelligence"),
            version=data.get("version", DEFAULT_SPEC_VERSION),
            description=data.get("description", ""),
            artifacts=artifacts,
            metadata=data.get("metadata", {}),
        )
