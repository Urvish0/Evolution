"""
Evolution Full-Spectrum Demonstration Lab

Exercises 100% of Evolution's capabilities:
1. 4 Specialized AI Agents (Litigation Attorney, Commercial Mediator, Tech Architect, Risk Officer)
2. Live Groq LPU inference with @evo.track
3. Pure Python SDK operations (Repository, Manifest, Artifacts, Framework Adapters)
4. LLM-as-a-Judge Semantic Evaluation with real quality scorecards
5. Full Go CLI command suite execution (init, commit, log, diff, branch, merge, replay, evaluate)
"""

from __future__ import annotations

import os
import shutil
import subprocess
import sys
from pathlib import Path
from dotenv import load_dotenv
from groq import Groq

import evolution as evo
from evolution.adapters import (
    AnthropicAdapter,
    CrewAIAdapter,
    LangChainAdapter,
    LlamaIndexAdapter,
    OpenAIAdapter,
)

if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

# 1. Environment & API Key
env_path = Path(__file__).parent / ".env"
load_dotenv(env_path)

api_key = os.environ.get("GROQ_API_KEY")
if not api_key:
    print("[ERROR] GROQ_API_KEY not found in playground/.env.")
    sys.exit(1)

client = Groq(api_key=api_key)

import stat

def remove_readonly(func, path, excinfo):
    os.chmod(path, stat.S_IWRITE)
    func(path)

ROOT_DIR = Path(__file__).resolve().parent.parent
SANDBOX_DIR = Path(__file__).resolve().parent / "sandbox_workspace"
if SANDBOX_DIR.exists():
    shutil.rmtree(SANDBOX_DIR, onerror=remove_readonly)
SANDBOX_DIR.mkdir(parents=True, exist_ok=True)

EVO_BIN = ROOT_DIR / "evo.exe"
BUILD_CMD = ["go", "build", "-o", str(EVO_BIN), "./cmd/evo"]


def run_cli(*args: str) -> str:
    """Helper to run Go CLI binary commands inside the sandbox."""
    res = subprocess.run(
        [str(EVO_BIN), *args],
        cwd=str(SANDBOX_DIR),
        capture_output=True,
        text=True,
    )
    if res.returncode != 0 and "conflict" not in res.stderr.lower():
        return f"[STDERR]: {res.stderr.strip()}"
    return res.stdout.strip()


def groq_judge_fn(system_prompt: str, user_prompt: str) -> str:
    """Judge model invocation for LLM-as-a-Judge semantic grading."""
    res = client.chat.completions.create(
        model="qwen/qwen3.8-27b",
        temperature=0.0,
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_prompt},
        ],
        response_format={"type": "json_object"},
    )
    return res.choices[0].message.content


