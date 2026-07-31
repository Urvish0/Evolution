# Intelligence Manifest Specification v0.1 (Draft)

> Defining the open standard for describing the complete operational state of an AI system.

**Status:** Draft  
**Version:** 0.1.0  
**Authors:** Urvish  
**Date:** 2026-07-31  

---

## 1. Purpose

The Intelligence Manifest is a declarative, version-controlled document that captures the **complete operational state** of an AI system at a point in time.

It answers the question: **"What exact configuration of prompts, memory, tools, model settings, and policies produced this AI behavior?"**

### 1.1 — Design Principles

1. **Declarative over Imperative** — The manifest describes *what* the AI system is, not *how* to build it.
2. **Framework-Agnostic** — Works with any AI framework (LangChain, LlamaIndex, custom agents, raw API calls).
3. **Content-Addressed** — Every artifact is referenced by its SHA-256 hash for integrity and deduplication.
4. **Human-Readable** — JSON primary format; easily inspectable without tooling.
5. **Incrementally Adoptable** — Only `version` and `name` are required. Everything else is optional.

---

## 2. Manifest Schema

The root manifest file is `evolution.manifest.json`, stored in the workspace root.

### 2.1 — Root Structure

```json
{
  "version": "0.1.0",
  "name": "legal-research-assistant",
  "description": "AI assistant specializing in legal case research and citation",
  "artifacts": {
    "prompts": [],
    "memory": [],
    "retrieval": [],
    "tools": [],
    "model_config": {},
    "policies": []
  },
  "metadata": {}
}
```

### 2.2 — Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | `string` | ✅ | Spec version this manifest conforms to (semver) |
| `name` | `string` | ✅ | Human-readable name of the intelligence |
| `description` | `string` | No | Brief description of what this AI system does |
| `artifacts` | `object` | No | Collection of typed artifacts (see §3) |
| `metadata` | `object` | No | Execution context and environment info (see §5) |

---

## 3. Artifact Types

Artifacts are the building blocks of an AI system's intelligence. Each artifact has a `type`, a `name`, and is content-addressed by its SHA-256 `hash`.

### 3.1 — Common Artifact Fields

Every artifact shares these base fields:

