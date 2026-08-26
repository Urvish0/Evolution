# Evolution Framework Integration Specification v1.0 (Stable)

> Defining the open standard and protocols for integrating AI frameworks with the Evolution Intelligence Manifest.

**Status:** Stable  
**Version:** 1.0.0  
**Authors:** Urvish  
**Date:** 2026-08-26  
**Related Specs:** [Intelligence Manifest Specification v1.0](intelligence-manifest-v1.0.md)  

---

## 1. Purpose & Scope

The **Framework Integration Specification** defines how AI orchestration frameworks, agent libraries, and LLM runtime systems interface with Evolution to capture, version, and reconstitute AI operational state.

### 1.1 — The Core Problem

Modern AI frameworks (such as LangChain, LlamaIndex, CrewAI, AutoGen, DSPy, and Haystack) maintain internal abstractions for prompts, tools, memory, retrievers, and model parameters. However:
- Every framework represents these components with proprietary, incompatible in-memory schemas.
- When an AI pipeline changes, there is no standardized way to serialize *what* changed at the intelligence level.
- Code version control (Git) tracks text diffs in source code, but cannot isolate whether a behavior shift was caused by a temperature tweak, a prompt template revision, or a tool schema modification.

### 1.2 — Mission of this Specification

This specification provides:
1. **Normalized Primitive Mapping** — Standard mapping rules from framework-specific entities to Evolution's 6 typed artifacts (`prompt`, `memory`, `retrieval`, `tool`, `model_config`, `policy`).
2. **Standard Adapter Lifecycle** — The 4-phase lifecycle for framework adapters: Introspection, Mapping, Telemetry Hooking, and Reconstitution.
3. **Bidirectional Protocol** — Formal interface requirements for **Export** (`Framework State → Manifest`) and **Import** (`Manifest → Framework State`).
4. **Third-Party Developer Guide** — Architectural rules and compliance criteria for creating new community adapters.

---

## 2. Universal Framework Primitive Mapping

Every adapter MUST translate framework-specific primitives into one of the six core artifact types defined in Intelligence Manifest Spec v1.0.

| Framework Primitive | Evolution Artifact Type | Key Attributes Captured | Example Sources |
|---|---|---|---|
| System/User Prompt Templates | `prompt` | `role`, `format`, `path`, `hash`, `description` | LangChain `PromptTemplate`, DSPy Signatures, Raw Messages |
| Context & Conversation Memory | `memory` | `strategy`, `max_tokens`, `path`, `hash` | `ConversationBufferMemory`, `SummaryMemory`, LangGraph Checkpointers |
| Vector Stores, Indices, RAG | `retrieval` | `source`, `chunk_size`, `top_k`, `path`, `hash` | LlamaIndex `VectorStoreIndex`, Pinecone, Chroma, Weaviate |
| Function Calling & External APIs | `tool` | `provider`, `auth_required`, `path`, `hash` | LangChain `BaseTool`, CrewAI Tools, OpenAI Functions |
| Model IDs & Inference Hyperparams | `model_config` | `model`, `provider`, `temperature`, `max_tokens`, `top_p` | `ChatOpenAI`, `ChatAnthropic`, Ollama, Bedrock |
| Guardrails, Moderation & Policies | `policy` | `enforcement` (`strict`/`warn`/`log`), `path`, `hash` | NeMo Guardrails, LlamaGuard, Custom Safety Rules |

---

## 3. The 4-Phase Adapter Lifecycle

An Evolution Framework Adapter operates across four distinct phases:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      EVOLUTION ADAPTER LIFECYCLE                         │
└─────────────────────────────────────────────────────────────────────────┘
                                   │
               ┌───────────────────▼───────────────────┐
               │        1. INTROSPECTION PHASE         │
               │ Deep inspection of runtime object     │
               │ tree using safe duck-typing.          │
               └───────────────────┬───────────────────┘
                                   │
               ┌───────────────────▼───────────────────┐
               │           2. MAPPING PHASE            │
               │ Normalization into 6 typed Artifacts  │
               │ + Git-compatible SHA-256 auto-hashing.│
               └───────────────────┬───────────────────┘
                                   │
               ┌───────────────────▼───────────────────┐
               │       3. TELEMETRY HOOK PHASE         │
               │ Intercept latency, tokens, inputs,    │
               │ and outputs into Execution records.   │
               └───────────────────┬───────────────────┘
                                   │
               ┌───────────────────▼───────────────────┐
               │       4. RECONSTITUTION PHASE         │
               │ Rehydrate runtime agents from a       │
               │ committed Intelligence Manifest.      │
               └───────────────────────────────────────┘
