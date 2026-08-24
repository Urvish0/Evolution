"""
Unit tests for framework adapters (LangChain, LlamaIndex, CrewAI, OpenAI, Anthropic).
"""

import unittest
from types import SimpleNamespace

from evolution import (
    from_anthropic,
    from_crewai,
    from_langchain,
    from_llamaindex,
    from_openai,
)
from evolution.models.artifacts import (
    MemoryArtifact,
    ModelConfigArtifact,
    PromptArtifact,
    RetrievalArtifact,
    ToolArtifact,
)


class TestFrameworkAdapters(unittest.TestCase):

    def test_langchain_adapter(self):
        # Mock LangChain AgentExecutor / Chain
        mock_chain = SimpleNamespace(
            llm=SimpleNamespace(model_name="gpt-4o", temperature=0.3, max_tokens=2048),
            prompt=SimpleNamespace(template="You are a legal advisor. Answer: {query}"),
            tools=[
                SimpleNamespace(name="case_search", description="Searches case law"),
                SimpleNamespace(name="statute_lookup", description="Looks up statutes"),
            ],
            memory=SimpleNamespace(__class__=type("ConversationBufferMemory", (), {})),
        )

        manifest = from_langchain(mock_chain, name="legal-chain")
        manifest.validate()

        self.assertEqual(manifest.name, "legal-chain")
        self.assertEqual(manifest.artifacts.model_config.model, "gpt-4o")
        self.assertEqual(manifest.artifacts.model_config.temperature, 0.3)
        self.assertEqual(len(manifest.artifacts.prompts), 1)
        self.assertEqual(len(manifest.artifacts.tools), 2)
        self.assertEqual(len(manifest.artifacts.memory), 1)
        self.assertEqual(manifest.artifacts.memory[0].strategy, "buffer_window")

    def test_llamaindex_adapter(self):
        # Mock LlamaIndex Index with Chroma vector store and retriever
        mock_index = SimpleNamespace(
            storage_context=SimpleNamespace(
                vector_store=type("ChromaVectorStore", (), {})()
            ),
            retriever=SimpleNamespace(similarity_top_k=8),
            service_context=SimpleNamespace(
                node_parser=SimpleNamespace(chunk_size=1024),
                llm=SimpleNamespace(model="gpt-4o-mini", temperature=0.1),
            ),
        )

        manifest = from_llamaindex(mock_index, name="legal-rag")
        manifest.validate()

        self.assertEqual(manifest.name, "legal-rag")
        self.assertEqual(len(manifest.artifacts.retrieval), 1)
        ret = manifest.artifacts.retrieval[0]
        self.assertEqual(ret.source, "chroma")
        self.assertEqual(ret.chunk_size, 1024)
        self.assertEqual(ret.top_k, 8)
        self.assertEqual(manifest.artifacts.model_config.model, "gpt-4o-mini")

    def test_crewai_adapter(self):
        # Mock CrewAI Crew with 2 agents and tasks
        mock_crew = SimpleNamespace(
            agents=[
                SimpleNamespace(
                    role="Legal Researcher",
                    goal="Find relevant precedent",
                    backstory="Experienced in appellate litigation.",
                    tools=[SimpleNamespace(name="westlaw_search", description="Search Westlaw")],
                    llm=SimpleNamespace(model="claude-3.5-sonnet"),
                ),
                SimpleNamespace(
                    role="Legal Drafter",
                    goal="Draft court briefs",
                    backstory="Specialized in clear statutory argument.",
                    tools=[],
                    llm=SimpleNamespace(model="gpt-4o"),
                ),
            ],
            tasks=[SimpleNamespace(description="Research", expected_output="Cases")],
            process="sequential",
        )

        manifest = from_crewai(mock_crew, name="appellate-crew")
        manifest.validate()

        self.assertEqual(manifest.name, "appellate-crew")
        self.assertEqual(len(manifest.artifacts.prompts), 2)
        self.assertEqual(len(manifest.artifacts.tools), 1)
        self.assertEqual(manifest.metadata["crewai"]["agent_count"], 2)
        self.assertEqual(manifest.metadata["crewai"]["process"], "sequential")

    def test_openai_adapter(self):
        payload = {
            "model": "gpt-4o",
            "temperature": 0.2,
            "messages": [
                {"role": "system", "content": "You are a legal assistant."},
                {"role": "user", "content": "Explain strict liability."},
            ],
            "tools": [
                {"type": "function", "function": {"name": "lookup_case", "description": "Lookup case"}}
            ],
        }

        manifest = from_openai(payload, name="openai-agent")
        manifest.validate()

        self.assertEqual(manifest.name, "openai-agent")
        self.assertEqual(manifest.artifacts.model_config.model, "gpt-4o")
        self.assertEqual(manifest.artifacts.model_config.provider, "openai")
        self.assertEqual(len(manifest.artifacts.prompts), 2)
        self.assertEqual(len(manifest.artifacts.tools), 1)

    def test_anthropic_adapter(self):
        payload = {
            "model": "claude-3.5-sonnet",
            "temperature": 0.5,
            "system": "You are an Anthropic AI assistant.",
            "tools": [{"name": "web_search", "description": "Search web"}],
        }

        manifest = from_anthropic(payload, name="anthropic-agent")
        manifest.validate()

        self.assertEqual(manifest.name, "anthropic-agent")
        self.assertEqual(manifest.artifacts.model_config.model, "claude-3.5-sonnet")
        self.assertEqual(manifest.artifacts.model_config.provider, "anthropic")
        self.assertEqual(len(manifest.artifacts.prompts), 1)
        self.assertEqual(len(manifest.artifacts.tools), 1)


if __name__ == "__main__":
    unittest.main()
