# tests/adapters/test_inbound.py
import pytest
import asyncio
from src.core.models.message import AgentRequest, AgentResponse
from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import MockPubSubPublisherAdapter
from src.adapters.inbound.pubsub_subscriber import BackgroundSubscriber
import json

@pytest.mark.asyncio
async def test_subscriber_processing():
    agent = ADKAgentAdapter()
    publisher = MockPubSubPublisherAdapter()
    subscriber = BackgroundSubscriber(agent_service=agent, publisher=publisher)

    # Simulate receiving a message
    req_json = json.dumps({
        "task_id": "tsk_123",
        "project_name": "Test",
        "github_url": "https://github.com/test/test",
        "submission_text": "text",
        "judging_rubric": "rubric",
        "scoring_criteria": [{"name": "Innovate", "weight": 1.0, "max_score": 10}]
    })

    await subscriber.process_raw_message("123", req_json)

    assert len(publisher.published_messages) == 1
    assert publisher.published_messages[0].task_id == "tsk_123"
