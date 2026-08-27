"""
Direct API adapters for OpenAI and Anthropic API payloads and clients.
"""

from __future__ import annotations

from typing import Any

from evolution.adapters.base import BaseAdapter
from evolution.models.artifacts import (
    ModelConfigArtifact,
    PromptArtifact,
    ToolArtifact,
)
from evolution.models.manifest import Manifest


class OpenAIAdapter(BaseAdapter):
    """Converts OpenAI request dictionaries or client configurations into an Intelligence Manifest."""

    @classmethod
    def from_openai(cls, obj: Any, name: str = "openai-system", description: str = "") -> Manifest:
        return cls().to_manifest(obj, name=name, description=description)

    def to_manifest(self, obj: Any, name: str = "openai-system", description: str = "") -> Manifest:
        manifest = Manifest(name=name, description=description or "OpenAI API intelligence state")
        params = obj if isinstance(obj, dict) else getattr(obj, "__dict__", {})

        # 1. Model configuration
        model = params.get("model", "gpt-4o")
        temp = params.get("temperature", 0.7)
        max_tokens = params.get("max_tokens") or params.get("max_completion_tokens")
        top_p = params.get("top_p")

        manifest.add_artifact(ModelConfigArtifact(
            name=f"openai-{model}",
            model=str(model),
            provider="openai",
            temperature=float(temp) if temp is not None else 0.7,
            max_tokens=int(max_tokens) if max_tokens is not None else None,
            top_p=float(top_p) if top_p is not None else None,
        ))

        # 2. Prompts from messages
        messages = params.get("messages", [])
        for idx, msg in enumerate(messages):
            if isinstance(msg, dict):
                role = msg.get("role", "user")
                content = str(msg.get("content", ""))
                if role in ("system", "user", "assistant", "few_shot"):
                    manifest.add_artifact(PromptArtifact(
                        name=f"openai-msg-{idx+1}-{role}",
                        role=role,
                        description=content[:200],
                    ))

        # 3. Tools
        tools = params.get("tools", [])
        for t in tools:
            if isinstance(t, dict):
                fn = t.get("function", {})
                t_name = fn.get("name") or t.get("name", "custom-tool")
                t_desc = fn.get("description", "")
                manifest.add_artifact(ToolArtifact(
                    name=t_name,
                    provider="openai",
                    description=t_desc,
                ))

        return manifest


class AnthropicAdapter(BaseAdapter):
    """Converts Anthropic request dictionaries or client configurations into an Intelligence Manifest."""

    @classmethod
    def from_anthropic(cls, obj: Any, name: str = "anthropic-system", description: str = "") -> Manifest:
        return cls().to_manifest(obj, name=name, description=description)

    def to_manifest(self, obj: Any, name: str = "anthropic-system", description: str = "") -> Manifest:
        manifest = Manifest(name=name, description=description or "Anthropic API intelligence state")
        params = obj if isinstance(obj, dict) else getattr(obj, "__dict__", {})

        # 1. Model configuration
        model = params.get("model", "claude-3.5-sonnet")
        temp = params.get("temperature", 0.7)
        max_tokens = params.get("max_tokens", 4096)

        manifest.add_artifact(ModelConfigArtifact(
            name=f"anthropic-{model}",
            model=str(model),
            provider="anthropic",
            temperature=float(temp) if temp is not None else 0.7,
            max_tokens=int(max_tokens) if max_tokens is not None else None,
        ))

        # 2. System prompt
        system = params.get("system")
        if system:
            manifest.add_artifact(PromptArtifact(
                name="anthropic-system-prompt",
                role="system",
                description=str(system)[:200],
            ))

        # 3. Tools
        tools = params.get("tools", [])
        for t in tools:
            if isinstance(t, dict):
                t_name = t.get("name", "custom-tool")
                t_desc = t.get("description", "")
                manifest.add_artifact(ToolArtifact(
                    name=t_name,
                    provider="anthropic",
                    description=t_desc,
                ))

        return manifest


def from_openai(params: dict[str, Any] | Any, name: str = "openai-system", description: str = "") -> Manifest:
    """Convenience function to extract an Intelligence Manifest from an OpenAI request payload."""
    return OpenAIAdapter().to_manifest(params, name=name, description=description)


def from_anthropic(params: dict[str, Any] | Any, name: str = "anthropic-system", description: str = "") -> Manifest:
    """Convenience function to extract an Intelligence Manifest from an Anthropic request payload."""
    return AnthropicAdapter().to_manifest(params, name=name, description=description)
