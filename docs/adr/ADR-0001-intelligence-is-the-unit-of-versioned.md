# ADR-0001: Intelligence Is the Unit of Versioning

## Status

Accepted

---

## Context

Traditional version control systems manage source code, while AI systems derive their behavior from multiple evolving components such as prompts, memory, retrieval strategies, model configurations, tools, and policies.

Versioning these components independently makes it difficult to understand the complete state of an AI system at any given point in time.

Evolution requires a versioning model that captures the entire operational state of an AI system rather than isolated artifacts.

---

## Decision

Evolution versions Intelligence as a whole.

Each version represents a complete snapshot of an AI system's operational state at a specific point in time.

This snapshot is represented as an Intelligence Commit.

---

## Consequences

### Benefits

- Complete reproducibility
- Reliable historical replay
- Simpler mental model
- Consistent version history

### Trade-offs

- Larger snapshots
- Additional storage requirements
- More sophisticated diffing mechanisms

---

## Alternatives Considered

### Version prompts only

Rejected because prompts do not fully determine AI behavior.

### Version every artifact independently

Rejected because developers ultimately need to understand the behavior of the complete system rather than isolated components.
