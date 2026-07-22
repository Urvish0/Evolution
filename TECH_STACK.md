# Technology Stack

## Purpose

This document defines the core technologies used to build Evolution.

The chosen stack prioritizes simplicity, performance, maintainability, portability, and long-term stability.

---

# Guiding Principles

- Keep dependencies minimal.
- Prefer standard libraries whenever possible.
- Build for long-term maintainability.
- Cross-platform support is a first-class requirement.
- Optimize for developer experience without sacrificing performance.

---

# Core Language

## Go

### Why Go?

- Cross-platform compilation
- Excellent CLI ecosystem
- Fast execution
- Simple concurrency model
- Strong standard library
- Easy distribution as a single binary
- Mature tooling and testing

---

# CLI

## Cobra

Purpose:

- Command parsing
- Help generation
- Command hierarchy

Examples:

- evo init
- evo commit
- evo replay

---

# Configuration

YAML

Reason:

- Human-readable
- Widely adopted
- Easy to edit

---

# Repository Metadata

JSON

Reason:

- Simple
- Portable
- Language agnostic
- Easy to inspect

---

# Storage

Local File System

Future:

- Pluggable storage backends

---

# Logging

Go slog

Reason:

- Standard library
- Structured logging
- Lightweight

---

# Testing

Go Testing Package

Additional tools may be introduced when necessary.

---

# Documentation

Markdown

Architecture Decision Records (ADRs)

GitHub Issues

---

# Version Control

Git

Evolution repositories remain compatible with Git repositories.

---

# Build System

Go Toolchain

Examples:

go build

go test

go fmt

go vet

---

# CI/CD

GitHub Actions

Initial workflows:

- Build
- Test
- Lint

---

# Code Quality

Formatting

- gofmt

Static Analysis

- go vet

Linting

- golangci-lint (future)

---

# Package Management

Go Modules

---

# Supported Platforms

- Linux
- macOS
- Windows

---

# Future Technologies

These are intentionally deferred until needed.

- SQLite
- PostgreSQL
- gRPC
- REST API
- Web UI
- Cloud Services
- Plugin Marketplace

---

# Technology Philosophy

Technology choices should support Evolution's goals of transparency, reproducibility, portability, and simplicity.

New technologies should only be introduced when they provide clear value and align with the project's architectural principles.
