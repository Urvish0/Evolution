# Getting Started with Evolution

> **Zero to working AI version control in 5 minutes.**

This tutorial walks you through installing Evolution, tracking your first AI agent, committing an intelligence snapshot, and running a semantic evaluation — all in pure Python with zero external dependencies.

---

## Prerequisites

- Python 3.10 or higher
- A terminal (PowerShell, bash, zsh)
- An LLM API key (OpenAI, Groq, or Anthropic) — *optional, works with mock responses too*

---

## Step 1: Install the SDK (30 seconds)

```bash
pip install evolution-sdk
```

Verify it works:

```bash
evo-py --version
# Output: evolution-sdk 0.8.5
```

You now have:
- `import evolution as evo` — the Python SDK for use inside your code.
- `evo-py` — a standalone CLI you can run in your terminal.

---

## Step 2: Initialize a Repository (30 seconds)

Create a project folder and initialize an Evolution repository:

```python
import evolution as evo

repo = evo.Repository.init("./my-ai-project", name="customer-support-bot")
```

This creates the `.evolution/` directory inside `my-ai-project/`, containing:

```
my-ai-project/
├── .evolution/
│   ├── HEAD                      # Points to current branch (main)
│   ├── branches/                 # Branch references
│   ├── commits/                  # Immutable commit objects
│   ├── objects/                  # Content-addressable blob storage
│   ├── executions/               # Recorded AI execution telemetry
│   └── evaluations/              # Semantic evaluation reports
└── evolution.manifest.json       # Intelligence Manifest (auto-generated)
```

Or from the terminal:

```bash
evo-py init ./my-ai-project --name "customer-support-bot"
```

---

## Step 3: Track Your First AI Agent (2 minutes)

Create a file called `agent.py`:

```python
import evolution as evo

# Open the repository
repo = evo.Repository.open("./my-ai-project")

# Define your agent with @evo.track
@evo.track(repo=repo, model="gpt-4o", temperature=0.2)
def support_agent(question: str):
    """You are a helpful customer support specialist.
    Always be polite, concise, and suggest relevant documentation links."""

    # Option A: Real OpenAI call (if you have an API key)
    # from openai import OpenAI
    # client = OpenAI()
    # return client.chat.completions.create(
    #     model="gpt-4o",
    #     temperature=0.2,
    #     messages=[
    #         {"role": "system", "content": support_agent.__doc__},
    #         {"role": "user", "content": question}
    #     ]
    # )

    # Option B: Simulated response (works without an API key)
    return {
        "content": f"To answer '{question}': Please check our docs at docs.example.com.",
        "usage": {"prompt_tokens": 45, "completion_tokens": 22}
    }


# Run the agent
print("Running agent...")
result = support_agent("How do I reset my password?")
print(f"Response: {result['content']}")
```

**What happens when you run `python agent.py`:**

1. `@evo.track` extracts the docstring as the system prompt artifact.
2. It records `model="gpt-4o"` and `temperature=0.2` as model config artifacts.
3. It measures execution latency with high-precision timers.
4. It introspects the return value to extract output text and token counts.
5. It saves an execution telemetry record to `.evolution/executions/`.
6. It updates `evolution.manifest.json` with the captured intelligence state.

All of this happens automatically. You wrote one decorator.

---

## Step 4: Commit the Intelligence Snapshot (30 seconds)

Add this line at the end of `agent.py`:

```python
commit = repo.commit("feat: baseline customer support agent with gpt-4o")
print(f"Committed: {commit.id}")
```

Or from the terminal:

```bash
evo-py -C ./my-ai-project status
```

Output:

```
On branch main
Manifest: customer-support-bot (v1.0.0) | Model: gpt-4o | Prompts: 1 | Tools: 0
Recorded Executions: 1
Evaluations: 0

nothing to commit, working tree clean
```

---

## Step 5: Experiment on a Branch (1 minute)

Now suppose you want to test a cheaper, faster model. Create a branch and modify the agent:

```python
# experiment.py
import evolution as evo

repo = evo.Repository.open("./my-ai-project")

@evo.track(repo=repo, model="llama-3.3-70b", temperature=0.7)
def support_agent(question: str):
    """You are a friendly, casual support agent.
    Use simple language and emoji when appropriate."""

    return {
        "content": f"Hey! For '{question}' — just hit the 'Forgot Password' button on the login page.",
        "usage": {"prompt_tokens": 38, "completion_tokens": 19}
    }

result = support_agent("How do I reset my password?")
repo.commit("exp: test llama-3.3-70b with casual tone")
```

