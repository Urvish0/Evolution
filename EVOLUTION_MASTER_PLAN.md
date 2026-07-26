# Evolution — Master Plan

> **Single Source of Truth**
>
> This document tracks every completed milestone, every known issue, and every future phase of Evolution.
> It replaces the need to cross-reference ROADMAP.md, IMPLEMENTATION_BLUEPRINT.md, CHANGELOG.md, or the Founder's Playbook for "what's done" and "what's next."
>
> Update this file whenever work is completed, a decision changes, or a new phase begins.

---

## Long-Term Goal

**Product:** Build the engineering platform that AI development is missing — a system that versions, replays, evaluates, and evolves the complete intelligence of AI systems, not just their source code.

**Career:** Establish Urvish as a top-tier AI Product & Systems Engineer by shipping a production-grade, open-source platform that demonstrates system design, distributed thinking, Go mastery, and product leadership.

**North Star:** *"Version Intelligence, Not Code."*

---

## Architecture Principles (Non-Negotiable)

| # | Principle | Meaning |
|---|-----------|---------|
| 1 | Local First | Works fully offline; cloud is optional |
| 2 | Intelligence is the Unit of Versioning | We version complete AI state, not individual prompts |
| 3 | Immutable History | Commits are never modified; new commits supersede old ones |
| 4 | Framework Agnostic | Evolution integrates with any AI framework; it replaces none |
| 5 | API First | CLI, SDK, and GUI all build on the same core interfaces |
| 6 | Extensible by Design | New artifact types, storage backends, and integrations without core changes |
| 7 | Deterministic Replay | Reproduce previous intelligence when sufficient context exists |
| 8 | Simplicity Over Complexity | No abstraction without measurable value |

---

## Core Domain Model

```
Workspace
  └── Repository
        └── Branch
              └── Intelligence Commit (immutable snapshot)
                    └── Artifacts (prompt, memory, retrieval, tools, model config, policies)
                          └── Execution (recorded run)
                                └── Evaluation (quality, safety, cost, performance)
                                      └── Deployment (promote to environment)
```

**Domain Rules:**
- Intelligence Commits are immutable.
- Every Artifact belongs to exactly one Intelligence Commit.
- Every Execution references exactly one Intelligence Commit.
- Every Evaluation references exactly one Execution.
- Branches maintain independent version histories.

---

## Technology Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Language | Go 1.26 | Cross-platform, single binary, excellent CLI ecosystem |
| CLI Framework | Cobra | Command hierarchy, help generation, flags |
| Serialization | JSON | Portable, human-readable, language-agnostic |
| Configuration | JSON (future: YAML) | Simple for current phase |
| Storage | Local filesystem | Local-first principle; pluggable backends later |
| Logging | Go slog | Standard library, structured |
| Testing | Go testing | Standard library |
| IDs | UUID v4 (google/uuid) | Globally unique commit identifiers |
| CI/CD | GitHub Actions (planned) | Build, test, lint |
| Linting | golangci-lint (planned) | Static analysis |

---

## Repository Layout (On-Disk)

```
.evolution/
├── config.json          # Repository metadata (version, id, created_at, default_branch)
├── HEAD                 # Current active branch name (e.g. "main")
├── branches/            # Branch JSON files mapping name → head commit ID
├── commits/             # Individual commit JSON files (id.json)
├── objects/             # (Future) Content-addressable blob/tree storage
└── artifacts/           # (Future) Versioned AI artifact storage
```

**Planned additions from REPOSITORY_FORMAT.md:**
- `executions/` — Execution history
- `evaluations/` — Evaluation results
- `refs/` — Branch and tag references

---

## Code Architecture

