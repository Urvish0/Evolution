# ADR-0006: Artifact Hash Auto-Computation at Commit Time

**Status:** Accepted  
**Date:** 2026-08-06  
**Deciders:** Urvish  

## Context

In the Intelligence Manifest Specification v0.1, the `hash` field on every artifact was marked as **required**. This meant developers had to manually compute SHA-256 hashes before adding artifacts to `evolution.manifest.json`.

During implementation (Phase 4.1–4.5), we discovered that requiring pre-computed hashes creates significant friction:

1. Developers must run a separate hashing tool before every manifest update.
2. Hash values become stale if the file content changes but the manifest is not updated.
3. It breaks the "incrementally adoptable" design principle — adding an artifact should be as simple as specifying a type, name, and path.

## Decision

**Make `hash` optional at authoring time. Auto-compute it at commit time.**

When `evo commit` runs and detects `evolution.manifest.json`:

1. For each artifact where `hash` is empty (`""`), Evolution reads the file at `path`.
2. Computes `SHA-256("blob <len>\0<content>")` (same format as blob storage).
3. Stores the artifact metadata JSON in `.evolution/artifacts/<type>/<hash>.json`.
4. Attaches the computed hash to the artifact in the Intelligence Commit.

If the file at `path` does not exist, the hash remains empty and the artifact is still recorded (metadata-only artifact).

## Consequences

### Positive

- **Zero-friction artifact registration:** `evo artifact add prompt my-prompt prompts/system.txt` — no hash needed.
- **Always-correct hashes:** Hashes are computed from actual file content at commit time, never stale.
- **Consistent with Git:** Git computes blob hashes automatically; users never manually hash files.

### Negative

- **Authoring-time hash unavailable:** Tools that want to verify artifact integrity before committing cannot rely on the manifest hash field.
- **Path dependency:** If `path` is incorrect or the file is missing, the hash cannot be computed.

### Mitigations

- `evo manifest validate` warns about empty hashes.
- `evo artifact list` displays `(unhashed)` for artifacts without computed hashes, making the state visible.

## Alternatives Considered

1. **Keep hash required:** Rejected — too much developer friction, violates incremental adoption principle.
2. **Compute hash on `evo artifact add`:** Rejected — the file may change between `add` and `commit`, causing stale hashes.
3. **Compute hash on `evo manifest validate`:** Considered but not adopted — validation should be read-only, not mutating.