```

### 3.1 — Phase 1: Introspection
The adapter defensively inspects the framework object tree without assuming hard package imports. It queries standard attributes (`llm`, `prompt`, `tools`, `retriever`, `memory`, `storage_context`, `agents`, `tasks`).

### 3.2 — Phase 2: Mapping & Auto-Hashing
The adapter builds an `evolution.Manifest` conforming to Spec v1.0. If artifact content is located on the filesystem, the adapter computes `SHA-256("blob <len>\0<content>")` for content-addressed deduplication.

### 3.3 — Phase 3: Telemetry Hooking
The adapter provides execution recording hooks. When an invocation runs, the adapter measures elapsed duration via high-resolution monotonic clocks (`time.perf_counter()`), parses token usage (`prompt_tokens`, `completion_tokens`, `total_tokens`), and writes the execution JSON to `.evolution/executions/<id>.json`.

### 3.4 — Phase 4: Reconstitution (Import)
The adapter parses a saved `evolution.manifest.json` and instantiates corresponding framework runtime classes (e.g., creating a `ChatOpenAI` instance with the recorded model and temperature).

---

## 4. Standard Adapter Protocol (Interface Definition)

Every framework adapter MUST conform to the `BaseAdapter` contract.

### 4.1 — Python Interface Protocol

```python
from abc import ABC, abstractmethod
from typing import Any
from evolution.models.manifest import Manifest

class BaseAdapter(ABC):
    """Abstract interface for all framework adapters."""

    @abstractmethod
    def to_manifest(
        self,
        obj: Any,
        name: str = "ai-intelligence",
        description: str = "",
    ) -> Manifest:
        """Exports a framework-specific object into an Evolution Manifest.
        
        Args:
            obj: Framework runtime object, chain, index, crew, or dictionary payload.
            name: Name of the intelligence system.
            description: Optional human-readable description.
            
        Returns:
            Manifest: Validated Evolution Intelligence Manifest conforming to Spec v1.0.
        """
        pass

    def from_manifest(self, manifest: Manifest) -> Any:
        """Optional: Reconstitutes a framework object from a Manifest.
        
        Args:
            manifest: Evolution Intelligence Manifest.
            
        Returns:
            Framework-specific runtime object or configuration dictionary.
        """
        raise NotImplementedError("Reconstitution is not implemented for this adapter.")
