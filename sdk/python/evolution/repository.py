"""
Repository operations for the Evolution Python SDK.
"""

from __future__ import annotations

import json
import shutil
import subprocess
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

from evolution.exceptions import (
    CommandExecutionError,
    ManifestNotFoundError,
    RepositoryAlreadyExistsError,
    RepositoryNotFoundError,
)
from evolution.models.evaluation import EvaluationResult
from evolution.models.execution import Execution
from evolution.models.manifest import MANIFEST_FILE_NAME, Manifest

EVOLUTION_DIR = ".evolution"
EXECUTIONS_DIR = "executions"
EVALUATIONS_DIR = "evaluations"


@dataclass
class CommitInfo:
    """Represents an Intelligence Commit summary."""
    id: str
    tree_id: str = ""
    parent_ids: list[str] = field(default_factory=list)
    author: str = ""
    timestamp: str = ""
    message: str = ""
    tags: list[str] = field(default_factory=list)
    metadata: dict[str, str] = field(default_factory=dict)


@dataclass
class RepoStatus:
    """Represents the current repository working tree status."""
    branch: str = "main"
    head_commit: str = ""
    is_clean: bool = True
    staged_files: list[str] = field(default_factory=list)
    modified_files: list[str] = field(default_factory=list)
    untracked_files: list[str] = field(default_factory=list)
    deleted_files: list[str] = field(default_factory=list)


