"""
Agent 1: Strict Corporate Legal Analyst
- Model: llama-3.3-70b-versatile (or llama-3.1-8b-instant)
- Temperature: 0.1 (Strict, deterministic)
- Persona: Aggressive liability analysis, statutory citations, no emotional language.
"""

from __future__ import annotations

import os
from typing import Any
from groq import Groq
import evolution as evo


def create_strict_analyst(repo: evo.Repository):
    api_key = os.environ.get("GROQ_API_KEY")
    client = Groq(api_key=api_key) if api_key else None

    @evo.track(
        repo=repo,
        name="strict-legal-analyst",
        model="qwen/qwen3.8-27b",
        provider="local",
        temperature=0.1,
    )
    def analyze_dispute(dispute_facts: str) -> dict[str, Any]:
        """Strict Corporate Legal Analyst. Identifies material breaches, statutory liabilities, and immediate legal remedies without compromise."""
        if client:
            chat_completion = client.chat.completions.create(
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are a Senior Corporate Litigation Attorney. "
                            "Analyze the dispute facts strictly under contract law. "
                            "Be concise, aggressive in protecting company rights, identify specific material breaches, "
                            "and demand immediate formal remedy. Use bullet points."
                        ),
                    },
                    {
                        "role": "user",
                        "content": dispute_facts,
                    },
                ],
                model="qwen/qwen3.8-27b",
                temperature=0.1,
                max_tokens=1500,
            )
            return chat_completion

        # Fallback if no key
        return {
            "choices": [{"message": {"content": "MATERIAL BREACH: Section 4.2 violated. Immediate notice of cure required."}}],
            "usage": {"prompt_tokens": 40, "completion_tokens": 30, "total_tokens": 70},
            "model": "llama-3.1-8b-instant",
        }

    return analyze_dispute
