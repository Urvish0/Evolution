# ADR-0002: Local-First Architecture

## Status

Accepted

---

## Context

Developers expect engineering tools to function without requiring continuous connectivity or cloud services.

A local-first architecture enables faster development workflows, offline usage, and seamless integration into existing development environments.

---

## Decision

Evolution repositories shall operate locally by default.

Cloud services will provide optional collaboration and synchronization capabilities but will not be required for core functionality.

---

## Consequences

### Benefits

- Offline development
- Faster operations
- Improved privacy
- Easier CI/CD integration

### Trade-offs

- Synchronization complexity
- Future collaboration challenges

---

## Alternatives Considered

### Cloud-first architecture

Rejected because it introduces unnecessary dependencies for local development.

### Fully managed SaaS

Rejected because Evolution should remain developer-centric and infrastructure-independent.
