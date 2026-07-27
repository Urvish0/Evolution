# Changelog

All notable changes to the Evolution platform are documented in this file.

## [v0.3.0] - Object Model & Content-Addressable Storage (In Progress)

### Added
- SHA-256 Object Hashing Engine (`HashContent`, `HashRaw`) with Git-style header prefix (`<type> <len>\0<content>`).
- Content-Addressable Blob Storage (`WriteBlob`, `ReadBlob`, `HasBlob`) with 2-character directory sharding (`.evolution/objects/xx/yyyy...`).
- Automatic Blob Content Deduplication (skips re-writing existing SHA-256 hashes).
- Merkle Tree Object Engine (`WriteTree`, `ReadTree`, `SerializeTree`, `ParseTree`).
- Recursive Merkle Tree workspace creation (`BuildTreeFromDirectory`).

---

## [v0.2.0] - Repository Engine

### Added
- User Configuration System (`evo config set user.name/user.email`, `evo config list`).
- Auto-attached commit `Author` identity from global user config (`~/.evolution/config.json`).
- Full Branch Management CLI (`evo branch`, `evo branch -n <name>`, `evo checkout <name>`, `evo branch -d <name>`).
- Dynamic active branch resolution via `.evolution/HEAD`.
- Enhanced Log CLI (`evo log --oneline`, `evo log -n <limit>`, 8-character short commit IDs, ANSI color output).
- Full commit chain DAG traversal walking parent pointers to genesis.
- Contextual error wrapping (`fmt.Errorf`) across all repository functions.
- Automated Go unit test suite (`config_test`, `commit_test`, `branch_test`, `init_test`, `log_test`, `hash_test`, `blob_test`, `tree_test`).
- GitHub Actions CI/CD workflow (`.github/workflows/ci.yml`) and build status badge.
- Comprehensive user CLI guide ([USAGE.md](USAGE.md)).

---

## [v0.1.0] - Bootstrap

### Added
- Cobra CLI foundation (`evo`).
- Version command (`evo version`).
- Repository initialization (`evo init`).
- Basic repository status (`evo status`).
- Initial local repository storage layout (`.evolution/`).
