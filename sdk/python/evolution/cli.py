"""
Command-line interface for Evolution Python SDK (`evo-py` / `evolution`).
Allows developers to manage Evolution repositories without requiring a Go runtime.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Sequence

import evolution as evo
from evolution.models.manifest import MANIFEST_FILE_NAME


def cmd_init(args: argparse.Namespace) -> int:
    target = Path(args.path).resolve()
    try:
        repo = evo.Repository.init(target, name=args.name)
        print(f"Initialized Evolution repository in {repo.evolution_dir}")
        return 0
    except Exception as e:
        print(f"[ERROR] Failed to initialize repository: {e}", file=sys.stderr)
        return 1


def cmd_status(args: argparse.Namespace) -> int:
    try:
        repo = evo.Repository.open(args.path)
        status = repo.status()
        print(f"On branch {status.branch}")
        
        manifest_path = repo.root / MANIFEST_FILE_NAME
        if manifest_path.is_file():
            try:
                m = repo.get_manifest()
                prompt_count = len(m.artifacts.prompts)
                tool_count = len(m.artifacts.tools)
                model_name = m.artifacts.model_config.model if m.artifacts.model_config else "None"
                print(f"Manifest: {m.name} (v{m.version}) | Model: {model_name} | Prompts: {prompt_count} | Tools: {tool_count}")
            except Exception as e:
                print(f"Manifest: invalid ({e})")
        else:
            print("Manifest: not found (run 'evo-py manifest init' or track an agent)")

        # Count executions
        exec_dir = repo.evolution_dir / "executions"
        exec_count = len(list(exec_dir.glob("*.json"))) if exec_dir.is_dir() else 0
        print(f"Recorded Executions: {exec_count}")

        # Count evaluations
        eval_dir = repo.evolution_dir / "evaluations"
        eval_count = len(list(eval_dir.glob("*.json"))) if eval_dir.is_dir() else 0
        print(f"Evaluations: {eval_count}")

        if status.is_clean:
            print("\nnothing to commit, working tree clean")
        return 0
    except Exception as e:
        print(f"[ERROR] {e}", file=sys.stderr)
        return 1


def cmd_log(args: argparse.Namespace) -> int:
    try:
        repo = evo.Repository.open(args.path)
        commits_dir = repo.evolution_dir / "commits"
        if not commits_dir.is_dir():
            print("No commits found.")
            return 0

        commit_files = sorted(commits_dir.glob("*.json"), key=lambda p: p.stat().st_mtime, reverse=True)
        if not commit_files:
            print("No commits found.")
            return 0

        limit = args.limit or len(commit_files)
        for cf in commit_files[:limit]:
            try:
                data = json.loads(cf.read_text(encoding="utf-8"))
                cid = data.get("id", cf.stem)
                msg = data.get("message", "")
                author = data.get("author", "unknown")
                ts = data.get("timestamp", "")
                meta = data.get("metadata", {})
                tags = meta.get("tags", [])
                tag_str = f" ({', '.join(tags)})" if tags else ""

                if args.oneline:
                    print(f"{cid[:8]}{tag_str} {msg}")
                else:
                    print(f"commit {cid}{tag_str}")
                    print(f"Author: {author}")
                    print(f"Date:   {ts}")
                    print(f"\n    {msg}\n")
            except Exception:
                continue
        return 0
    except Exception as e:
        print(f"[ERROR] {e}", file=sys.stderr)
        return 1


def cmd_manifest(args: argparse.Namespace) -> int:
    try:
        repo = evo.Repository.open(args.path)
        if args.manifest_action == "show":
            m = repo.get_manifest()
            print(json.dumps(m.to_dict(), indent=2))
            return 0
        elif args.manifest_action == "validate":
            m = repo.get_manifest()
            m.validate()
            print(f"Manifest '{m.name}' is valid according to Intelligence Manifest Spec v1.0.")
            return 0
        elif args.manifest_action == "init":
            m = evo.Manifest(name=args.name or "my-ai-intelligence", description="Standard Evolution Manifest")
            repo.save_manifest(m)
            print(f"Created standard {MANIFEST_FILE_NAME} in {repo.root}")
            return 0
        else:
            print("Usage: evo-py manifest [show|validate|init]")
            return 1
    except Exception as e:
        print(f"[ERROR] {e}", file=sys.stderr)
        return 1


def cmd_execution(args: argparse.Namespace) -> int:
    try:
        repo = evo.Repository.open(args.path)
        exec_dir = repo.evolution_dir / "executions"

        if args.exec_action == "list":
            if not exec_dir.is_dir():
                print("No executions recorded.")
                return 0

            files = sorted(exec_dir.glob("*.json"), key=lambda p: p.stat().st_mtime, reverse=True)
            if not files:
                print("No executions recorded.")
                return 0

            print(f"{'ID':<10} {'COMMIT':<10} {'TOKENS':<10} {'DURATION':<12} {'TIMESTAMP'}")
            print("-" * 65)
            for f in files:
                try:
                    d = json.loads(f.read_text(encoding="utf-8"))
                    eid = d.get("id", f.stem)[:8]
                    cid = str(d.get("commit_id", "none"))[:8]
                    tok = d.get("tokens", {}).get("total_tokens", 0)
                    dur = f"{d.get('duration_ms', 0)}ms"
                    ts = d.get("timestamp", "")
                    print(f"{eid:<10} {cid:<10} {tok:<10} {dur:<12} {ts}")
                except Exception:
                    continue
            return 0
        elif args.exec_action == "show":
            if not args.exec_id:
                print("Please provide an execution ID: evo-py execution show <ID>")
                return 1
            exec_obj = repo.get_execution(args.exec_id)
            print(json.dumps(exec_obj.to_dict(), indent=2))
            return 0
        else:
            print("Usage: evo-py execution [list|show <ID>]")
            return 1
    except Exception as e:
        print(f"[ERROR] {e}", file=sys.stderr)
        return 1


def cmd_evaluate(args: argparse.Namespace) -> int:
    try:
        repo = evo.Repository.open(args.path)
        eval_dir = repo.evolution_dir / "evaluations"

        if not eval_dir.is_dir():
            print("No evaluations recorded yet.")
            return 0

        files = sorted(eval_dir.glob("*.json"), key=lambda p: p.stat().st_mtime, reverse=True)
        if not files:
            print("No evaluations recorded yet.")
            return 0

        print(f"{'ID':<10} {'COMMIT':<10} {'SCORE':<10} {'TIMESTAMP'}")
        print("-" * 50)
        for f in files:
            try:
                d = json.loads(f.read_text(encoding="utf-8"))
                eid = d.get("id", f.stem)[:8]
                cid = str(d.get("commit_id", "head"))[:8]
                score = f"{d.get('overall_score', 0.0) * 100:.1f}%"
                ts = d.get("timestamp", "")
                print(f"{eid:<10} {cid:<10} {score:<10} {ts}")
            except Exception:
                continue
        return 0
    except Exception as e:
        print(f"[ERROR] {e}", file=sys.stderr)
        return 1


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="evo-py",
        description="Evolution Python SDK CLI — AI-Native Version Control Platform",
    )
    parser.add_argument("--version", action="version", version=f"evolution-sdk {evo.__version__}")
    parser.add_argument("-C", "--path", default=".", help="Run as if evo-py was started in <path> (default: .)")

    subparsers = parser.add_subparsers(dest="command", help="Command to execute")

    # init
    init_parser = subparsers.add_parser("init", help="Initialize a new Evolution repository")
    init_parser.add_argument("path", nargs="?", default=".", help="Directory to initialize (default: .)")
    init_parser.add_argument("--name", default="ai-intelligence", help="Intelligence system name")

    # status
    subparsers.add_parser("status", help="Show working tree status and intelligence summary")

    # log
    log_parser = subparsers.add_parser("log", help="Show intelligence commit history")
    log_parser.add_argument("--oneline", action="store_true", help="Display one line per commit")
    log_parser.add_argument("-n", "--limit", type=int, default=None, help="Limit number of commits")

    # manifest
    manifest_parser = subparsers.add_parser("manifest", help="Inspect and validate intelligence manifests")
    manifest_parser.add_argument("manifest_action", choices=["show", "validate", "init"], default="show", nargs="?")
    manifest_parser.add_argument("--name", default=None, help="Manifest name when initializing")

    # execution
    exec_parser = subparsers.add_parser("execution", help="Inspect recorded execution telemetry")
    exec_parser.add_argument("exec_action", choices=["list", "show"], default="list", nargs="?")
    exec_parser.add_argument("exec_id", nargs="?", default=None, help="Execution ID for show")

    # evaluate
    subparsers.add_parser("evaluate", help="Inspect evaluation scores and reports")

    return parser


def main(argv: Sequence[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)

    if not args.command:
        parser.print_help()
        return 0

    dispatch = {
        "init": cmd_init,
        "status": cmd_status,
        "log": cmd_log,
        "manifest": cmd_manifest,
        "execution": cmd_execution,
        "evaluate": cmd_evaluate,
    }

    handler = dispatch.get(args.command)
    if handler:
        return handler(args)

    parser.print_help()
    return 1


if __name__ == "__main__":
    sys.exit(main())
