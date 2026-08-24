"""
Manifest validator conforming to Intelligence Manifest Specification v1.0.
"""

from __future__ import annotations

import re
from typing import TYPE_CHECKING

from evolution.exceptions import ManifestValidationError

if TYPE_CHECKING:
    from evolution.models.manifest import Manifest

SEMVER_PATTERN = re.compile(r"^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$")
VALID_ARTIFACT_TYPES = {"prompt", "memory", "retrieval", "tool", "model_config", "policy"}


def validate_manifest(manifest: Manifest) -> None:
    """Validates that a Manifest instance strictly adheres to Spec v1.0.
    Raises ManifestValidationError if any violations are found.
    """
    errors: list[str] = []

    # 1. Version validation
    if not manifest.version or not manifest.version.strip():
        errors.append("manifest 'version' is required and cannot be empty")
    elif not SEMVER_PATTERN.match(manifest.version):
        errors.append(f"manifest 'version' '{manifest.version}' is not valid semantic versioning (e.g. 1.0.0)")

    # 2. Name validation
    if not manifest.name or not manifest.name.strip():
        errors.append("manifest 'name' is required and cannot be empty")

    # 3. Artifact validations
    for art in manifest.artifacts.all():
        if not art.type:
            errors.append(f"artifact '{art.name}' has missing type")
        elif art.type not in VALID_ARTIFACT_TYPES:
            errors.append(f"artifact '{art.name}' has invalid type '{art.type}' (must be one of {sorted(VALID_ARTIFACT_TYPES)})")

        if not art.name or not art.name.strip():
            errors.append(f"found artifact of type '{art.type}' with empty name")

    # 4. Model Config specific validation
    if manifest.artifacts.model_config:
        mc = manifest.artifacts.model_config
        if not mc.model or not mc.model.strip():
            errors.append("model_config artifact requires a non-empty 'model' field")

    if errors:
        raise ManifestValidationError(
            f"Manifest validation failed with {len(errors)} error(s)",
            errors=errors,
        )