```
Evolution/
├── cmd/evo/main.go                    # Entry point
├── internal/
│   ├── cli/                           # Cobra command handlers
│   │   ├── root.go                    # Base "evo" command
│   │   ├── init.go                    # evo init
│   │   ├── status.go                  # evo status
│   │   ├── commit.go                  # evo commit -m
│   │   ├── log.go                     # evo log
│   │   └── version.go                 # evo version
│   ├── repository/                    # Core versioning engine
│   │   ├── constants.go               # RepositoryDir, ConfigFile, HeadFile, DefaultBranch
│   │   ├── layout.go                  # BranchesDir, CommitsDir, ObjectsDir, ArtifactsDir
│   │   ├── config.go                  # Config struct, NewConfig, LoadConfig, writeConfig
│   │   ├── repository.go              # Repository struct, OpenRepository
│   │   ├── branch.go                  # Branch struct, NewBranch, LoadBranch, writeBranch, UpdateBranch
│   │   ├── commit.go                  # Commit struct, NewCommit, LoadCommit, WriteCommit, CreateCommit
│   │   ├── status.go                  # Status struct, GetStatus
│   │   ├── log.go                     # Log (returns commit history)
│   │   ├── init.go                    # Exists, Init
│   │   └── files.go                   # ensureNotInitialized, createRepositoryDirectory, writeHead
│   └── version/
│       └── version.go                 # const Version = "0.1.0"
├── docs/                              # Design documentation
├── go.mod                             # Module: github.com/Urvish0/evolution
└── go.sum
```

---

# History — What Has Been Completed

## Phase 0 — Foundation ✅

**Goal:** Design the product before writing production code.

| Task | Status |
|------|--------|
| Initialize Git repository | ✅ Done |
| Create documentation structure | ✅ Done |
| Product Requirement Document (PRD) | ✅ Done |
| Architecture Decision Records (ADR-0001 through ADR-0005) | ✅ Done |
| Domain Model | ✅ Done |
| Architecture Principles | ✅ Done |
| Engineering Decisions log | ✅ Done |
| Glossary | ✅ Done |
| System Overview | ✅ Done |
| Repository Format spec | ✅ Done |
| CLI Specification (draft) | ✅ Done |
| SDK Specification (draft) | ✅ Done |
| API Specification (draft) | ✅ Done |
| Storage Model (draft) | ✅ Done |
| Evaluation Engine spec (draft) | ✅ Done |
| Replay Engine spec (draft) | ✅ Done |
| Extensibility spec (draft) | ✅ Done |
| Security Model (draft) | ✅ Done |
| AI Pair Programmer Charter | ✅ Done |
| Founder's Playbook (evolution-bible.md) | ✅ Done |
| Project Blueprint | ✅ Done |
| Technology Stack document | ✅ Done |
| Roadmap | ✅ Done |
| Implementation Blueprint | ✅ Done |

**Learning Achieved:** Product thinking, documentation-driven development, ADR authoring, domain modeling.

---

## Phase 1 — Bootstrap ✅ (Partially)

**Goal:** Create the Evolution CLI and local repository.

| Task | Status |
|------|--------|
| Initialize Go module | ✅ Done |
| Create project structure (cmd/evo, internal/cli, internal/repository) | ✅ Done |
| Install Cobra CLI framework | ✅ Done |
| `evo version` command | ✅ Done |
| `evo help` (auto-generated by Cobra) | ✅ Done |
| `evo init` — creates `.evolution/` directory structure | ✅ Done |
| `evo status` — shows branch, commit count, working tree | ✅ Done |
| `evo commit -m "message"` — creates commit with UUID, parent, timestamp | ✅ Done |
| `evo log` — displays commit history | ✅ Done (partial — see Known Issues) |
| Configure GitHub Actions CI | ❌ Not started |
| Configure linting (golangci-lint) | ❌ Not started |
| Global CLI flags (--verbose, --quiet) | ❌ Not started |
| `.gitignore` hardening (build artifacts, binaries) | ❌ Not started |

**Learning Achieved:** Go modules, packages, structs, JSON marshalling, file I/O, error handling, Cobra CLI, UUID generation.

---

# Known Issues & Technical Debt

> These must be resolved before starting Phase 2.

### 1. `evo log` Only Shows HEAD Commit
**File:** `internal/repository/log.go`
**Problem:** `Log()` returns `[]Commit{commit}` — only the single HEAD commit. It does not walk the `Parent` chain.
**Fix:** Implement a loop that follows `commit.Parent` until it reaches an empty string (genesis commit).

### 2. Branch Functions Ignore the `name` Parameter
**File:** `internal/repository/branch.go`
**Problem:** `LoadBranch(name)` reads `DefaultBranch` instead of `name`. `writeBranch()` writes to `DefaultBranch` instead of `branch.Name`. This makes multi-branch support impossible.
**Fix:** Replace `DefaultBranch` with the function parameter in both `LoadBranch` and `writeBranch`.

