# tests/adapters/test_outbound.py
import pytest
from unittest.mock import patch, MagicMock
from src.core.models.message import AgentRequest, AgentResponse, ScoringCriteria
from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import MockPubSubPublisherAdapter, PubSubPublisherAdapter

@pytest.mark.asyncio
async def test_adk_agent_adapter():
    agent = ADKAgentAdapter()
    req = AgentRequest(
        task_id="tsk_1",
        project_name="Test Project",
        github_url="https://github.com/test/test",
        submission_text="Here is my code",
        judging_rubric="Be nice",
        scoring_criteria=[ScoringCriteria(name="Innovation", weight=0.5)]
    )
    res = await agent.process_message(req)
    assert res.task_id == "tsk_1"
    # Basic check to ensure we get a status
    assert res.status in ["success", "error"]

@pytest.mark.asyncio
async def test_mock_publisher():
    pub = MockPubSubPublisherAdapter()
    res = AgentResponse(task_id="tsk_1", status="success")
    await pub.publish(res)
    assert len(pub.published_messages) == 1

@pytest.mark.asyncio
@patch('src.adapters.outbound.pubsub_publisher.pubsub_v1.PublisherClient')
async def test_real_publisher_no_attribute_errors(mock_publisher_client_cls):
    # Setup mock publisher
    mock_publisher_instance = MagicMock()
    mock_publisher_client_cls.return_value = mock_publisher_instance
    mock_future = MagicMock()
    mock_future.result.return_value = "msg_123"
    mock_publisher_instance.publish.return_value = mock_future
    mock_publisher_instance.topic_path.return_value = "projects/test/topics/test"

    pub = PubSubPublisherAdapter(project_id="test-project", topic_id="test-topic")
    
    # Create response model
    res = AgentResponse(
        task_id="tsk_123",
        status="success",
        overall_comments="Tested properly."
    )
    
    # This should run without throwing any exceptions like AttributeError
    await pub.publish(res)
    
    # Verify the mock was called correctly
    mock_publisher_instance.publish.assert_called_once()
    args, kwargs = mock_publisher_instance.publish.call_args
    assert args[0] == "projects/test/topics/test"
    assert b"tsk_123" in kwargs['data']
