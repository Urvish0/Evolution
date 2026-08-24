# Evolution

> Version intelligence, not code.

[![CI](https://github.com/Urvish0/Evolution/actions/workflows/ci.yml/badge.svg)](https://github.com/Urvish0/Evolution/actions/workflows/ci.yml)

Evolution is an AI engineering platform that introduces version control for intelligence.

Instead of versioning only source code, Evolution versions the entire lifecycle of an AI system — prompts, memory, retrieval, tool usage, evaluations, policies, and deployments.

---

## Quick Start

```bash
git clone https://github.com/Urvish0/Evolution.git
cd Evolution

go run ./cmd/evo version
go run ./cmd/evo init
go run ./cmd/evo status
go run ./cmd/evo commit -m "initial intelligence"
go run ./cmd/evo log
```

---

## Problem

AI systems are driven by dynamic components — prompts, memory, retrieval, models, tools, and policies — that evolve independently and influence behavior unpredictably.

Teams struggle to:

- Understand why an agent behaves differently
- Identify which change introduced a regression
- Reproduce previous AI behavior
- Compare intelligence across versions
- Deploy AI updates with confidence

Evolution exists to reduce that uncertainty.

---

## Core Philosophy

Evolution does **not** version prompts, memories, or workflows individually.

Evolution versions **Intelligence** — the complete operational state of an AI system at a point in time.

Everything else is a versioned artifact within that intelligence.

---

## Features

| Feature | Description | Status |
|---------|-------------|--------|
| `evo init` | Initialize an Evolution repository | ✅ |
| `evo status` | Show repository status and branch | ✅ |
| `evo add` | Stage files or directories for commit | ✅ |
| `evo commit` | Create an Intelligence Commit | ✅ |
| `evo log` | View commit history (with `--oneline`, `-n`, colors) | ✅ |
| `evo config` | Manage user identity (`user.name`, `user.email`) | ✅ |
| `evo branch` | List, create (`-n`), and delete (`-d`) branches | ✅ |
| `evo checkout` | Switch active branch | ✅ |
| `evo version` | Display CLI version | ✅ |
| `evo diff` | Compare working tree against HEAD | ✅ |
| `evo manifest` | Manage `evolution.manifest.json` (`init`, `validate`, `show`) | ✅ |
| `evo artifact` | Manage AI artifacts (`add`, `list`, `show`, `diff`) | ✅ |
| `evo restore` | Restore files or working tree from HEAD commit | ✅ |
| `evo merge` | Three-way branch merging with conflict marker injection | ✅ |
| `evo record` | Record an AI execution run bound to HEAD commit | ✅ |
| `evo execution` | List and inspect recorded AI executions | ✅ |
| `evo replay` | Reconstruct and replay historical AI states | ✅ |
| `evo evaluate` | Run quality evaluations (Performance, Cost, Safety, Correctness) | ✅ |

---

## Architecture

```
Developer → CLI / SDK → Core Engine → Evolution Repository → Replay & Evaluation
```

**Tech Stack:** Go, Cobra CLI, JSON serialization, local filesystem storage.

See [EVOLUTION_MASTER_PLAN.md](EVOLUTION_MASTER_PLAN.md) for the complete roadmap, architecture, and phase details.

---

## Project Structure

```
Evolution/
├── cmd/evo/             # CLI entry point
├── internal/
│   ├── cli/             # Cobra command handlers
│   ├── repository/      # Core versioning engine
│   └── version/         # Version metadata
├── spec/                # Intelligence Manifest Specification
├── docs/                # Design documentation
├── .agents/             # AI collaboration rules
├── DOCUMENTATION.md     # Official Technical Manual & Technical Reference
├── USAGE.md             # Complete CLI usage guide
└── EVOLUTION_MASTER_PLAN.md  # Single source of truth
```

---

## Documentation

| Document | Purpose |
|----------|---------|
| [Official Manual](DOCUMENTATION.md) | Complete technical manual, architecture, CLI reference, and end-to-end tutorials |
| [CLI Usage Guide](USAGE.md) | Quick reference for CLI commands and workflows |
| [Master Plan](EVOLUTION_MASTER_PLAN.md) | Roadmap, progress, architecture — single source of truth |
| [PRD](docs/PRD.md) | Product requirements and vision |
| [Domain Model](docs/DOMAIN_MODEL.md) | Core entities and relationships |
| [Architecture Principles](docs/ARCHITECTURE_PRINCIPLES.md) | Design constraints |
| [CLI Spec](docs/CLI_SPEC.md) | Command definitions |
| [ADRs](docs/adr/) | Architecture Decision Records |
| [Intelligence Manifest Spec v1.0](spec/intelligence-manifest-v1.0.md) | Open standard for describing AI system state (stable) |
| [JSON Schema](spec/schema/evolution-manifest.schema.json) | Formal JSON Schema for manifest validation (v1.0) |

---

## Current Status

> **Phase 7 — Evaluation Engine & Spec v1.0** (complete)

AI evaluation framework with pluggable evaluators (performance, cost, safety, correctness), cross-commit evaluation comparison, regression detection with CI/CD quality gates (`--fail-under`, `--max-drop`, `--require-safety`), and Intelligence Manifest Specification v1.0 (Stable) with execution and evaluation schemas. Ready for v0.7.0 tag.

---

## License

Apache-2.0
