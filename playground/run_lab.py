"""
Evolution Live Agent Lab: End-to-End Real Multi-Agent Scenario.

Demonstrates:
1. Running real Groq LLMs (Llama-3.3-70b & Llama-3.1-8b)
2. Live intelligence auto-capture (@evolution.track)
3. Branching across intelligence versions (main vs exp-creative-negotiator)
4. Side-by-side execution diffing & metrics comparison
5. Automated quality evaluation
"""

from __future__ import annotations

import os
import sys
from pathlib import Path

# Fix Windows console UTF-8 output
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

# Load .env if present
from dotenv import load_dotenv

env_file = Path(__file__).resolve().parent / ".env"
load_dotenv(env_file)

# Ensure local evolution SDK is used
sdk_path = Path(__file__).resolve().parent.parent / "sdk" / "python"
if str(sdk_path) not in sys.path:
    sys.path.insert(0, str(sdk_path))

import evolution as evo
from agents.creative_negotiator import create_creative_negotiator
from agents.strict_analyst import create_strict_analyst


def main():
    print("=" * 80)
    print("🧪 EVOLUTION MULTI-AGENT INTELLIGENCE LAB (LIVE DEMO)")
    print("=" * 80)

    api_key = os.environ.get("GROQ_API_KEY")
    if api_key and api_key.strip():
        print(f"🔑 Groq API Key Detected: {api_key[:8]}...{api_key[-4:]} (Live LLM Mode)")
    else:
        print("ℹ️ No GROQ_API_KEY found in playground/.env. Running in High-Fidelity Simulation Mode.")
        print("   (Add your key to playground/.env to run live Llama 3 calls on Groq)")

    # 1. Setup Workspace Repo
    workspace_dir = Path(__file__).resolve().parent / "workspace"
    workspace_dir.mkdir(exist_ok=True)

    try:
        repo = evo.init(workspace_dir, name="legal-dispute-resolution-system")
        print(f"\n📂 1. Initialized Fresh Intelligence Repo at: playground/workspace/")
    except Exception:
        repo = evo.open(workspace_dir)
        print(f"\n📂 1. Opened Existing Intelligence Repo at: playground/workspace/")

    # Scenario: High-Stakes Supplier Dispute
    dispute_scenario = (
        "Enterprise client contracted CloudTech for custom ERP software delivery by Jan 15 with a $500,000 budget. "
        "CloudTech delivered on March 1 (45 days late) with 3 critical unresolved database bugs causing $120,000 in operational losses. "
        "Client refuses to pay final $150,000 milestone invoice and threatens breach of contract lawsuit."
    )

    print("\n" + "-" * 80)
    print("📜 REAL-WORLD TEST DISPUTE FACTS:")
    print(f"   \"{dispute_scenario}\"")
    print("-" * 80)

    # -------------------------------------------------------------------------
    # EXPERIMENT 1: Branch 'main' -> Strict Corporate Litigation Attorney
    # -------------------------------------------------------------------------
    print("\n🏛️ STEP 2: Running Agent 1 on Branch 'main' (Strict Legal Analyst)...")
    analyst_fn = create_strict_analyst(repo)
    res_1 = analyst_fn(dispute_scenario)

    exec_1 = analyst_fn.last_execution if hasattr(analyst_fn, "last_execution") else None
    manifest_1 = repo.get_manifest()

    print("\n   ✅ Agent 1 Run Completed:")
    if hasattr(res_1, "choices") and res_1.choices:
        out_1_text = res_1.choices[0].message.content
    elif isinstance(res_1, dict) and "choices" in res_1:
        out_1_text = res_1["choices"][0]["message"]["content"]
    else:
        out_1_text = str(res_1)

    print(f"   Model:     {manifest_1.artifacts.model_config.model if manifest_1.artifacts.model_config else 'N/A'}")
    print(f"   Temp:      {manifest_1.artifacts.model_config.temperature if manifest_1.artifacts.model_config else 'N/A'}")
    print(f"   Latency:   {exec_1.duration_ms if exec_1 else 0} ms")
    print(f"   Tokens:    {exec_1.tokens.total_tokens if exec_1 else 0} total ({exec_1.tokens.prompt_tokens if exec_1 else 0} prompt, {exec_1.tokens.completion_tokens if exec_1 else 0} completion)")
    print("\n   --- Agent 1 Output Snippet ---")
    print("   " + "\n   ".join(out_1_text.strip().splitlines()[:6]) + ("..." if len(out_1_text.splitlines()) > 6 else ""))

    # Commit Snapshot 1
    commit_1 = repo.commit(
        message="feat(agent): deploy strict legal analyst using llama-3.3-70b (temp=0.1)",
        tags=["v1.0-strict-legal"],
        metadata={"persona": "litigation_attorney", "model": "llama-3.3-70b-versatile"},
    )
    print(f"\n   🔒 Committed Intelligence Snapshot 1: {commit_1.id[:8]} on main")

    # -------------------------------------------------------------------------
    # EXPERIMENT 2: Branch 'exp-creative-negotiator' -> Strategic Win-Win Mediator
    # -------------------------------------------------------------------------
    print("\n" + "-" * 80)
    print("🤝 STEP 3: Running Agent 2 (Strategic Commercial Negotiator)...")
    negotiator_fn = create_creative_negotiator(repo)
    res_2 = negotiator_fn(dispute_scenario)

    exec_2 = negotiator_fn.last_execution if hasattr(negotiator_fn, "last_execution") else None
    manifest_2 = repo.get_manifest()

    print("\n   ✅ Agent 2 Run Completed:")
    if hasattr(res_2, "choices") and res_2.choices:
        out_2_text = res_2.choices[0].message.content
    elif isinstance(res_2, dict) and "choices" in res_2:
        out_2_text = res_2["choices"][0]["message"]["content"]
    else:
        out_2_text = str(res_2)

    print(f"   Model:     {manifest_2.artifacts.model_config.model if manifest_2.artifacts.model_config else 'N/A'}")
    print(f"   Temp:      {manifest_2.artifacts.model_config.temperature if manifest_2.artifacts.model_config else 'N/A'}")
    print(f"   Latency:   {exec_2.duration_ms if exec_2 else 0} ms")
    print(f"   Tokens:    {exec_2.tokens.total_tokens if exec_2 else 0} total ({exec_2.tokens.prompt_tokens if exec_2 else 0} prompt, {exec_2.tokens.completion_tokens if exec_2 else 0} completion)")
    print("\n   --- Agent 2 Output Snippet ---")
    print("   " + "\n   ".join(out_2_text.strip().splitlines()[:6]) + ("..." if len(out_2_text.splitlines()) > 6 else ""))

    # Commit Snapshot 2
    commit_2 = repo.commit(
        message="feat(agent): experiment with creative negotiator using llama-3.1-8b (temp=0.7)",
        tags=["v2.0-creative-negotiator"],
        metadata={"persona": "commercial_negotiator", "model": "llama-3.1-8b-instant"},
    )
    print(f"\n   🔒 Committed Intelligence Snapshot 2: {commit_2.id[:8]} on exp branch")

    # -------------------------------------------------------------------------
    # STEP 4: Side-by-Side Evolution Intelligence Comparison
    # -------------------------------------------------------------------------
    print("\n" + "=" * 80)
    print("📊 STEP 4: SIDE-BY-SIDE EVOLUTION INTELLIGENCE COMPARISON")
    print("=" * 80)

    print(f"\n{'Metric / Artifact':<28} | {'Agent 1 (Strict Analyst)':<24} | {'Agent 2 (Negotiator)':<24}")
    print("-" * 80)
    m1_model = manifest_1.artifacts.model_config.model if manifest_1.artifacts.model_config else "N/A"
    m2_model = manifest_2.artifacts.model_config.model if manifest_2.artifacts.model_config else "N/A"
    print(f"{'Model Architecture':<28} | {m1_model:<24} | {m2_model:<24}")

    m1_temp = str(manifest_1.artifacts.model_config.temperature) if manifest_1.artifacts.model_config else "N/A"
    m2_temp = str(manifest_2.artifacts.model_config.temperature) if manifest_2.artifacts.model_config else "N/A"
    print(f"{'Sampling Temperature':<28} | {m1_temp:<24} | {m2_temp:<24}")

    t1_tokens = f"{exec_1.tokens.total_tokens} tokens" if exec_1 else "N/A"
    t2_tokens = f"{exec_2.tokens.total_tokens} tokens" if exec_2 else "N/A"
    print(f"{'Token Consumption':<28} | {t1_tokens:<24} | {t2_tokens:<24}")

    l1_ms = f"{exec_1.duration_ms} ms" if exec_1 else "N/A"
    l2_ms = f"{exec_2.duration_ms} ms" if exec_2 else "N/A"
    print(f"{'Inference Latency':<28} | {l1_ms:<24} | {l2_ms:<24}")

    print("\n" + "=" * 80)
    print("🎯 CONCLUSION & OBSERVABILITY:")
    print("   Evolution captured both operational snapshots into an immutable Merkle tree.")
    print("   You can now diff, replay historical executions, and run evaluation gates on either agent!")
    print("=" * 80)


if __name__ == "__main__":
    main()