def main():
    print("=" * 80)
    print("EVOLUTION FULL-SPECTRUM MASTER LAB & TESTBED")
    print("=" * 80)

    # 0. Build latest Go CLI binary
    print("\n[STEP 0] Building latest Go CLI binary (evo.exe)...")
    subprocess.run(BUILD_CMD, cwd=str(ROOT_DIR), check=True)
    print("  -> Built evo.exe successfully.")

    # 1. Initialize Repo via Python SDK & Go CLI
    print("\n[STEP 1] Initializing Repository...")
    repo = evo.init(SANDBOX_DIR, name="enterprise-incident-response-system")
    print(f"  -> Python SDK initialized repo at: {SANDBOX_DIR}")

    # Set user config via CLI
    run_cli("config", "set", "user.name", "Urvish")
    run_cli("config", "set", "user.email", "prajapatiurvish712@gmail.com")
    cfg_out = run_cli("config", "list")
    print(f"  -> CLI Configured:\n     {cfg_out}")

    # Dispute Scenario
    dispute_facts = (
        "Enterprise client contracted CloudTech for custom ERP software delivery by Jan 15 with a $500,000 budget. "
        "CloudTech delivered on March 1 (45 days late) with 3 critical unresolved database bugs causing $120,000 in operational losses. "
        "Client refuses to pay final $150,000 milestone invoice and threatens lawsuit."
    )

    # =========================================================================
    # AGENT 1: Strict Corporate Litigation Attorney (Branch: main)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[AGENT 1] Strict Corporate Litigation Attorney (qwen/qwen3.8-27b, temp=0.1)")
    print("=" * 80)

    @evo.track(repo=repo, name="litigation-attorney", model="qwen/qwen3.8-27b", temperature=0.1)
    def run_attorney(facts: str):
        """Strict Corporate Litigation Attorney. Identifies breaches, demands remedy, asserts claims."""
        return client.chat.completions.create(
            model="qwen/qwen3.8-27b",
            temperature=0.1,
            max_tokens=1500,
            messages=[
                {"role": "system", "content": "You are a Senior Corporate Litigation Attorney. Analyze dispute facts strictly under contract law. Demand cure and assert statutory set-off."},
                {"role": "user", "content": facts},
            ],
        )

    res1 = run_attorney(dispute_facts)
    exec1 = getattr(run_attorney, "last_execution")
    print(f"  -> Execution Latency: {exec1.duration_ms} ms | Tokens: {exec1.tokens.total_tokens}")

    commit1 = repo.commit(
        message="feat(agent): deploy strict legal analyst snapshot",
        tags=["v1.0-legal"],
        metadata={"agent": "litigation_attorney", "temp": "0.1"},
    )
    print(f"  -> Committed Snapshot 1: {commit1.id[:8]} on main")

    # =========================================================================
    # AGENT 2: Strategic Commercial Negotiator (Branch: exp-negotiator)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[AGENT 2] Strategic Commercial Mediator (qwen/qwen3.6-27b, temp=0.7)")
    print("=" * 80)

    # Create & checkout branch via CLI
    run_cli("branch", "exp-negotiator")
    run_cli("checkout", "exp-negotiator")
    repo = evo.open(SANDBOX_DIR)

    @evo.track(repo=repo, name="commercial-mediator", model="qwen/qwen3.6-27b", temperature=0.7)
    def run_mediator(facts: str):
        """Strategic Commercial Negotiator. Proposes creative win-win solutions and preserves business partnerships."""
        return client.chat.completions.create(
            model="qwen/qwen3.6-27b",
            temperature=0.7,
            max_tokens=1500,
            messages=[
                {"role": "system", "content": "You are a Master Commercial Negotiator. Propose win-win restructuring to avoid costly litigation."},
                {"role": "user", "content": facts},
            ],
        )

    res2 = run_mediator(dispute_facts)
    exec2 = getattr(run_mediator, "last_execution")
    print(f"  -> Execution Latency: {exec2.duration_ms} ms | Tokens: {exec2.tokens.total_tokens}")

    commit2 = repo.commit(
        message="feat(agent): deploy commercial mediator snapshot",
        tags=["v2.0-mediator"],
        metadata={"agent": "commercial_mediator", "temp": "0.7"},
    )
    print(f"  -> Committed Snapshot 2: {commit2.id[:8]} on exp-negotiator")

    # =========================================================================
    # AGENT 3: Lead Software Systems Architect (Branch: exp-tech-remediation)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[AGENT 3] Lead Software Systems Architect (qwen/qwen3.8-27b, temp=0.2)")
    print("=" * 80)

    run_cli("branch", "exp-tech-remediation")
    run_cli("checkout", "exp-tech-remediation")
    repo = evo.open(SANDBOX_DIR)

    @evo.track(repo=repo, name="tech-architect", model="qwen/qwen3.8-27b", temperature=0.2)
    def run_architect(facts: str):
        """Lead Systems Architect. Diagnoses technical failures, database concurrency bugs, and proposes remediation SLA."""
        return client.chat.completions.create(
            model="qwen/qwen3.8-27b",
            temperature=0.2,
            max_tokens=1500,
            messages=[
                {"role": "system", "content": "You are a Principal Software Architect. Focus exclusively on technical root-cause remediation, database patch validation, and uptime SLAs."},
                {"role": "user", "content": facts},
            ],
        )

    res3 = run_architect(dispute_facts)
    exec3 = getattr(run_architect, "last_execution")
    print(f"  -> Execution Latency: {exec3.duration_ms} ms | Tokens: {exec3.tokens.total_tokens}")

    commit3 = repo.commit(
        message="feat(agent): deploy technical remediation architect snapshot",
        tags=["v3.0-tech-arch"],
        metadata={"agent": "tech_architect", "temp": "0.2"},
    )
    print(f"  -> Committed Snapshot 3: {commit3.id[:8]} on exp-tech-remediation")

    # =========================================================================
    # AGENT 4: Chief Financial & Risk Officer (Branch: exp-risk-audit)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[AGENT 4] Chief Financial & Risk Officer (qwen/qwen3.8-27b, temp=0.3)")
    print("=" * 80)

    run_cli("branch", "exp-risk-audit")
    run_cli("checkout", "exp-risk-audit")
    repo = evo.open(SANDBOX_DIR)

    @evo.track(repo=repo, name="risk-officer", model="qwen/qwen3.8-27b", temperature=0.3)
    def run_risk_officer(facts: str):
        """Chief Financial & Risk Officer. Quantifies balance-sheet exposure, insurance indemnification, and liability caps."""
        return client.chat.completions.create(
            model="qwen/qwen3.8-27b",
            temperature=0.3,
            max_tokens=1500,
            messages=[
                {"role": "system", "content": "You are a Chief Financial & Risk Officer (CFRO). Quantify financial loss exposure ($120k vs $150k escrow), legal cost-benefit ratio, and balance sheet protection."},
                {"role": "user", "content": facts},
            ],
        )

    res4 = run_risk_officer(dispute_facts)
    exec4 = getattr(run_risk_officer, "last_execution")
    print(f"  -> Execution Latency: {exec4.duration_ms} ms | Tokens: {exec4.tokens.total_tokens}")

    commit4 = repo.commit(
        message="feat(agent): deploy CFRO risk audit snapshot",
        tags=["v4.0-risk-audit"],
        metadata={"agent": "risk_officer", "temp": "0.3"},
    )
    print(f"  -> Committed Snapshot 4: {commit4.id[:8]} on exp-risk-audit")

    # =========================================================================
    # STEP 5: Framework Adapters Verification (Zero-Dependency Introspection)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[STEP 5] Universal Framework Adapters Verification")
    print("=" * 80)

    # Test dummy objects for each framework
    class DummyChain:
        first = "Summarize: {input}"
        middle = ["format_tool", "search_tool"]
        last = {"model": "gpt-4o", "temperature": 0.2}

    class DummyIndex:
        index_struct = {"chunk_size": 512, "similarity_top_k": 5}
        vector_store = "Pinecone"

    class DummyCrew:
        agents = [{"name": "researcher", "role": "Fact Checker"}]
        tasks = [{"name": "audit_task", "description": "Verify claims"}]

    m_lc = LangChainAdapter.from_langchain(DummyChain(), name="langchain-adapter-demo")
    m_li = LlamaIndexAdapter.from_llamaindex(DummyIndex(), name="llamaindex-adapter-demo")
    m_cr = CrewAIAdapter.from_crewai(DummyCrew(), name="crewai-adapter-demo")
    m_oa = OpenAIAdapter.from_openai({"model": "gpt-4o", "usage": {"prompt_tokens": 50, "completion_tokens": 100}}, name="openai-demo")
    m_an = AnthropicAdapter.from_anthropic({"model": "claude-3-5-sonnet-20241022", "usage": {"input_tokens": 60, "output_tokens": 120}}, name="anthropic-demo")

    print(f"  -> LangChain Adapter:  Generated manifest with {len(m_lc.artifacts.prompts)} prompts, {len(m_lc.artifacts.tools)} tools.")
    print(f"  -> LlamaIndex Adapter: Generated retrieval config (chunk_size={m_li.artifacts.retrieval[0].chunk_size if m_li.artifacts.retrieval else 0}).")
    print(f"  -> CrewAI Adapter:    Generated {len(m_cr.artifacts.tools)} task tools.")
    print(f"  -> Direct API Adapters: Successfully converted OpenAI ({m_oa.artifacts.model_config.model if m_oa.artifacts.model_config else 'N/A'}) & Anthropic ({m_an.artifacts.model_config.model if m_an.artifacts.model_config else 'N/A'}).")

    # =========================================================================
    # STEP 6: Semantic Evaluation (LLM-as-a-Judge)
    # =========================================================================
    print("\n" + "=" * 80)
    print("[STEP 6] Running Impartial LLM Judge on All 4 Agents (qwen/qwen3.8-27b)")
    print("=" * 80)

    judge = evo.SemanticEvaluator(judge_fn=groq_judge_fn, judge_model="qwen/qwen3.8-27b (Groq LPU)")

    agents = [
        ("Agent 1 (Litigation Attorney)", exec1, "Strict Litigation Attorney"),
        ("Agent 2 (Commercial Mediator)", exec2, "Commercial Mediator"),
        ("Agent 3 (Tech Architect)", exec3, "Technical Solutions Architect"),
        ("Agent 4 (Risk & Finance Officer)", exec4, "Chief Risk & Finance Officer"),
    ]

    reports = []
    for name, exc, role in agents:
        print(f"  Evaluating {name}...")
        rep = judge.evaluate(
            inputs=exc.inputs,
            outputs=exc.outputs,
            agent_prompt=f"Role: {role}. Address enterprise software dispute facts.",
            execution_id=exc.id,
            commit_id=exc.commit_id,
        )
        reports.append((name, rep))
        repo.save_evaluation(
            evo.EvaluationResult(
                commit_id=exc.commit_id,
                execution_id=exc.id,
                overall_score=rep.overall_score,
                scores={k: evo.EvaluationScore(name=k, score=v.score, details=v.reasoning) for k, v in rep.dimensions.items()},
            )
        )

    # Print 4-Agent Scorecard
    print("\n" + "=" * 80)
    print("4-AGENT SEMANTIC EVALUATION SCORECARD")
    print("=" * 80)
    print(f"{'Dimension':<24} | {'Agent 1 (Legal)':<15} | {'Agent 2 (Mediate)':<17} | {'Agent 3 (Tech)':<15} | {'Agent 4 (Risk)':<15}")
    print("-" * 95)

    all_dims = ["accuracy", "helpfulness", "instruction_following", "safety"]
    for dim in all_dims:
        scores = []
        for _, rep in reports:
            d = rep.dimensions.get(dim)
            scores.append(f"{d.raw_score:.1f}/10 ({d.score*100:.0f}%)" if d else "N/A")
        print(f"{dim.replace('_', ' ').title():<24} | {scores[0]:<15} | {scores[1]:<17} | {scores[2]:<15} | {scores[3]:<15}")

    print("-" * 95)
    overall_strs = [f"{rep.overall_score*100:.1f}%" for _, rep in reports]
    print(f"{'OVERALL QUALITY SCORE':<24} | {overall_strs[0]:<15} | {overall_strs[1]:<17} | {overall_strs[2]:<15} | {overall_strs[3]:<15}")
    print("=" * 95)

    # =========================================================================
    # STEP 7: Comprehensive Go CLI Command Suite Execution
    # =========================================================================
    print("\n" + "=" * 80)
    print("[STEP 7] Comprehensive Go CLI Commands Execution Suite")
    print("=" * 80)

    print("\n--- [evo status] ---")
    print(run_cli("status"))

    print("\n--- [evo branch] ---")
    print(run_cli("branch"))

    print("\n--- [evo log --oneline] ---")
    print(run_cli("log", "--oneline"))

    print("\n--- [evo diff main exp-tech-remediation] ---")
    print(run_cli("diff", "main", "exp-tech-remediation")[:500] + "...\n[diff output truncated for readability]")

    print("\n--- [evo manifest show] ---")
    print(run_cli("manifest", "show"))

    print("\n--- [evo execution list] ---")
    print(run_cli("execution", "list"))

    print("\n--- [evo replay (State Reconstruction of Commit on 'main')] ---")
    print(run_cli("replay", "main"))

    print("\n--- [evo evaluate --compare main exp-negotiator] ---")
    print(run_cli("evaluate", "--compare", "main", "exp-negotiator"))

    print("\n" + "=" * 80)
    print("FULL-SPECTRUM MASTER LAB EXECUTION COMPLETED SUCCESSFULLY!")
    print("=" * 80)


if __name__ == "__main__":
    main()
