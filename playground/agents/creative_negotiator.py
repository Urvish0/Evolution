"""
Agent 2: Strategic Commercial Negotiator
- Model: llama-3.1-8b-instant
- Temperature: 0.7 (Creative, exploratory)
- Persona: Win-win mediation, relationship preservation, commercial restructuring.
"""

from __future__ import annotations

import os
from typing import Any
from groq import Groq
import evolution as evo


def create_creative_negotiator(repo: evo.Repository):
    api_key = os.environ.get("GROQ_API_KEY")
    client = Groq(api_key=api_key) if api_key else None

    @evo.track(
        repo=repo,
        name="creative-negotiator",
        model="qwen/qwen3.6-27b",
        provider="local",
        temperature=0.7,
    )
    def negotiate_dispute(dispute_facts: str) -> dict[str, Any]:
        """Strategic Commercial Negotiator. Preserves business partnerships, proposes win-win restructuring, and resolves disputes without litigation."""
        if client:
            chat_completion = client.chat.completions.create(
                messages=[
                    {
                        "role": "system",
                        "content": (
                            "You are a Master Commercial Negotiator and Mediator. "
                            "Analyze the dispute with empathy and commercial pragmatism. "
                            "Propose creative win-win solutions that preserve the business relationship, "
                            "avoid costly litigation, and restructure the agreement constructively."
                        ),
                    },
                    {
                        "role": "user",
                        "content": dispute_facts,
                    },
                ],
                model="qwen/qwen3.6-27b",
                temperature=0.7,
                max_tokens=1500,
            )
            return chat_completion

        # Fallback if no key
        return {
            "choices": [{"message": {"content": "COMMERCIAL PROPOSAL: Restructure delivery milestones and offer 10% future credit to preserve relationship."}}],
            "usage": {"prompt_tokens": 40, "completion_tokens": 35, "total_tokens": 75},
            "model": "llama-3.1-8b-instant",
        }

    return negotiate_dispute
