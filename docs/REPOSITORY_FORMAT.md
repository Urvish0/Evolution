# Repository Format

## Purpose

This document defines the directory structure and on-disk layout of an Evolution repository.

> **Status:** Draft

---

## Repository Structure

.evolution/

├── config/
├── branches/
├── commits/
├── artifacts/
├── executions/
├── evaluations/
├── refs/
└── HEAD

---

## Directory Overview

| Directory   | Purpose                     |
| ----------- | --------------------------- |
| config      | Repository configuration    |
| branches    | Branch metadata             |
| commits     | Intelligence Commit objects |
| artifacts   | Versioned artifacts         |
| executions  | Execution history           |
| evaluations | Evaluation results          |
| refs        | Branch and tag references   |
| HEAD        | Current active branch       |

---

## Design Goals

- Human-readable where practical
- Immutable commit history
- Efficient artifact lookup
- Extensible repository structure
