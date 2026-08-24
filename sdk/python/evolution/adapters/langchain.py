"""
LangChain framework adapter for extracting Intelligence Manifests.
"""

from __future__ import annotations

from typing import Any

from evolution.adapters.base import BaseAdapter
from evolution.models.artifacts import (
    MemoryArtifact,
    ModelConfigArtifact,
    PromptArtifact,
    ToolArtifact,
)
from evolution.models.manifest import Manifest


class LangChainAdapter(BaseAdapter):
    """Converts LangChain chains, agents, prompt templates, and tools into an Intelligence Manifest."""

    def to_manifest(self, obj: Any, name: str = "langchain-agent", description: str = "") -> Manifest:
        manifest = Manifest(name=name, description=description or "LangChain application state")

        # 1. Inspect LLM / Model
        llm = getattr(obj, "llm", None) or getattr(obj, "llm_chain", None)
        if llm is None and (hasattr(obj, "model_name") or hasattr(obj, "model")):
            llm = obj

        if llm is not None:
            model_name = getattr(llm, "model_name", None) or getattr(llm, "model", "gpt-4o")
            temperature = getattr(llm, "temperature", 0.7)
            max_tokens = getattr(llm, "max_tokens", None) or getattr(llm, "max_output_tokens", None)

            # Determine provider
            provider = "openai"
            cls_name = llm.__class__.__name__.lower()
            if "anthropic" in cls_name or "claude" in str(model_name).lower():
                provider = "anthropic"
            elif "google" in cls_name or "gemini" in str(model_name).lower():
                provider = "google"
            elif "mistral" in cls_name:
                provider = "mistral"
            elif "cohere" in cls_name:
                provider = "cohere"

            manifest.add_artifact(ModelConfigArtifact(
                name="langchain-llm",
                model=str(model_name),
                provider=provider,
                temperature=float(temperature) if temperature is not None else 0.7,
                max_tokens=int(max_tokens) if max_tokens is not None else None,
            ))

        # 2. Inspect Prompts
        prompt = getattr(obj, "prompt", None) or getattr(obj, "prompt_template", None)
        if prompt is None and hasattr(obj, "messages"):
            prompt = obj

        if prompt is not None:
            prompt_text = ""
            if hasattr(prompt, "template"):
                prompt_text = str(prompt.template)
            elif hasattr(prompt, "messages"):
                msg_strs = []
                for m in prompt.messages:
                    role = getattr(m, "role", "system")
                    content = getattr(m, "content", getattr(m, "template", str(m)))
                    msg_strs.append(f"[{role}] {content}")
                prompt_text = "\n".join(msg_strs)
            else:
                prompt_text = str(prompt)

            manifest.add_artifact(PromptArtifact(
                name="langchain-prompt",
                role="system",
                description=prompt_text[:200] if prompt_text else "LangChain prompt template",
            ))

        # 3. Inspect Tools
        tools = getattr(obj, "tools", []) or []
        for t in tools:
            tool_name = getattr(t, "name", str(t))
            tool_desc = getattr(t, "description", "")
            manifest.add_artifact(ToolArtifact(
                name=tool_name,
                provider="langchain",
                description=tool_desc,
            ))

        # 4. Inspect Memory
        memory = getattr(obj, "memory", None)
        if memory is not None:
            mem_cls = memory.__class__.__name__.lower()
            strategy = "buffer_window"
            if "summary" in mem_cls:
                strategy = "summary"
            elif "vector" in mem_cls:
                strategy = "vector"

            manifest.add_artifact(MemoryArtifact(
                name="langchain-memory",
                strategy=strategy,
                max_tokens=getattr(memory, "max_token_limit", None),
            ))

        return manifest


def from_langchain(obj: Any, name: str = "langchain-agent", description: str = "") -> Manifest:
    """Convenience function to extract an Intelligence Manifest from a LangChain object."""
    return LangChainAdapter().to_manifest(obj, name=name, description=description)
