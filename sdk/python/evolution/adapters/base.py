"""
Base protocol and abstract class for framework adapters.
"""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any

from evolution.models.manifest import Manifest


class BaseAdapter(ABC):
    """Abstract base class for converting framework-specific AI configurations
    into standard Intelligence Manifests (Spec v1.0).
    """

    @abstractmethod
    def to_manifest(self, obj: Any, name: str = "ai-intelligence", description: str = "") -> Manifest:
        """Converts a framework-specific object into an Evolution Manifest."""
        pass
