"""
Live LLM-as-a-Judge Semantic Evaluation Demo for Evolution.

Uses Groq LPU with qwen/qwen3.8-27b as an impartial judge to score Agent 1
(Strict Legal Analyst) vs Agent 2 (Creative Negotiator) across 4 dimensions:
1. Accuracy
2. Helpfulness
3. Instruction Following
4. Safety & Guardrails
"""

from __future__ import annotations

import os
import sys
from pathlib import Path
from dotenv import load_dotenv
from groq import Groq

import evolution as evo

# Force UTF-8 on Windows console
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

# 1. Load environment variables
env_path = Path(__file__).parent / ".env"
load_dotenv(env_path)

api_key = os.environ.get("GROQ_API_KEY")
if not api_key:
    print("⚠️ GROQ_API_KEY not found in playground/.env.")
    sys.exit(1)

client = Groq(api_key=api_key)


def groq_judge_fn(system_prompt: str, user_prompt: str) -> str:
    """Invokes Groq LPU model as an impartial judge."""
    response = client.chat.completions.create(
        model="qwen/qwen3.8-27b",
        temperature=0.0,  # Deterministic judging
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_prompt},
        ],
        response_format={"type": "json_object"},
    )
    return response.choices[0].message.content


def main():
    print("=" * 80)
    print("⚖️  EVOLUTION SEMANTIC EVALUATION: LLM-AS-A-JUDGE (GROQ LPU)")
    print("=" * 80)

    repo_dir = Path(__file__).parent / "workspace"
    repo = evo.Repository.open(repo_dir)

    # Initialize Semantic Evaluator with Groq judge
    judge = evo.SemanticEvaluator(
        judge_fn=groq_judge_fn,
        judge_model="qwen/qwen3.8-27b (Groq LPU)",
    )

    # Import both agents
    from agents.strict_analyst import create_strict_analyst
    from agents.creative_negotiator import create_creative_negotiator

    strict_analyst = create_strict_analyst(repo)
    creative_negotiator = create_creative_negotiator(repo)

    dispute_facts = (
        "Enterprise client contracted CloudTech for custom ERP software delivery by Jan 15 with a $500,000 budget. "
        "CloudTech delivered on March 1 (45 days late) with 3 critical unresolved database bugs causing $120,000 "
        "in operational losses. Client refuses to pay final $150,000 milestone invoice and threatens lawsuit."
    )

    print("\n🚀 Executing Agent 1 (Strict Litigation Analyst)...")
    res1 = strict_analyst(dispute_facts)
    exec1 = getattr(strict_analyst, "last_execution", None)

    print("🚀 Executing Agent 2 (Creative Strategic Negotiator)...")
    res2 = creative_negotiator(dispute_facts)
    exec2 = getattr(creative_negotiator, "last_execution", None)

    if not exec1 or not exec2:
        print("❌ Could not capture executions.")
        return

    print("\n" + "=" * 80)
    print("🧑‍⚖️  RUNNING IMPARTIAL LLM JUDGE EVALUATION...")
    print("=" * 80)

    print("\nEvaluating Agent 1...")
    report1 = judge.evaluate(
        inputs=exec1.inputs,
        outputs=exec1.outputs,
        agent_prompt="You are a Senior Corporate Litigation Attorney. Analyze dispute facts strictly under contract law.",
        execution_id=exec1.id,
        commit_id=exec1.commit_id,
    )
    repo.save_evaluation(
        evo.EvaluationResult(
            commit_id=exec1.commit_id or "head-main",
            execution_id=exec1.id,
            overall_score=report1.overall_score,
            scores={
                k: evo.EvaluationScore(name=k, score=v.score, details=v.reasoning)
                for k, v in report1.dimensions.items()
            },
        )
    )

    print("Evaluating Agent 2...")
    report2 = judge.evaluate(
        inputs=exec2.inputs,
        outputs=exec2.outputs,
        agent_prompt="You are a Master Commercial Negotiator and Mediator. Propose win-win solutions.",
        execution_id=exec2.id,
        commit_id=exec2.commit_id,
    )
    repo.save_evaluation(
        evo.EvaluationResult(
            commit_id=exec2.commit_id or "head-exp",
            execution_id=exec2.id,
            overall_score=report2.overall_score,
            scores={
                k: evo.EvaluationScore(name=k, score=v.score, details=v.reasoning)
                for k, v in report2.dimensions.items()
            },
        )
    )

    # Print Scorecard
    print("\n" + "=" * 80)
    print("📊 SIDE-BY-SIDE SEMANTIC JUDGMENT SCORECARD")
    print("=" * 80)
    print(f"Judge Model: {report1.judge_model}\n")

    header = f"{'Dimension':<25} | {'Agent 1 (Strict Analyst)':<24} | {'Agent 2 (Negotiator)':<24}"
    print(header)
    print("-" * len(header))

    all_dims = sorted(set(list(report1.dimensions.keys()) + list(report2.dimensions.keys())))
    for dim_name in all_dims:
        d1 = report1.dimensions.get(dim_name)
        d2 = report2.dimensions.get(dim_name)

        s1_str = f"{d1.raw_score:.1f}/10 ({d1.score*100:.0f}%)" if d1 else "N/A"
        s2_str = f"{d2.raw_score:.1f}/10 ({d2.score*100:.0f}%)" if d2 else "N/A"

        print(f"{dim_name.replace('_', ' ').title():<25} | {s1_str:<24} | {s2_str:<24}")

    print("-" * len(header))
    print(f"{'OVERALL SCORE':<25} | {report1.overall_score*100:.1f}%{'':<19} | {report2.overall_score*100:.1f}%")
    print("=" * 80)

    print("\n📝 Judge Summary for Agent 1 (Strict Analyst):")
    print(f"   \"{report1.judge_reasoning}\"")

    print("\n📝 Judge Summary for Agent 2 (Creative Negotiator):")
    print(f"   \"{report2.judge_reasoning}\"")

    print("\n🔒 Evaluations permanently persisted in .evolution/evaluations/ directory.")
    print("=" * 80)


if __name__ == "__main__":
    main()
