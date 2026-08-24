"""
Evolution framework adapters package.
"""

from evolution.adapters.base import BaseAdapter
from evolution.adapters.crewai import CrewAIAdapter, from_crewai
from evolution.adapters.direct import (
    AnthropicAdapter,
    OpenAIAdapter,
    from_anthropic,
    from_openai,
)
from evolution.adapters.langchain import LangChainAdapter, from_langchain
from evolution.adapters.llamaindex import LlamaIndexAdapter, from_llamaindex

__all__ = [
    "AnthropicAdapter",
    "BaseAdapter",
    "CrewAIAdapter",
    "LangChainAdapter",
    "LlamaIndexAdapter",
    "OpenAIAdapter",
    "from_anthropic",
    "from_crewai",
    "from_langchain",
    "from_llamaindex",
    "from_openai",
]
