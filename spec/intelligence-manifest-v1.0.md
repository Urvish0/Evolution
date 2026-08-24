# Intelligence Manifest Specification v1.0 (Stable)

> Defining the open standard for describing the complete operational state of an AI system.

**Status:** Stable  
**Version:** 1.0.0  
**Authors:** Urvish  
**Date:** 2026-08-24  
**Supersedes:** v0.2.0  

---

## 1. Purpose

The Intelligence Manifest is a declarative, version-controlled document that captures the **complete operational state** of an AI system at a point in time.

It answers the question: **"What exact configuration of prompts, memory, tools, model settings, policies, executions, and evaluations produced this AI behavior?"**

### 1.1 — Design Principles

1. **Declarative over Imperative** — The manifest describes *what* the AI system is, not *how* to build it.
2. **Framework-Agnostic** — Works with any AI framework (LangChain, LlamaIndex, custom agents, raw API calls).
3. **Content-Addressed** — Every artifact is referenced by its SHA-256 hash for integrity and deduplication.
4. **Human-Readable** — JSON primary format; easily inspectable without tooling.
5. **Incrementally Adoptable** — Only `version` and `name` are required. Everything else is optional.
6. **Auto-Computed Integrity** — If `hash` is omitted or empty, Evolution auto-computes it from file content at commit time.
7. **Observable** — Execution recordings and evaluation scores are first-class citizens. *(New in v1.0)*

### 1.2 — Changes from v0.2

| Change | v0.2 | v1.0 |
|--------|------|------|
| Execution schema | Undefined | Formal `execution` object with inputs, outputs, tokens, duration |
| Evaluation schema | Undefined | Formal `evaluation` object with scored evaluators and regression rules |
| Spec status | Validated | Stable |
| `model_config.provider` enum | 4 values | Extended with `mistral`, `cohere`, `aws_bedrock` |
| Regression rules | Undefined | `regression_rules` object with `min_score`, `max_drop`, `require_safety` |

---

## 2. Manifest Schema

The root manifest file is `evolution.manifest.json`, stored in the workspace root.

### 2.1 — Root Structure