### 3. `GetStatus()` Always Reports "Clean"
**File:** `internal/repository/status.go`
**Problem:** `Clean` is hardcoded to `true` and `Commits` is hardcoded to `0`. No actual file comparison or commit counting occurs.
**Fix:** Count commits by walking the log. Compare workspace files against last commit snapshot (requires Phase 3 object model).

### 4. `NewBranch()` Ignores Its Parameter
**File:** `internal/repository/branch.go`
**Problem:** `NewBranch(name)` uses `DefaultBranch` instead of the provided `name`.
**Fix:** Use the `name` parameter.

### 5. Commit Author Is Always Empty
**File:** `internal/repository/commit.go`
**Problem:** `NewCommit()` sets `Author: ""`. There is no user configuration system.
**Fix:** Add `evo config` command to set user identity, then read it during commit creation.

### 6. Empty Placeholder File
**File:** `internal/repository/load.go`
**Problem:** Contains only `package repository`. Either implement workspace loading logic or delete the file.

### 7. `.gitignore` Is Incomplete
**File:** `.gitignore`
**Problem:** Does not ignore Go build artifacts (`main`, `*.exe`, `evo`, `evo.exe`).
**Fix:** Add standard Go ignore patterns.

### 8. Compiled Binary in Repo Root
**File:** `main` (root directory)
**Problem:** A compiled binary named `main` exists in the project root and should be gitignored and deleted.

---

# Future Phases — The Road to a Polished Product

---

## Phase 2 — Repository Engine

**Goal:** Make the versioning engine production-quality. Fix all known bugs and implement complete Git-like local operations.

**Version Target:** v0.2.0

### 2.1 — Bug Fixes & Cleanup
- [ ] Fix `LoadBranch()` to use `name` parameter instead of `DefaultBranch`
- [ ] Fix `writeBranch()` to use `branch.Name` instead of `DefaultBranch`
- [ ] Fix `NewBranch()` to use the provided `name` parameter
- [ ] Implement commit chain traversal in `Log()` (walk Parent pointers)
- [ ] Count commits dynamically in `GetStatus()`
- [ ] Delete or implement `load.go`
- [ ] Remove compiled `main` binary from repo root
- [ ] Harden `.gitignore` (add `*.exe`, `main`, `/evo`, build dirs)
- [ ] Add error wrapping with `fmt.Errorf("context: %w", err)` across all functions

### 2.2 — User Configuration
- [ ] Design `evo config` command (`evo config set user.name "Urvish"`)
- [ ] Create `internal/repository/user.go` — UserConfig struct
- [ ] Store user config in `.evolution/config.json` or a separate `user.json`
- [ ] Populate `Author` field on commit creation from config
- [ ] Add `evo config list` to display current configuration

### 2.3 — Branch Operations
- [ ] `evo branch` — list all branches
- [ ] `evo branch <name>` — create a new branch from current HEAD
- [ ] `evo checkout <branch>` — switch active branch (update HEAD file)
- [ ] `evo branch -d <name>` — delete a branch
- [ ] Update `evo status` to display current branch from HEAD file (not just config)

### 2.4 — Enhanced Log & History
- [ ] Walk full commit DAG (follow Parent chain to genesis)
- [ ] `evo log --oneline` — compact single-line format
- [ ] `evo log -n <count>` — limit number of commits shown
- [ ] Display short commit IDs (first 8 characters)
- [ ] Add colored terminal output (ANSI codes or a library like `fatih/color`)

### 2.5 — Testing Foundation
- [ ] `internal/repository/config_test.go` — test NewConfig, LoadConfig, writeConfig
- [ ] `internal/repository/commit_test.go` — test NewCommit, WriteCommit, LoadCommit
- [ ] `internal/repository/branch_test.go` — test NewBranch, LoadBranch, UpdateBranch
- [ ] `internal/repository/init_test.go` — test Init, Exists
- [ ] `internal/repository/log_test.go` — test Log with multi-commit chains
- [ ] Set up test helpers (temp directory creation/cleanup)
- [ ] Configure GitHub Actions to run `go test ./...` on every push

