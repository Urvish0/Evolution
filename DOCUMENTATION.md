# Evolution Platform — Official Documentation & Technical Reference

> **Version:** v0.3.0 / v0.4.0 (In Progress)  
> **Core Mission:** *"Version Intelligence, Not Code."*  
> **Specification:** Intelligence Manifest Specification v0.1  

---

## Table of Contents

1. [Overview & Philosophy](#1-overview--philosophy)
2. [Core Architecture & Concepts](#2-core-architecture--concepts)
   - [2.1 Content-Addressable Storage (Blobs)](#21-content-addressable-storage-blobs)
   - [2.2 Merkle Trees](#22-merkle-trees)
   - [2.3 Staging Area (Index)](#23-staging-area-index)
   - [2.4 Working Tree Comparison & Content Diffs](#24-working-tree-comparison--content-diffs)
   - [2.5 Artifact Model](#25-artifact-model)
   - [2.6 Intelligence Manifest Specification v0.1](#26-intelligence-manifest-specification-v01)
3. [CLI Reference Manual](#3-cli-reference-manual)
   - [`evo init`](#evo-init)
   - [`evo status`](#evo-status)
   - [`evo config`](#evo-config)
   - [`evo add`](#evo-add)
   - [`evo commit`](#evo-commit)
   - [`evo log`](#evo-log)
   - [`evo diff`](#evo-diff)
   - [`evo branch`](#evo-branch)
   - [`evo checkout`](#evo-checkout)
   - [`evo manifest`](#evo-manifest)
   - [`evo version`](#evo-version)
4. [End-to-End Example Use Cases](#4-end-to-end-example-use-cases)
   - [Use Case 1: Initializing and Versioning an AI Legal Assistant](#use-case-1-initializing-and-versioning-an-ai-legal-assistant)
   - [Use Case 2: Parallel Branching for Model Configuration Experiments](#use-case-2-parallel-branching-for-model-configuration-experiments)
   - [Use Case 3: Line-by-Line Content Diffing on System Prompts](#use-case-3-line-by-line-content-diffing-on-system-prompts)
5. [Repository On-Disk Layout](#5-repository-on-disk-layout)

---

## 1. Overview & Philosophy

Evolution is an AI-native version control platform. Traditional version control systems like Git version source code files. Evolution versions **Intelligence** — the complete operational state of an AI system at any point in time.

An AI system's behavior is governed by dynamic, interdependent components:
- **System Prompts & Templates** (`prompt`)
- **Memory & Conversation History Strategies** (`memory`)
- **Knowledge Retrieval & Vector DB Settings** (`retrieval`)
- **External Tools & API Definitions** (`tool`)
- **Model Configurations** (`model_config` - temperature, max tokens, provider)
- **Safety Policies & Guardrails** (`policy`)

If any of these change, the AI behaves differently. Evolution captures all of them as **typed, content-addressed artifacts** bundled into immutable **Intelligence Commits**.

---

## 2. Core Architecture & Concepts

### 2.1 Content-Addressable Storage (Blobs)
All file content stored in Evolution is hashed using **SHA-256** with Git-style header prefixes (`"blob <len>\0<content>"`).
- Object paths are sharded using the first 2 characters of the hash: `.evolution/objects/xx/yyyy...`
- **Automatic Deduplication:** If two files have identical content, they produce the exact same SHA-256 hash and are written to disk only once.
- **Immutability:** Object files are stored with read-only permissions (`0444`).

### 2.2 Merkle Trees
Directories are represented as **Tree Objects**. Each tree entry contains `[mode, type, hash, filename]`.
- Because child tree and blob hashes are embedded into parent tree objects, a single 64-character Root Tree SHA-256 hash cryptographically guarantees the integrity of an entire workspace folder structure.

### 2.3 Staging Area (Index)
The Staging Area (`.evolution/index`) is a preparation zone.
- Running `evo add <file>` stages specific files or directory trees into the index.
- When `evo commit` runs, it builds a Merkle tree snapshot directly from the staged entries and clears the index.

### 2.4 Working Tree Comparison & Content Diffs
Evolution compares current workspace files against the HEAD committed tree snapshot:
- **Untracked:** Files on disk not in HEAD commit and not staged.
- **Modified:** Tracked files whose content SHA-256 hash differs from the committed version.
- **Deleted:** Files present in HEAD commit but removed from disk.
- **Content Diffs:** `evo diff` uses the **Longest Common Subsequence (LCS)** algorithm to render line-by-line `+` (added) and `-` (deleted) unified diffs between committed blobs and current disk files.

### 2.5 Artifact Model
Evolution implements a typed artifact model (`Artifact` interface in Go). Every AI component belongs to one of 6 core types matching the Intelligence Manifest Spec v0.1:
1. `prompt` — System prompts, templates, few-shot examples
2. `memory` — Buffer window, summary, vector, or graph memory strategies
3. `retrieval` — Vector store settings, chunk sizes, top-k retrieval rules
4. `tool` — External API and function calling definitions
5. `model_config` — LLM model provider, temperature, max output tokens, top-p
6. `policy` — Safety rules and output enforcement levels

### 2.6 Intelligence Manifest Specification v0.1
The open standard defining how AI system state is declared in `evolution.manifest.json`.
- Required fields: `version`, `name`
- Contains arrays of typed artifacts and environment metadata.
- Automatically detected, validated, and attached to commits when running `evo commit`.

---

## 3. CLI Reference Manual

### `evo init`
Initializes a new `.evolution/` repository in the current working directory.

```bash
evo init
```
**Output:**
```text
Initialized empty Evolution repository in .evolution/
```

---

### `evo status`
Displays repository health, active branch, commit count, and working tree changes.

```bash
evo status
```
**Example Output:**
```text
On branch main
Commits: 3

Changes to be committed:
  staged:   prompts/system.txt

Changes not staged for commit:
  modified: config/model.json

Untracked files:
  scratch.py
```

---

### `evo config`
Gets or sets developer identity stored globally in `~/.evolution/config.json`.

```bash
# Set identity
evo config set user.name "Urvish"
evo config set user.email "urvish@example.com"

# View current identity
evo config list
```

---

### `evo add`
Stages files or directory hierarchies into `.evolution/index`.

```bash
# Stage a specific file
evo add prompts/system.txt

# Stage a folder recursively
evo add tools/

# Stage all changes in working directory
evo add .
```

---

### `evo commit`
Creates an immutable Intelligence Commit snapshot. Automatically checks for `evolution.manifest.json` and attaches declared artifacts, tag labels, and execution metadata.

```bash
# Basic commit
evo commit -m "Updated system prompt and temperature settings"

# Commit with tags and custom execution metadata
evo commit -m "Promoted prompt to production" --tag production --tag v1.0 --meta model=gpt-4o --meta temp=0.2
```
**Flags:**
- `-m, --message <string>` (Required): Commit description message.
- `-t, --tag <string>` (Optional, repeatable): Tag labels attached to the commit.
- `--meta <key=value>` (Optional, repeatable): Custom execution metadata key-value pairs.

---

### `evo log`
Displays commit history starting from HEAD back to the initial (genesis) commit.

```bash
# Standard history view
evo log

# Compact single-line view
evo log --oneline

# Limit output to N commits
evo log -n 5
```
**Example Output:**
```text
commit ab318534
Tree:   e0f78b91
Author: Urvish <urvish@example.com>
Date:   2026-08-06T12:30:00Z

    Updated system prompt and temperature settings
```

---

### `evo diff`
Shows line-by-line unified content diffs (`--- a/` / `+++ b/`) for modified tracked files between the working directory and HEAD commit.

```bash
evo diff
```
**Example Output:**
```text
--- a/prompts/system.txt
+++ b/prompts/system.txt
@@ -1 +1 @@
-System Prompt: You are a helpful assistant.
+System Prompt: You are a legal research AI. Always cite sources.
```

---

### `evo branch`
Lists, creates, renames, or deletes branches. Renders rich branch metadata (short head commit ID, commit count, and last commit message).

```bash
# List all branches with rich details (* indicates active branch)
evo branch

# Create a new branch
evo branch -n feature-claude-model

# Rename current branch (or 'evo branch -m old-name new-name')
evo branch -m exp-claude-sonnet

# Delete an inactive branch (active branch protected)
evo branch -d exp-claude-sonnet
```

---

### `evo checkout`
Switches active branch and automatically restores physical workspace files on disk to match the branch's HEAD Merkle Tree snapshot. Supports `-b / --create` flag and detached HEAD commit checkout.

```bash
# Switch to an existing branch (restores physical files)
evo checkout main

# Create a new branch and switch to it immediately
evo checkout -b exp-prompt-v2

# Checkout a specific commit directly (detached HEAD mode)
evo checkout ed523a4d
```

---

### `evo restore`
Discards working directory modifications for a specified file and restores its content directly from the HEAD commit snapshot.

```bash
evo restore prompts/system.txt
```

---

### `evo manifest`
Manages the `evolution.manifest.json` file conforming to Spec v0.1.

```bash
# Generate starter manifest file
evo manifest init --name "legal-agent" --description "Legal Research Assistant"

# Validate manifest schema against Spec v0.1
evo manifest validate

# Display current workspace manifest JSON
evo manifest show
```

---

### `evo artifact`
Registers, lists, inspects, and performs semantic diffs on typed AI artifacts (`prompt`, `memory`, `retrieval`, `tool`, `model_config`, `policy`).

```bash
# Register an artifact in evolution.manifest.json
evo artifact add prompt sys-prompt prompts/system.txt

# List all AI artifacts attached to HEAD commit
evo artifact list

# Display raw JSON for a stored artifact
evo artifact show prompt b671051c

# Perform high-level semantic diff between two commits
evo artifact diff <commit1_id> <commit2_id>
```

---

### `evo version`
Displays CLI version information.

```bash
evo version
```

---

## 4. End-to-End Example Use Cases

### Use Case 1: Initializing and Versioning an AI Legal Assistant

```bash
# 1. Setup global identity
evo config set user.name "Urvish"
evo config set user.email "urvish@example.com"

# 2. Initialize repo
evo init

# 3. Create starter manifest
evo manifest init --name "legal-assistant" --description "AI Legal Research System"

# 4. Create prompt file
mkdir prompts
echo "System Prompt: You are a legal research AI." > prompts/system.txt

# 5. Stage and commit initial state
evo add .
evo commit -m "Initial legal assistant configuration"

# 6. Verify commit history and attached tree
evo log
```

---

### Use Case 2: Parallel Branching for Model Configuration Experiments

```bash
# 1. Create and switch to experimental branch for Claude 3.5 Sonnet
evo branch -n exp-claude-sonnet
evo checkout exp-claude-sonnet

# 2. Update model config in manifest or file
# (Modify config/model.json to use claude-3.5-sonnet)

# 3. Commit experimental configuration
evo add .
evo commit -m "Experiment: Switched to Claude 3.5 Sonnet with temperature 0.2"

# 4. Switch back to main branch to compare
evo checkout main
evo log --oneline
```

---

### Use Case 3: Line-by-Line Content Diffing on System Prompts

```bash
# 1. Modify prompt text in working directory
echo "System Prompt: You are a senior legal researcher. Always cite federal case law." > prompts/system.txt

# 2. Run evo diff to view precise line changes
evo diff

# Output will highlight exact line additions/deletions in red/green!
```

---

## 5. Repository On-Disk Layout

Evolution repositories maintain a clean, local-first storage layout under `.evolution/`:

```text
.evolution/
├── HEAD                  # Contains active branch name (e.g. "main")
├── config.json           # Repository configuration
├── index                 # Staging area JSON index
├── branches/             # Branch pointer files
│   ├── main              # Stores current commit ID for main branch
│   └── exp-claude-sonnet # Stores commit ID for experimental branch
├── commits/              # Intelligence Commit JSON snapshots
│   └── <commit_id>.json  # Holds message, author, timestamp, tree_hash, artifacts, metadata
├── objects/              # Content-Addressable Storage (SHA-256 sharded)
│   ├── b6/
│   │   └── 71051cad...   # Blob object (file content)
│   └── aa/
│       └── 0237644a...   # Tree object (directory structure)
└── artifacts/            # Typed Artifact JSON objects
    ├── prompt/
    ├── model_config/
    ├── retrieval/
    ├── tool/
    ├── memory/
    └── policy/
```