```json
{
  "version": "1.0.0",
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

### 2.2 — Required Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | `string` | ✅ | Semver string (e.g., `"1.0.0"`) |
| `name` | `string` | ✅ | Human-readable intelligence name |
| `description` | `string` | ❌ | Brief description |
| `artifacts` | `object` | ❌ | Collection of typed AI artifacts |
| `metadata` | `object` | ❌ | Execution context and environment metadata |

---

## 3. Artifact Types

Evolution recognizes six artifact types, each representing a distinct component of an AI system's operational state.

### 3.1 — Prompt Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"prompt"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `path` | `string` | ✅ | File path relative to workspace root |
| `hash` | `string` | ❌ | SHA-256 content hash (auto-computed if empty) |
| `role` | `"system" \| "user" \| "assistant" \| "few_shot"` | ❌ | Prompt role in conversation |
| `format` | `"text" \| "template" \| "jinja2" \| "mustache"` | ❌ | Template format |

### 3.2 — Memory Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"memory"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `path` | `string` | ✅ | File path |
| `strategy` | `"buffer_window" \| "summary" \| "vector" \| "graph"` | ❌ | Memory strategy |
| `max_tokens` | `integer` | ❌ | Token budget for memory window |

### 3.3 — Retrieval Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"retrieval"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `path` | `string` | ✅ | File path |
| `source` | `"pinecone" \| "chroma" \| "weaviate" \| "local" \| "elasticsearch"` | ❌ | Vector store provider |
| `chunk_size` | `integer` | ❌ | Document chunk size |
| `top_k` | `integer` | ❌ | Number of results to retrieve |

### 3.4 — Tool Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"tool"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `path` | `string` | ✅ | File path |
| `provider` | `string` | ❌ | Tool provider name |
| `auth_required` | `boolean` | ❌ | Whether authentication is required |

### 3.5 — Model Config Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"model_config"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `model` | `string` | ✅ | Model identifier (e.g., `"gpt-4o"`, `"claude-3.5-sonnet"`) |
| `provider` | `"openai" \| "anthropic" \| "google" \| "local" \| "mistral" \| "cohere" \| "aws_bedrock"` | ❌ | LLM provider |
| `temperature` | `number` (0–2) | ❌ | Sampling temperature |
| `max_tokens` | `integer` | ❌ | Maximum output tokens |
| `top_p` | `number` (0–1) | ❌ | Nucleus sampling parameter |

### 3.6 — Policy Artifact

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | `"policy"` | ✅ | Artifact type discriminator |
| `name` | `string` | ✅ | Artifact identifier |
| `path` | `string` | ✅ | File path |
| `enforcement` | `"strict" \| "warn" \| "log"` | ❌ | Enforcement level |

---

## 4. Execution Schema *(New in v1.0)*

An **Execution** records a single AI system invocation at a specific commit snapshot.

### 4.1 — Execution Object

```json
{
  "id": "92c91ed0-f152-44ed-88e0-daf5c27cf6e2",
  "commit_id": "58e45843-c4e7-447c-8194-86b7bc16130d",
  "inputs": "What is corporate law?",
  "outputs": "Corporate law governs company formation and governance.",
  "duration_ms": 280,
  "tokens": {
    "prompt_tokens": 45,
    "completion_tokens": 30,
    "total_tokens": 75
  },
  "timestamp": "2026-08-10T19:41:21+05:30",
  "metadata": {}
}
```

### 4.2 — Execution Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | `string` (UUID) | ✅ | Unique execution identifier |
| `commit_id` | `string` (UUID) | ✅ | Intelligence Commit this execution is bound to |
| `inputs` | `string` | ✅ | User input or query |
| `outputs` | `string` | ✅ | AI-generated response |
| `duration_ms` | `integer` | ✅ | Response latency in milliseconds |
| `tokens` | `object` | ✅ | Token consumption metrics |
| `tokens.prompt_tokens` | `integer` | ✅ | Input token count |
| `tokens.completion_tokens` | `integer` | ✅ | Output token count |
| `tokens.total_tokens` | `integer` | ✅ | Total token count |
| `timestamp` | `string` (RFC 3339) | ✅ | When the execution occurred |
| `metadata` | `object` | ❌ | Custom key-value execution metadata |

### 4.3 — Storage

Executions are stored as individual JSON files in `.evolution/executions/<execution_id>.json`.

---

## 5. Evaluation Schema *(New in v1.0)*

An **Evaluation** scores an execution against one or more quality evaluators.

### 5.1 — Evaluation Result Object

```json
{
  "id": "a9b6cc69-5a67-466d-b98f-bb5ddffb12c1",
  "commit_id": "58e45843-c4e7-447c-8194-86b7bc16130d",
  "execution_id": "92c91ed0-f152-44ed-88e0-daf5c27cf6e2",
  "overall_score": 1.0,
  "scores": {
    "performance": {
      "name": "performance",
      "score": 1.0,
      "unit": "latency_score",
      "details": "Latency: 280ms (Target: 1000ms)"
    },
    "cost": {
      "name": "cost",
      "score": 1.0,
      "unit": "USD",
      "details": "Est. Cost: $0.000412 (75 tokens)"
    },
    "safety": {
      "name": "safety",
      "score": 1.0,
      "unit": "pass_rate",
      "details": "Passed safety guardrails"
    },
    "correctness": {
      "name": "correctness",
      "score": 1.0,
      "unit": "quality_score",
      "details": "Output response generated successfully"
    }
  },
  "timestamp": "2026-08-10T20:02:16+05:30"
}
```

### 5.2 — Evaluation Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | `string` (UUID) | ✅ | Unique evaluation identifier |
| `commit_id` | `string` (UUID) | ✅ | Intelligence Commit this evaluation belongs to |
| `execution_id` | `string` (UUID) | ✅ | Execution being evaluated |
| `overall_score` | `number` (0.0–1.0) | ✅ | Aggregate normalized quality score |
| `scores` | `object` | ✅ | Map of evaluator name → score detail |
| `timestamp` | `string` (RFC 3339) | ✅ | When the evaluation was generated |

### 5.3 — Evaluation Score Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | ✅ | Evaluator name (e.g., `"performance"`) |
| `score` | `number` (0.0–1.0) | ✅ | Normalized score |
| `unit` | `string` | ❌ | Measurement unit |
| `details` | `string` | ❌ | Human-readable detail |

### 5.4 — Built-in Evaluators

| Evaluator | Measures | Score Logic |
|-----------|----------|-------------|
| `performance` | Response latency vs SLA target | `1.0` if within target, proportional decay beyond |
| `cost` | Token cost estimation (USD) | `1.0` for < $0.001, `0.5` for < $0.01, `0.2` beyond |
| `safety` | Forbidden keyword detection in outputs | `1.0` if clean, `0.0` if violations detected |
| `correctness` | Output completeness and golden match | `1.0` for exact match, `0.75` for partial, `0.0` for empty |

### 5.5 — Storage

Evaluations are stored as individual JSON files in `.evolution/evaluations/<evaluation_id>.json`.

---

## 6. Regression Rules *(New in v1.0)*

Regression rules define quality gate constraints for CI/CD integration.

### 6.1 — Regression Rule Object

```json
{
  "min_score": 0.80,
  "max_drop": 0.05,
  "require_safety": true
}
```

| Field | Type | Description |
|-------|------|-------------|
| `min_score` | `number` (0.0–1.0) | Fail if overall quality score is below this threshold |
| `max_drop` | `number` (0.0–1.0) | Fail if score drop between baseline and candidate exceeds this ratio |
| `require_safety` | `boolean` | Fail if safety evaluator score is not `1.0` |

### 6.2 — CI/CD Integration

When regression rules are violated, `evo evaluate` exits with **status code 1**, blocking CI/CD pipelines.

```yaml
# GitHub Actions example
- name: AI Quality Gate
  run: evo evaluate --fail-under 0.80 --require-safety