```

---

## 5. Reference Implementations

Evolution provides canonical reference adapters across the primary AI archetypes:

### 5.1 — Orchestration Pattern: `LangChainAdapter`

- **Target Entities**: `Chain`, `AgentExecutor`, `RunnableSequence`, `PromptTemplate`.
- **Extraction Rules**:
  - LLM model & temperature $\rightarrow$ `ModelConfigArtifact`
  - Prompt templates $\rightarrow$ `PromptArtifact`
  - Tool list $\rightarrow$ `ToolArtifact` collection
  - Memory buffers $\rightarrow$ `MemoryArtifact`

```python
import evolution as evo
# Export LangChain chain to Manifest
manifest = evo.from_langchain(legal_agent_chain, name="legal-chain")
manifest.validate()
```

### 5.2 — Retrieval/RAG Pattern: `LlamaIndexAdapter`

- **Target Entities**: `VectorStoreIndex`, `QueryEngine`, `Retriever`, `StorageContext`.
- **Extraction Rules**:
  - Vector store backend (`chroma`, `pinecone`, `weaviate`, `local`) $\rightarrow$ `RetrievalArtifact.source`
  - Chunk size & overlap $\rightarrow$ `RetrievalArtifact.chunk_size`
  - Similarity top-k $\rightarrow$ `RetrievalArtifact.top_k`
  - Embedded LLM settings $\rightarrow$ `ModelConfigArtifact`

```python
import evolution as evo
# Export LlamaIndex index to Manifest
manifest = evo.from_llamaindex(contract_index, name="contract-rag")
manifest.validate()
```

### 5.3 — Multi-Agent Swarm Pattern: `CrewAIAdapter`

- **Target Entities**: `Crew`, `Agent`, `Task`.
- **Extraction Rules**:
  - Each Agent role, goal, and backstory $\rightarrow$ distinct `PromptArtifact`
  - Agent tools $\rightarrow$ deduplicated `ToolArtifact` collection
  - Agent LLMs $\rightarrow$ `ModelConfigArtifact`
  - Crew process (`sequential`/`hierarchical`) & agent counts $\rightarrow$ `metadata.crewai`

```python
import evolution as evo
# Export multi-agent Crew to Manifest
manifest = evo.from_crewai(appellate_crew, name="appellate-crew")
manifest.validate()
```

### 5.4 — Direct Provider Pattern: `OpenAIAdapter` & `AnthropicAdapter`

- **Target Entities**: Direct API request payloads or SDK client objects.
- **Extraction Rules**:
  - `model`, `temperature`, `max_tokens` $\rightarrow$ `ModelConfigArtifact`
  - `messages` / `system` $\rightarrow$ `PromptArtifact`
  - `tools` / `functions` $\rightarrow$ `ToolArtifact`

```python
import evolution as evo
manifest = evo.from_openai({"model": "gpt-4o", "messages": [...]})
manifest = evo.from_anthropic({"model": "claude-3.5-sonnet", "system": "..."})
```

---

## 6. Third-Party Contributor Guide

Community contributors developing new adapters (e.g., for DSPy, AutoGen, Semantic Kernel, Haystack, or custom in-house agent loops) MUST adhere to the following rules:

### 6.1 — Architectural Guardrails

1. **Zero Hard Runtime Dependencies**:
   - Adapters MUST NOT add heavy third-party framework libraries to `dependencies` in `pyproject.toml`.
   - Use safe introspection (`getattr`, `hasattr`, duck-typing) or optional imports (`try: import ... except ImportError:`).
2. **Fail Gracefully**:
   - If an unrecognized object is passed, the adapter MUST extract whatever attributes are discernable and populate generic fallbacks rather than crashing.
3. **Spec v1.0 Compliance**:
   - All generated manifests MUST pass `manifest.validate()` (valid semver, non-empty name, allowed artifact enum types).
4. **Deterministic Auto-Hashing**:
   - If file paths are referenced, hashes MUST use the standard Git blob format: `SHA-256("blob <len>\0<content>")`.

### 6.2 — Adapter Compliance Checklist

Before submitting a new framework adapter to the Evolution repository:

- [ ] Adapter inherits from `BaseAdapter` and implements `to_manifest()`.
- [ ] Export convenience function is provided (e.g., `from_autogen()`, `from_dspy()`).
- [ ] Function exported in `evolution/adapters/__init__.py` and top-level `evolution/__init__.py`.
- [ ] Unit tests added in `tests/test_adapters.py` with mock objects covering all extracted artifact types.
- [ ] All tests pass without requiring third-party framework packages to be pre-installed.
- [ ] Documentation updated in `DOCUMENTATION.md` and `sdk/python/README.md`.

---

## 7. Versioning & Evolution

This specification follows **Semantic Versioning**:
- **MAJOR (v1.0 $\rightarrow$ v2.0)**: Breaking changes to the adapter interface or artifact mapping requirements.
- **MINOR (v1.0 $\rightarrow$ v1.1)**: Introduction of new supported framework primitives or optional lifecycle hooks.
- **PATCH (v1.0 $\rightarrow$ v1.0.1)**: Documentation fixes, mapping clarifications, and example additions.
