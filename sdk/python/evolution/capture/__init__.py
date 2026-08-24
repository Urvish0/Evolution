"""
Automatic intelligence capture and execution recording package.
"""

from evolution.capture.introspect import (
    extract_docstring_prompt,
    extract_inputs_from_args,
    extract_llm_response,
    extract_model_config_from_kwargs,
)
from evolution.capture.recorder import RecordContextManager, record
from evolution.capture.tracker import track

__all__ = [
    "RecordContextManager",
    "extract_docstring_prompt",
    "extract_inputs_from_args",
    "extract_llm_response",
    "extract_model_config_from_kwargs",
    "record",
    "track",
]
