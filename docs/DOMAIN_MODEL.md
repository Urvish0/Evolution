# Domain Model

## Purpose

The Domain Model defines the core entities of Evolution, their relationships, and the rules governing how they interact.

It is implementation-independent and serves as the conceptual foundation of the platform.

> **Status:** Draft

---

# Core Entities

## Workspace

Top-level organizational container for one or more repositories.

Responsible for shared configuration and resources.

---

## Repository

The primary unit for managing and versioning AI intelligence.

Contains branches, commits, artifacts, and execution history.

---

## Branch

Represents an independent line of intelligence development.

Allows experimentation without affecting other branches.

---

## Intelligence

The complete operational behavior of an AI system.

Composed of multiple versioned artifacts.

---

## Intelligence Commit

An immutable snapshot of Intelligence at a point in time.

Represents the primary unit of version history.

---

## Artifact

A versioned component that contributes to Intelligence.

Examples include prompts, memory, retrieval configuration, model settings, tools, and policies.

---

## Execution

A recorded run of an Intelligence Commit.

Captures inputs, outputs, metadata, and execution context.

---

## Evaluation

Measures the quality, correctness, safety, cost, or performance of an Execution.

---

## Deployment

Represents promoting an Intelligence Commit into an environment.

---

# Relationships

- A Workspace contains multiple Repositories.
- A Repository contains multiple Branches.
- A Branch references a sequence of Intelligence Commits.
- An Intelligence Commit references one or more Artifacts.
- Executions are produced from Intelligence Commits.
- Evaluations assess Executions.
- Deployments promote Intelligence Commits.

---

# Domain Rules

- Intelligence Commits are immutable.
- Every Artifact belongs to an Intelligence Commit.
- Every Execution references exactly one Intelligence Commit.
- Every Evaluation references exactly one Execution.
- Branches maintain independent version histories.

---

# Future Entities

- Organization
- Team
- User
- Environment
- Registry
- Policy
- Review
- Release
