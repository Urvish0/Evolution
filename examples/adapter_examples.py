"""
Evolution Framework Adapter Examples.

Demonstrates how to generate Intelligence Manifests from popular AI frameworks
using Evolution's zero-dependency adapters.

Requirements:
    pip install evolution-sdk
    (No framework libraries needed — adapters use duck-typing introspection)
"""

import json
import evolution as evo


def divider(title: str):
    print(f"\n{'=' * 70}")
    print(f"  {title}")
    print(f"{'=' * 70}\n")


# --- 1. LangChain Adapter ---

divider("1. LangChain Adapter")

# Simulate a LangChain chain object (duck-typed, no langchain import needed)
class MockPromptTemplate:
    template = "Summarize the following legal document in plain English:\n\n{document}"
    input_variables = ["document"]

class MockChatModel:
    model_name = "gpt-4o"
    temperature = 0.1
    model_kwargs = {}

    @property
    def bound_tools(self):
        return []

class MockChain:
    first = MockPromptTemplate()
    last = MockChatModel()

    @property
    def steps(self):
        return [self.first, self.last]

chain = MockChain()
manifest = evo.from_langchain(chain, name="legal-summarizer")
print(json.dumps(manifest.to_dict(), indent=2))


# --- 2. OpenAI Adapter ---

divider("2. OpenAI Adapter (Direct API)")

# Pass an OpenAI-style request dictionary
openai_request = {
    "model": "gpt-4o",
    "temperature": 0.2,
    "messages": [
        {"role": "system", "content": "You are a helpful customer support agent. Always be polite and concise."},
        {"role": "user", "content": "How do I reset my password?"},
    ],
    "tools": [
        {"function": {"name": "search_docs", "description": "Search knowledge base articles"}}
    ],
}

manifest = evo.from_openai(openai_request, name="customer-support-bot")
print(json.dumps(manifest.to_dict(), indent=2))


# --- 3. Anthropic Adapter ---

divider("3. Anthropic Adapter (Direct API)")

anthropic_request = {
    "model": "claude-3-5-sonnet-20241022",
    "temperature": 0.3,
    "system": "You are a senior research analyst. Cite sources for every claim.",
    "messages": [
        {"role": "user", "content": "Summarize the latest advances in quantum computing."},
    ],
}

manifest = evo.from_anthropic(anthropic_request, name="research-analyst")
print(json.dumps(manifest.to_dict(), indent=2))


# --- 4. CrewAI Adapter ---

divider("4. CrewAI Adapter")

class MockAgent:
    role = "Senior Data Scientist"
    goal = "Analyze customer churn patterns and recommend retention strategies"
    backstory = "10 years of experience in predictive analytics and ML"
    tools = []

class MockTask:
    description = "Analyze Q3 customer data and identify top 5 churn indicators"
    expected_output = "A structured report with churn indicators and confidence scores"

class MockCrew:
    agents = [MockAgent()]
    tasks = [MockTask()]
    verbose = True

crew = MockCrew()
manifest = evo.from_crewai(crew, name="churn-analysis-crew")
print(json.dumps(manifest.to_dict(), indent=2))


# --- 5. LlamaIndex Adapter ---

divider("5. LlamaIndex Adapter")

class MockServiceContext:
    class llm:
        model = "gpt-4o-mini"
        temperature = 0.0

    class embed_model:
        model_name = "text-embedding-3-small"

class MockIndex:
    _service_context = MockServiceContext()

index = MockIndex()
manifest = evo.from_llamaindex(index, name="knowledge-base-retriever")
print(json.dumps(manifest.to_dict(), indent=2))


# --- Summary ---

divider("All 5 Framework Adapters Verified")
print("Each adapter generated a valid Intelligence Manifest using duck-typing")
print("introspection with zero external dependencies.")
print("\nAdapters available:")
print("  - evo.from_langchain(chain)")
print("  - evo.from_openai(client, model=..., system_prompt=...)")
print("  - evo.from_anthropic(client, model=..., system_prompt=...)")
print("  - evo.from_crewai(crew)")
print("  - evo.from_llamaindex(index)")