### 2.6 — CI/CD Pipeline
- [ ] Create `.github/workflows/ci.yml`
- [ ] Run `go build`, `go test`, `go vet` on push and PR
- [ ] Add `golangci-lint` step
- [ ] Add build status badge to README

**Learning Goals:** Error wrapping, interfaces, table-driven tests, GitHub Actions, CI/CD, DAG traversal, ANSI terminal formatting.

**Interview Talking Points:** *"Implemented a directed acyclic graph for commit history traversal. Built a CI pipeline from scratch with automated testing and linting."*

---

## Phase 3 — Object Model & Content-Addressable Storage

**Goal:** Implement Git-style content-addressable storage so that commits reference hashed snapshots of actual file/artifact content.

**Version Target:** v0.3.0

### 3.1 — Hashing Engine
- [ ] Implement SHA-256 hashing for file content
- [ ] Create `internal/repository/hash.go` — `HashContent(data []byte) string`
- [ ] Store objects as `.evolution/objects/<hash>` files
- [ ] Implement object type prefix (blob, tree, commit)

### 3.2 — Blob Storage
- [ ] `internal/repository/blob.go` — store raw file content by hash
- [ ] Write blob to `objects/` directory using first 2 chars as subdirectory
- [ ] Read blob by hash
- [ ] Deduplication — identical content produces identical hash

### 3.3 — Tree Objects
- [ ] `internal/repository/tree.go` — represent directory snapshots
- [ ] Tree entries: `[mode, type, hash, filename]`
- [ ] Hash the tree itself for immutable directory representation
- [ ] Recursive tree building from a directory path

### 3.4 — Snapshot Engine
- [ ] `evo snapshot` or integrate into `evo commit`
- [ ] Walk workspace directory, create blobs for every file
- [ ] Build tree objects from directory structure
- [ ] Link root tree hash to commit object
- [ ] Update Commit struct to include `TreeHash string`

### 3.5 — Staging Area
- [ ] Implement an index/staging area (`.evolution/index`)
- [ ] `evo add <file>` — stage files for commit
- [ ] `evo add .` — stage all changes
- [ ] Compare index against HEAD tree to determine changes

### 3.6 — Working Tree Comparison
- [ ] Compare current workspace files against last committed tree
- [ ] Detect: new files, modified files, deleted files
- [ ] Update `GetStatus()` to show real modified/untracked/deleted files
- [ ] `evo diff` — show changes between working tree and HEAD

**Learning Goals:** SHA-256 hashing, content-addressable storage, Merkle trees, file diffing algorithms, staging/index concepts, io.Reader patterns.

**Interview Talking Points:** *"Built a content-addressable object store with SHA-256 deduplication. Implemented Merkle tree snapshots for workspace state tracking."*

---

## Phase 4 — Intelligence Commits

**Goal:** Extend basic commits into full Intelligence Commits that capture the complete AI system state.

**Version Target:** v0.4.0

### 4.1 — Artifact Model
- [ ] Define `Artifact` interface with `Type()`, `Hash()`, `Serialize()` methods
- [ ] Implement artifact types: Prompt, Memory, Retrieval, ToolConfig, ModelConfig, Policy
- [ ] Create `internal/repository/artifact.go`
- [ ] Store artifacts in `.evolution/artifacts/<type>/<hash>`

### 4.2 — Intelligence Commit Schema
- [ ] Extend Commit struct to become an Intelligence Commit
- [ ] Fields: ID, Parent, Message, Author, Timestamp, TreeHash, Artifacts map, Metadata
- [ ] Each commit references a collection of typed, hashed artifacts
- [ ] `evo commit` bundles all staged artifacts into the commit

### 4.3 — Artifact CLI Commands
- [ ] `evo artifact add <type> <path>` — register an artifact
- [ ] `evo artifact list` — list artifacts in current commit
- [ ] `evo artifact show <hash>` — display artifact content
- [ ] `evo artifact diff <commit1> <commit2>` — compare artifacts between commits

### 4.4 — Metadata & Context
- [ ] Capture execution metadata with each commit (model name, temperature, token limits)
- [ ] Add environment info (OS, Go version, Evolution version)
- [ ] Add tags/labels for commits

