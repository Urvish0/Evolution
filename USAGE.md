# Evolution CLI — Usage Guide

> Quick reference for all commands and workflows supported by the `evo` CLI.

---

## Installation & Setup

### Option 1: Run directly with Go
```bash
go run ./cmd/evo <command>
```

### Option 2: Install globally as `evo` binary
```bash
go install ./cmd/evo
```
*Note: Ensure `%USERPROFILE%\go\bin` (Windows) or `$GOPATH/bin` (Linux/macOS) is in your system `PATH` to run `evo` directly from any terminal.*

---

## Command Reference

### 1. User Identity & Configuration (`evo config`)

Set developer identity (stored globally in `~/.evolution/config.json`). This identity is automatically attached as `Author` to every commit you make.

```bash
# Set user name
evo config set user.name "Urvish"

# Set user email
evo config set user.email "urvish@example.com"

# View current configuration
evo config list
```

---

### 2. Repository Initialization (`evo init`)

Initializes a local `.evolution/` version control repository in the current directory.

```bash
evo init
```

*Creates `.evolution/` folder structure containing config, default branch (`main`), commit store, object store, and artifacts directory.*

---

### 3. Repository Status (`evo status`)

Displays the current repository state: initialized status, active branch, commit count, and working tree health.

```bash
evo status
```

---

### 4. Staging Files (`evo add`)

Stages specific files or directories into the Staging Area (`.evolution/index`) before committing.

```bash
# Stage a specific file
evo add prompt.txt

# Stage a directory recursively
evo add tools/

# Stage all files in current directory
evo add .
```

---

### 5. Creating Intelligence Commits (`evo commit`)

Creates an immutable commit snapshot. Automatically captures staged files from the Staging Area (`.evolution/index`) as a Merkle tree of hashed Blobs (`.evolution/objects/xx/yyyy...`), auto-detects `evolution.manifest.json`, and links attached artifacts, tags (`--tag`), and execution metadata (`--meta key=val`).

```bash
# Basic commit
evo commit -m "Initial intelligence setup"

# Commit with tags and custom execution metadata
evo commit -m "Deployed v1 prompt" --tag production --tag v1.0 --meta model=gpt-4o --meta temp=0.2
```

---

### 6. Viewing History (`evo log`)

Inspects commit history from HEAD back to the initial (genesis) commit.

```bash
# Standard history view (with short 8-char commit IDs and colors)
evo log

# Compact single-line format
evo log --oneline

# Limit output to N commits
evo log -n 5

# Combine flags
evo log --oneline -n 3
```

---

### 7. Branch Management (`evo branch` & `evo checkout`)

Create, list, switch, and delete branches for parallel intelligence experimentation.

```bash
# List all local branches (* indicates current active branch)
evo branch

# Create a new branch
evo branch -n feature-prompt
# OR
evo branch feature-prompt

# Switch active branch
evo checkout feature-prompt

# Delete a branch
evo branch -d feature-prompt
```

---

### 8. Working Tree Comparison (`evo diff`)

Compares the current workspace files against the last committed Merkle tree snapshot to detect new, modified, and deleted files.

```bash
# Show all changes vs HEAD
evo diff
```

*Output shows `[staged]`, `[modified]`, `[deleted]`, and `[new]` file labels with color coding.*

---

### 9. Intelligence Manifest (`evo manifest`)

Manage the `evolution.manifest.json` file conforming to the Intelligence Manifest Specification v0.1.

```bash
# Generate starter manifest file
evo manifest init --name "my-ai-agent" --description "Legal assistant AI"

# Validate manifest against Spec v0.1 schema
evo manifest validate

# Display active workspace manifest JSON
evo manifest show
```

---

### 10. Typed Artifacts (`evo artifact`)

Register, list, inspect, and perform high-level semantic diffs on typed AI artifacts.

```bash
# Register an artifact in evolution.manifest.json
evo artifact add prompt sys-prompt prompts/system.txt

# List all AI artifacts attached to the current HEAD commit
evo artifact list

# Display JSON content of a stored artifact
evo artifact show prompt b671051c

# Perform high-level semantic diff between two commits
evo artifact diff <commit1_id> <commit2_id>
```

---

### 11. Version & Help

```bash
# Display CLI version
evo version

# Display help for any command
evo --help
evo commit --help
evo artifact --help
evo manifest --help
```

---

## Complete Example Workflow

```bash
# 1. Setup user identity
evo config set user.name "Urvish"
evo config set user.email "urvish@example.com"

# 2. Initialize repository
evo init

# 3. Create initial commit
evo commit -m "v1 prompt setup"

# 4. Create and switch to experimental branch
evo branch -n exp-claude-model
evo checkout exp-claude-model

# 5. Commit experimental changes
evo commit -m "Switched model to Claude 3.5 Sonnet"

# 6. View single-line history
evo log --oneline

# 7. Switch back to main branch
evo checkout main

# 8. Check status
evo status
```
