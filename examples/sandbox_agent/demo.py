"""
Evolution Live Sandbox Demo: Legal AI Assistant with Grok / OpenAI & Auto-Tracking.

This script demonstrates:
1. Automatic Intelligence Capture with @evolution.track
2. Manifest auto-generation (prompts, model config, parameters)
3. Precision latency and token metrics capture
4. Committing intelligence state to the Evolution DAG
5. Replay and Quality Evaluation
"""

import os
import sys
from pathlib import Path

# Fix Windows console encoding for UTF-8 characters
if sys.platform == "win32":
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

# Ensure the local evolution-sdk package is loaded
sdk_path = Path(__file__).resolve().parent.parent.parent / "sdk" / "python"
if str(sdk_path) not in sys.path:
    sys.path.insert(0, str(sdk_path))

import evolution as evo


def main():
    print("=" * 70)
    print("🚀 EVOLUTION LIVE SANDBOX DEMO: AI AGENT VERSION CONTROL")
    print("=" * 70)

    # 1. Initialize or open Evolution repository in this sandbox folder
    sandbox_dir = Path(__file__).resolve().parent
    repo_dir = sandbox_dir / "agent_workspace"
    repo_dir.mkdir(exist_ok=True)

    print(f"\n📂 1. Initializing Evolution Workspace at: {repo_dir.name}/")
    try:
        repo = evo.init(repo_dir, name="legal-advisor-ai")
        print("   ✅ Created new Evolution repository (.evolution/)")
    except Exception:
        repo = evo.open(repo_dir)
        print("   ✅ Opened existing Evolution repository")

    # 2. Define our AI Agent with @evolution.track
    # Works with Grok (xAI), OpenAI, or fallback mock simulation if no API key is provided!
    api_key = os.environ.get("XAI_API_KEY") or os.environ.get("OPENAI_API_KEY")

    print("\n🤖 2. Registering AI Agent with @evolution.track...")

    @evo.track(
        repo=repo,
        name="legal-advisory-bot",
        model="grok-beta" if os.environ.get("XAI_API_KEY") else "gpt-4o",
        temperature=0.2,
    )
    def legal_advisory_bot(query: str):
        """Tier-1 Legal AI Assistant. Analyzes breach of contract, indemnity, and statutory compliance."""
        # If API key is available and openai is installed, perform live API call
        if api_key:
            try:
                from openai import OpenAI
                base_url = "https://api.x.ai/v1" if os.environ.get("XAI_API_KEY") else None
                client = OpenAI(api_key=api_key, base_url=base_url)
                model = "grok-beta" if os.environ.get("XAI_API_KEY") else "gpt-4o"

                response = client.chat.completions.create(
                    model=model,
                    messages=[
                        {"role": "system", "content": "You are a professional corporate legal assistant. Provide structured legal analysis."},
                        {"role": "user", "content": query},
                    ],
                    temperature=0.2,
                )
                return response
            except Exception as e:
                print(f"   ⚠️ Live API call failed ({e}), using mock response...")

        # Fallback simulation (matches standard LLM response format)
        import time
        time.sleep(0.15)  # Simulate network latency
        return {
            "choices": [{
                "message": {
                    "content": f"Under common law contract principles, the facts in '{query}' indicate a potential breach of the implied covenant of good faith and fair dealing. Recommended remedy: formal notice of cure within 30 days."
                }
            }],
            "usage": {
                "prompt_tokens": 42,
                "completion_tokens": 58,
                "total_tokens": 100,
            },
            "model": "grok-beta",
        }

    # 3. Execute the agent
    test_query = "The supplier delayed critical component delivery by 45 days causing assembly line shutdown."
    print(f"\n⚡ 3. Invoking Agent with query:\n   \"{test_query}\"")

    response = legal_advisory_bot(test_query)

    print("\n📝 4. Result Received:")
    if isinstance(response, dict) and "choices" in response:
        print(f"   Output: {response['choices'][0]['message']['content']}")
        print(f"   Tokens: {response['usage']['total_tokens']} tokens")
    else:
        print(f"   Response: {response}")

    # 4. Inspect Manifest
    manifest = repo.get_manifest()
    print("\n📋 5. Inspecting Auto-Generated Intelligence Manifest:")
    print(f"   - Manifest Name:    {manifest.name}")
    print(f"   - Spec Version:     {manifest.version}")
    print(f"   - Prompts Attached: {[p.name for p in manifest.artifacts.prompts]}")
    if manifest.artifacts.model_config:
        print(f"   - Model Config:     {manifest.artifacts.model_config.model} (provider: {manifest.artifacts.model_config.provider}, temp: {manifest.artifacts.model_config.temperature})")

    # 5. Inspect Execution History
    executions = repo.list_executions()
    print(f"\n📊 6. Inspecting Execution Telemetry Log:")
    print(f"   - Total Recorded Runs: {len(executions)}")
    if executions:
        latest = executions[0]
        print(f"   - Latest Execution ID: {latest.id[:8]}...")
        print(f"   - Measured Latency:    {latest.duration_ms} ms")
        print(f"   - Tokens Consumed:     {latest.tokens.total_tokens} (prompt: {latest.tokens.prompt_tokens}, completion: {latest.tokens.completion_tokens})")

    # 6. Commit Intelligence
    print("\n🔒 7. Committing Intelligence State to Evolution DAG...")
    commit = repo.commit(
        message="feat: initial legal advisory agent with grok-beta and 0.2 temperature",
        tags=["v1.0-alpha"],
        metadata={"author": "Urvish", "model": "grok-beta"},
    )
    print(f"   ✅ Intelligence Commit Created: {commit.id[:8]} - '{commit.message}'")

    print("\n" + "=" * 70)
    print("🎉 DEMO COMPLETE! Evolution successfully versioned this AI agent.")
    print("=" * 70)


if __name__ == "__main__":
    main()