**Learning Goals:** Go interfaces, polymorphism, type assertions, builder patterns, schema evolution.

**Interview Talking Points:** *"Designed an extensible artifact system using Go interfaces. Each Intelligence Commit captures the complete operational state of an AI system including prompts, memory, and model configuration."*

---

## Phase 5 — Branch & Merge Engine

**Goal:** Support parallel intelligence development with branching and merging.

**Version Target:** v0.5.0

### 5.1 — Branch Management
- [ ] Solidify branch creation, deletion, renaming
- [ ] `evo branch rename <old> <new>`
- [ ] Branch metadata (created_at, created_from)

### 5.2 — Checkout & Restore
- [ ] `evo checkout <branch>` — switch branches and restore workspace to that branch's HEAD tree
- [ ] `evo checkout <commit-id>` — detached HEAD mode
- [ ] `evo restore <file>` — restore a single file from HEAD

### 5.3 — Merge
- [ ] `evo merge <branch>` — merge another branch into current
- [ ] Three-way merge strategy (common ancestor, ours, theirs)
- [ ] Conflict detection and resolution for text artifacts
- [ ] Merge commits (two parent references)

### 5.4 — Diff Engine
- [ ] `evo diff` — working tree vs HEAD
- [ ] `evo diff <commit1> <commit2>` — compare any two commits
- [ ] `evo diff --artifacts` — compare artifact snapshots
- [ ] Line-by-line diff output with unified diff format

**Learning Goals:** Graph algorithms (LCA — lowest common ancestor), three-way merge, conflict resolution, diff algorithms (Myers diff).

**Interview Talking Points:** *"Implemented a three-way merge engine with conflict detection. Built a diff system supporting both file content and AI artifact comparison."*

---

## Phase 6 — Replay Engine

**Goal:** Reproduce historical AI executions deterministically.

**Version Target:** v0.6.0

### 6.1 — Execution Recording
- [ ] Define `Execution` struct (input, output, commit ref, timing, metadata)
- [ ] Store executions in `.evolution/executions/`
- [ ] Link execution to its Intelligence Commit

### 6.2 — Replay Command
- [ ] `evo replay <commit-id>` — restore Intelligence Commit state
- [ ] Reconstruct artifacts (prompt, memory, retrieval, tools, model config)
- [ ] Display what the AI system looked like at that point in time
- [ ] Compare replay output against original execution

### 6.3 — Execution Comparison
- [ ] `evo replay --compare <commit1> <commit2>` — side-by-side comparison
- [ ] Highlight behavioral differences between versions

**Learning Goals:** Goroutines, context.Context, channels, process execution, I/O piping.

**Interview Talking Points:** *"Built a deterministic replay engine that reconstructs historical AI system states and enables behavioral comparison across versions."*

---

## Phase 7 — Evaluation Engine

**Goal:** Measure and compare AI system quality across versions.

**Version Target:** v0.7.0

### 7.1 — Evaluation Framework
- [ ] Define `Evaluator` interface
- [ ] Built-in evaluators: Correctness, Performance, Cost, Safety
- [ ] Custom evaluator support (user-defined scripts/functions)
- [ ] Store evaluation results in `.evolution/evaluations/`

### 7.2 — Evaluation CLI
- [ ] `evo evaluate <commit-id>` — run evaluations against a commit
- [ ] `evo evaluate --compare <commit1> <commit2>` — compare evaluations
- [ ] `evo evaluate --report` — generate human-readable evaluation report

### 7.3 — Regression Detection
- [ ] Automatic regression detection between commits
- [ ] Threshold-based pass/fail (e.g. "accuracy dropped by more than 5%")
- [ ] Integration with CI/CD for automated evaluation gates

**Learning Goals:** Strategy pattern, plugin architecture, metrics collection, reporting.

**Interview Talking Points:** *"Designed an evaluation platform with pluggable evaluators and automated regression detection for AI systems."*

---

## Phase 8 — Python SDK

**Goal:** Expose Evolution programmatically for integration into AI workflows.

**Version Target:** v0.8.0

