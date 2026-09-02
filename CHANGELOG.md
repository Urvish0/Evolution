# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.8.5] - 2026-08-27

### Added
- **Python CLI (`evo-py` / `python -m evolution`)**: Pure-Python CLI entry point allowing developers to manage repositories, inspect manifests, list executions, and view evaluations without requiring Go installed.
- **LLM-as-a-Judge Semantic Evaluation (`SemanticEvaluator`)**: Implemented structured multi-dimensional evaluation engine (Accuracy, Helpfulness, Instruction Following, Safety) with qualitative reasoning and weighted scoring.
- **Live 4-Agent Benchmark Testbed (`playground/full_spectrum_demo.py`)**: End-to-end sandbox runner demonstrating 4 specialized agents on Groq LPU with side-by-side judging.
- **Vector Workflow Diagram (`docs/assets/workflow.svg`)**: Dark-mode terminal visualization showing the 4-step workflow: Capture $\rightarrow$ Commit $\rightarrow$ Diff $\rightarrow$ Evaluate.
- **SPDX MIT License & PyPI Landing Page**: Added `sdk/python/LICENSE` and `sdk/python/README.md`.

### Changed
- Refactored framework adapters to provide `@classmethod from_*` constructors (`LangChainAdapter`, `LlamaIndexAdapter`, `CrewAIAdapter`, `OpenAIAdapter`, `AnthropicAdapter`).
- Cleaned and standardized all CLI outputs to follow an engineering-first format.

---

## [0.8.0] - 2026-08-27

### Added
- **Official PyPI Release (`evolution-sdk`)**: Initial release on the Python Package Index ([pypi.org/project/evolution-sdk/](https://pypi.org/project/evolution-sdk/)) with **zero runtime dependencies**.
- **Automatic Intelligence Capture (`@evo.track`)**: Non-invasive decorator extracting docstring prompts, model parameters, execution tokens, and hardware latency.
- **Universal Framework Adapters**: Zero-dependency introspection for LangChain, LlamaIndex, CrewAI, OpenAI, and Anthropic.
- **Precision Telemetry (`evo.record`)**: Context manager for fine-grained execution latency and token logging.
- **Open Specifications (v1.0 Stable)**:
  - Intelligence Manifest Specification v1.0 (`evolution.manifest.json`).
  - Framework Integration Specification v1.0.
  - Formal JSON Schema validation.

---

## [0.7.0] - 2026-08-15

### Added
- **Core Engine (Go)**:
  - Content-Addressable Storage (CAS) with type-prefixed SHA-256 (`blob <size>\0<content>`).
  - Merkle Tree directory snapshots with root tree integrity verification.
  - Directed Acyclic Graph (DAG) commit engine.
  - Three-Way Branch Merging with Lowest Common Ancestor (LCA) ancestor detection.
  - Longest Common Subsequence (LCS) unified line-diff engine.
  - Replay engine for historical AI state reconstruction.
  - Pluggable quality evaluation gates with CI regression thresholds.
- Comprehensive Go CLI command suite (`evo init`, `evo commit`, `evo log`, `evo diff`, `evo branch`, `evo checkout`, `evo merge`, `evo replay`, `evo evaluate`).
