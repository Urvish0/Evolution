# Architecture Decision Records (ADRs)

## Purpose

Architecture Decision Records (ADRs) document the significant architectural decisions made throughout the development of Evolution.

Rather than relying on tribal knowledge or historical discussions, ADRs provide a permanent record of why important decisions were made, the alternatives that were considered, and the trade-offs involved.

They help ensure architectural consistency as the project evolves and make it easier for future contributors to understand the reasoning behind the system's design.

---

## What belongs in an ADR?

An ADR should capture decisions that have a long-term impact on the architecture or direction of Evolution.

Examples include:

- System architecture
- Core domain concepts
- Storage strategy
- Repository structure
- Versioning model
- Deployment strategy
- API design
- SDK design
- Integration strategy

Implementation details should **not** be documented as ADRs.

---

## ADR Lifecycle

Each ADR has one of the following statuses:

- Proposed
- Accepted
- Superseded
- Deprecated

Accepted ADRs become part of Evolution's architectural foundation.

If a decision changes in the future, a new ADR should supersede the previous one rather than modifying history.

---

## Naming Convention

```
ADR-XXXX-short-title.md
```

Examples:

```
ADR-0001-intelligence-is-versioned.md
ADR-0002-local-first.md
ADR-0003-immutable-intelligence-commits.md
```

---

## Creating a New ADR

1. Copy `TEMPLATE.md`.
2. Assign the next sequential ADR number.
3. Fill in the Context, Decision, Consequences, and Alternatives sections.
4. Submit the ADR for review.
5. Mark it as Accepted once approved.

---

## Guiding Principle

Architecture decisions should be difficult to change but easy to understand.

Every accepted ADR should clearly answer:

> Why was this decision made?