### 8.1 — SDK Design
- [ ] Define public API surface (Repository, Commit, Artifact, Execution, Evaluation)
- [ ] Python package: `evolution-sdk`
- [ ] Repository operations: init, open, commit, log, checkout
- [ ] Artifact operations: add, list, get, diff

### 8.2 — Integration Patterns
- [ ] Decorator-based commit capture (`@evolution.track`)
- [ ] Context manager for execution recording
- [ ] Framework integrations (LangChain, LlamaIndex, CrewAI)

**Learning Goals:** Python packaging, API design, cross-language interop, gRPC/subprocess communication.

**Interview Talking Points:** *"Designed a Python SDK that integrates with major AI frameworks, enabling automatic intelligence versioning through decorators and context managers."*

---

## Phase 9 — REST API & Server

**Goal:** Enable remote access and prepare for collaboration features.

**Version Target:** v0.9.0

### 9.1 — API Server
- [ ] HTTP server using Go standard library or Chi router
- [ ] REST endpoints: repositories, commits, artifacts, executions, evaluations
- [ ] API versioning (v1)
- [ ] Authentication (API keys, JWT)

### 9.2 — Remote Sync
- [ ] `evo remote add <name> <url>` — configure remote server
- [ ] `evo push` — upload commits and objects to remote
- [ ] `evo pull` — download commits and objects from remote
- [ ] Conflict resolution for concurrent pushes

**Learning Goals:** HTTP servers in Go, REST API design, authentication, networking, protocol design.

**Interview Talking Points:** *"Built a REST API server with authentication and designed a remote sync protocol for distributed AI version control."*

---

## Phase 10 — Collaboration & Review

**Goal:** Enable team workflows for AI intelligence changes.

**Version Target:** v1.0.0

### 10.1 — Review Workflow
- [ ] `evo review create` — open a review for pending changes
- [ ] `evo review approve/reject` — review actions
- [ ] Review comments and discussions
- [ ] Approval gates before deployment

### 10.2 — Organizations & Access Control
- [ ] Organization and team management
- [ ] Role-based access control (RBAC)
- [ ] Audit logs for all operations

**Learning Goals:** Authorization systems, RBAC design, audit logging, collaborative workflows.

---

## Phase 11 — Deployment Pipeline

**Goal:** Safely promote intelligence versions to production environments.

**Version Target:** v1.1.0

### 11.1 — Deployment Management
- [ ] `evo deploy <commit-id> --env <environment>`
- [ ] Canary deployments
- [ ] Automatic rollback on evaluation failure
- [ ] Deployment history and audit trail

### 11.2 — Container Integration
- [ ] Docker image generation from Intelligence Commit
- [ ] Kubernetes manifest generation
- [ ] Helm chart support

**Learning Goals:** Docker, Kubernetes, Helm, deployment strategies, infrastructure as code.

---

## Phase 12 — Observability

**Goal:** Provide runtime visibility into AI system behavior.

**Version Target:** v1.2.0

- [ ] OpenTelemetry integration
- [ ] Prometheus metrics export
- [ ] Grafana dashboards
- [ ] Structured logging with correlation IDs
- [ ] Tracing across intelligence executions

**Learning Goals:** Distributed tracing, metrics instrumentation, observability platforms.

---

## Phase 13 — Platform & Ecosystem

**Goal:** Transform Evolution into a full platform with extensibility and ecosystem support.

**Version Target:** v2.0.0

### 13.1 — Plugin System
- [ ] Plugin architecture for custom artifact types
- [ ] Evaluation provider plugins
- [ ] Storage backend plugins (S3, GCS, Azure Blob)

### 13.2 — Cloud Platform
- [ ] Multi-tenancy
- [ ] Web dashboard
- [ ] Billing and usage tracking
- [ ] Enterprise features (SSO, compliance)

### 13.3 — Ecosystem Integrations
- [ ] VS Code Extension
- [ ] GitHub App (PR-based intelligence reviews)
- [ ] Desktop application
- [ ] MCP (Model Context Protocol) integration
- [ ] gRPC API for high-performance integrations

**Learning Goals:** Plugin architectures, multi-tenant systems, web development, cloud infrastructure, product management at scale.

---

# Accepted Architecture Decision Records

