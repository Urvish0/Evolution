# Evolution Intelligence Manifest Specification

> The open standard for describing the complete operational state of an AI system.

## Overview

The Intelligence Manifest Specification defines a declarative, version-controlled format for capturing **what** an AI system is at any point in time — its prompts, memory, retrieval sources, tools, model configuration, policies, execution recordings, and evaluation scores.

**Current Version:** 1.0.0 (Stable)

## Quick Start

Create an `evolution.manifest.json` in your project root:

```json
{
  "version": "1.0.0",
  "name": "my-ai-assistant"
}
```

That's it. Only `version` and `name` are required. Everything else is optional and incrementally adoptable.

## Specification Documents

| Version | Status | Document |
|---------|--------|----------|
| **1.0.0** | **Stable** | [intelligence-manifest-v1.0.md](intelligence-manifest-v1.0.md) |
| 0.2.0 | Superseded | [intelligence-manifest-v0.2.md](intelligence-manifest-v0.2.md) |
| 0.1.0 | Superseded | [intelligence-manifest-v0.1.md](intelligence-manifest-v0.1.md) |

## JSON Schema

A machine-readable JSON Schema for automated validation is available at:

```
spec/schema/evolution-manifest.schema.json
```

### Usage with VS Code

Add the following to your workspace `.vscode/settings.json`:

```json
{
  "json.schemas": [
    {
      "fileMatch": ["evolution.manifest.json"],
      "url": "./spec/schema/evolution-manifest.schema.json"
    }
  ]
}
```

### Usage with JSON Schema Validators

```bash
# Example: validate using ajv-cli
npx ajv-cli validate -s spec/schema/evolution-manifest.schema.json -d evolution.manifest.json
```

## Key Concepts

### Artifact Types

The spec defines six artifact types representing distinct components of an AI system:

| Type | Purpose |
|------|---------|
| `prompt` | System, user, and few-shot prompt templates |
| `memory` | Conversation memory configurations |
| `retrieval` | Vector store and retrieval source configurations |
| `tool` | External tool and API integrations |
| `model_config` | LLM provider, model, and parameter settings |
| `policy` | Safety, compliance, and behavioral constraints |

### Execution Recording *(v1.0)*

Captures individual AI system invocations with inputs, outputs, latency, and token usage.

### Evaluation Scoring *(v1.0)*

Scores executions against pluggable evaluators (performance, cost, safety, correctness) with normalized 0.0–1.0 scores.

### Regression Rules *(v1.0)*

CI/CD quality gates: `min_score`, `max_drop`, and `require_safety` constraints that fail builds on quality regression.

## Design Principles

1. **Declarative over Imperative** — Describes *what*, not *how*.
2. **Framework-Agnostic** — Works with LangChain, LlamaIndex, custom agents, raw API calls.
3. **Content-Addressed** — SHA-256 hashes for integrity and deduplication.
4. **Human-Readable** — JSON format, inspectable without tooling.
5. **Incrementally Adoptable** — Start with two fields, add more as needed.
6. **Auto-Computed Integrity** — Hashes computed at commit time if omitted.
7. **Observable** — Execution and evaluation are first-class citizens.

## Contributing

This specification is developed as part of the [Evolution](https://github.com/Urvish0/Evolution) project. For questions, feedback, or proposals:

1. Open an issue on the Evolution repository.
2. Reference the spec version you're discussing.
3. For breaking changes, propose via RFC with migration path.

## License

This specification is released under the same license as the Evolution project.
