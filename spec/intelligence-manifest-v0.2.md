# Intelligence Manifest Specification v0.2 (Validated)

> Defining the open standard for describing the complete operational state of an AI system.

**Status:** Validated  
**Version:** 0.2.0  
**Authors:** Urvish  
**Date:** 2026-08-06  
**Supersedes:** v0.1.0  

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
6. **Auto-Computed Integrity** — If `hash` is omitted or empty, Evolution auto-computes it from file content at commit time. *(New in v0.2)*

### 1.2 — Changes from v0.1

| Change | v0.1 | v0.2 |
|--------|------|------|
| `hash` field | Required, must be pre-computed | Optional at authoring time; auto-computed at commit time |
| `metadata` schema | Free-form object | Structured: `environment` map + `tags` array |
| Commit metadata | Undefined | Auto-captures OS, architecture, Go version, custom key-value pairs |
| Validation rules | Informal | Formal JSON Schema + programmatic validation |
| `model_config` cardinality | Single object only | Single object (unchanged, but clarified in schema) |
| Spec status | Draft | Validated (implementation-proven) |

---

## 2. Manifest Schema

The root manifest file is `evolution.manifest.json`, stored in the workspace root.

### 2.1 — Root Structure

```json
{
  "version": "0.2.0",
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
  "metadata": {
    "environment": {},
    "tags": []
  }
}
```

### 2.2 — Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | `string` | ✅ | Spec version this manifest conforms to (semver) |
| `name` | `string` | ✅ | Human-readable name of the intelligence |
| `description` | `string` | No | Brief description of what this AI system does |
| `artifacts` | `object` | No | Collection of typed artifacts (see §3) |
| `metadata` | `object` | No | Structured execution context (see §5) |

---

## 3. Artifact Types

Artifacts are the building blocks of an AI system's intelligence. Each artifact has a `type`, a `name`, and is content-addressed by its SHA-256 `hash`.

### 3.1 — Common Artifact Fields

Every artifact shares these base fields:

```json
{
  "type": "<artifact_type>",
  "name": "<human_readable_name>",
  "hash": "<sha256_hash_or_empty>",
  "path": "<relative_file_path>",
  "description": "<optional_description>"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `string` | ✅ | One of: `prompt`, `memory`, `retrieval`, `tool`, `model_config`, `policy` |
| `name` | `string` | ✅ | Unique identifier within the manifest |
| `hash` | `string` | No | SHA-256 content hash. If empty, auto-computed from `path` at commit time |
| `path` | `string` | ✅ | Relative file path in the workspace |
| `description` | `string` | No | Human-readable description of the artifact's purpose |

> **v0.2 Change:** `hash` is no longer required at authoring time. When `evo commit` runs, Evolution reads the file at `path`, computes `SHA-256("blob <len>\0<content>")`, and stores the hash. This means developers can add artifacts without manually computing hashes.

### 3.2 — Prompt

A prompt artifact represents a system prompt, user prompt template, or few-shot example collection.

**Type identifier:** `"prompt"`

```json
{
  "type": "prompt",
  "name": "system-prompt",
  "hash": "",
  "path": "prompts/system.txt",
  "description": "Main system prompt defining assistant persona and behavior",
  "role": "system",
  "format": "text"
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `role` | `string` | No | Prompt role: `"system"`, `"user"`, `"assistant"`, `"few_shot"` |
| `format` | `string` | No | Content format: `"text"`, `"template"`, `"jinja2"`, `"mustache"` |

### 3.3 — Memory

A memory artifact represents persistent context, conversation history, or knowledge the AI system retains across interactions.

**Type identifier:** `"memory"`

```json
{
  "type": "memory",
  "name": "conversation-buffer",
  "hash": "",
  "path": "memory/conversation_buffer.json",
  "description": "Rolling conversation history (last 50 turns)",
  "strategy": "buffer_window",
  "max_tokens": 4096
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `strategy` | `string` | No | Memory strategy: `"buffer_window"`, `"summary"`, `"vector"`, `"graph"` |
| `max_tokens` | `integer` | No | Maximum token budget for this memory |

### 3.4 — Retrieval

A retrieval artifact represents a data source the AI system queries for information (RAG configurations, vector stores, search indices).

**Type identifier:** `"retrieval"`

```json
{
  "type": "retrieval",
  "name": "case-law-index",
  "hash": "",
  "path": "retrieval/case_law_config.json",
  "description": "Vector index over 10k federal case law documents",
  "source": "pinecone",
  "chunk_size": 512,
  "top_k": 5
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `source` | `string` | No | Backend: `"pinecone"`, `"chroma"`, `"weaviate"`, `"local"`, `"elasticsearch"` |
| `chunk_size` | `integer` | No | Document chunk size in tokens |
| `top_k` | `integer` | No | Number of results returned per query |

### 3.5 — Tool Configuration

A tool artifact represents an external capability the AI system can invoke (API calls, function calls, plugins).

**Type identifier:** `"tool"`

```json
{
  "type": "tool",
  "name": "web-search",
  "hash": "",
  "path": "tools/web_search.json",
  "description": "Google Search API integration for real-time web queries",
  "provider": "google",
  "auth_required": true
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `provider` | `string` | No | Tool provider or API service name |
| `auth_required` | `boolean` | No | Whether the tool requires authentication (default: `false`) |

### 3.6 — Model Configuration

The model configuration artifact captures the LLM or AI model settings used by the system.

**Type identifier:** `"model_config"`

Unlike other artifacts, `model_config` is a **single object** (not an array), because a system runs on one primary model configuration at a time.

```json
{
  "type": "model_config",
  "name": "primary-model",
  "hash": "",
  "path": "config/model.json",
  "model": "gpt-4o",
  "provider": "openai",
  "temperature": 0.7,
  "max_tokens": 4096,
  "top_p": 1.0
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `model` | `string` | ✅ | Model identifier (e.g., `"gpt-4o"`, `"claude-3.5-sonnet"`) |
| `provider` | `string` | No | Provider: `"openai"`, `"anthropic"`, `"google"`, `"local"` |
| `temperature` | `float` | No | Sampling temperature (default: `0.7`) |
| `max_tokens` | `integer` | No | Maximum output tokens (default: `4096`) |
| `top_p` | `float` | No | Nucleus sampling parameter (default: `1.0`) |

### 3.7 — Policy

A policy artifact represents safety rules, guardrails, output constraints, or behavioral policies governing the AI system.

**Type identifier:** `"policy"`

```json
{
  "type": "policy",
  "name": "content-safety",
  "hash": "",
  "path": "policies/content_safety.json",
  "description": "Content moderation and safety guardrails",
  "enforcement": "strict"
}
```

| Extended Field | Type | Required | Description |
|----------------|------|----------|-------------|
| `enforcement` | `string` | No | Enforcement level: `"strict"`, `"warn"`, `"log"` (default: `"strict"`) |

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

The `metadata` object captures execution context and environment information.

### 5.1 — Structured Schema *(New in v0.2)*

```json
{
  "metadata": {
    "environment": {
      "os": "windows",
      "arch": "amd64",
      "go_version": "go1.24.5",
      "evolution_version": "0.4.0"
    },
    "tags": ["production", "v1.0"]
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `metadata.environment` | `map[string]string` | Key-value pairs for runtime context. Auto-populated: `os`, `arch`, `go_version`. Custom values via `--meta key=value`. |
| `metadata.tags` | `string[]` | Categorical labels for commits. Set via `--tag <label>`. |

### 5.2 — Commit-Level Metadata

When `evo commit` runs, Evolution automatically captures:

| Key | Source | Example |
|-----|--------|---------|
| `os` | `runtime.GOOS` | `"windows"`, `"linux"`, `"darwin"` |
| `arch` | `runtime.GOARCH` | `"amd64"`, `"arm64"` |
| `go_version` | `runtime.Version()` | `"go1.24.5"` |
| Custom keys | `--meta key=value` | `"model=gpt-4o"`, `"temp=0.2"` |

---

## 6. Validation Rules *(New in v0.2)*

Evolution validates manifests programmatically using `ValidateManifest()`:

### 6.1 — Required Field Validation

- `version` MUST be a non-empty string conforming to semver.
- `name` MUST be a non-empty string.

### 6.2 — Artifact Type Validation

- `artifacts.*.type` MUST be one of: `prompt`, `memory`, `retrieval`, `tool`, `model_config`, `policy`.
- If `artifacts.model_config` is present, `model_config.model` MUST be a non-empty string.

### 6.3 — Hash Auto-Computation

- If `hash` is empty (`""`) and `path` points to an existing file, Evolution computes `SHA-256("blob <len>\0<content>")` at commit time.
- If `path` does not exist, `hash` remains empty and the artifact is still attached to the commit.

### 6.4 — Formal JSON Schema

A machine-readable JSON Schema is available at `spec/schema/evolution-manifest.schema.json` for third-party tooling integration.

---

## 7. Spec Versioning

The Intelligence Manifest Specification follows [Semantic Versioning](https://semver.org/):

- **MAJOR** — Breaking changes to required fields or fundamental schema structure.
- **MINOR** — New optional fields, new artifact types, backward-compatible additions.
- **PATCH** — Clarifications, typo fixes, example corrections.

### Version History

| Version | Date | Status | Notes |
|---------|------|--------|-------|
| 0.1.0 | 2026-07-31 | Superseded | Initial draft. Core schema, 6 artifact types, metadata. |
| 0.2.0 | 2026-08-06 | Validated | Auto-hash computation, structured metadata, formal JSON Schema, validation rules. |

---

## 8. Future Considerations (v0.3+)

The following are out of scope for v0.2 but planned for future versions:

- **Evaluation artifacts** — Test suites, benchmarks, and golden datasets.
- **Deployment artifacts** — Endpoint configurations, scaling rules, A/B test splits.
- **Dependency graph** — Explicit relationships between artifacts (e.g., "this prompt requires this retrieval source").
- **Inheritance** — Manifests that extend or override a base manifest.
- **Signing** — Cryptographic signatures for manifest integrity verification.
- **Remote references** — Artifacts stored externally (S3, GCS) referenced by URL + hash.
- **Multi-model support** — Multiple `model_config` entries for multi-agent systems.
