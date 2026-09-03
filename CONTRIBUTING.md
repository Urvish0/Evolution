# Contributing to Evolution

Thank you for your interest in contributing to Evolution. This document provides guidelines and information for contributors.

---

## How to Contribute

### Reporting Bugs

Open a [GitHub Issue](https://github.com/Urvish0/Evolution/issues/new) with:
- A clear, descriptive title.
- Steps to reproduce the behavior.
- Expected vs. actual behavior.
- Your environment: OS, Go version, Python version, `evolution-sdk` version.

### Suggesting Features

Open a [GitHub Issue](https://github.com/Urvish0/Evolution/issues/new) describing:
- The problem your feature solves.
- How it relates to versioning AI intelligence (prompts, models, memory, tools, policies).
- Any proposed API or CLI interface.

### Submitting Code

1. **Fork** the repository and create your branch from `main`.
2. **Write tests** for any new functionality.
3. **Follow existing code style** (see Standards below).
4. **Open a Pull Request** with a clear description of your changes.

---

## Development Setup

### Go Core Engine

```bash
# Clone the repository
git clone https://github.com/Urvish0/Evolution.git
cd Evolution

# Build the CLI
go build -o evo.exe ./cmd/evo

# Run Go tests (53 tests)
go test -v ./...
```

### Python SDK

```bash
# Create a virtual environment
python -m venv .venv
.venv\Scripts\activate    # Windows
source .venv/bin/activate # macOS/Linux

# Install in editable mode with dev dependencies
pip install -e sdk/python[dev]

# Run Python tests (31 tests)
pytest sdk/python/tests -v
```

---

## Code Standards

### Go

- Follow standard Go conventions (`gofmt`, `go vet`).
- Every exported function must have a doc comment.
- Error wrapping with context: `fmt.Errorf("operation: %w", err)`.
- Small, focused functions organized by domain.

### Python

- Type hints on all function signatures (Python 3.10+ syntax).
- Docstrings on all public functions and classes.
- Zero external runtime dependencies in the core SDK. Test dependencies go in `[project.optional-dependencies]`.
- Follow PEP 8 conventions.

---

## Project Structure

```
Evolution/
├── cmd/evo/                 # Go CLI entry point (main.go)
├── internal/
│   ├── cli/                 # Cobra command handlers
│   ├── repository/          # Core versioning engine (CAS, Merkle, Merge, Diff)
│   └── version/             # Build version metadata
├── sdk/python/
│   ├── evolution/           # Python SDK source
│   │   ├── adapters/        # Framework adapters (LangChain, LlamaIndex, CrewAI)
│   │   ├── capture/         # @track decorator, introspection, recorder
│   │   ├── models/          # Data models (Manifest, Execution, Evaluation)
│   │   ├── cli.py           # evo-py CLI implementation
│   │   ├── evaluators.py    # LLM-as-a-Judge semantic evaluator
│   │   └── repository.py    # Repository operations
│   └── tests/               # Python unit tests
├── spec/                    # Intelligence Manifest Specification v1.0
├── docs/                    # Documentation and guides
├── examples/                # Runnable example scripts
└── playground/              # Live multi-agent testing lab
```

---

## Areas for Contribution

### High-Impact Areas

- **New Framework Adapters**: Add support for additional AI frameworks (AutoGen, DSPy, Haystack). See `sdk/python/evolution/adapters/` for the pattern.
- **CLI Commands**: Extend `evo-py` with additional commands. See `sdk/python/evolution/cli.py`.
- **Evaluation Dimensions**: Add new scoring dimensions to the `SemanticEvaluator` (e.g., coherence, citation accuracy).
- **Documentation**: Improve tutorials, examples, and inline docstrings.

### Writing a Framework Adapter

Framework adapters live in `sdk/python/evolution/adapters/`. Each adapter:

1. Extends `BaseAdapter` from `evolution.adapters.base`.
2. Implements `to_manifest(obj, name) -> Manifest` to convert a framework object into an Intelligence Manifest.
3. Uses duck-typing introspection (no external imports) to extract prompts, model configs, and tools.
4. Provides both a module-level function (`from_X(obj)`) and a classmethod (`XAdapter.from_X(obj)`).

Example skeleton:

```python
from evolution.adapters.base import BaseAdapter
from evolution.models.manifest import Manifest

class MyFrameworkAdapter(BaseAdapter):
    """Adapter for MyFramework pipelines."""

    @classmethod
    def from_my_framework(cls, pipeline, name: str = "my-pipeline") -> Manifest:
        manifest = Manifest(name=name)
        # Extract prompts, model config, tools from pipeline
        # Add them to manifest.artifacts
        return manifest

def from_my_framework(pipeline, name: str = "my-pipeline") -> Manifest:
    return MyFrameworkAdapter.from_my_framework(pipeline, name=name)
```

---

## Running the Full Test Suite

```bash
# Go tests (53 tests)
go test -v ./...

# Python tests (31 tests)
pytest sdk/python/tests -v

# Total: 84 automated tests
```

All tests must pass before submitting a pull request.

---

## License

By contributing to Evolution, you agree that your contributions will be licensed under the [MIT License](LICENSE).