```

---

## 7. Metadata Schema

### 7.1 — Manifest Metadata

```json
{
  "metadata": {
    "environment": {
      "os": "windows",
      "arch": "amd64",
      "go_version": "go1.24.5"
    },
    "tags": ["production", "v2.1"]
  }
}
```

When `evo commit` runs, Evolution automatically captures:

| Key | Source | Example |
|-----|--------|---------|
| `os` | `runtime.GOOS` | `"windows"`, `"linux"`, `"darwin"` |
| `arch` | `runtime.GOARCH` | `"amd64"`, `"arm64"` |
| `go_version` | `runtime.Version()` | `"go1.24.5"` |
| Custom keys | `--meta key=value` | `"model=gpt-4o"`, `"temp=0.2"` |

---

## 8. Validation Rules

Evolution validates manifests programmatically using `ValidateManifest()`:

### 8.1 — Required Field Validation

- `version` MUST be a non-empty string conforming to semver.
- `name` MUST be a non-empty string.

### 8.2 — Artifact Type Validation

- `artifacts.*.type` MUST be one of: `prompt`, `memory`, `retrieval`, `tool`, `model_config`, `policy`.
- If `artifacts.model_config` is present, `model_config.model` MUST be a non-empty string.

### 8.3 — Hash Auto-Computation

- If `hash` is empty (`""`) and `path` points to an existing file, Evolution computes `SHA-256("blob <len>\0<content>")` at commit time.
- If `path` does not exist, `hash` remains empty and the artifact is still attached to the commit.

### 8.4 — Formal JSON Schema

A machine-readable JSON Schema is available at `spec/schema/evolution-manifest.schema.json` for third-party tooling integration.

---

## 9. Spec Versioning

The Intelligence Manifest Specification follows [Semantic Versioning](https://semver.org/):

- **MAJOR** — Breaking changes to required fields or fundamental schema structure.
- **MINOR** — New optional fields, new artifact types, backward-compatible additions.
- **PATCH** — Clarifications, typo fixes, example corrections.

### Version History

| Version | Date | Status | Notes |
|---------|------|--------|-------|
| 0.1.0 | 2026-07-31 | Superseded | Initial draft. Core schema, 6 artifact types, metadata. |
| 0.2.0 | 2026-08-06 | Superseded | Auto-hash computation, structured metadata, formal JSON Schema, validation rules. |
| 1.0.0 | 2026-08-24 | **Stable** | Execution schema, evaluation schema, regression rules, extended provider enum, CI/CD integration. |

---

## 10. Future Considerations (v1.1+)

The following are out of scope for v1.0 but planned for future versions:

- **Deployment artifacts** — Endpoint configurations, scaling rules, A/B test splits.
- **Dependency graph** — Explicit relationships between artifacts.
- **Inheritance** — Manifests that extend or override a base manifest.
- **Signing** — Cryptographic signatures for manifest integrity verification.
- **Remote references** — Artifacts stored externally (S3, GCS) referenced by URL + hash.
- **Multi-model support** — Multiple `model_config` entries for multi-agent systems.
- **Python SDK auto-capture** — Automatic manifest generation from Python frameworks.