```json
{
  "type": "<artifact_type>",
  "name": "<human_readable_name>",
  "hash": "<sha256_hash>",
  "path": "<relative_file_path>",
  "description": "<optional_description>"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `string` | ✅ | One of the defined artifact types below |
| `name` | `string` | ✅ | Unique identifier within the manifest |
| `hash` | `string` | ✅ | SHA-256 content hash (auto-computed by Evolution) |
| `path` | `string` | ✅ | Relative file path in the workspace |
| `description` | `string` | No | Human-readable description of the artifact's purpose |

### 3.2 — Prompt

A prompt artifact represents a system prompt, user prompt template, or few-shot example collection.

**Type identifier:** `"prompt"`

```json
{
  "type": "prompt",
  "name": "system-prompt",
  "hash": "b671051cadc3b0806ab998726cb95055a713983b1e94605e9bc3d88689d21a0e",
  "path": "prompts/system.txt",
  "description": "Main system prompt defining assistant persona and behavior",
  "role": "system",
  "format": "text"
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `role` | `string` | Prompt role: `"system"`, `"user"`, `"assistant"`, `"few_shot"` |
| `format` | `string` | Content format: `"text"`, `"template"`, `"jinja2"`, `"mustache"` |

### 3.3 — Memory

A memory artifact represents persistent context, conversation history, or knowledge the AI system retains across interactions.

**Type identifier:** `"memory"`

```json
{
  "type": "memory",
  "name": "conversation-buffer",
  "hash": "a42f891b2c3d...",
  "path": "memory/conversation_buffer.json",
  "description": "Rolling conversation history (last 50 turns)",
  "strategy": "buffer_window",
  "max_tokens": 4096
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `strategy` | `string` | Memory strategy: `"buffer_window"`, `"summary"`, `"vector"`, `"graph"` |
| `max_tokens` | `integer` | Maximum token budget for this memory |

### 3.4 — Retrieval

A retrieval artifact represents a data source the AI system queries for information (RAG configurations, vector stores, search indices).

**Type identifier:** `"retrieval"`

```json
{
  "type": "retrieval",
  "name": "case-law-index",
  "hash": "5f8c3d7e1a2b...",
  "path": "retrieval/case_law_config.json",
  "description": "Vector index over 10k federal case law documents",
  "source": "pinecone",
  "chunk_size": 512,
  "top_k": 5
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `source` | `string` | Backend: `"pinecone"`, `"chroma"`, `"weaviate"`, `"local"`, `"elasticsearch"` |
| `chunk_size` | `integer` | Document chunk size in tokens |
| `top_k` | `integer` | Number of results returned per query |

### 3.5 — Tool Configuration

A tool artifact represents an external capability the AI system can invoke (API calls, function calls, plugins).

**Type identifier:** `"tool"`

```json
{
  "type": "tool",
  "name": "web-search",
  "hash": "f910ab3c8d2e...",
  "path": "tools/web_search.json",
  "description": "Google Search API integration for real-time web queries",
  "provider": "google",
  "auth_required": true
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `provider` | `string` | Tool provider or API service name |
| `auth_required` | `boolean` | Whether the tool requires authentication |

### 3.6 — Model Configuration

The model configuration artifact captures the LLM or AI model settings used by the system.

**Type identifier:** `"model_config"`

Unlike other artifacts, `model_config` is a single object (not an array), because a system runs on one primary model at a time.

```json
{
  "type": "model_config",
  "name": "primary-model",
  "hash": "e3d4f5a6b7c8...",
  "path": "config/model.json",
  "model": "gpt-4o",
  "provider": "openai",
  "temperature": 0.7,
  "max_tokens": 4096,
  "top_p": 1.0
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `model` | `string` | Model identifier (e.g., `"gpt-4o"`, `"claude-3.5-sonnet"`) |
| `provider` | `string` | Provider: `"openai"`, `"anthropic"`, `"google"`, `"local"` |
| `temperature` | `float` | Sampling temperature |
| `max_tokens` | `integer` | Maximum output tokens |
| `top_p` | `float` | Nucleus sampling parameter |

### 3.7 — Policy

A policy artifact represents safety rules, guardrails, output constraints, or behavioral policies governing the AI system.

**Type identifier:** `"policy"`

```json
{
  "type": "policy",
  "name": "content-safety",
  "hash": "c1d2e3f4a5b6...",
  "path": "policies/content_safety.json",
  "description": "Content moderation and safety guardrails",
  "enforcement": "strict"
}
```

| Extended Field | Type | Description |
|----------------|------|-------------|
| `enforcement` | `string` | Enforcement level: `"strict"`, `"warn"`, `"log"` |

---

## 4. Serialization

### 4.1 — Primary Format

**JSON** is the primary serialization format. All manifest files MUST be valid JSON.

File naming convention: `evolution.manifest.json`

### 4.2 — Alternative Format

**YAML** is supported as an alternative for human authoring. Evolution tooling converts YAML to JSON internally.

File naming convention: `evolution.manifest.yaml`

### 4.3 — Encoding

All manifest files MUST be encoded in UTF-8 without BOM.

---

## 5. Metadata

The `metadata` object captures execution context and environment information. All fields are optional.

```json
{
  "metadata": {
    "environment": {
      "os": "windows",
      "evolution_version": "0.3.0",
      "go_version": "1.24.5"
    },
    "tags": ["production", "v2-experiment"],
    "created_by": "Urvish <urvish@example.com>",
    "created_at": "2026-07-31T14:00:00+05:30"
  }
}
```

---

## 6. Spec Versioning

The Intelligence Manifest Specification follows [Semantic Versioning](https://semver.org/):

- **MAJOR** — Breaking changes to required fields or fundamental schema structure.
- **MINOR** — New optional fields, new artifact types, backward-compatible additions.
- **PATCH** — Clarifications, typo fixes, example corrections.

### Version History

| Version | Date | Notes |
|---------|------|-------|
| 0.1.0 | 2026-07-31 | Initial draft. Core schema, 6 artifact types, metadata. |

---

## 7. Example Manifests

### 7.1 — Simple Chatbot

```json
{
  "version": "0.1.0",
  "name": "support-chatbot",
  "description": "Customer support chatbot for SaaS product",
  "artifacts": {
    "prompts": [
      {
        "type": "prompt",
        "name": "system-prompt",
        "hash": "a1b2c3d4...",
        "path": "prompts/system.txt",
        "role": "system",
        "format": "text"
      }
    ],
    "model_config": {
      "type": "model_config",
      "name": "primary-model",
      "hash": "e5f6a7b8...",
      "path": "config/model.json",
      "model": "gpt-4o-mini",
      "provider": "openai",
      "temperature": 0.3,
      "max_tokens": 1024
    }
  }
}
```

### 7.2 — RAG Application

```json
{
  "version": "0.1.0",
  "name": "legal-research-assistant",
  "description": "AI assistant that retrieves and cites federal case law",
  "artifacts": {
    "prompts": [
      {
        "type": "prompt",
        "name": "system-prompt",
        "hash": "b671051c...",
        "path": "prompts/system.txt",
        "role": "system",
        "format": "text"
      },
      {
        "type": "prompt",
        "name": "citation-template",
        "hash": "c3d4e5f6...",
        "path": "prompts/citation_template.txt",
        "role": "user",
        "format": "template"
      }
    ],
    "retrieval": [
      {
        "type": "retrieval",
        "name": "case-law-index",
        "hash": "5f8c3d7e...",
        "path": "retrieval/case_law.json",
        "source": "pinecone",
        "chunk_size": 512,
        "top_k": 5
      }
    ],
    "model_config": {
      "type": "model_config",
      "name": "primary-model",
      "hash": "d4e5f6a7...",
      "path": "config/model.json",
      "model": "gpt-4o",
      "provider": "openai",
      "temperature": 0.2,
      "max_tokens": 4096
    },
    "policies": [
      {
        "type": "policy",
        "name": "citation-required",
        "hash": "f8a9b0c1...",
        "path": "policies/citation_policy.json",
        "enforcement": "strict"
      }
    ]
  }
}
```

### 7.3 — Autonomous Agent

```json
{
  "version": "0.1.0",
  "name": "research-agent",
  "description": "Autonomous agent that plans, searches, and synthesizes research reports",
  "artifacts": {
    "prompts": [
      {
        "type": "prompt",
        "name": "planner-prompt",
        "hash": "a1a1a1a1...",
        "path": "prompts/planner.txt",
        "role": "system",
        "format": "text"
      },
      {
        "type": "prompt",
        "name": "synthesizer-prompt",
        "hash": "b2b2b2b2...",
        "path": "prompts/synthesizer.txt",
        "role": "system",
        "format": "text"
      }
    ],
    "memory": [
      {
        "type": "memory",
        "name": "task-memory",
        "hash": "c3c3c3c3...",
        "path": "memory/task_buffer.json",
        "strategy": "summary",
        "max_tokens": 8192
      }
    ],
    "tools": [
      {
        "type": "tool",
        "name": "web-search",
        "hash": "d4d4d4d4...",
        "path": "tools/search.json",
        "provider": "google",
        "auth_required": true
      },
      {
        "type": "tool",
        "name": "code-executor",
        "hash": "e5e5e5e5...",
        "path": "tools/code_exec.json",
        "provider": "local",
        "auth_required": false
      }
    ],
    "model_config": {
      "type": "model_config",
      "name": "primary-model",
      "hash": "f6f6f6f6...",
      "path": "config/model.json",
      "model": "claude-3.5-sonnet",
      "provider": "anthropic",
      "temperature": 0.5,
      "max_tokens": 8192
    },
    "policies": [
      {
        "type": "policy",
        "name": "safety-guardrails",
        "hash": "a7a7a7a7...",
        "path": "policies/safety.json",
        "enforcement": "strict"
      },
      {
        "type": "policy",
        "name": "cost-limits",
        "hash": "b8b8b8b8...",
        "path": "policies/cost.json",
        "enforcement": "warn"
      }
    ]
  },
  "metadata": {
    "environment": {
      "os": "linux",
      "evolution_version": "0.4.0"
    },
    "tags": ["autonomous", "research", "multi-step"]
  }
}
```

---

## 8. Future Considerations (v0.2+)

The following are out of scope for v0.1 but planned for future versions:

- **Evaluation artifacts** — Test suites, benchmarks, and golden datasets.
- **Deployment artifacts** — Endpoint configurations, scaling rules, A/B test splits.
- **Dependency graph** — Explicit relationships between artifacts (e.g., "this prompt requires this retrieval source").
- **Inheritance** — Manifests that extend or override a base manifest.
- **Signing** — Cryptographic signatures for manifest integrity verification.
- **Remote references** — Artifacts stored externally (S3, GCS) referenced by URL + hash.
