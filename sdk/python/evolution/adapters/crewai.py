"""
CrewAI framework adapter for extracting multi-agent Intelligence Manifests.
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


class CrewAIAdapter(BaseAdapter):
    """Converts CrewAI Crews, Agents, and Tasks into an Intelligence Manifest."""

    @classmethod
    def from_crewai(cls, obj: Any, name: str = "crewai-agent", description: str = "") -> Manifest:
        return cls().to_manifest(obj, name=name, description=description)

    def to_manifest(self, obj: Any, name: str = "crewai-agent", description: str = "") -> Manifest:
        manifest = Manifest(name=name, description=description or "CrewAI multi-agent system state")

        agents = getattr(obj, "agents", []) or []
        if not agents and hasattr(obj, "role"):  # If a single Agent was passed
            agents = [obj]

        seen_models = set()
        seen_tools = set()

        for idx, agent in enumerate(agents):
            role = getattr(agent, "role", f"agent-{idx+1}")
            goal = getattr(agent, "goal", "")
            backstory = getattr(agent, "backstory", "")

            # 1. Agent role prompt artifact
            prompt_content = f"Role: {role}\nGoal: {goal}\nBackstory: {backstory}".strip()
            manifest.add_artifact(PromptArtifact(
                name=f"{role.lower().replace(' ', '-')}-prompt",
                role="system",
                description=prompt_content[:300],
            ))

            # 2. Agent tools
            agent_tools = getattr(agent, "tools", []) or []
            for t in agent_tools:
                t_name = getattr(t, "name", str(t))
                if t_name not in seen_tools:
                    seen_tools.add(t_name)
                    manifest.add_artifact(ToolArtifact(
                        name=t_name,
                        provider="crewai",
                        description=getattr(t, "description", ""),
                    ))

            # 3. Agent LLM
            llm = getattr(agent, "llm", None)
            if llm is not None and str(llm) not in seen_models:
                seen_models.add(str(llm))
                model_str = getattr(llm, "model", getattr(llm, "model_name", str(llm)))
                provider = "openai"
                if "claude" in str(model_str).lower():
                    provider = "anthropic"
                elif "gemini" in str(model_str).lower():
                    provider = "google"

                manifest.add_artifact(ModelConfigArtifact(
                    name=f"model-{role.lower().replace(' ', '-')}",
                    model=str(model_str),
                    provider=provider,
                ))

        # Capture crew structure in metadata
        tasks = getattr(obj, "tasks", []) or []
        manifest.metadata["crewai"] = {
            "agent_count": len(agents),
            "task_count": len(tasks),
            "process": str(getattr(obj, "process", "sequential")),
        }

        return manifest


def from_crewai(obj: Any, name: str = "crewai-multi-agent", description: str = "") -> Manifest:
    """Convenience function to extract an Intelligence Manifest from a CrewAI Crew or Agent."""
    return CrewAIAdapter().to_manifest(obj, name=name, description=description)
