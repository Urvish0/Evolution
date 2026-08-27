# Evolution Python SDK (`evolution-sdk`)

> **"Version Intelligence, Not Code."**

[![PyPI Version](https://img.shields.io/badge/pypi-v0.8.0-blue.svg)](https://pypi.org/project/evolution-sdk/)
[![Python Versions](https://img.shields.io/badge/python-3.10%20%7C%203.11%20%7C%203.12%20%7C%203.13-blue)](https://pypi.org/project/evolution-sdk/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/Urvish0/Evolution/blob/main/LICENSE)
[![Zero Dependencies](https://img.shields.io/badge/dependencies-0%20runtime-brightgreen.svg)](#features)

The official Python SDK for **[Evolution](https://github.com/Urvish0/Evolution)** — an open-source, AI-native version control platform.

While traditional version control systems (like Git) track line-by-line changes to text files, **Evolution tracks and versions Intelligence**: the complete operational state of an AI system (system prompts, sampling parameters, memory strategies, vector store retrieval configurations, tool schemas, and guardrails).

---

## Key Features

- **Zero Runtime Dependencies**: Pure Python implementation with zero third-party dependencies. Ultra-fast, lightweight, and won't conflict with your environment.
- **Automatic Intelligence Capture**: Use the `@evolution.track` decorator to automatically extract docstring prompts, model hyperparameters, token consumption, and latency.
- **Precision Execution Telemetry**: Record execution traces, inputs, outputs, and hardware latency using `@track` or the `evolution.record` context manager.
- **Universal Framework Adapters**: Seamlessly introspect and export standard Intelligence Manifests from **LangChain**, **LlamaIndex**, **CrewAI**, and direct OpenAI / Anthropic client calls.
- **Git-Compatible CAS Blobs**: Content-Addressable Storage with auto-computed SHA-256 Merkle hashes.
- **Intelligence Manifest Spec v1.0 Compliant**: Native support for creating, validating, and saving standard `evolution.manifest.json` schemas.

---

## Installation

Install `evolution-sdk` via `pip` or `uv`:

```bash
pip install evolution-sdk
```

Or using `uv`:

```bash
uv add evolution-sdk
```

---

## Quickstart in 30 Seconds

### 1. Initialize or Open an Evolution Repository

```python
import evolution as evo

# Initialize a new Evolution repository in the current directory
repo = evo.Repository.init("./my_ai_agent")

# Or open an existing repository
# repo = evo.Repository.open("./my_ai_agent")
```

---

### 2. Automatic Intelligence Tracking with `@evo.track`

Decorate your AI agent functions. Evolution automatically extracts the docstring as the system prompt, registers model configs, measures latency, and logs execution tokens without modifying your business logic:

```python
from openai import OpenAI
import evolution as evo

client = OpenAI()
repo = evo.Repository.init("./my_agent")

@evo.track(
    repo=repo,
    name="customer-support-agent",
    model="gpt-4o",
    temperature=0.2,
)
def handle_customer_inquiry(user_message: str):
    """Senior Customer Support Specialist. Always maintain professional tone and provide step-by-step guidance."""
    response = client.chat.completions.create(
        model="gpt-4o",
        temperature=0.2,
        messages=[
            {"role": "system", "content": "You are a Senior Customer Support Specialist."},
            {"role": "user", "content": user_message}
        ]
    )
    return response

# Run your agent
result = handle_customer_inquiry("How do I reset my API key?")

# Commit the captured intelligence snapshot
repo.commit("feat: initial customer support intelligence snapshot")
```

---

### 3. Precision Latency & Execution Recording with `record()`

For fine-grained telemetry or manual pipelines, use the `record` context manager:

```python
import evolution as evo

repo = evo.Repository.open(".")

with evo.record(repo, metadata={"phase": "testing"}) as ctx:
    # Set inputs
    ctx.set_inputs("User asked: Explain Merkle trees.")
    
    # Execute your LLM call
    response_text = "A Merkle tree is a cryptographic tree structure..."
    
    # Record outputs and token metrics
    ctx.set_outputs(response_text)
    ctx.set_tokens(prompt_tokens=42, completion_tokens=128)

print(f"Recorded execution with duration: {ctx.execution.duration_ms} ms")
```

---

## 🔌 Framework Adapters

Evolution provides plug-and-play adapters to export standard Intelligence Manifests directly from your existing AI framework abstractions:

### LangChain Adapter
```python
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate
from evolution.adapters import LangChainAdapter

prompt = ChatPromptTemplate.from_template("Summarize the following legal document: {doc}")
llm = ChatOpenAI(model="gpt-4o", temperature=0.1)
chain = prompt | llm

manifest = LangChainAdapter.from_langchain(chain, name="legal-summarizer")
print(manifest.to_dict())
```

### LlamaIndex Adapter
```python
from llama_index.core import VectorStoreIndex
from evolution.adapters import LlamaIndexAdapter

manifest = LlamaIndexAdapter.from_llamaindex(
    index=my_vector_index,
    name="financial-rag-agent",
    similarity_top_k=5,
)
```

### CrewAI Adapter
```python
from crewai import Agent, Crew, Task
from evolution.adapters import CrewAIAdapter

manifest = CrewAIAdapter.from_crewai(my_crew, name="research-crew")
```

---

## 📋 The Intelligence Manifest Format (`evolution.manifest.json`)

Evolution stores your complete AI system state in a portable, framework-agnostic manifest following the **[Intelligence Manifest Specification v1.0](https://github.com/Urvish0/Evolution/blob/main/spec/intelligence-manifest-v1.0.md)**:

```json
{
  "version": "1.0.0",
  "name": "legal-dispute-resolution-system",
  "description": "AI system powered by Evolution version control",
  "artifacts": {
    "prompts": [
      {
        "type": "prompt",
        "name": "strict-legal-analyst-prompt",
        "description": "Extracted docstring prompt",
        "role": "system",
        "format": "text"
      }
    ],
    "model_config": {
      "type": "model_config",
      "name": "analyst-model",
      "model": "qwen/qwen3.8-27b",
      "provider": "local",
      "temperature": 0.1
    }
  }
}
```

---

## Running Tests

To run the SDK test suite:

```bash
pytest sdk/python/tests/ -v
```

---

## License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

## Links

- **GitHub Repository:** [https://github.com/Urvish0/Evolution](https://github.com/Urvish0/Evolution)
- **Specification:** [Intelligence Manifest Spec v1.0](https://github.com/Urvish0/Evolution/blob/main/spec/intelligence-manifest-v1.0.md)
- **Documentation:** [Technical Manual](https://github.com/Urvish0/Evolution/blob/main/DOCUMENTATION.md)
