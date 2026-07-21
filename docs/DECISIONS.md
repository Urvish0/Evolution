# Engineering Decisions

> This document captures the key engineering insights that shaped Evolution.
>
> Unlike Architecture Decision Records (ADRs), these notes are lightweight observations, design realizations, and product principles discovered throughout development.
>
> Every entry should answer:
>
> - What did we learn?
> - Why does it matter?
> - How does it influence the architecture?

---

# 2026-07-21

## Evolution versions Intelligence, not code

One of the earliest and most important realizations was that Evolution is **not Git for AI**.

Git versions source code.

Evolution versions **Intelligence**.

An Intelligence Commit captures the complete state of an AI system at a point in time, including prompts, memory, retrieval, tools, evaluations, policies, and metadata.

Everything else is simply an artifact of that intelligence.

**Impact**

This became the central abstraction of the entire platform.

---

## Memory is an artifact

Initially, memory appeared to be one of the core objects.

Further analysis showed that memory should instead be treated as one versioned artifact within an Intelligence Commit.

Memory should never exist independently of the intelligence that produced or consumed it.

**Impact**

- Simplifies architecture.
- Makes replay deterministic.
- Keeps all artifacts synchronized.

---

## Replay is a first-class primitive

Traditional software debugging uses logs.

AI systems require something stronger.

Evolution should be able to replay historical intelligence exactly as it originally executed.

Replay should include:

- Prompt
- Memory
- Retrieval
- Tool calls
- Model configuration
- Policies
- Inputs
- Outputs

**Impact**

Replay becomes a foundational capability rather than an optional debugging feature.

---

## Product drives technology

Technology choices should never be made because they are popular.

Every dependency must solve a clearly identified product problem.

Examples:

- Kafka only when asynchronous event streaming is required.
- Neo4j only if relational storage cannot efficiently model the Decision Graph.
- Kubernetes only when deployment complexity justifies orchestration.

**Impact**

Prevents unnecessary complexity and encourages intentional architecture.

---

## Architecture before implementation

Evolution will be designed before it is built.

The goal is to reduce uncertainty through documentation, discussion, and design before writing production code.

Implementation should validate the architecture—not discover it.

**Impact**

This establishes Phase 0 as a formal engineering stage.

---

## Evolution is infrastructure

Evolution is not:

- an AI agent framework
- a prompt management tool
- a memory library
- an observability platform

Instead, Evolution provides the engineering infrastructure that manages the lifecycle of AI systems regardless of the framework they use.

**Impact**

This defines the product's long-term positioning and prevents unnecessary feature expansion.

---

## Build for the next decade

Short-term convenience should never compromise long-term maintainability.

Whenever there is a trade-off between speed and architecture, architecture takes priority unless there is a compelling product reason to decide otherwise.

**Impact**

This principle guides engineering decisions across the entire project.

---

# Future Entries

This document will continue to grow as Evolution evolves.

Only add entries when a decision significantly changes how the product is understood or built.
