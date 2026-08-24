# Evolution Python SDK (`evolution-sdk`)

> **"Version Intelligence, Not Code."**

The official Python SDK for [Evolution](https://github.com/Urvish0/Evolution) — the AI-native version control platform.

## Features

- 📦 **Manifest Management**: Programmatically generate, inspect, validate, and export `evolution.manifest.json` conforming to **Intelligence Manifest Spec v1.0**.
- 🗃️ **Typed Artifacts**: Native Python classes for Prompts, Memory, Retrieval, Tools, Model Configurations, and Policies with auto-SHA-256 computation.
- 🌳 **Repository Operations**: Init, open, stage, commit, log, diff, and checkout directly from Python code.
- ⚡ **Zero External Dependencies**: Built with Python 3.10+ standard library.

## Quick Start

### Installation

```bash
pip install evolution-sdk
```

### 1. Initialize or Open a Repository

```python
import evolution as evo

# Initialize a new Evolution repository
repo = evo.init("./my-ai-project")

# Or open an existing one
repo = evo.open("./my-ai-project")
```

### 2. Manage the Intelligence Manifest

```python
from evolution import (
    Manifest,
    PromptArtifact,
    ModelConfigArtifact,
    ToolArtifact,
)

# Create a new manifest
manifest = Manifest(name="legal-assistant", description="Legal Research Agent")

# Add AI artifacts
manifest.add_artifact(PromptArtifact(
    name="system-prompt",
    path="prompts/system.txt",
    role="system"
))

manifest.add_artifact(ModelConfigArtifact(
    name="primary-llm",
    model="gpt-4o",
    provider="openai",
    temperature=0.2,
    max_tokens=2048
))

# Save to disk as evolution.manifest.json
manifest.save("./my-ai-project")

# Validate against v1.0 Specification
manifest.validate()
```

### 3. Commit Intelligence State

```python
# Stage and commit your AI configuration
commit = repo.commit(
    message="feat: optimize legal retrieval prompts and model temperature",
    tags=["v1.0-release"],
    metadata={"experiment": "exp-42"}
)

print(f"Committed: {commit.id[:8]} - {commit.message}")
```

### 4. Automatic Intelligence Capture with `@evolution.track`

Decorate any LLM or agent function to automatically capture system prompts from docstrings, model configurations, and execution telemetry:

```python
import evolution as evo

@evo.track(model="gpt-4o", temperature=0.2)
def run_legal_agent(query: str):
    """Specialized AI assistant for analyzing corporate legal liability."""
    # Your LLM call (OpenAI, Anthropic, or custom)
    return client.chat.completions.create(
        model="gpt-4o",
        messages=[{"role": "user", "content": query}]
    )

# Calling the function automatically records inputs, outputs, latency, tokens,
# and attaches prompt/model config to evolution.manifest.json!
response = run_legal_agent("Review this indemnity clause...")
```

### 5. Precision Execution Tracing with Context Manager

```python
import evolution as evo

repo = evo.open("./my-ai-project")

with evo.record(repo, inputs="What is promissory estoppel?") as rec:
    response = client.chat.completions.create(...)
    rec.set_output(response)
    rec.set_metadata("experiment", "temperature-ab-test")
```

### 6. Framework Adapters (LangChain, LlamaIndex, CrewAI, OpenAI, Anthropic)

Convert existing agent and RAG setups directly into version-controlled Intelligence Manifests:

```python
import evolution as evo

# --- LangChain ---
manifest = evo.from_langchain(chain, name="legal-chain")
manifest.save("./my-ai-project")

# --- LlamaIndex ---
manifest = evo.from_llamaindex(index, name="legal-rag")
manifest.save("./my-ai-project")

# --- CrewAI ---
manifest = evo.from_crewai(crew, name="appellate-crew")
manifest.save("./my-ai-project")

# --- Raw OpenAI / Anthropic Payloads ---
manifest = evo.from_openai({"model": "gpt-4o", "messages": [...]})
manifest = evo.from_anthropic({"model": "claude-3.5-sonnet", "system": "..."})
```

## License

MIT License. Developed as part of the [Evolution](https://github.com/Urvish0/Evolution) platform.
