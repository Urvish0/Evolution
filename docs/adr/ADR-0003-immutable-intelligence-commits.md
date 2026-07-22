# ADR-0003: Intelligence Commits Are Immutable

## Status

Accepted

---

## Context

Reliable replay and debugging require historical versions to remain unchanged.

Mutable history makes it impossible to guarantee reproducibility.

---

## Decision

Once created, an Intelligence Commit cannot be modified.

Any change to intelligence results in the creation of a new commit.

---

## Consequences

### Benefits

- Reliable history
- Deterministic replay
- Strong audit trail
- Simpler reasoning

### Trade-offs

- Additional storage
- More commits over time

---

## Alternatives Considered

### Mutable commits

Rejected because modifying history undermines reproducibility and trust.
