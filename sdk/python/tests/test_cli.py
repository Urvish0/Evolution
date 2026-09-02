"""
Unit tests for the Evolution Python CLI (`evo-py`).
"""

from __future__ import annotations

import io
from contextlib import redirect_stdout
from pathlib import Path

from evolution.cli import main


def test_cli_version(capsys):
    try:
        main(["--version"])
    except SystemExit as e:
        assert e.code == 0
    captured = capsys.readouterr()
    assert "evolution-sdk" in captured.out or "evolution-sdk" in captured.err


def test_cli_init_and_status(tmp_path: Path, capsys):
    repo_dir = tmp_path / "test_repo"

    # 1. init
    rc = main(["init", str(repo_dir), "--name", "cli-test-agent"])
    assert rc == 0
    captured = capsys.readouterr()
    assert "Initialized Evolution repository" in captured.out
    assert (repo_dir / ".evolution").is_dir()

    # 2. status
    rc = main(["-C", str(repo_dir), "status"])
    assert rc == 0
    captured = capsys.readouterr()
    assert "On branch main" in captured.out
    assert "Recorded Executions: 0" in captured.out


def test_cli_manifest_operations(tmp_path: Path, capsys):
    repo_dir = tmp_path / "manifest_repo"
    main(["init", str(repo_dir), "--name", "manifest-agent"])
    capsys.readouterr()

    # 1. manifest show
    rc = main(["-C", str(repo_dir), "manifest", "show"])
    assert rc == 0
    captured = capsys.readouterr()
    assert '"name": "manifest-agent"' in captured.out

    # 2. manifest validate
    rc = main(["-C", str(repo_dir), "manifest", "validate"])
    assert rc == 0
    captured = capsys.readouterr()
    assert "is valid according to Intelligence Manifest Spec v1.0" in captured.out


def test_cli_execution_and_evaluate_empty(tmp_path: Path, capsys):
    repo_dir = tmp_path / "eval_repo"
    main(["init", str(repo_dir)])
    capsys.readouterr()

    # execution list
    rc = main(["-C", str(repo_dir), "execution", "list"])
    assert rc == 0
    captured = capsys.readouterr()
    assert "No executions recorded" in captured.out

    # evaluate
    rc = main(["-C", str(repo_dir), "evaluate"])
    assert rc == 0
    captured = capsys.readouterr()
    assert "No evaluations recorded" in captured.out
