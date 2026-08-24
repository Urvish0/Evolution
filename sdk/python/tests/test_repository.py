"""
Unit tests for Repository operations, execution recording, and evaluation storage.
"""

import tempfile
import unittest
from pathlib import Path

from evolution.exceptions import RepositoryAlreadyExistsError, RepositoryNotFoundError
from evolution.models.artifacts import PromptArtifact
from evolution.models.evaluation import EvaluationResult, EvaluationScore
from evolution.repository import Repository


class TestRepository(unittest.TestCase):

    def test_repo_init_and_open(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir, name="init-test")
            self.assertTrue(repo.is_valid)

            # Re-init should fail
            with self.assertRaises(RepositoryAlreadyExistsError):
                Repository.init(tmpdir)

            # Open existing
            opened = Repository.open(tmpdir)
            self.assertTrue(opened.is_valid)

    def test_manifest_management_via_repo(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir, name="agent-repo")
            manifest = repo.get_manifest()
            self.assertEqual(manifest.name, "agent-repo")

            # Add artifact and save through repo
            manifest.add_artifact(PromptArtifact(name="agent-sys", path="prompts/sys.txt"))
            repo.save_manifest(manifest)

            # Reload and verify
            reloaded = repo.get_manifest()
            self.assertIsNotNone(reloaded.get_artifact("agent-sys"))

    def test_execution_recording_and_retrieval(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir)

            exec_record = repo.record_execution(
                inputs="Analyze this contract",
                outputs="Contract is valid and binding.",
                duration_ms=180,
                prompt_tokens=50,
                completion_tokens=20,
                metadata={"model": "gpt-4o"},
            )

            self.assertIsNotNone(exec_record.id)
            self.assertEqual(exec_record.tokens.total_tokens, 70)

            # Retrieve by ID
            loaded = repo.get_execution(exec_record.id)
            self.assertEqual(loaded.inputs, "Analyze this contract")
            self.assertEqual(loaded.duration_ms, 180)

            # List executions
            all_execs = repo.list_executions()
            self.assertEqual(len(all_execs), 1)

    def test_evaluation_persistence(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir)

            eval_res = EvaluationResult(
                commit_id="commit-123",
                execution_id="exec-456",
                overall_score=0.95,
                scores={
                    "performance": EvaluationScore(name="performance", score=1.0, details="Fast"),
                    "safety": EvaluationScore(name="safety", score=0.90, details="Passed"),
                },
            )

            repo.save_evaluation(eval_res)

            # Retrieve
            loaded = repo.get_evaluation(eval_res.id)
            self.assertEqual(loaded.overall_score, 0.95)
            self.assertIn("performance", loaded.scores)
            self.assertEqual(loaded.scores["performance"].score, 1.0)


if __name__ == "__main__":
    unittest.main()
