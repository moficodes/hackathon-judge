# tests/adapters/test_inbound.py
import pytest
import asyncio
from src.core.models.message import AgentResponse
from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import MockPubSubPublisherAdapter
from src.adapters.inbound.pubsub_subscriber import BackgroundSubscriber

@pytest.mark.asyncio
async def test_subscriber_processing():
    agent = ADKAgentAdapter()
    publisher = MockPubSubPublisherAdapter()
    subscriber = BackgroundSubscriber(agent_service=agent, publisher=publisher)
    
    # Simulate receiving a message
    await subscriber.process_raw_message("123", "Hello World")
    
    assert len(publisher.published_messages) == 1
    assert publisher.published_messages[0].original_message_id == "123"
    assert "Hello World" in publisher.published_messages[0].response_text
