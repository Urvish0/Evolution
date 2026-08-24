"""
Decorator-based telemetry and intelligence tracking.
"""

from __future__ import annotations

import asyncio
import functools
import inspect
import time
from pathlib import Path
from typing import Any, Callable

from evolution.capture.introspect import (
    extract_docstring_prompt,
    extract_inputs_from_args,
    extract_llm_response,
    extract_model_config_from_kwargs,
)
from evolution.models.artifacts import ModelConfigArtifact
from evolution.models.execution import Execution
from evolution.repository import Repository


def track(
    _func: Callable | None = None,
    *,
    name: str | None = None,
    repo: Repository | Path | str | None = None,
    model: str | None = None,
    provider: str = "openai",
    temperature: float | None = None,
    auto_manifest: bool = True,
    metadata: dict[str, Any] | None = None,
):
    """Decorator to automatically capture intelligence artifacts and record executions.

    Can be used as `@track` or `@track(...)`.
    """
    def decorator(func: Callable) -> Callable:
        # Resolve target repository
        target_repo: Repository | None
        if isinstance(repo, Repository):
            target_repo = repo
        elif repo is not None:
            try:
                target_repo = Repository.open(repo)
            except Exception:
                target_repo = None
        else:
            try:
                target_repo = Repository.open(".")
            except Exception:
                target_repo = None

        func_name = name or func.__name__

        # Auto-update manifest with docstring prompt and model config if enabled
        if auto_manifest and target_repo is not None:
            try:
                manifest = target_repo.get_manifest()
                updated = False

                # 1. Capture docstring as prompt artifact
                prompt_art = extract_docstring_prompt(func, name=f"{func_name}-prompt")
                if prompt_art:
                    manifest.add_artifact(prompt_art)
                    updated = True

                # 2. Capture explicit model config if specified
                if model:
                    mc_art = ModelConfigArtifact(
                        name=f"{func_name}-model",
                        model=model,
                        provider=provider,
                        temperature=temperature,
                    )
                    manifest.add_artifact(mc_art)
                    updated = True

                if updated:
                    target_repo.save_manifest(manifest, auto_hash=True)
            except Exception:
                pass

        if inspect.iscoroutinefunction(func):
            @functools.wraps(func)
            async def async_wrapper(*args, **kwargs):
                inputs_str = extract_inputs_from_args(func, args, kwargs)
                start_time = time.perf_counter()
                exec_meta = dict(metadata or {})
                exec_meta["function"] = func.__name__

                try:
                    result = await func(*args, **kwargs)
                    elapsed = time.perf_counter() - start_time
                    duration_ms = max(1, int(elapsed * 1000))

                    output_text, tokens, detected_model = extract_llm_response(result)
                    if detected_model:
                        exec_meta["model"] = detected_model
                    elif model:
                        exec_meta["model"] = model

                    if target_repo is not None:
                        exec_obj = target_repo.record_execution(
                            inputs=inputs_str,
                            outputs=output_text,
                            duration_ms=duration_ms,
                            prompt_tokens=tokens.prompt_tokens,
                            completion_tokens=tokens.completion_tokens,
                            metadata=exec_meta,
                        )
                        setattr(async_wrapper, "last_execution", exec_obj)

                    return result
                except Exception as exc:
                    elapsed = time.perf_counter() - start_time
                    duration_ms = max(1, int(elapsed * 1000))
                    exec_meta["error"] = str(exc)
                    exec_meta["error_type"] = type(exc).__name__

                    if target_repo is not None:
                        target_repo.record_execution(
                            inputs=inputs_str,
                            outputs=f"Error: {exc}",
                            duration_ms=duration_ms,
                            metadata=exec_meta,
                        )
                    raise

            return async_wrapper

        else:
            @functools.wraps(func)
            def sync_wrapper(*args, **kwargs):
                inputs_str = extract_inputs_from_args(func, args, kwargs)
                start_time = time.perf_counter()
                exec_meta = dict(metadata or {})
                exec_meta["function"] = func.__name__

                # Check if model passed in runtime kwargs
                detected_kwarg_mc = extract_model_config_from_kwargs(kwargs)
                if detected_kwarg_mc and auto_manifest and target_repo is not None:
                    try:
                        manifest = target_repo.get_manifest()
                        manifest.add_artifact(detected_kwarg_mc)
                        target_repo.save_manifest(manifest, auto_hash=True)
                    except Exception:
                        pass

                try:
                    result = func(*args, **kwargs)
                    elapsed = time.perf_counter() - start_time
                    duration_ms = max(1, int(elapsed * 1000))

                    output_text, tokens, detected_model = extract_llm_response(result)
                    if detected_model:
                        exec_meta["model"] = detected_model
                    elif model:
                        exec_meta["model"] = model

                    if target_repo is not None:
                        exec_obj = target_repo.record_execution(
                            inputs=inputs_str,
                            outputs=output_text,
                            duration_ms=duration_ms,
                            prompt_tokens=tokens.prompt_tokens,
                            completion_tokens=tokens.completion_tokens,
                            metadata=exec_meta,
                        )
                        setattr(sync_wrapper, "last_execution", exec_obj)

                    return result
                except Exception as exc:
                    elapsed = time.perf_counter() - start_time
                    duration_ms = max(1, int(elapsed * 1000))
                    exec_meta["error"] = str(exc)
                    exec_meta["error_type"] = type(exc).__name__

                    if target_repo is not None:
                        target_repo.record_execution(
                            inputs=inputs_str,
                            outputs=f"Error: {exc}",
                            duration_ms=duration_ms,
                            metadata=exec_meta,
                        )
                    raise

            return sync_wrapper

    if _func is None:
        return decorator
    return decorator(_func)
