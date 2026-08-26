# Evolution Sandbox Agent Demo

A runnable demonstration of how AI engineers use Evolution to automatically track, version, and evaluate AI agents.

## Quick Run

```bash
# Run the sandbox demo
python examples/sandbox_agent/demo.py
```

### With a Real Grok or OpenAI Key (Optional)

```bash
# For Grok (xAI):
export XAI_API_KEY="your-xai-key"
python examples/sandbox_agent/demo.py

# Or for OpenAI:
export OPENAI_API_KEY="your-openai-key"
python examples/sandbox_agent/demo.py
```

*(If no API key is provided, the script runs in high-fidelity mock mode so you can test all Evolution features locally without spending API credits).*

## What This Demo Demonstrates

1. **`@evolution.track`**: Introspects docstrings as system prompts, captures model parameters, and measures real-time latency.
2. **`evolution.manifest.json`**: Auto-generated and updated without manual JSON editing.
3. **Execution Telemetry**: Execution logs stored in `.evolution/executions/`.
4. **DAG Commit**: State committed to the Evolution Merkle tree.
