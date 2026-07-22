# ADR-0005: CLI Before GUI

## Status

Accepted

---

## Context

Developer tools benefit from automation, scripting, and CI/CD integration.

A CLI provides a stable foundation that can later power graphical interfaces.

---

## Decision

The first user interface for Evolution will be a command-line interface.

Future graphical interfaces will build upon the same underlying APIs.

---

## Consequences

### Benefits

- Faster development
- Easier automation
- Better testing
- Stable API design

### Trade-offs

- Less accessible for non-technical users
- Slower visual adoption

---

## Alternatives Considered

### GUI-first development

Rejected because it slows iteration and tightly couples the interface with early architectural decisions.
