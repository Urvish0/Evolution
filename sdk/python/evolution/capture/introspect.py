"""
Introspection and duck-typing utilities for LLM responses and function signatures.
"""

from __future__ import annotations

import inspect
from typing import Any, Callable

from evolution.models.artifacts import ModelConfigArtifact, PromptArtifact
from evolution.models.execution import TokenUsage


def extract_inputs_from_args(func: Callable, args: tuple, kwargs: dict) -> str:
    """Formats function arguments into a readable input representation."""
    sig = inspect.signature(func)
    try:
        bound = sig.bind(*args, **kwargs)
        bound.apply_defaults()
        arg_dict = bound.arguments
    except Exception:
        arg_dict = {"args": args, "kwargs": kwargs}

    # If only one string argument, return it directly
    if len(arg_dict) == 1:
        val = next(iter(arg_dict.values()))
        if isinstance(val, str):
            return val

    # Otherwise format key-value pairs
    items = []
    for k, v in arg_dict.items():
        if k in ("self", "cls"):
            continue
        items.append(f"{k}={v!r}")
    return ", ".join(items) if items else "(no input arguments)"


def extract_docstring_prompt(func: Callable, name: str | None = None) -> PromptArtifact | None:
    """Extracts function docstring as a system prompt artifact if present."""
    doc = inspect.getdoc(func)
    if doc and doc.strip():
        art_name = name or f"{func.__name__}-prompt"
        return PromptArtifact(
            name=art_name,
            role="system",
            description=f"Extracted docstring prompt for {func.__name__}",
        )
    return None


def extract_model_config_from_kwargs(kwargs: dict[str, Any], name: str | None = None) -> ModelConfigArtifact | None:
    """Detects model and generation parameters passed in kwargs."""
    model = kwargs.get("model") or kwargs.get("model_name")
    if not model or not isinstance(model, str):
        return None

    provider = kwargs.get("provider", "openai")
    if provider not in ("openai", "anthropic", "google", "local", "mistral", "cohere", "aws_bedrock"):
        provider = "openai"

    temperature = kwargs.get("temperature")
    if temperature is not None:
        try:
            temperature = float(temperature)
        except (ValueError, TypeError):
            temperature = 0.7

    return ModelConfigArtifact(
        name=name or f"model-{model}",
        model=model,
        provider=provider,
        temperature=temperature,
        max_tokens=kwargs.get("max_tokens"),
        top_p=kwargs.get("top_p"),
    )


def extract_llm_response(result: Any) -> tuple[str, TokenUsage, str | None]:
    """Duck-types standard LLM response objects (OpenAI, Anthropic, LangChain, or dicts)
    to extract:
      - output text (str)
      - token usage (TokenUsage)
      - model name (str | None)
    """
    output_text = ""
    tokens = TokenUsage()
    model_name: str | None = None

    if result is None:
        return "", tokens, None

    # If it's already a simple string
    if isinstance(result, str):
        return result, tokens, None

    # Check for dict response
    if isinstance(result, dict):
        # 1. Output extraction
        if "choices" in result and isinstance(result["choices"], list) and len(result["choices"]) > 0:
            first = result["choices"][0]
            if isinstance(first, dict):
                msg = first.get("message", {})
                output_text = msg.get("content", "") if isinstance(msg, dict) else str(first.get("text", ""))
        elif "content" in result:
            if isinstance(result["content"], list) and len(result["content"]) > 0:
                first = result["content"][0]
                output_text = first.get("text", "") if isinstance(first, dict) else str(first)
            else:
                output_text = str(result["content"])
        elif "output" in result:
            output_text = str(result["output"])
        elif "text" in result:
            output_text = str(result["text"])

        # 2. Token usage extraction
        usage = result.get("usage", {})
        if isinstance(usage, dict):
            tokens.prompt_tokens = usage.get("prompt_tokens") or usage.get("input_tokens") or 0
            tokens.completion_tokens = usage.get("completion_tokens") or usage.get("output_tokens") or 0
            tokens.total_tokens = usage.get("total_tokens") or (tokens.prompt_tokens + tokens.completion_tokens)

        # 3. Model name
        model_name = result.get("model")

        if not output_text:
            output_text = str(result)
        return output_text, tokens, model_name

    # Check for Object response (OpenAI ChatCompletion, Anthropic Message, etc.) via attributes
    # 1. Output extraction
    if hasattr(result, "choices") and result.choices:
        first = result.choices[0]
        if hasattr(first, "message") and hasattr(first.message, "content"):
            output_text = str(first.message.content or "")
        elif hasattr(first, "text"):
            output_text = str(first.text or "")
    elif hasattr(result, "content"):
        c = result.content
        if isinstance(c, list) and len(c) > 0:
            first = c[0]
            output_text = getattr(first, "text", str(first))
        else:
            output_text = str(c)
    elif hasattr(result, "output"):
        output_text = str(result.output)
    elif hasattr(result, "text"):
        output_text = str(result.text)

    # 2. Token usage extraction
    if hasattr(result, "usage") and result.usage:
        u = result.usage
        tokens.prompt_tokens = getattr(u, "prompt_tokens", getattr(u, "input_tokens", 0)) or 0
        tokens.completion_tokens = getattr(u, "completion_tokens", getattr(u, "output_tokens", 0)) or 0
        tokens.total_tokens = getattr(u, "total_tokens", tokens.prompt_tokens + tokens.completion_tokens) or (tokens.prompt_tokens + tokens.completion_tokens)

    # 3. Model name
    if hasattr(result, "model") and result.model:
        model_name = str(result.model)

    if not output_text:
        output_text = str(result)

    return output_text, tokens, model_name
