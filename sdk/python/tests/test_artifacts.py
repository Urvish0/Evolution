"""
Unit tests for Evolution SDK typed artifacts and content hashing.
"""

import unittest
from pathlib import Path
import tempfile

from evolution.models.artifacts import (
    PromptArtifact,
    MemoryArtifact,
    RetrievalArtifact,
    ToolArtifact,
    ModelConfigArtifact,
    PolicyArtifact,
    artifact_from_dict,
    compute_blob_hash,
)


class TestArtifacts(unittest.TestCase):

    def test_compute_blob_hash(self):
        content = b"Hello Evolution!"
        h = compute_blob_hash(content)
        self.assertEqual(len(h), 64)
        # Verify deterministic
        self.assertEqual(h, compute_blob_hash(content))

    def test_prompt_artifact_serialization(self):
        prompt = PromptArtifact(
            name="system-prompt",
            path="prompts/sys.txt",
            role="system",
            format="template",
            description="Main system prompt",
        )
        d = prompt.to_dict()
        self.assertEqual(d["type"], "prompt")
        self.assertEqual(d["name"], "system-prompt")
        self.assertEqual(d["role"], "system")
        self.assertEqual(d["format"], "template")

        # Roundtrip
        reconstructed = artifact_from_dict(d)
        self.assertIsInstance(reconstructed, PromptArtifact)
        self.assertEqual(reconstructed.name, "system-prompt")
        self.assertEqual(reconstructed.role, "system")

    def test_model_config_artifact(self):
        mc = ModelConfigArtifact(
            name="primary-llm",
            model="claude-3.5-sonnet",
            provider="anthropic",
            temperature=0.3,
            max_tokens=4096,
        )
        d = mc.to_dict()
        self.assertEqual(d["type"], "model_config")
        self.assertEqual(d["model"], "claude-3.5-sonnet")
        self.assertEqual(d["provider"], "anthropic")
        self.assertEqual(d["temperature"], 0.3)

        reconstructed = artifact_from_dict(d)
        self.assertIsInstance(reconstructed, ModelConfigArtifact)
        self.assertEqual(reconstructed.model, "claude-3.5-sonnet")

    def test_auto_hash_file_computation(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            tmp_path = Path(tmpdir)
            file_path = tmp_path / "prompt.txt"
            file_path.write_bytes(b"You are an AI assistant.")

            prompt = PromptArtifact(name="p1", path="prompt.txt")
            self.assertEqual(prompt.hash, "")

            computed = prompt.compute_hash(workspace_root=tmp_path)
            self.assertEqual(len(computed), 64)
            self.assertEqual(prompt.hash, computed)


if __name__ == "__main__":
    unittest.main()
