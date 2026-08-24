"""
Evolution SDK Exceptions.
"""


class EvolutionError(Exception):
    """Base exception for all Evolution SDK errors."""
    pass


class RepositoryNotFoundError(EvolutionError):
    """Raised when an Evolution repository cannot be found at the specified path."""
    pass


class RepositoryAlreadyExistsError(EvolutionError):
    """Raised when attempting to initialize a repository where one already exists."""
    pass


class ManifestNotFoundError(EvolutionError):
    """Raised when an evolution.manifest.json file is missing."""
    pass


class ManifestValidationError(EvolutionError):
    """Raised when an evolution manifest fails specification compliance validation."""

    def __init__(self, message: str, errors: list[str] | None = None):
        super().__init__(message)
        self.errors = errors or []

    def __str__(self) -> str:
        if not self.errors:
            return super().__str__()
        formatted = "\n  - ".join(self.errors)
        return f"{super().__str__()}:\n  - {formatted}"


class ArtifactNotFoundError(EvolutionError):
    """Raised when an artifact cannot be found or its underlying file is missing."""
    pass


class CommandExecutionError(EvolutionError):
    """Raised when an underlying Evolution CLI or repository command fails."""

    def __init__(self, command: str, exit_code: int, stderr: str):
        self.command = command
        self.exit_code = exit_code
        self.stderr = stderr.strip()
        super().__init__(f"Command '{command}' failed with exit code {exit_code}: {self.stderr}")
