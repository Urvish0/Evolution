"""
Unit tests for SemanticEvaluator (LLM-as-a-Judge) and evaluation reporting.
"""

from __future__ import annotations

import json
from evolution.evaluators import DimensionScore, EvaluationReport, SemanticEvaluator
from evolution.models.execution import Execution


def test_parse_clean_json():
    evaluator = SemanticEvaluator()
    judge_json = {
        "accuracy": {"score": 9, "reasoning": "Highly accurate and grounded."},
        "helpfulness": {"score": 8, "reasoning": "Direct and well-structured."},
        "instruction_following": {"score": 10, "reasoning": "Followed all constraints."},
        "safety": {"score": 10, "reasoning": "Zero safety violations."},
        "summary": "Excellent performance overall.",
    }

    dims, summary = evaluator.parse_judge_response(json.dumps(judge_json))

    assert summary == "Excellent performance overall."
    assert "accuracy" in dims
    assert dims["accuracy"].raw_score == 9.0
    assert dims["accuracy"].score == 0.9
    assert dims["accuracy"].reasoning == "Highly accurate and grounded."
    assert dims["instruction_following"].score == 1.0


def test_parse_fenced_markdown_json():
    evaluator = SemanticEvaluator()
    fenced_output = """Here is my evaluation of the agent response:

```json
{
  "accuracy": {"score": 7, "reasoning": "Minor assumption made."},
  "helpfulness": {"score": 9, "reasoning": "Very thorough."},
  "instruction_following": {"score": 8, "reasoning": "Followed format."},
  "safety": {"score": 10, "reasoning": "Safe."},
  "summary": "Solid legal advice."
}
```

Hope this helps!"""

    dims, summary = evaluator.parse_judge_response(fenced_output)

    assert summary == "Solid legal advice."
    assert dims["accuracy"].score == 0.7
    assert dims["helpfulness"].score == 0.9
    assert dims["instruction_following"].score == 0.8


def test_evaluation_with_mock_judge():
    mock_response = {
        "accuracy": {"score": 9, "reasoning": "Statutory facts correct."},
        "helpfulness": {"score": 9, "reasoning": "Actionable notice letter."},
        "instruction_following": {"score": 10, "reasoning": "Maintained strict tone."},
        "safety": {"score": 10, "reasoning": "Compliant."},
        "summary": "Senior-tier legal output.",
    }

    def mock_judge_fn(sys_prompt: str, user_prompt: str):
        assert "AGENT SYSTEM PROMPT" in user_prompt
        assert "USER INPUT" in user_prompt
        return json.dumps(mock_response)

    evaluator = SemanticEvaluator(
        judge_fn=mock_judge_fn,
        judge_model="mock-qwen-judge",
    )

    report = evaluator.evaluate(
        inputs="Enterprise dispute with CloudTech",
        outputs="Formal Notice of Material Breach...",
        agent_prompt="You are a Senior Litigation Attorney.",
        execution_id="exec-12345",
        commit_id="commit-abcde",
    )

    assert report.evaluator == "semantic-llm-judge"
    assert report.judge_model == "mock-qwen-judge"
    assert report.execution_id == "exec-12345"
    assert report.commit_id == "commit-abcde"
    # Weighted average: (0.9 + 0.9 + 1.0 + 1.0) / 4 = 0.95
    assert report.overall_score == 0.95
    assert report.judge_reasoning == "Senior-tier legal output."


def test_custom_weights():
    # Weigh accuracy 3x higher than other dimensions
    weights = {
        "accuracy": 3.0,
        "helpfulness": 1.0,
        "instruction_following": 1.0,
        "safety": 1.0,
    }

    mock_response = {
        "accuracy": {"score": 10, "reasoning": "Flawless."},     # 1.0 * 3 = 3.0
        "helpfulness": {"score": 5, "reasoning": "Average."},      # 0.5 * 1 = 0.5
        "instruction_following": {"score": 5, "reasoning": "Ok."}, # 0.5 * 1 = 0.5
        "safety": {"score": 10, "reasoning": "Safe."},            # 1.0 * 1 = 1.0
    }
    # Total weight = 6.0, Weighted sum = 5.0 -> 5.0 / 6.0 = 0.8333

    def mock_judge(s: str, u: str):
        return mock_response

    evaluator = SemanticEvaluator(judge_fn=mock_judge, weights=weights)
    report = evaluator.evaluate(inputs="in", outputs="out")

    assert report.overall_score == 0.8333


def test_evaluate_execution_helper():
    mock_response = {
        "accuracy": {"score": 8, "reasoning": "Grounded."},
        "helpfulness": {"score": 8, "reasoning": "Clear."},
        "instruction_following": {"score": 8, "reasoning": "Good."},
        "safety": {"score": 8, "reasoning": "Safe."},
        "summary": "Good run.",
    }

    evaluator = SemanticEvaluator(
        judge_fn=lambda s, u: mock_response,
        judge_model="test-judge",
    )

    execution = Execution(
        id="exec-999",
        commit_id="commit-888",
        inputs="Explain quantum state",
        outputs="A quantum state represents...",
    )

    report = evaluator.evaluate_execution(execution, agent_prompt="Be clear.")
    assert report.execution_id == "exec-999"
    assert report.commit_id == "commit-888"
    assert report.overall_score == 0.8
    report_dict = report.to_dict()
    assert report_dict["evaluator"] == "semantic-llm-judge"
    assert isinstance(report_dict["dimensions"]["accuracy"], dict)
