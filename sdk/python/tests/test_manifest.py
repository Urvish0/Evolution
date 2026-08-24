"""
Unit tests for Evolution Manifest operations and validation.
"""

import json
import tempfile
import unittest
from pathlib import Path

from evolution.exceptions import ManifestValidationError
from evolution.models.artifacts import ModelConfigArtifact, PromptArtifact
from evolution.models.manifest import Manifest


class TestManifest(unittest.TestCase):

    def test_manifest_creation_and_serialization(self):
        m = Manifest(name="test-agent", description="An agent for testing")
        m.add_artifact(PromptArtifact(name="sys", path="sys.txt", role="system"))
        m.add_artifact(ModelConfigArtifact(name="mc", model="gpt-4o", provider="openai"))

        self.assertEqual(len(m.list_artifacts()), 2)

        data = m.to_dict()
        self.assertEqual(data["version"], "1.0.0")
        self.assertEqual(data["name"], "test-agent")
        self.assertIn("prompts", data["artifacts"])
        self.assertIn("model_config", data["artifacts"])

    def test_manifest_validation_pass(self):
        m = Manifest(name="valid-agent", version="1.0.0")
        m.add_artifact(ModelConfigArtifact(name="mc", model="gpt-4o"))
        # Should not raise
        m.validate()

    def test_manifest_validation_failures(self):
        # Empty name
        m1 = Manifest(name="", version="1.0.0")
        with self.assertRaises(ManifestValidationError):
            m1.validate()

        # Invalid semver
        m2 = Manifest(name="test", version="invalid-version")
        with self.assertRaises(ManifestValidationError):
            m2.validate()

        # Model config with empty model
        m3 = Manifest(name="test", version="1.0.0")
        m3.add_artifact(ModelConfigArtifact(name="mc", model=""))
        with self.assertRaises(ManifestValidationError):
            m3.validate()

    def test_manifest_save_and_load(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            m = Manifest(name="persisted-agent", description="Testing disk persistence")
            m.add_artifact(PromptArtifact(name="prompt-1", path="prompts/p1.txt"))

            saved_path = m.save(tmpdir)
            self.assertTrue(saved_path.is_file())

            loaded = Manifest.load(tmpdir)
            self.assertEqual(loaded.name, "persisted-agent")
            self.assertEqual(len(loaded.list_artifacts()), 1)
            self.assertEqual(loaded.list_artifacts()[0].name, "prompt-1")


if __name__ == "__main__":
    unittest.main()