| ADR | Title | Status |
|-----|-------|--------|
| ADR-0001 | Intelligence Is the Unit of Versioning | Accepted |
| ADR-0002 | Local First | Accepted |
| ADR-0003 | Immutable Intelligence Commits | Accepted |
| ADR-0004 | Framework Agnostic | Accepted |
| ADR-0005 | CLI First | Accepted |

---

# Key Engineering Decisions

| Date | Decision | Impact |
|------|----------|--------|
| 2026-07-21 | Intelligence, not code, is what we version | Central abstraction of the platform |
| 2026-07-21 | Memory is an artifact, not a core object | Simplifies architecture, enables deterministic replay |
| 2026-07-21 | Replay is a first-class primitive | Foundational capability, not optional debugging |
| 2026-07-21 | Product drives technology | No Kafka/K8s/Neo4j until the product needs them |
| 2026-07-21 | Architecture before implementation | Phase 0 is a formal engineering stage |
| 2026-07-21 | Evolution is infrastructure | Not a framework, not a library — infrastructure |
| 2026-07-21 | Build for the next decade | Architecture over speed when there's a trade-off |

---

# Document Index

> Reference for all existing documentation in the project.

| Document | Path | Purpose |
|----------|------|---------|
| README | `README.md` | Public project overview |
| **Master Plan** | **`EVOLUTION_MASTER_PLAN.md`** | **This file — single source of truth** |
| Changelog | `CHANGELOG.md` | Release notes |
| Agent Rules | `.agents/AGENTS.md` | AI collaboration rules for the project |
| PRD | `docs/PRD.md` | Product requirements and vision |
| Domain Model | `docs/DOMAIN_MODEL.md` | Entity definitions and relationships |
| Architecture Principles | `docs/ARCHITECTURE_PRINCIPLES.md` | Non-negotiable design constraints |
| Engineering Decisions | `docs/DECISIONS.md` | Lightweight decision log |
| Glossary | `docs/GLOSSARY.md` | Term definitions |
| Repository Format | `docs/REPOSITORY_FORMAT.md` | On-disk directory structure |
| CLI Specification | `docs/CLI_SPEC.md` | Command definitions |
| API Specification | `docs/API_SPEC.md` | REST API design (draft) |
| SDK Specification | `docs/SDK_SPEC.md` | SDK design (draft) |
| Evaluation Engine | `docs/EVALUATION_ENGINE.md` | Evaluation design (draft) |
| Replay Engine | `docs/REPLAY_ENGINE.md` | Replay design (draft) |
| ADR Index | `docs/adr/README.md` | ADR process and index |

### Removed Documents (merged into Master Plan)

The following documents were removed during the professional cleanup. Their content has been absorbed into `EVOLUTION_MASTER_PLAN.md` and `.agents/AGENTS.md`:

- `ROADMAP.md` — phases merged into Master Plan
- `IMPLEMENTATION_BLUEPRINT.md` — task tracking merged into Master Plan
- `TECH_STACK.md` — tech stack table merged into Master Plan
- `invisible.md` — personal scratch notes
- `docs/evolution-bible.md` — vision, strategy, and learning goals merged into Master Plan
- `docs/PROJECT_BLUEPRINT.md` — core concepts merged into Master Plan
- `docs/AI_PAIR_PROGRAMMER_CHARTER.md` — rules moved to `.agents/AGENTS.md`
- `docs/SYSTEM_OVERVIEW.md` — covered by Domain Model and Master Plan
- `docs/STORAGE_MODEL.md` — covered by Repository Format
- `docs/EXTENSIBILITY.md` — covered by Architecture Principles
- `docs/SECURITY_MODEL.md` — principles noted in Architecture Principles
- `internal/repository/load.go` — empty placeholder file deleted

---

# How to Use This Document

1. **Before starting work:** Check this file for the next unchecked task.
2. **After completing work:** Mark the task `[x]` and add any new decisions or bugs discovered.
3. **When making architectural decisions:** Add a row to the Key Engineering Decisions table and create an ADR if the decision is significant.
4. **When starting a new phase:** Add the date and update the version target.
5. **In interviews:** Reference the "Interview Talking Points" section for each phase to articulate what you built and why.

---

> *"Don't version prompts. Version intelligence."*