class Repository:
    """Interface to an Evolution intelligence repository."""

    def __init__(self, root_path: Path | str = "."):
        self.root = Path(root_path).resolve()
        self.evolution_dir = self.root / EVOLUTION_DIR
        self._evo_bin = shutil.which("evo")

    @classmethod
    def init(cls, path: Path | str = ".", name: str = "ai-intelligence") -> Repository:
        """Initializes a new Evolution repository at the target path."""
        target = Path(path).resolve()
        evo_dir = target / EVOLUTION_DIR

        if evo_dir.is_dir():
            raise RepositoryAlreadyExistsError(f"Evolution repository already exists at {target}")

        target.mkdir(parents=True, exist_ok=True)
        repo = cls(target)

        if repo._evo_bin:
            repo._run_cli(["init"])
        else:
            # Native fallback initialization
            evo_dir.mkdir(parents=True, exist_ok=True)
            (evo_dir / "objects").mkdir(exist_ok=True)
            (evo_dir / "refs" / "heads").mkdir(parents=True, exist_ok=True)
            (evo_dir / EXECUTIONS_DIR).mkdir(exist_ok=True)
            (evo_dir / EVALUATIONS_DIR).mkdir(exist_ok=True)
            (evo_dir / "HEAD").write_text("ref: refs/heads/main\n", encoding="utf-8")

        # Create starter manifest if not present
        manifest_path = target / MANIFEST_FILE_NAME
        if not manifest_path.is_file():
            m = Manifest(name=name)
            m.save(target)

        return repo

    @classmethod
    def open(cls, path: Path | str = ".") -> Repository:
        """Opens an existing Evolution repository."""
        target = Path(path).resolve()
        # Search upward if not directly in repo root
        cur = target
        while cur != cur.parent:
            if (cur / EVOLUTION_DIR).is_dir():
                return cls(cur)
            cur = cur.parent

        raise RepositoryNotFoundError(f"No Evolution repository found at or above {target}")

    @property
    def is_valid(self) -> bool:
        """Checks if this repository directory contains an initialized .evolution structure."""
        return self.evolution_dir.is_dir()

    def _run_cli(self, args: list[str]) -> str:
        """Runs the evo CLI in the workspace root if available."""
        if not self._evo_bin:
            raise CommandExecutionError("evo", 127, "Evolution CLI binary 'evo' not found in PATH")

        cmd = [self._evo_bin] + args
        try:
            result = subprocess.run(
                cmd,
                cwd=str(self.root),
                capture_output=True,
                text=True,
                check=False,
            )
        except Exception as e:
            raise CommandExecutionError(" ".join(cmd), 1, str(e))

        if result.returncode != 0:
            raise CommandExecutionError(" ".join(cmd), result.returncode, result.stderr)

        return result.stdout.strip()

    # --- Manifest Operations ---

    def get_manifest(self) -> Manifest:
        """Loads and returns the current workspace Intelligence Manifest."""
        manifest_path = self.root / MANIFEST_FILE_NAME
        if not manifest_path.is_file():
            raise ManifestNotFoundError(f"Manifest not found in repository at {manifest_path}")
        return Manifest.load(manifest_path)

    def save_manifest(self, manifest: Manifest, auto_hash: bool = True) -> Path:
        """Saves a manifest to the repository workspace, optionally auto-hashing artifacts."""
        if auto_hash:
            manifest.compute_hashes(workspace_root=self.root)
        manifest.validate()
        return manifest.save(self.root)

    # --- Git-Like VCS Operations ---

    def add(self, *paths: str) -> None:
        """Stages file paths for the next commit."""
        if self._evo_bin:
            args = ["add"] + list(paths)
            self._run_cli(args)

    def commit(
        self,
        message: str,
        author: str | None = None,
        tags: list[str] | None = None,
        metadata: dict[str, str] | None = None,
        auto_stage_all: bool = True,
    ) -> CommitInfo:
        """Creates an Intelligence Commit capturing the current intelligence state."""
        # Auto-compute hashes on manifest before commit if manifest exists
        manifest_path = self.root / MANIFEST_FILE_NAME
        if manifest_path.is_file():
            try:
                m = self.get_manifest()
                self.save_manifest(m, auto_hash=True)
            except Exception:
                pass

        if auto_stage_all and self._evo_bin:
            self._run_cli(["add", "."])

        if self._evo_bin:
            args = ["commit", "-m", message]
            if author:
                args += ["--author", author]
            if tags:
                for t in tags:
                    args += ["--tag", t]
            if metadata:
                for k, v in metadata.items():
                    args += ["--meta", f"{k}={v}"]

            out = self._run_cli(args)
            # Parse commit ID from output or read current branch head
            commit_id = "unknown"
            for line in out.splitlines():
                if "[" in line and "]" in line:
                    parts = line.split("[")[1].split("]")[0].split()
                    if len(parts) >= 2:
                        commit_id = parts[1]
                elif line.startswith("commit "):
                    commit_id = line.split()[1]

            if commit_id == "unknown":
                head_file = self.root / ".evolution" / "HEAD"
                if head_file.is_file():
                    branch_name = head_file.read_text(encoding="utf-8").strip()
                    branch_file = self.root / ".evolution" / "branches" / f"{branch_name}.json"
                    if branch_file.is_file():
                        try:
                            bdata = json.loads(branch_file.read_text(encoding="utf-8"))
                            commit_id = bdata.get("head", "unknown")
                        except Exception:
                            pass

            return CommitInfo(id=commit_id, message=message, author=author or "", tags=tags or [], metadata=metadata or {})

        # Fallback simulation when running in pure-python mode without evo binary
        fake_id = "py-" + message[:8].replace(" ", "_")
        return CommitInfo(id=fake_id, message=message, author=author or "", tags=tags or [], metadata=metadata or {})

    def status(self) -> RepoStatus:
        """Returns the current working tree and branch status."""
        import re
        if self._evo_bin:
            out = self._run_cli(["status"])
            status = RepoStatus()
            for line in out.splitlines():
                clean_line = re.sub(r"\x1b\[[0-9;]*m", "", line)
                if clean_line.startswith("On branch "):
                    status.branch = clean_line.replace("On branch ", "").strip()
                elif "nothing to commit, working tree clean" in clean_line:
                    status.is_clean = True
            return status
        return RepoStatus(branch="main", is_clean=True)

    def diff(self, rev1: str | None = None, rev2: str | None = None) -> str:
        """Renders unified diff between revisions or working tree."""
        if self._evo_bin:
            args = ["diff"]
            if rev1:
                args.append(rev1)
            if rev2:
                args.append(rev2)
            return self._run_cli(args)
        return ""

    def checkout(self, target: str, file_path: str | None = None) -> str:
        """Checks out a branch, commit, or restores a file snapshot."""
        if self._evo_bin:
            args = ["checkout", target]
            if file_path:
                args += ["--", file_path]
            return self._run_cli(args)
        return ""

    # --- Execution Recording Operations ---

    def record_execution(
        self,
        inputs: str,
        outputs: str,
        duration_ms: int = 0,
        prompt_tokens: int = 0,
        completion_tokens: int = 0,
        commit_id: str | None = None,
        metadata: dict[str, Any] | None = None,
    ) -> Execution:
        """Records an AI system execution linked to the current HEAD commit snapshot."""
        from evolution.models.execution import TokenUsage

        if not commit_id:
            # Attempt to resolve current HEAD commit
            head_file = self.evolution_dir / "HEAD"
            if head_file.is_file():
                head_ref = head_file.read_text().strip()
                if head_ref.startswith("ref: "):
                    ref_path = self.evolution_dir / head_ref[5:]
                    if ref_path.is_file():
                        commit_id = ref_path.read_text().strip()
                    else:
                        commit_id = "uncommitted"
                else:
                    commit_id = head_ref
            else:
                commit_id = "uncommitted"

        exec_obj = Execution(
            commit_id=commit_id,
            inputs=inputs,
            outputs=outputs,
            duration_ms=duration_ms,
            tokens=TokenUsage(prompt_tokens=prompt_tokens, completion_tokens=completion_tokens),
            metadata=metadata or {},
        )

        self.save_execution(exec_obj)
        return exec_obj

    def save_execution(self, execution: Execution) -> Path:
        """Saves an execution record to .evolution/executions/<id>.json."""
        exec_dir = self.evolution_dir / EXECUTIONS_DIR
        exec_dir.mkdir(parents=True, exist_ok=True)
        file_path = exec_dir / f"{execution.id}.json"
        file_path.write_text(json.dumps(execution.to_dict(), indent=2) + "\n", encoding="utf-8")
        return file_path

    def get_execution(self, execution_id: str) -> Execution:
        """Loads an execution record by ID."""
        file_path = self.evolution_dir / EXECUTIONS_DIR / f"{execution_id}.json"
        if not file_path.is_file():
            # Try prefix match
            exec_dir = self.evolution_dir / EXECUTIONS_DIR
            if exec_dir.is_dir():
                matches = list(exec_dir.glob(f"{execution_id}*.json"))
                if len(matches) == 1:
                    file_path = matches[0]
                elif len(matches) > 1:
                    raise ValueError(f"Ambiguous execution ID prefix '{execution_id}'")

        if not file_path.is_file():
            raise FileNotFoundError(f"Execution '{execution_id}' not found")

        data = json.loads(file_path.read_text(encoding="utf-8"))
        return Execution.from_dict(data)

    def list_executions(self) -> list[Execution]:
        """Lists all recorded executions in reverse chronological order."""
        exec_dir = self.evolution_dir / EXECUTIONS_DIR
        if not exec_dir.is_dir():
            return []

        executions: list[Execution] = []
        for file in exec_dir.glob("*.json"):
            try:
                data = json.loads(file.read_text(encoding="utf-8"))
                executions.append(Execution.from_dict(data))
            except Exception:
                continue

        executions.sort(key=lambda e: e.timestamp, reverse=True)
        return executions

    # --- Evaluation Operations ---

    def save_evaluation(self, evaluation: EvaluationResult) -> Path:
        """Saves an evaluation result to .evolution/evaluations/<id>.json."""
        eval_dir = self.evolution_dir / EVALUATIONS_DIR
        eval_dir.mkdir(parents=True, exist_ok=True)
        file_path = eval_dir / f"{evaluation.id}.json"
        file_path.write_text(json.dumps(evaluation.to_dict(), indent=2) + "\n", encoding="utf-8")
        return file_path

    def get_evaluation(self, evaluation_id: str) -> EvaluationResult:
        """Loads an evaluation report by ID."""
        file_path = self.evolution_dir / EVALUATIONS_DIR / f"{evaluation_id}.json"
        if not file_path.is_file():
            raise FileNotFoundError(f"Evaluation '{evaluation_id}' not found")

        data = json.loads(file_path.read_text(encoding="utf-8"))
        return EvaluationResult.from_dict(data)

    def list_evaluations(self) -> list[EvaluationResult]:
        """Lists all recorded evaluations."""
        eval_dir = self.evolution_dir / EVALUATIONS_DIR
        if not eval_dir.is_dir():
            return []

        evaluations: list[EvaluationResult] = []
        for file in eval_dir.glob("*.json"):
            try:
                data = json.loads(file.read_text(encoding="utf-8"))
                evaluations.append(EvaluationResult.from_dict(data))
            except Exception:
                continue

        evaluations.sort(key=lambda e: e.timestamp, reverse=True)
        return evaluations