---

## Step 6: View History and Diff (30 seconds)

From the terminal, inspect what you have built:

```bash
# View commit history
evo-py -C ./my-ai-project log --oneline

# Output:
# b3f1a2c4 exp: test llama-3.3-70b with casual tone
# a7d92e81 feat: baseline customer support agent with gpt-4o
```

With the Go CLI (if installed), you can run semantic diffs:

```bash
# See what changed between the two commits
evo diff main exp-casual-model
```

Output:

```diff
--- a/evolution.manifest.json (main)
+++ b/evolution.manifest.json (exp-casual-model)
-  "model": "gpt-4o",
+  "model": "llama-3.3-70b",
-  "temperature": 0.2,
+  "temperature": 0.7,
-  "prompt": "You are a helpful customer support specialist..."
+  "prompt": "You are a friendly, casual support agent..."
```

---

## Step 7: Evaluate Quality with LLM-as-a-Judge (1 minute)

Use the built-in semantic evaluator to score your agent's output:

```python
import evolution as evo

evaluator = evo.SemanticEvaluator(
    api_key="your-groq-or-openai-api-key",
    judge_model="llama-3.3-70b-versatile",  # or "gpt-4o"
    base_url="https://api.groq.com/openai/v1",
)

report = evaluator.evaluate(
    prompt="How do I reset my password?",
    response="Hey! Just hit the 'Forgot Password' button on the login page.",
    instruction="You are a helpful customer support specialist.",
)

print(f"Overall Score: {report.overall_score:.0%}")
for dim in report.scores:
    print(f"  {dim.dimension}: {dim.score}/10 — {dim.reasoning}")
```

Output:

```
Overall Score: 85%
  accuracy: 8/10 — Correct information about password reset flow.
  helpfulness: 9/10 — Direct and actionable response.
  instruction_following: 7/10 — Tone is casual, not professional as instructed.
  safety: 10/10 — No harmful content.
```

---

## Step 8: Inspect Everything from the Terminal

```bash
# View execution telemetry
evo-py -C ./my-ai-project execution list

# Output:
# ID         COMMIT     TOKENS     DURATION     TIMESTAMP
# -----------------------------------------------------------------
# 4a2c9f10   a7d92e81   67         142ms        2026-09-03T21:10:00

# View evaluation reports
evo-py -C ./my-ai-project evaluate

# Validate your manifest against the Intelligence Manifest Spec v1.0
evo-py -C ./my-ai-project manifest validate
# Output: Manifest 'customer-support-bot' is valid according to Intelligence Manifest Spec v1.0.
```

---

## What You Just Built

In under 5 minutes, you:

1. **Initialized** an Evolution repository with a typed Intelligence Manifest.
2. **Tracked** an AI agent with a single decorator (`@evo.track`), automatically capturing its prompt, model config, tokens, and latency.
3. **Committed** an immutable, content-addressed intelligence snapshot to a Merkle DAG.
4. **Branched** to test an alternative model and prompt configuration.
5. **Diffed** two intelligence states to see semantic changes (model swap, prompt drift, parameter mutation).
6. **Evaluated** output quality with an LLM-as-a-Judge scoring rubric across accuracy, helpfulness, instruction-following, and safety.
7. **Inspected** execution telemetry and evaluation reports from the terminal.

---

## Next Steps

| Goal | Resource |
|------|----------|
| Full CLI command reference | [USAGE.md](../USAGE.md) |
| Technical architecture and internals | [DOCUMENTATION.md](../DOCUMENTATION.md) |
| Intelligence Manifest Specification v1.0 | [spec/intelligence-manifest-v1.0.md](../spec/intelligence-manifest-v1.0.md) |
| Framework adapters (LangChain, CrewAI, etc.) | [sdk/python/README.md](../sdk/python/README.md) |
| LLM-as-a-Judge evaluation guide | [playground/evaluate_agents.py](../playground/evaluate_agents.py) |
| Contributing to Evolution | [CONTRIBUTING.md](../CONTRIBUTING.md) |

---

## Questions or Feedback?

- Open an issue on [GitHub](https://github.com/Urvish0/Evolution/issues)
- Start a discussion in [GitHub Discussions](https://github.com/Urvish0/Evolution/discussions)
