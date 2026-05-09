# tests/adapters/test_outbound.py
import pytest
from src.core.models.message import AgentRequest, AgentResponse
from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import MockPubSubPublisherAdapter

@pytest.mark.asyncio
async def test_adk_agent_adapter():
    agent = ADKAgentAdapter()
    req = AgentRequest(message_id="1", text="test")
    res = await agent.process_message(req)
    assert res.original_message_id == "1"
    # Basic check to ensure it returns something
    assert len(res.response_text) > 0

@pytest.mark.asyncio
async def test_mock_publisher():
    pub = MockPubSubPublisherAdapter()
    res = AgentResponse(original_message_id="1", response_text="test")
    await pub.publish(res)
    assert len(pub.published_messages) == 1
