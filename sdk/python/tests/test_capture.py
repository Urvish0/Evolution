"""
Unit tests for automatic intelligence capture (@track and record context manager).
"""

import asyncio
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace

from evolution import Repository, record, track
from evolution.capture.introspect import extract_llm_response


class TestCapture(unittest.TestCase):

    def test_extract_llm_response_dict(self):
        # OpenAI style dict
        openai_dict = {
            "choices": [{"message": {"content": "This is legal advice."}}],
            "usage": {"prompt_tokens": 15, "completion_tokens": 10, "total_tokens": 25},
            "model": "gpt-4o",
        }
        text, tokens, model = extract_llm_response(openai_dict)
        self.assertEqual(text, "This is legal advice.")
        self.assertEqual(tokens.total_tokens, 25)
        self.assertEqual(model, "gpt-4o")

        # Anthropic style dict
        anthropic_dict = {
            "content": [{"text": "Claude response"}],
            "usage": {"input_tokens": 20, "output_tokens": 30},
            "model": "claude-3.5-sonnet",
        }
        text, tokens, model = extract_llm_response(anthropic_dict)
        self.assertEqual(text, "Claude response")
        self.assertEqual(tokens.prompt_tokens, 20)
        self.assertEqual(tokens.completion_tokens, 30)
        self.assertEqual(tokens.total_tokens, 50)
        self.assertEqual(model, "claude-3.5-sonnet")

    def test_extract_llm_response_objects(self):
        # Mock OpenAI ChatCompletion object
        mock_obj = SimpleNamespace(
            choices=[SimpleNamespace(message=SimpleNamespace(content="Object output"))],
            usage=SimpleNamespace(prompt_tokens=100, completion_tokens=50, total_tokens=150),
            model="gpt-4o-mini",
        )
        text, tokens, model = extract_llm_response(mock_obj)
        self.assertEqual(text, "Object output")
        self.assertEqual(tokens.total_tokens, 150)
        self.assertEqual(model, "gpt-4o-mini")

    def test_record_context_manager(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir)

            with record(repo, inputs="What is negligence?") as rec:
                rec.set_output("Negligence is a failure to exercise appropriate care.")
                rec.set_tokens(prompt_tokens=12, completion_tokens=18)
                rec.set_metadata("source", "black_law_dict")

            self.assertIsNotNone(rec.execution)
            self.assertEqual(rec.execution.inputs, "What is negligence?")
            self.assertEqual(rec.execution.tokens.total_tokens, 30)
            self.assertGreaterEqual(rec.execution.duration_ms, 1)

            # Check persisted
            executions = repo.list_executions()
            self.assertEqual(len(executions), 1)
            self.assertEqual(executions[0].id, rec.execution.id)

    def test_track_decorator_sync(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir)

            @track(repo=repo, model="gpt-4o", temperature=0.2)
            def analyze_case(case_name: str, jurisdiction: str = "US"):
                """Analyzes appellate case law precedents."""
                return {
                    "choices": [{"message": {"content": f"Analysis for {case_name} in {jurisdiction}"}}],
                    "usage": {"prompt_tokens": 40, "completion_tokens": 60},
                }

            result = analyze_case("Marbury v. Madison")
            self.assertIn("choices", result)

            # Check execution was recorded
            executions = repo.list_executions()
            self.assertEqual(len(executions), 1)
            self.assertIn("Marbury v. Madison", executions[0].inputs)
            self.assertIn("Analysis for Marbury v. Madison", executions[0].outputs)
            self.assertEqual(executions[0].tokens.total_tokens, 100)

            # Check manifest auto-updated with docstring prompt and model config
            manifest = repo.get_manifest()
            self.assertIsNotNone(manifest.artifacts.model_config)
            self.assertEqual(manifest.artifacts.model_config.model, "gpt-4o")
            self.assertEqual(len(manifest.artifacts.prompts), 1)
            self.assertEqual(manifest.artifacts.prompts[0].name, "analyze_case-prompt")

    def test_track_decorator_async(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            repo = Repository.init(tmpdir)

            @track(repo=repo, name="async-agent", model="claude-3.5-sonnet")
            async def async_agent(prompt: str):
                """Async legal query assistant."""
                await asyncio.sleep(0.01)
                return f"Answer for: {prompt}"

            result = asyncio.run(async_agent("Statute of limitations"))
            self.assertEqual(result, "Answer for: Statute of limitations")

            executions = repo.list_executions()
            self.assertEqual(len(executions), 1)
            self.assertEqual(executions[0].outputs, "Answer for: Statute of limitations")


if __name__ == "__main__":
    unittest.main()
